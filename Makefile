PACKAGE_NAME := microsvc
PROJECT_NAME := microsvc-base
PROJECT_DIR := ~/go/src/github.com/hathbanger/$(PROJECT_NAME)

ARCH := $(shell uname -m)

ENV ?= dev

GHE_TOKEN ?=

GIT_COMMIT := $(shell git rev-parse --short HEAD)

GOVERSION := 1.11.5
GOPATH ?= $(shell go env GOPATH)
GOMAXPROCS ?= 4
GO111MODULE ?= on
LD_FLAGS ?= \
	-s \
	-w \
	-extldflags "-static" \
	-X $(PROJECT_DIR)/cmd.Name=$(PROJECT_NAME) \
	-X $(PROJECT_DIR)/cmd.GitCommit=$(GIT_COMMIT) \
	-X $(PROJECT_DIR)/cmd.ARCH=$(ARCH)

PORT ?=

SERVICE_DIR := /opt/atu/$(PROJECT_NAME)
RELEASE_DIR := /go/src/github.com/hathbanger/$(PROJECT_NAME)/build/release
RELEASE_NAME := $(PROJECT_NAME)

FILE=templates/test.txt
VARIABLE=`echo $(FILE)`

USER := $(shell id -u)

all: usage
usage: Makefile
	@echo
	@echo "$(PROJECT_NAME) supports the following:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo
.PHONY: usage

help: usage

## build - builds binaries
build:
	@echo "[INFO]: building binaries with LD_FLAGS \"$(LD_FLAGS)\""
	go build -a -v \
		-ldflags '$(LD_FLAGS)' \
		-o $(RELEASE_DIR)/$(RELEASE_NAME) \
		cmd/$(PROJECT_NAME).go
.PHONY: build

## bootstrap - bootstraps the cuurent system for development
bootstrap:
	@echo "[INFO]: installing prerequisites"
	@/usr/bin/ruby -e $$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)
	@brew install gnu-getopt
	@brew install docker-compose
	@brew install jq
	xcode-select --install
	@curl -s -L https://git.io/vp6lP | sh -s -- -b $(GOPATH)/binary
	@echo "\n- - - - - - - -"
	@echo "!!!ATTN add the following line to your desired .profile"
	@echo 'export PATH=/usr/local/opt/coreutils/libexec/gnubin:usr/local/opt/gnu-getopt/bin:$$PATH'
	@echo "\n_ _ _ _ _ _ _ _"
	@echo "[INFO]: finished bootstrap."
.PHONY: bootstrap

## cd - clones down and updates the deployment info
cd:
	@echo "[INFO]: about to clone son"

.PHONY: cd

## deps - installs and vendors dependencies
deps:
	@echo "[INFO] installing dependencies"
	rm -f go.sum
	@export GO111MODULE=on
	go mod download; \
	go mod vendor

## image - build docker image
image:
	docker build \
		--build-arg SERVICE_NAME=$(PROJECT_NAME) \
		--build-arg SERVICE_PORT=$(PORT) \
		--build-arg SERVICE_DIR=$(SERVICE_DIR) \
		--build-arg RELEASE_DIR=$(RELEASE_DIR) \
		-f build/Dockerfile .

## install - installs binary
install:
	go install -v
.PHONY: install

# endpoint - creates endpoint
endpoint: templates
	# @echo "[INFO] - parsing flags: ENDPOINT => $(ENDPOINT)"
	# @echo "[INFO] - editing for uppercase"
	# grep -rl Bar templates | xargs sed -i.bak "s/Bar/$(ENDPOINT)/g"
	# @echo "[INFO] - editing for lowercase"
	# @grep -rl bar templates | xargs sed -i.bak "s/bar/$(ENDPOINT2)/g"
	# @echo "[INFO] - adding to instrumenting.go"
	# @cat templates/instrumenting.txt >> pkg/$(PACKAGE_NAME)/instrumenting.go
	# @echo "[INFO] - creating special service file and adding package name: $(ENDPOINT2)"
	# @echo "package $(PACKAGE_NAME)" >> pkg/$(PACKAGE_NAME)/$(ENDPOINT2).go
	# @echo "[INFO] - adding to service file: $(ENDPOINT2).go"
	# @cat templates/service.txt >> pkg/$(PACKAGE_NAME)/$(ENDPOINT2).go
	# @cat templates/decodeRequest.txt >> pkg/$(PACKAGE_NAME)/transport.go
	# @cat templates/endpoints.txt >> pkg/$(PACKAGE_NAME)/endpoints.go
	# @cp templates/models/$(ENDPOINT2).txt  pkg/$(PACKAGE_NAME)/models/$(ENDPOINT2).go
	sed -i.bak "s/replace me/$$(cat templates/test.txt)/g" test/service_test.go
	rm -rf test/fakes
	go generate ./...
	mv pkg/$(PACKAGE_NAME)/$(PACKAGE_NAME)fakes test/fakes
	make templates


## templates - fetches templates
templates:
	@echo "[INFO] - pre clean - templates"
	rm -rf templates
	@echo "[INFO] - pre clean - base"
	rm -rf microsvc-base-temp
	@echo "[INFO] - cloning microsvc-base for svc templates"
	@git clone git@github.com:hathbanger/microsvc-base.git microsvc-base-temp
	@echo "[INFO] - extracting templates"
	@cp -r microsvc-base-temp/templates .
	@echo "[INFO] - deleting temp microsvc"
	@rm -rf microsvc-base-temp
.PHONY: templates

## test - tests binary
test:
	go test ./...
.PHONY: test

## watch - watch the local code for changes and rebuilds the test container
watch:
	docker-compose  --verbose -f $(PROJECT_DIR)/deployments/local/docker-compose.local.yml \
		build --no-cache \
		--build-arg ENV=$(ENV) \
		--build-arg PORT=$(PORT) \
		--build-arg SVC=$(PROJECT_NAME)

	docker-compose -f $(PROJECT_DIR)/deployments/local/docker-compose.local.yml up
.PHONY: watch
