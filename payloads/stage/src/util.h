#pragma once

#include <stddef.h>
#include <sys/types.h>

#ifdef _WIN32
#include <windows.h>

#define bzero(b, len) ZeroMemory(b, len)
void fprintfsock(SOCKET sock, const char *fmt, ...);
#endif

int  randint(int min, int max);
int  name_to_qname(char *qname, char **name);
bool parse_addr(char *addr, struct sockaddr_in *ret);

void rot13(char *str);
char rot13_char(char c);

bool copy_to_buffer(void *buffer, void *src, size_t size, ssize_t *total, ssize_t *used);
bool copy_from_buffer(void *dst, void *buffer, size_t size, ssize_t *total, ssize_t *used);
