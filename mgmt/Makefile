.PHONY: default clean alpine

default: ./mgmt

clean:
	rm -f mgmt

./mgmt: *.go Makefile go.mod
	go build .

# Build for alpine, see https://www.cloudbees.com/blog/building-minimal-docker-containers-for-go-applications
alpine:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o mgmt .
