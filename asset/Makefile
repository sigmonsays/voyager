all:
	go generate .

dev:
	go generate -tags dev .
	go install -a -tags dev github.com/sigmonsays/voyager/...

clean:
	rm -vf *.gen.go
