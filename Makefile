# Variable
BUILD_DIR="build"
APP_NAME="dm_exporter"
TAG_VERSION=$(if $(shell git describe --abbrev=0 --tags),$(shell git describe --abbrev=0 --tags),"")
APP_VERSION=$(shell cat VERSION)
VERSION=$(if $(TAG_VERSION),$(TAG_VERSION),$(APP_VERSION))
LDFLAGS += -s -w
LDFLAGS += -X dm_exporter/global.Version_=${VERSION}

.PHONY: help build clean init debug docker-build demo

# Function
default: help
help:
	@echo "usage: make <option>"
	@echo "options and effects:"
	@echo "    help   : Show help"
	@echo "    build  : Build the binary of this project for linux/amd64,linux/arm64 platform"
	@echo "	debug  : Build the binary of this project for current platform"
build: clean
	@go work sync
	@CGO_ENABLED=0 go build -o ${BUILD_DIR}/${APP_NAME} -ldflags '$(LDFLAGS)' ${APP_NAME}
	@echo "build completed"
docker-build:
	@docker build -t ${APP_NAME}:${VERSION} .
build_all: clean init _build_linux
	@echo "build completed"
build_linux: clean init _build_linux_amd64 _build_linux_arm64
	@echo "build completed"
_build_linux: _build_linux_amd64 _build_linux_arm64
_build_linux_amd64:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BUILD_DIR}/${APP_NAME}__${VERSION}-linux-amd64 -ldflags '$(LDFLAGS)'
_build_linux_arm64:
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o ${BUILD_DIR}/${APP_NAME}__${VERSION}-linux-arm64 -ldflags '$(LDFLAGS)'
init:
	@go mod tidy
	@echo "tidy completed"
clean:
	@rm -rf build/
	@echo "clean completed"