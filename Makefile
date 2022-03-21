.PHONY: default build install run uninstall upgrade test

SRS_PREFIX=/usr/local/srs-terraform
__REAL_INSTALL=$(DESTDIR)$(SRS_PREFIX)

default:
	@echo "Usage: make build|install|upgrade|test"
	@echo "     build       Build the project, npm install and build the ui"
	@echo "     install     Copy files for installer"
	@echo "     upgrade     Build for upgrade"
	@echo "     test     	Run tests"

build:
	@cd ffmpeg && npm install
	@cd platform && npm install
	@cd hooks && npm install
	@cd mgmt && npm install
	@cd mgmt/ui && npm install
	@cd releases && npm install
	@cd tencent && npm install
	@cd mgmt && npm run build

upgrade:
	@cd mgmt && npm install
	@cd mgmt/ui && npm install
	@cd mgmt && npm run upgrade

install:
	@mkdir -p $(__REAL_INSTALL)
	@rm -rf $(__REAL_INSTALL)/mgmt $(__REAL_INSTALL)/js-core
	@ln -sf `pwd`/mgmt $(__REAL_INSTALL)/mgmt
	@ln -sf `pwd`/js-core $(__REAL_INSTALL)/js-core
	@rm -rf $(__REAL_INSTALL)/usr
	@cp -rf usr $(__REAL_INSTALL)/usr
	@sed -i "s|/usr/local/srs-terraform|$(SRS_PREFIX)|g" $(__REAL_INSTALL)/usr/lib/systemd/system/srs-terraform.service

uninstall:
	@echo "rmdir $(SRS_PREFIX)"
	@rm -rf $(SRS_PREFIX)

test:
	@cd ffmpeg && npm test
	@cd platform && npm test
	@cd hooks && npm test
	@cd mgmt && npm test
	@cd mgmt/ui && npm test
	@cd releases && npm test
	@cd tencent && npm test

