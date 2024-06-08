HEADERS = $(wildcard src/*.h)
CSRCS   = $(wildcard src/*.c)
CFLAGS  = -O3
CC      = gcc

all: stage

stage: $(CSRCS) $(HEADERS)
	$(CC) $(CSRCS) $(CFLAGS) -o $@

clean:
	rm stage

format:
	clang-format -i -style=file $(HEADERS) $(CSRCS)

.PHONY: clean format
