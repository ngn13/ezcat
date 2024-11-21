#pragma once

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#ifdef _WIN32
#include <winsock2.h>
#else
#include <sys/socket.h>
#endif

#define PACKET_MAX_SIZE (255 + 4 + 3)
#define PACKET_VERSION 0

typedef enum {
  PACKET_TYPE_REQ = 0,
  PACKET_TYPE_RES = 1,
} packet_type_t;

typedef struct {
  struct {
    uint32_t session; // agent session
    uint16_t work_id; // ID of the work assoicated with this packet
    uint8_t  flags;   // version (3 bits), 1 (type), 4 (command)
    uint8_t  size;    // data size
  } header;
  char *data;
} packet_t;

#define packet_version(p) ((p)->header.flags >> 5)
#define packet_type(p) ((p)->header.flags >> 4) & 1)
#define packet_cmd(p) ((p)->header.flags & 0b1111)

void packet_free(packet_t *packet);

void packet_set_flags(packet_t *packet, uint8_t type, uint8_t cmd);
void packet_set_data(packet_t *packet, char *data, uint8_t size);

bool packet_send(packet_t *packet, int s);
bool packet_recv(packet_t *packet, int s);
