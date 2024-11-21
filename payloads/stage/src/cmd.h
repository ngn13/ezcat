#pragma once
#include "agent.h"
#include <stdbool.h>

typedef enum {
  CMD_SUCCESS  = 0,
  CMD_REGISTER = 1,
  CMD_KILL     = 2,
  CMD_RUN      = 3,
} cmd_t;

bool cmd_handle(agent_t *agent);
bool cmd_register(agent_t *agent);
bool cmd_send(agent_t *agent, cmd_t cmd);
