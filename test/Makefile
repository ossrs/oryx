.PHONY: default clean alpine build

default: build

build: ./srs-stack.test

clean:
	rm -f srs-stack.test

./srs-stack.test: *.go Makefile go.mod
	go test -mod=vendor -c -o srs-stack.test .
