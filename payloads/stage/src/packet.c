#ifdef _WIN32
#include <winsock2.h>
#else
#include <netinet/in.h>
#include <sys/socket.h>
#endif
#include <errno.h>
#include <stdbool.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "log.h"
#include "packet.h"
#include "util.h"

void packet_alloc(dns_packet_t *p) {
  p->questions = malloc(sizeof(dns_qd_t) * p->header.qdcount);
  p->answers   = malloc(sizeof(dns_rr_t) * p->header.ancount);

  // these are not really needed (we dont use them for c2 communcation), but might as well
  // yknow just in case...
  p->authorities = malloc(sizeof(dns_rr_t) * p->header.nscount);
  p->additionals = malloc(sizeof(dns_rr_t) * p->header.arcount);
}

void packet_free(dns_packet_t *p) {
  for (int i = 0; i < p->header.qdcount; i++) {
    if (p->questions[i].qname == NULL)
      continue;

    char **cur = p->questions[i].qname;
    int    j   = 0;

    for (; cur[j] != NULL; j++)
      free(cur[j]);
    free(cur);
  }
  free(p->questions);

  for (int i = 0; i < p->header.ancount; i++) {
    free(p->answers[i].name);
    free(p->answers[i].rdata);
  }
  free(p->answers);

  for (int i = 0; i < p->header.nscount; i++) {
    free(p->authorities[i].name);
    free(p->authorities[i].rdata);
  }
  free(p->authorities);

  for (int i = 0; i < p->header.arcount; i++) {
    free(p->authorities[i].name);
    free(p->authorities[i].rdata);
  }
  free(p->additionals);
}

bool packet_copy_rr(dns_rr_t *rr, char *buffer, ssize_t *buffer_total, ssize_t *buffer_used) {
  for (char *c = buffer + *buffer_used; *c != 0; c++) {
    (*buffer_used)++;
    if (*buffer_used >= *buffer_total) {
      debug("failed to pass the qname for the buffer");
      return false;
    }
  }
  (*buffer_used)++;

  rr->name = NULL;
  copy_from_buffer(&rr->type, buffer, sizeof(rr->type), buffer_total, buffer_used);
  copy_from_buffer(&rr->class, buffer, sizeof(rr->class), buffer_total, buffer_used);
  copy_from_buffer(&rr->ttl, buffer, sizeof(rr->ttl), buffer_total, buffer_used);

  if (!copy_from_buffer(&rr->rdlength, buffer, sizeof(rr->rdlength), buffer_total, buffer_used)) {
    debug("failed to read the rdlength from the buffer");
    return false;
  }

  rr->type  = ntohs(rr->type);
  rr->class = ntohs(rr->class);
  rr->ttl   = ntohl(rr->ttl);

  rr->rdlength = ntohs(rr->rdlength);

  if (rr->rdlength > *buffer_total - *buffer_used) {
    debug("bad rdlength (%d)", rr->rdlength);
    return false;
  }

  // pass the TXT length field
  if (rr->type == 16) {
    (*buffer_used)++;
    rr->rdlength--;
  }

  // alloc memory for the actual data
  rr->rdata = malloc(rr->rdlength);

  if (!copy_from_buffer(rr->rdata, buffer, rr->rdlength, buffer_total, buffer_used)) {
    debug("failed to read the rdata from the buffer");
    free(rr->rdata);
    return false;
  }

  return true;
}

bool packet_send(int s, struct sockaddr *addr, dns_packet_t *packet) {
  ssize_t      total = UDP_LIMIT, used = 0;
  char         buffer[total];
  dns_packet_t copy;

  copy.header.id      = htons(packet->header.id);
  copy.header.flags   = htons(packet->header.flags);
  copy.header.qdcount = htons(packet->header.qdcount);
  copy.header.ancount = htons(packet->header.ancount);
  copy.header.nscount = htons(packet->header.nscount);
  copy.header.arcount = htons(packet->header.arcount);

  if (!copy_to_buffer(buffer, &copy.header, sizeof(copy.header), &total, &used)) {
    debug("failed to copy header to the buffer (possible overflow)");
    return false;
  }

  for (int i = 0; i < packet->header.qdcount; i++) {
    dns_qd_t *cur = &packet->questions[0];
    dns_qd_t  curcp;

    char qname[NAME_LIMIT + 1];
    int  qsize = name_to_qname(qname, cur->qname);
    if (qsize <= 0) {
      debug("failed to convert name to qname");
      return false;
    }

    curcp.qtype  = htons(cur->qtype);
    curcp.qclass = htons(cur->qclass);

    copy_to_buffer(buffer, qname, qsize, &total, &used);
    copy_to_buffer(buffer, &curcp.qtype, sizeof(curcp.qtype), &total, &used);
    copy_to_buffer(buffer, &curcp.qclass, sizeof(curcp.qclass), &total, &used);
  }

  if (sendto(s, buffer, used, 0, addr, sizeof(struct sockaddr_in)) <= 0) {
    debug("sendto failed: %s", strerror(errno));
    return false;
  }

  return true;
}

bool packet_recv(int s, struct sockaddr *addr, dns_packet_t *packet) {
  ssize_t total = UDP_LIMIT, used = 0;
  char    buffer[total];

  unsigned int addrlen = sizeof(struct sockaddr_in); // socklen_t
  bool         ret     = false;

#ifdef _WIN32
  if ((total = recvfrom(s, buffer, UDP_LIMIT, 0, addr, &addrlen)) == SOCKET_ERROR) {
#else
  if ((total = recvfrom(s, buffer, UDP_LIMIT, 0, addr, &addrlen)) <= 0) {
#endif
    debug("recvfrom failed: %s", strerror(errno));
    return ret;
  }

  if (!copy_from_buffer(&packet->header, buffer, sizeof(packet->header), &total, &used)) {
    debug("failed to copy buffer to the head (possible overflow)");
    return ret;
  }

  packet->header.qdcount = ntohs(packet->header.qdcount);
  packet->header.ancount = ntohs(packet->header.ancount);
  packet->header.nscount = ntohs(packet->header.nscount);
  packet->header.arcount = ntohs(packet->header.arcount);

  packet_alloc(packet);
  dns_qd_t question; // temporary, just used for "sizeof"

  for (int i = 0; i < packet->header.qdcount; i++) {
    // we dont actually need to parse the questions
    // they are not required for the c2 communcation

    // however they are "in the way" of the answers section
    // sooo lets just quickly go over them, shall we?
    packet->questions[i].qname = NULL;

    for (char *c = buffer + used; *c != 0; c++) {
      used++;
      if (used >= total) {
        debug("failed to pass the qname for the question %d", i);
        goto end;
      }
    }
    used++;

    used += sizeof(question.qtype) + sizeof(question.qclass);
    if (used > total) {
      debug("pakcet is to small for the rest of the question %d", i);
      goto end;
    }
  }

  for (int i = 0; i < packet->header.ancount; i++) {
    dns_rr_t *cur = &packet->answers[i];
    if (!packet_copy_rr(cur, buffer, &total, &used)) {
      debug("failed to read the %d. answer", i);
      goto end;
    }
  }

  for (int i = 0; i < packet->header.nscount; i++) {
    dns_rr_t *cur = &packet->authorities[i];
    if (!packet_copy_rr(cur, buffer, &total, &used)) {
      debug("failed to read the %d. authority", i);
      goto end;
    }
  }

  for (int i = 0; i < packet->header.arcount; i++) {
    dns_rr_t *cur = &packet->additionals[i];
    if (!packet_copy_rr(cur, buffer, &total, &used)) {
      debug("failed to read the %d. additional", i);
      goto end;
    }
  }

  ret = true;

end:
  if (!ret)
    packet_free(packet);
  return ret;
}
