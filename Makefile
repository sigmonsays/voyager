all:
	$(MAKE) -C filetype
	$(MAKE) -C asset

dev:
	go get github.com/tv42/becky
	go get golang.org/x/tools/cmd/stringer
