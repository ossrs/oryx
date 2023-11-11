.PHONY: default clean build
.PHONY: npm lint test build-cn build-en

default: build

build: npm build build-cn build-en

npm: lint
	npm install

# Don't run lint in GitHub Actions, because it crashes the build.
# See https://docs.github.com/en/actions/learn-github-actions/variables#default-environment-variables
lint:
	@if [ -z "$(GITHUB_ACTION)" ] && [ -z "$(SRS_NO_LINT)" ]; then npm run lint; fi

test: lint
	npm run test

build: build-cn npm

build-cn: npm
	env PUBLIC_URL=/mgmt REACT_APP_LOCALE=zh BUILD_PATH=build/zh npm run build

build-en: npm
	env PUBLIC_URL=/mgmt REACT_APP_LOCALE=en BUILD_PATH=build/en npm run build

clean:
	rm -rf ui/build