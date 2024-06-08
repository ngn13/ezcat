#ifdef _WIN32
#include <winsock2.h>
#else
#include <arpa/inet.h>
#include <netinet/in.h>
#endif
#include <stdarg.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <strings.h>
#include <unistd.h>

#include "log.h"
#include "packet.h"
#include "util.h"

int randint(int min, int max) {
  return min + rand() % (max + 1 - min);
}

int name_to_qname(char *qname, char **name) {
  int len = 0, indx = 0;
  bzero(qname, NAME_LIMIT + 1);

  for (int i = 0; name[i] != NULL; i++) {
    len = strlen(name[i]);

    if (len > LABEL_LIMIT) {
      // label is too large
      return -1;
    }

    if (indx + len + 1 > NAME_LIMIT) {
      // name is too large
      return -1;
    }

    qname[indx++] = len;
    memcpy(qname + indx, name[i], len);
    indx += len;
  }

  qname[indx] = 0;
  return indx + 1;
}

bool parse_addr(char *addr, struct sockaddr_in *ret) {
  char   *host = NULL, port[6], *c = NULL;
  ssize_t indx = 0;

  for (c = addr; *c != ':'; c++) {
    if (*c == 0)
      return false;

    if (NULL == host)
      host = malloc(indx + 1);
    else
      host = realloc(host, indx + 1);

    host[indx++] = *c;
    host[indx]   = 0;
  }

  indx = 0;
  c++;

  for (; *c != 0; c++) {
    if (indx + 1 >= 6) {
      free(host);
      return false;
    }

    port[indx++] = *c;
    port[indx]   = 0;
  }

  ret->sin_addr.s_addr = inet_addr(host);
  free(host);

  if (atoi(port) > UINT16_MAX || atoi(port) <= 0)
    return false;

  ret->sin_port = htons(atoi(port));
  return true;
}

char rot13_char(char c) {
  if (c >= 'A' && c <= 'Z')
    return 'A' + ((c - 'A' + 13) % 26);
  else if (c >= 'a' && c <= 'z')
    return 'a' + ((c - 'a' + 13) % 26);
  return c;
}

void rot13(char *rot) {
  int len = strlen(rot);
  for (int i = 0; i < len; i++)
    rot[i] = rot13_char(rot[i]);
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
