all:
	$(MAKE) -C filetype
	$(MAKE) -C asset

dev:
	$(MAKE) -C asset dev
