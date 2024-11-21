#pragma once
#include "packet.h"
#include <stdint.h>

typedef struct {
  uint32_t session;
  uint16_t job_id;
  int      socket;
} agent_t;

bool agent_send(agent_t *agent, packet_t *packet);
bool agent_recv(agent_t *agent, packet_t *packet);

bool agent_connect(agent_t *agent);
void agent_disconnect(agent_t *agent);
