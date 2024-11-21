#include "cmd.h"
#include "agent.h"
#include "packet.h"
#include "util.h"

bool cmd_handle(agent_t *agent) {
  bool     ret = true;
  packet_t packet;

  packet_set_flags(&packet, PACKET_TYPE_REQ, CMD_ASK);
  packet_set_data(&packet, NULL, 0);

  if (!agent_send(agent, &packet)) {
    debug("failed to send the ask command");
    goto end;
  }

  packet_free(&packet);

  if (!agent_recv(agent, &packet)) {
    debug("failed to receive a packet");
    goto end;
  }

  switch (packet_cmd(&packet)) {
  case CMD_AUTH:
    debug("server asked for registeration");
    if (!cmd_register(agent)) {
      debug("registration failed");
      ret = false;
    }
    break;

  case CMD_NONE:
    debug("asked for a job, but there is nothing to do");
    break;

  case CMD_INFO:
    debug("received an info command");
    cmd_info_handler(agent, &packet);
    break;

  case CMD_RUN:
    debug("received a run command");
    cmd_run_handler(agent, &packet);
    break;

  case CMD_KILL:
    debug("received a kill command");
    cmd_success(agent, "success", 0);
    ret = false;
    break;

  default:
    debug("received an unknown command: %d", packet_cmd(&packet));
    break;
  }

end:
  packet_free(&packet);
  return ret;
}
