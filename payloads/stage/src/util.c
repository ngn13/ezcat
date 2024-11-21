#ifdef _WIN32
// clang-format off
#include <winsock2.h>
#include <ws2tcpip.h>
#include <ws2def.h>
#include <windows.h>
// clang-format on
#else
#include <arpa/inet.h>
#include <netdb.h>
#include <netinet/in.h>
#include <sys/socket.h>
#endif

#include <errno.h>
#include <stdarg.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <strings.h>
#include <sys/types.h>
#include <unistd.h>

#include "packet.h"
#include "util.h"

void print_debug(const char *func, const char *msg, ...) {
  if (!STAGE_DEBUG)
    return;

  va_list args;
  va_start(args, msg);

  printf("[debug] (%s) ", func);
  vprintf(msg, args);
  printf("\n");

  va_end(args);
}

void print_debug_dump(const char *func, char *buf, uint64_t size) {
  if (!STAGE_DEBUG)
    return;

  uint64_t i          = 0;
  int      title_size = printf("------------- %s dumping %p -------------\n", func, buf);

  for (i = 0; i < size; i++) {
    if (i % 10 == 0 || i == 0)
      printf(i == 0 ? "0x%08lx: " : "\n0x%08lx: ", i);

    printf("0x%02x ", (unsigned char)buf[i]);

    if (size - 1 == i)
      printf("\n");
  }

  for (i = 0; i < title_size - 1; i++) // -1 for \n
    printf("-");
  printf("\n");
}

int randint(int min, int max) {
  return min + rand() % (max + 1 - min);
}

void *copy_to(void *dst, void *src, uint64_t size) {
  memcpy(dst, src, size);
  return dst + size;
}

void *copy_from(void *dst, void *src, uint64_t size) {
  memcpy(dst, src, size);
  return src + size;
}

#ifdef _WIN32
// https://ao.ngn.tf/questions/8776103/how-to-redirect-fprintf-output-to-c-socket
void fprintfsock(SOCKET sock, const char *fmt, ...) {
  va_list args;
  va_start(args, fmt);

  int  len = vsnprintf(0, 0, fmt, args);
  char buf[len + 1];
  bzero(buf, len + 1);

  va_start(args, fmt);
  vsnprintf(buf, len + 1, fmt, args);

  send(sock, buf, len, 0);
}
#endif

bool resolve(struct sockaddr *saddr, char *addr, uint16_t port) {
  if (NULL == saddr || NULL == addr) {
    errno = EINVAL;
    return false;
  }

  struct addrinfo *res = NULL, *cur = NULL;
  bool             ret = false;

  if (getaddrinfo(addr, NULL, NULL, &res) != 0)
    goto end;

  if (NULL == res)
    goto end;

  for (cur = res; cur != NULL; cur = cur->ai_next) {
    if (AF_INET == cur->ai_family || AF_INET6 == cur->ai_family)
      break;
  }

  if (NULL == cur)
    goto end;

  memcpy(saddr, cur->ai_addr, sizeof(struct sockaddr));

  if (port == 0) {
    ret = true;
    goto end;
  }

  switch (saddr->sa_family) {
  case AF_INET:
    ((struct sockaddr_in *)saddr)->sin_port = htons(port);
    break;

  case AF_INET6:
    ((struct sockaddr_in6 *)saddr)->sin6_port = htons(port);
    break;

  default:
    errno = EPROTONOSUPPORT;
    goto end;
  }

  ret = true;
end:
  freeaddrinfo(res);
  return ret;
}

#ifndef _WIN32
char *get_distro() {
  char    *line = NULL, *distro = NULL, *c = NULL;
  uint64_t line_size = 0;
  FILE    *distf     = NULL;

  if ((distf = fopen("/etc/os-release", "r")) == NULL)
    return NULL;

  if (getline(&line, &line_size, distf) <= 1)
    goto fail;

  for (c = line;; c++) {
    if (*c == 0)
      goto fail;

    if (c == line)
      continue;

    if (*c == '"' && *(c - 1) == '=') {
      distro = strdup(++c);
      break;
    }
  }

  for (c = distro; *c != '"'; c++)
    if (*c == 0)
      goto fail;

  *c = 0;
  goto end;

fail:
  free(distro);
  distro = NULL;
end:
  free(line);
  return distro;
}
#endif

bool parse_addr(char *addr, char **host, uint16_t *port) {
  char *save = NULL, *cur = NULL;
  int   _port = 0;

  *host = NULL;
  *port = 0;

  if ((cur = strtok_r(addr, ":", &save)) == NULL) {
    debug("failed to obtain the hostname");
    return false;
  }

  if ((*host = strdup(cur)) == NULL) {
    debug("failed to allocate memory for the hostname");
    return false;
  }

  if ((cur = strtok_r(NULL, ":", &save)) == NULL) {
    debug("failed to obtain the port");
    return false;
  }

  if ((_port = atoi(cur)) > UINT16_MAX || _port <= 0) {
    debug("invalid port number");
    return false;
  }

  *port = _port;
  return true;
}
