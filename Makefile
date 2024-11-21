format:
	$(MAKE) -C server format
	$(MAKE) -C payloads/stage format
	cd app && npm run format

.PHONY: format
