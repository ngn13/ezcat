GOSRCS = $(shell find . -type f -name '*.go')
CGO_ENABLED = 1

all: ezcat.elf

ezcat.elf: $(GOSRCS)
	CGO_ENABLED=$(CGO_ENABLED) go build -o $@

format:
	gofmt -s -w .

clean:
	rm *.elf

debug:
	EZCAT_DEBUG=1 ./ezcat.elf

.PHONY: format clean debug
