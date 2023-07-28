.PHONY: default build install run uninstall upgrade test npm help clean

PREFIX ?= /usr/local/srs-cloud
__REAL_INSTALL = $(DESTDIR)$(PREFIX)

default: build

help:
	@echo "Usage: make build|install|test"
	@echo "     default     Show help and quit"
	@echo "     build       Build the project, npm install and build the ui"
	@echo "     install     Copy files for installer"
	@echo "     test     	Run tests"

build:
	make -C platform
	make -C ui
	make -C releases

clean:
	make -C platform clean
	make -C ui clean
	make -C releases clean

install:
ifeq ($(shell pwd), $(__REAL_INSTALL))
	@echo "Ignore install for $(__REAL_INSTALL)"
else
	rm -rf $(__REAL_INSTALL)
	mkdir -p $(__REAL_INSTALL)/mgmt $(__REAL_INSTALL)/platform $(__REAL_INSTALL)/ui
	cp -rf mgmt $(__REAL_INSTALL)/mgmt
	cp -rf ui/build $(__REAL_INSTALL)/ui/build
	cp -f platform/platform $(__REAL_INSTALL)/platform/platform
	cp -f platform/bootstrap $(__REAL_INSTALL)/platform/bootstrap
	cp -rf platform/auto $(__REAL_INSTALL)/platform/auto
	cp -rf platform/containers $(__REAL_INSTALL)/platform/containers
	rm -rf $(__REAL_INSTALL)/platform/containers/objs/*
	cp -rf usr $(__REAL_INSTALL)/usr
	sed -i "s|/usr/local/srs-cloud|$(PREFIX)|g" $(__REAL_INSTALL)/usr/lib/systemd/system/srs-cloud.service
endif

uninstall:
ifeq ($(shell pwd), $(__REAL_INSTALL))
	@echo "Ignore uninstall for $(__REAL_INSTALL)"
else
	rm -rf $(__REAL_INSTALL)
endif

test:
	cd platform && go test ./...
	cd releases && go test ./...
	cd ui && npm run test
