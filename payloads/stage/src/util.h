#pragma once

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <sys/types.h>

#ifdef _WIN32
#include <windows.h>

#define bzero(b, len) ZeroMemory(b, len)
void fprintfsock(SOCKET sock, const char *fmt, ...);
#endif

#define debug(...) print_debug(__func__, __VA_ARGS__);

void print_debug(const char *func, const char *msg, ...);
int  randint(int min, int max);
bool copy_to_buffer(void *buffer, void *src, size_t size, ssize_t *total, ssize_t *used);
bool copy_from_buffer(void *dst, void *buffer, size_t size, ssize_t *total, ssize_t *used);
bool resolve(struct sockaddr *saddr, char *addr, uint16_t port);
