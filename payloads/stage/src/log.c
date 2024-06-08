#include <stdarg.h>
#include <stdio.h>

#include "config.h"

void debug(const char *msg, ...) {
  if (!DEBUG)
    return;

  va_list args;
  va_start(args, msg);

  printf("[debug] ");
  vprintf(msg, args);
  printf("\n");

  va_end(args);
}
