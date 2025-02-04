#include "../cmd.h"
#include "../util.h"

bool cmd_failure(agent_t *agent, char *data, uint8_t data_size) {
  packet_t packet;
  bool     ret = false;

  packet_set_flags(&packet, PACKET_TYPE_REQ, CMD_FAILURE);
  packet_set_data(&packet, data, data_size);

  if (!(ret = agent_send(agent, &packet)))
    debug("failed to send the failure command");

  packet_free(&packet);
  return ret;
}
