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
