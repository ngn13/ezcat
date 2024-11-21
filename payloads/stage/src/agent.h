#pragma once
#include "packet.h"
#include <stdint.h>

typedef struct {
  uint32_t session;
  int      socket;
} agent_t;

#define agent_recv(a, p) packet_recv(p, a->socket)
#define agent_send(a, p) packet_send(p, a->socket)

bool agent_connect(agent_t *agent);
void agent_disconnect(agent_t *agent);
