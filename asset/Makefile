
FLAGS = -pkg=asset -ignore=bindata.go -ignore=Makefile


all: deps
	go-bindata $(FLAGS) *

dev: deps
	go-bindata $(FLAGS) -debug *

deps:
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/elazarl/go-bindata-assetfs/...

clean:
	rm -vf *.gen.go
