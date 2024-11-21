HEADERS = $(wildcard src/*.h)
CSRCS   = $(wildcard src/*.c)
CFLAGS  = -O3
CC      = gcc

# options
STAGE_ID = $(shell python3 -c 'print("0"*32, end="")')
STAGE_SERVER_HOST = 127.0.0.1
STAGE_SERVER_PORT = 5567
STAGE_DEBUG = 1

all: stage

stage: $(CSRCS) $(HEADERS)
	$(CC) $(CSRCS) $(CFLAGS)                       \
		-DSTAGE_ID=\"$(STAGE_ID)\"                   \
		-DSTAGE_SERVER_PORT=$(STAGE_SERVER_PORT)     \
		-DSTAGE_SERVER_HOST=\"$(STAGE_SERVER_HOST)\" \
		-DSTAGE_DEBUG=$(STAGE_DEBUG)                 \
		-o $@

clean:
	rm stage

format:
	clang-format -i -style=file $(HEADERS) $(CSRCS)

.PHONY: clean format
