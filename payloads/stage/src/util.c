#ifdef _WIN32
#include <winsock2.h>
#else
#include <arpa/inet.h>
#include <netinet/in.h>
#endif
#include <errno.h>
#include <netdb.h>
#include <stdarg.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <strings.h>
#include <sys/socket.h>
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

int randint(int min, int max) {
  return min + rand() % (max + 1 - min);
}

bool copy_to_buffer(void *buffer, void *src, size_t size, ssize_t *total, ssize_t *used) {
  if (*used + size > *total)
    return false;

  if (*used == 0)
    bzero(buffer, *total);

  memcpy(buffer + *used, src, size);
  *used += size;
  return true;
}

bool copy_from_buffer(void *dst, void *buffer, size_t size, ssize_t *total, ssize_t *used) {
  if (*used + size > *total)
    return false;

  if (*used == 0)
    bzero(dst, size);

  memcpy(dst, buffer + *used, size);
  *used += size;
  return true;
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

  if (getaddrinfo(addr, NULL, NULL, &res) < 0)
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
