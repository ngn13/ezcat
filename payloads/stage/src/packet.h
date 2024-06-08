#pragma once

#include <stddef.h>
#include <stdint.h>
#ifdef _WIN32
#include <winsock2.h>
#else
#include <sys/socket.h>
#endif

#define REQ_LABEL_LEN 4
#define LABEL_LIMIT 63
#define NAME_LIMIT 255
#define UDP_LIMIT 512

typedef struct dns_header {
  uint16_t id;
  uint16_t flags;
  uint16_t qdcount;
  uint16_t ancount;
  uint16_t nscount;
  uint16_t arcount;
} dns_header_t;

typedef struct dns_qd {
  char   **qname;
  uint16_t qtype;
  uint16_t qclass;
} dns_qd_t;

typedef struct dns_rr {
  char    *name;
  uint16_t type;
  uint16_t class;
  uint32_t ttl;
  uint16_t rdlength;
  char    *rdata;
} dns_rr_t;

typedef struct dns_packet {
  dns_header_t header;
  dns_qd_t    *questions;
  dns_rr_t    *answers;
  dns_rr_t    *authorities;
  dns_rr_t    *additionals;
} dns_packet_t;

void packet_alloc(dns_packet_t *p);
void packet_free(dns_packet_t *p);
bool packet_send(int s, struct sockaddr *addr, dns_packet_t *packet);
bool packet_recv(int s, struct sockaddr *addr, dns_packet_t *packet);
