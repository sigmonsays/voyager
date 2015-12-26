all:
	$(MAKE) -C filetype
	$(MAKE) -C asset

dev:
	go get -u google.golang.org/grpc
	go get -u github.com/golang/protobuf/protoc-gen-go                                                                                                                            

	$(MAKE) -C asset dev
