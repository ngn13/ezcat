#pragma once

#include <stdint.h>
#ifdef _WIN32
#include <winsock2.h>
#else
#include <sys/socket.h>
#endif

#define ARGUMENT_LEN 63
#define SESSION_LEN 31
#define ID_LEN 31

extern char session[SESSION_LEN + 1];

typedef enum req_codes {
  REQ_REGISTER = 'R',
  REQ_CONTINUE = 'C',
  REQ_WORK     = 'W',
  REQ_DONE     = 'D',
  REQ_FAIL     = 'F',
} req_codes_t;

typedef enum res_codes {
  RES_NOTNOW = 'N',
  RES_FAIL   = 'F',
  RES_OK     = 'K',

  RES_INFO = 'I',
  RES_KILL = 'D',
  RES_RUN  = 'R',
} res_codes_t;

typedef struct req {
  req_codes_t code;
  char       *arg;
  bool        end;
} req_t;

typedef struct res {
  res_codes_t code;
  char        arg[ARGUMENT_LEN];
} res_t;

bool client_setup_socket(int s);
bool client_send(int s, struct sockaddr *addr, req_t *req);
bool client_recv(int s, struct sockaddr *addr, res_t *res);
bool client_request(int s, struct sockaddr *addr, req_codes_t code, char *arg);

bool client_register(int s, struct sockaddr *addr);
bool client_info(int s, struct sockaddr *addr);
bool client_run(int s, struct sockaddr *addr, char *target);
bool client_kill(int s, struct sockaddr *addr);
