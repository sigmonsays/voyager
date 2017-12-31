all:
	$(MAKE) -C filetype
	$(MAKE) -C asset

dev:
	go get -u google.golang.org/grpc
	go get -u github.com/golang/protobuf/protoc-gen-go
	go get -u golang.org/x/tools/cmd/stringer

	$(MAKE) -C asset dev

server:
	voyager -log trace
