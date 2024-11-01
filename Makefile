.PHONY: default build build-no-ui install run uninstall upgrade test npm help clean

PREFIX ?= /usr/local/oryx
__REAL_INSTALL = $(DESTDIR)$(PREFIX)

default: build

help:
	@echo "Usage: make build|install|utest"
	@echo "     default     Show help and quit"
	@echo "     build       Build the project, npm install and build the ui"
	@echo "     install     Copy files for installer"
	@echo "     test     	Run tests"

build: build-no-ui
	make -C ui

build-no-ui:
	make -C platform
	make -C test
	make -C releases

clean:
	make -C platform clean
	make -C ui clean
	make -C test clean
	make -C releases clean

install:
ifeq ($(shell pwd), $(__REAL_INSTALL))
	@echo "Ignore install for $(__REAL_INSTALL)"
else
	rm -rf $(__REAL_INSTALL)
	mkdir -p $(__REAL_INSTALL)/mgmt $(__REAL_INSTALL)/platform $(__REAL_INSTALL)/ui
	cp -f mgmt/bootstrap $(__REAL_INSTALL)/mgmt/bootstrap
	cp -rf ui/build $(__REAL_INSTALL)/ui/build
	cp -f platform/platform $(__REAL_INSTALL)/platform/platform
	cp -f platform/bootstrap $(__REAL_INSTALL)/platform/bootstrap
	cp -f platform/bootstrap.origin.cluster $(__REAL_INSTALL)/platform/bootstrap.origin.cluster
	cp -rf platform/auto $(__REAL_INSTALL)/platform/auto
	cp -rf platform/containers $(__REAL_INSTALL)/platform/containers
	(cd platform && cp -P dvr objs record upload vlive vod transcript dub $(__REAL_INSTALL)/platform)
	rm -rf $(__REAL_INSTALL)/platform/containers/objs/*
	rm -rf $(__REAL_INSTALL)/platform/containers/data/*
	cp -rf usr $(__REAL_INSTALL)/usr
	sed -i "s|/usr/local/oryx|$(PREFIX)|g" $(__REAL_INSTALL)/usr/lib/systemd/system/oryx.service
endif

uninstall:
ifeq ($(shell pwd), $(__REAL_INSTALL))
	@echo "Ignore uninstall for $(__REAL_INSTALL)"
else
	rm -rf $(__REAL_INSTALL)
endif

test:
	cd platform && go test -v ./...
	cd releases && go test -v ./...
	cd test && go test -v -check-api-secret=false -test.run TestSystem_Empty ./...
	make -C ui test
