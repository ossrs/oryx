.PHONY: default clean

default: ./releases

clean:
	rm -f releases

./releases: *.go Makefile go.mod
	go build .

