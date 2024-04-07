.PHONY: default clean alpine build

default: build

build: ./oryx.test

clean:
	rm -f oryx.test

./oryx.test: *.go Makefile go.mod
	go test -mod=vendor -c -o oryx.test .
