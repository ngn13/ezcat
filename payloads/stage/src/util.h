#pragma once

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <sys/types.h>

#ifdef _WIN32
#include <windows.h>

#define bzero(b, len) ZeroMemory(b, len)
void fprintfsock(SOCKET sock, const char *fmt, ...);
#else
#include <sys/socket.h>

char *get_distro();
#endif

#define debug(...) print_debug(__func__, __VA_ARGS__);
#define dump(b, s) print_debug_dump(__func__, b, s)

void print_debug_dump(const char *func, char *buf, uint64_t size);
void print_debug(const char *func, const char *msg, ...);

bool resolve(struct sockaddr *saddr, char *addr, uint16_t port);
bool parse_addr(char *addr, char **host, uint16_t *port);

int randint(int min, int max);

void *copy_to(void *dst, void *src, uint64_t size);
void *copy_from(void *dst, void *src, uint64_t size);
