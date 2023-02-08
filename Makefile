.PHONY: default build install run uninstall upgrade test npm help clean

SRS_PREFIX=/usr/local/srs-cloud
__REAL_INSTALL=$(DESTDIR)$(SRS_PREFIX)

default: help

help:
	@echo "Usage: make build|install|upgrade|test"
	@echo "     default     Show help and quit"
	@echo "     build       Build the project, npm install and build the ui"
	@echo "     install     Copy files for installer"
	@echo "     upgrade     Build for upgrade"
	@echo "     test     	Run tests"

build:
	cd releases && make
	cd mgmt && make
	cd platform && make build

clean:
	cd releases && make clean
	cd mgmt && make clean
	cd platform && make clean

upgrade:
	@echo "ignore for upgrade"

ifeq ($(shell pwd), $(__REAL_INSTALL))
install:
	@echo "Install ok for $(__REAL_INSTALL)"
else
install:
	mkdir -p $(__REAL_INSTALL)
	ln -sf `pwd`/mgmt $(__REAL_INSTALL)/mgmt
	rm -rf $(__REAL_INSTALL)/usr
	cp -rf usr $(__REAL_INSTALL)/usr
	sed -i "s|/usr/local/srs-cloud|$(SRS_PREFIX)|g" $(__REAL_INSTALL)/usr/lib/systemd/system/srs-cloud.service
endif

uninstall:
	@echo "rmdir $(SRS_PREFIX)"
	rm -rf $(SRS_PREFIX)

test:
	cd platform && go test ./...
	cd mgmt && go test ./...
	cd releases && go test ./...

