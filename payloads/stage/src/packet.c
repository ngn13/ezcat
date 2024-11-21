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
  packet->header.flags |= 4 << (type & 1);
  packet->header.flags |= 5 << (PACKET_VERSION & 0b111);
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
  packet->data        = malloc(size);
  memcpy(packet->data, data, size);
}

bool packet_send(packet_t *packet, int s) {
  ssize_t  total = PACKET_MAX_SIZE, used = 0;
  char     buffer[total];
  packet_t copy;

  copy.header.flags   = htons(packet->header.flags);
  copy.header.session = htonl(packet->header.session);
  copy.header.work_id = htons(packet->header.work_id);
  copy.header.size    = htons(packet->header.size);

  if (!copy_to_buffer(buffer, &copy.header, sizeof(copy.header), &total, &used)) {
    debug("failed to copy header to the buffer (possible overflow)");
    return false;
  }

  if (copy.header.size != 0 && !copy_to_buffer(buffer, copy.data, copy.header.size, &total, &used)) {
    debug("failed to copy data to the buffer (possible overflow)");
    return false;
  }

  if (send(s, buffer, used, 0) <= 0) {
    debug("send failed: %s", strerror(errno));
    return false;
  }

  return true;
}

bool packet_recv(packet_t *packet, int s) {
  ssize_t total = PACKET_MAX_SIZE, used = 0;
  char    buffer[total];

#ifdef _WIN32
  if ((total = recvfrom(s, buffer, UDP_LIMIT, 0, addr, &addrlen)) == SOCKET_ERROR) {
#else
  if ((total = recv(s, buffer, PACKET_MAX_SIZE, 0)) <= 0) {
#endif
    debug("recv failed: %s", strerror(errno));
    return false;
  }

  if (!copy_from_buffer(&packet->header, buffer, sizeof(packet->header), &total, &used)) {
    debug("failed to copy buffer to the header (possible overflow)");
    return false;
  }

  packet->header.flags   = ntohs(packet->header.flags);
  packet->header.session = ntohl(packet->header.session);
  packet->header.work_id = ntohs(packet->header.work_id);
  packet->header.size    = htons(packet->header.size);

  if (packet_version(packet) != PACKET_VERSION) {
    debug("version mismatch");
    return false;
  }

  if ((packet->data = malloc(packet->header.size)) == NULL) {
    debug("failed to allocate data buffer (size: %u)", packet->header.size);
    return false;
  }

  if (!copy_from_buffer(packet->data, buffer, packet->header.size, &total, &used)) {
    debug("failed to copy buffer to the data buffer (possible overflow)");
    free(packet->data);
    packet->data = NULL;
    return false;
  }

  return true;
}
