.PHONY: default clean alpine build

default: build

build: ./srs-cloud.test

clean:
	rm -f srs-cloud.test

./srs-cloud.test: *.go Makefile go.mod
	go test -mod=vendor -c -o srs-cloud.test .
