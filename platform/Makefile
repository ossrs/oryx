.PHONY: default clean alpine ui-install ui-build ui-build-cn ui-build-en ui-build-docker ui-clean build help

default: ./platform

help:
	@echo "Usage: make build|clean"
	@echo "     default     Build the platform only"
	@echo "     build       Build the platform and ui"
	@echo "     clean       Cleanup the platform and ui"
	@echo "     ui-install	Use npm install for dependencies for ui"
	@echo "     ui-build	Use npm run build for building ui"
	@echo "     ui-clean	Cleanup the ui"

build: ./platform ui-install ui-build

clean: ui-clean
	rm -f platform

./platform: *.go Makefile go.mod
	go build .

# Build for alpine, see https://www.cloudbees.com/blog/building-minimal-docker-containers-for-go-applications
alpine:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o platform .

ui-install:
	cd ui && npm install

ui-build:
	cd ui && env PUBLIC_URL=/mgmt REACT_APP_LOCALE=zh BUILD_PATH=build/zh npm run build

ui-build-cn:
	cd ui && env PUBLIC_URL=/mgmt REACT_APP_LOCALE=zh BUILD_PATH=build/zh npm run build

ui-build-en:
	cd ui && env PUBLIC_URL=/mgmt REACT_APP_LOCALE=en BUILD_PATH=build/en npm run build

ui-build-docker:
	cd ui && npm run build

ui-clean:
	rm -rf platform/ui/build