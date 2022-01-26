.PHONY: default build install run uninstall upgrade

SRS_PREFIX=/usr/local/srs-terraform
__REAL_INSTALL=$(DESTDIR)$(SRS_PREFIX)

default: run

run:
	@cd mgmt && npm start

build:
	@cd mgmt && npm install
	@cd mgmt && npm run build

upgrade:
	@cd mgmt && npm install
	@cd mgmt && npm run upgrade

install:
	@mkdir -p $(__REAL_INSTALL)
	@rm -rf $(__REAL_INSTALL)/mgmt
	@ln -sf `pwd`/mgmt $(__REAL_INSTALL)/mgmt
	@ls -lh $(__REAL_INSTALL)/mgmt
	@rm -rf $(__REAL_INSTALL)/usr
	@cp -rf usr $(__REAL_INSTALL)/usr
	@sed -i "s|/usr/local/srs-terraform|$(SRS_PREFIX)|g" $(__REAL_INSTALL)/usr/lib/systemd/system/srs-terraform.service

uninstall:
	@echo "rmdir $(SRS_PREFIX)"
	@rm -rf $(SRS_PREFIX)
