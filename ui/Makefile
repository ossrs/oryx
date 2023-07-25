.PHONY: default clean build
.PHONY: npm build-cn build-en

default: build

build: npm build build-cn build-en

npm:
	npm install

build: build-cn

build-cn:
	env PUBLIC_URL=/mgmt REACT_APP_LOCALE=zh BUILD_PATH=build/zh npm run build

build-en:
	env PUBLIC_URL=/mgmt REACT_APP_LOCALE=en BUILD_PATH=build/en npm run build

clean:
	rm -rf ui/build