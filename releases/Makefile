.PHONY: default clean alpine

default: ./releases

clean:
	rm -f releases

./releases: *.go Makefile go.mod
	go build .

# Build for alpine, see https://www.cloudbees.com/blog/building-minimal-docker-containers-for-go-applications
alpine:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o releases .
