#include "cmd.h"
#include "agent.h"
#include "packet.h"
#include "util.h"

bool cmd_handle(agent_t *agent) {
  packet_t packet;

  if (!agent_recv(agent, &packet)) {
    debug("failed to receive a packet");
    return false;
  }

  switch (packet_cmd(&packet)) {
  case CMD_RUN:
    break;

  case CMD_KILL:
    return false;

  default:
    debug("unknown command: %d", packet_cmd(&packet));
    break;
  }

  return true;
}

bool cmd_register(agent_t *agent) {
  packet_t packet;

  packet_set_flags(&packet, PACKET_TYPE_REQ, CMD_REGISTER);
  packet_set_data(&packet, STAGE_ID, 0);

  if (!agent_send(agent, &packet)) {
    debug("failed to send the register command");
    return false;
  }

  packet_free(&packet);

  if (!agent_recv(agent, &packet)) {
    debug("failed to receive the register command result");
    return false;
  }

  agent->session = packet.header.session;
  debug("registered with session %u", agent->session);

  return true;
}
