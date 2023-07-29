.PHONY: default clean alpine build

default: build

build: ./platform

clean:
	rm -f platform

./platform: *.go Makefile go.mod
	go build -mod=vendor .

# Build for alpine, see https://www.cloudbees.com/blog/building-minimal-docker-containers-for-go-applications
alpine:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o platform .
