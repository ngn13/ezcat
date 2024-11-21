#include "../cmd.h"
#include "../util.h"

bool cmd_register(agent_t *agent) {
  bool     ret = false;
  packet_t packet;

  packet_set_flags(&packet, PACKET_TYPE_REQ, CMD_REGISTER);
  packet_set_data(&packet, STAGE_ID, 0);

  if (!agent_send(agent, &packet)) {
    debug("failed to send the register command");
    goto end;
  }

  packet_free(&packet);

  if (!agent_recv(agent, &packet)) {
    debug("failed to receive the register command result");
    goto end;
  }

  switch (packet_cmd(&packet)) {
  case CMD_SUCCESS:
    agent->session = packet.header.session;
    debug("registered with session %u", agent->session);

    ret = true;
    goto end;

  case CMD_FAILURE:
    debug("failed to register: %s", packet.data);
    goto end;
  }

  debug("received an unknown response command, register failed");

end:
  packet_free(&packet);
  return ret;
}
