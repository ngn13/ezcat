#ifdef _WIN32
#include <winsock2.h>
#else
#include <netinet/in.h>
#include <sys/socket.h>
#endif
#include <errno.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "packet.h"
#include "util.h"

void packet_free(packet_t *packet) {
  free(packet->data);
  packet->header.size = 0;
  packet->data        = NULL;
}

void packet_set_flags(packet_t *packet, uint8_t type, uint8_t cmd) {
  packet->header.flags = 0;

  packet->header.flags |= (cmd & 0b1111);
  packet->header.flags |= (type & 1) << 4;
  packet->header.flags |= (PACKET_VERSION & 0b111) << 5;
}

void packet_set_data(packet_t *packet, char *data, uint8_t size) {
  if (size == 0 && NULL == data) {
  no_alloc:
    packet->data        = NULL;
    packet->header.size = 0;
    return;
  }

  if (size == 0) {
    ssize_t _size = strlen(data);

    if (_size > UINT8_MAX)
      size = UINT8_MAX;
    else if (_size == 0)
      goto no_alloc;
    else
      size = _size;
  }

  packet->header.size = size;
  packet->data        = malloc(size + 1);

  bzero(packet->data, size + 1);
  memcpy(packet->data, data, size);
}

bool packet_send(packet_t *packet, int s) {
  char buffer[PACKET_MAX_SIZE], *bufp = buffer;

  packet->header.session = htonl(packet->header.session);
  packet->header.job_id  = htons(packet->header.job_id);

  bufp = copy_to(bufp, &packet->header, sizeof(packet->header));
  bufp = copy_to(bufp, packet->data, packet->header.size);

  dump(buffer, bufp - buffer);

  if (send(s, buffer, bufp - buffer, 0) <= 0) {
    debug("send failed: %s", strerror(errno));
    return false;
  }

  return true;
}

bool packet_recv(packet_t *packet, int s) {
  char    buffer[PACKET_MAX_SIZE], *bufp = buffer;
  int64_t size = 0;

#ifdef _WIN32
  if ((size = recv(s, buffer, PACKET_MAX_SIZE, 0)) == SOCKET_ERROR) {
#else
  if ((size = recv(s, buffer, PACKET_MAX_SIZE, 0)) <= 0) {
#endif
    debug("recv failed: %s", strerror(errno));
    return false;
  }

  dump(buffer, size);

  bufp = copy_from(&packet->header, bufp, sizeof(packet->header));

  packet->header.session = ntohl(packet->header.session);
  packet->header.job_id  = ntohs(packet->header.job_id);

  if (packet_version(packet) != PACKET_VERSION) {
    debug("version mismatch");
    return false;
  }

  if ((packet->data = malloc(packet->header.size + 1)) == NULL) {
    debug("failed to allocate data buffer (size: %u)", packet->header.size);
    return false;
  }

  bzero(packet->data, packet->header.size + 1);
  bufp = copy_from(packet->data, bufp, packet->header.size);

  return true;
}
