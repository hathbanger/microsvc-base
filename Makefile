# - - -
# Standard Makefile Variables
# - - -

# PROJECT_NAME - name of the project
PROJECT_NAME := microsvc-base

# PACKAGE_NAME - name of the service's package (can only be one)
PACKAGE_NAME := microsvc

# GITHUB_TARGET - the url for the desired github
GITHUB_TARGET := github.com

# USER_NAME - user name in github
USER_NAME := hathbanger

# PROJECT_DIR - directory for project
PROJECT_DIR := ~/go/src/$(GITHUB_TARGET)/$(USER_NAME)/$(PROJECT_NAME)

# CONCOURSE_TARGET - url for concourse
CONCOURSE_TARGET := localhost

# CONSUL_HTTP_ADDR - the consul address
CONSUL_HTTP_ADDR := localhost

# DEBUG - enable debug logging in service
DEBUG ?= true

# DOCKER_IMAGE - url for docker image
DOCKER_IMAGE := hub.docker.com/hathbanger/$(PROJECT_NAME)

# ENV - environment
ENV ?= dev

# GITHUB_TOKEN - github token
GITHUB_TOKEN ?=

# PORT - port for service
PORT ?= 8080

# VAULT_ADDR - address for vault
VAULT_ADDR ?= 

# VAULT_TOKEN - TOKEN
VAULT_TOKEN ?= $(shell vault login -tls-skip-verify -address=$(VAULT_ADDR) -token-only -method=github token="$(GITHUB_TOKEN)")

# - - - - - - - - -
# Default Config Variables
# - - - - - - - - -

# Arch - the architecture provided by uname -m
ARCH := $(shell uname -m)

# CONSUL_TEMPLATE_RELEASE - the version of the consul template
CONSUL_TEMPLATE_RELEASE := 0.20.0

# CONSUL_AUTH
CONSUL_CLIENT_CERT ?=
CONSUL_CLIENT_KEY ?=
CONSUL_CACERT ?=
CONSUL_HTTP_TOKEN := $(shell export VAULT_TOKEN=$(VAULT_DEV_ROOT_TOKEN_ID);vault kv get -tls-skip-verify -address=$(VAULT_ADDR) -field token secret/my-secret)

# GIT_COMMIT - the git commit
GIT_COMMIT := $(shell git rev-parse --short HEAD 2> /dev/null)

# GOLANG SPECIFICS
GO_RELEASE := 1.11.5
GOPATH ?= $(shell go env GOPATH)
GOMAXPROCS ?= 4
GO111MODULE ?= on
GOTAGS ?=

# KERNAL - the kernal name as provided by uname -s
KERNAL := $(shell uname -s)

# RELEASE_DIR - the output diretory for the binary builds.
# also the output diretor in the Dockerfile
RELEASE_DIR := build/release

# SERVICE_DIR - the directory to run the binary from in the container
SERVICE_DIR := /opt/$(USER_NAME)/$(PROJECT_NAME)

# USER - the running USER
USER := $(shell id -u)

# SED - the sed command to use
ifeq ($(KERNAL), Darwin)
SED := sed -i ""
else
SED := sed -i
endif

# CI - triggers CI specific build instructions.
# if Ci is set to any value, CI is assumed.
CI ?=

# CI_TARGET - defines the lcoation in CI for working directory
CI_TARGET := /tmp/build/$(shell ls /tmp/build 2> /dev/null)

# VERSION - the version as read from the branch version file VERSION
ifdef CI
VERSION := $(shell cat $(CI_TARGET)/version/version)
else
VERSION := $(shell git show version:version 2> /dev/null)
endif

LD_FLAGS ?= \
	-s \
	-w \
	-linkmode external \
	-extldflags "-static" \
	-X Z$(GITHUB_TARGET)/$(USER_NAME)/pkg/microsvc.Name=$(PROJECT_NAME)

# RELEASE_NAME - release nae of the binary build
RELEASE_NAME := $(PROJECT_NAME)-$(VERSION)

# VAULT_CMD - the path to the vault binary
VAULT_CMD := $(shell which vault)

# VAULT_RELEASE - version of hashicorp vault
VAULT_RELEASE := 1.1.1


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
build: lint test
	ls -la .
	echo $${PWD}
	echo "[INFO]: building binaries with LD_FLAGS \"$(LD_FLAGS)\""
	go build -a \
		-ldflags '$(LD_FLAGS)' \
		-o $(RELEASE_DIR)/$(RELEASE_NAME) \
		-tags "$(GOTAGS)" \
		cmd/$(PROJECT_NAME).go
ifdef CI
	mv $(RELEASE_DIR)/$(RELEASE_NAME) $(CI_TARGET)/resource-$(PROJECT_NAME)/
endif
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

# Check - checks if required variables are set
check:
	@:$(call check_defined, CONSUL_HTTP_ADDR)
	@:$(call check_defined, VAULT_ADDR)
	@:$(call check_defined, VAULT_TOKEN)
	@:$(call check_defined, GITHUB_TOKEN)
.PHONY: check

## clean - clean out all build artifactions and config
clean:
	@echo "[INFO]: Removing local build artifacts"
	@rm -rf configs/certs
	@rm -rf vendors
	@rm -f  configs/config.json
	@rm -f $(RELEASE_DIR)/$(RELEASE_NAME)
.PHONY: clean

## config - will assemble the config.json from Consul and Vault
config: check prerequisites
    @echo "[INFO]: parsing consul certificates and token"
    @echo "[INFO]: building local configuration file"
		export CONSUL_HTTP_ADDR=$(CONSUL_HTTP_ADDR); \
		export CONSUL_HTTP_TOKEN=$(CONSUL_HTTP_TOKEN); \
		export CONSUL_HTTP_SSL_VERIFY=0; \
		export VAULT_ADDR=$(VAULT_ADDR); \
		export VAULT_TOKEN=$(VAULT_TOKEN); \
		export ENV=$(ENV); \
		export PORT=$(PORT); \
		consul-template \
			-template="configs/config.ctmpl:configs/config.json" \
			-template="tools/etc/rsyslog.d/00-elasticsearch.ctmpl:tools/etc/rsyslog.d/00-elasticsearch.conf" -once
ifdef CI
    cp $(CI_TARGET)/gopath/src/$(GITHUB_TARGET)/$(USER_NAME)/$(PROJECT_NAME)/configs/config.json $(CI_TARGET)/configs/
    cp -R $(CI_TARGET)/gopath/src/$(GITHUB_TARGET)/$(USER_NAME)/$(PROJECT_NAME)/configs/certs $(CI_TARGET)/configs/
endif
.PHONY: config

## deps - installs and vendors dependencies
deps:
	@echo "[INFO] installing dependencies"
	rm -f go.sum
	@export GO111MODULE=$(GO111MODULE); \
	go mod download; \
	go mod vendor

## deploy - deploys the microservice to the $(ENV)
deploy: custom
ifdef CI
    mv $(CI_TARGET)/gopath/src/$(GITHUB_TARGET)/$(USER_NAME	)/$(PROJECT_NAME)/deployments/$(ENV)/manifest.yml $(CI_TARGET)/release_output/manifest.yml
    mv $(CI_TARGET)/gopath/src/$(GITHUB_TARGET)/$(USER_NAME)/$(PROJECT_NAME)/tools/signal_handler.sh $(CI_TARGET)/release_output/signal_handler.sh
    mv $(CI_TARGET)/resource-$(PROJECT_NAME)/$(PROJECT_NAME)-* $(CI_TARGET)/release_output/$(PROJECT_NAME)
    mv $(CI_TARGET)/configs/config.json $(CI_TARGET)/release_output/
    mv $(CI_TARGET)/configs/certs $(CI_TARGET)/release_output/
    @chmod +x $(CI_TARGET)/release_output/$(PROJECT_NAME)
endif


# endpoint - creates endpoint
endpoint: 
	go get github.com/ChimeraCoder/gojson/gojson
	git clone https://github.com/hathbanger/microsvc-generator
	cd microsvc-generator && make endpoint 
	rm -rf microsvc-generator
	#@make fakes
	#@make templates
	@make fmt

# fakes - creates fakes
fakes:
	@echo "[INFO] - generating fakes"
	@rm -rf test/fakes
	@cat tools/templates/tools.txt > tools/tools.go
	@GO111MODULE=off go get -u github.com/myitcv/gobin; \
	go generate -v ./...
	@echo "[INFO] - moving resources"
	@mv pkg/$(PACKAGE_NAME)/$(PACKAGE_NAME)fakes test/fakes
	@echo "package tools" > tools/tools.go
	@echo "// contents auto generated" >> tools/tools.go
	@echo "[INFO] - done generating fakes"

## fmt - will execute go fmt
fmt:
	go fmt ./...

## fly - will fly -t $(CONCOURSE_TARGET) set-pipeline
fly:
	fly -t $(CONCOURSE_TARGET) set-pipeline --pipeline=$(PROJECT_NAME) --config=build/ci/pipeline.yml --non-interactive
	fly -t $(CONCOURSE_TARGET) unpause-pipeline --pipeline=$(PROJECT_NAME)
.PHONY: fly

## image - will build the docker image
image:
ifdef CI
	@echo "SERVICE_NAME=$(PROJECT_NAME)" > $(CI_TARGET)/image-build-args/args
	@echo "PORT=$(PORT)" >> $(CI_TARGET)/image-build-args/args
	@echo "SERVICE_DIR=$(SERVICE_DIR)" >> $(CI_TARGET)/image-build-args/args
	@echo "RELEASE_DIR=$(RELEASE_DIR)" >> $(CI_TARGET)/image-build-args/args
	@echo "CONSUL_TEMPLATE_RELEASE=$(CONSUL_TEMPLATE_RELEASE)" >> $(CI_TARGET)/image-build-args/args
	@echo "VAULT_RELEASE=$(VAULT_RELEASE)" >> $(CI_TARGET)/image-build-args/args
else
    @docker build \
        --build-arg SERVICE_NAME=$(PROJECT_NAME) \
        --build-arg PORT=$(PORT) \
        --build-arg SERVICE_DIR=$(SERVICE_DIR) \
        --build-arg RELEASE_DIR=$(RELEASE_DIR) \
        --build-arg CONSUL_TEMPLATE_RELEASE=$(CONSUL_TEMPLATE_RELEASE) \
        --build-arg VAULT_RELEASE=$(VAULT_RELEASE) \
        -f build/Dockerfile --rm -t $(DOCKER_IMAGE):$(VERSION) .
endif

## install - will install the binary
install:
	go install -v

## integration - will run the service integration tests
integration:
	@:$(call check_defined, AUTH_URL)
	@:$(call check_defined, TOKEN_URL)
	@:$(call check_defined, CLIENT_ID)
	@:$(call check_defined, CLIENT_SECRET)
	@:$(call check_defined, TARGET_URL)
	@:$(call check_defined, USERNAME)
	@:$(call check_defined, PASSWORD)
	@cd test/_integration_test; go test -v; cd ../../
.PHONY: integration

## lint - will lint the code
lint:
	@if [ ! -f $(GOPATH)/bin/golangci-lint ]; then \
        wget -O - https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(go env GOPATH)/bin v1.16.0; \
    fi
	golangci-lint run -D unused -D errcheck --deadline=10m ./...

## pipeline - will fetch pipeline
pipeline:
	@rm -rf $(PROJECT_DIR)/build/ci
	git clone https://$(GITHUB_TARGET)/$(USER_NAME)/ci-template
	rm -rf ci-template/.git
	mv ci-template/ $(PROJECT_DIR)/build/ci
	@grep -rl service-foo build/ci | xargs $(SED) "s/service-foo/$(PROJECT_NAME)/g"
.PHONY: pipeline

## push - will push the docker image to the docker repository
push:
	docker push $(DOCKER_IMAGE)

## release - will create a release suitable for github
release:
ifdef CI
    mkdir $(CI_TARGET)/release
    mv $(CI_TARGET)/resource-$(PROJECT_NAME)/$(PROJECT_NAME)-* $(CI_TARGET)/release/$(PROJECT_NAME)
    tar -czvf $(VERSION).tar.gz $(PROJECT_NAME)
    rm $(PROJECT_NAME)
    cp $(CI_TARGET)/gopath/src/$(GITHUB_TARGET)/$(USER_NAME)/$(PROJECT_NAME)/.git/ref $(CI_TARGET)/release/commitsh
    # TODO: Make Release Notes.
endif
.PHONY: release

# sedtest - test
sedtest:
	PATTERN='// decodeRequest.txt' ./templater.awk templates/decodeRequest.txt pkg/$(PACKAGE_NAME)/transport.go
	PATTERN='// transport.txt' ./templater.awk templates/transport.txt pkg/$(PACKAGE_NAME)/transport.go
	PATTERN='// endpoints.txt' ./templater.awk templates/endpoints.txt pkg/$(PACKAGE_NAME)/endpoints.go
	PATTERN='// instrumenting.txt' ./templater.awk templates/instrumenting.txt pkg/$(PACKAGE_NAME)/instrumenting.go
	PATTERN='// test.txt' ./templater.awk templates/test.txt test/service_test.go


## templates - fetches templates
templates:
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

## watch - will watch the local code for changes and rebuild testing container
watch:
	echo $(CONSUL_HTTP_TOKEN)
	@docker-compose -f $(PROJECT_DIR)/deployments/local/docker-compose.local.yml \
        build --no-cache \
        --build-arg CONSUL_HTTP_ADDR=$(CONSUL_HTTP_ADDR) \
        --build-arg CONSUL_HTTP_TOKEN=$(CONSUL_HTTP_TOKEN) \
        --build-arg CONSUL_TEMPLATE_RELEASE=$(CONSUL_TEMPLATE_RELEASE) \
        --build-arg VAULT_ADDR=$(VAULT_ADDR) \
        --build-arg VAULT_TOKEN=$(VAULT_TOKEN) \
        --build-arg ENV=$(ENV) \
        --build-arg PORT=$(PORT) \
        --build-arg SVC=$(PROJECT_NAME)
	@echo "WOO"
	@export PORT=$(PORT); \
    docker-compose -f $(PROJECT_DIR)/deployments/local/docker-compose.local.yml up
.PHONY: watch

# prerequsites - installs any prerequisite software on the CI container
prerequisites:
ifdef CI
    @apt-get update
    @apt-get install unzip -y
    @curl -s https://releases.hashicorp.com/consul-template/$(CONSUL_TEMPLATE_RELEASE)/consul-template_$(CONSUL_TEMPLATE_RELEASE)_linux_amd64.tgz | tar -C /usr/local/bin -zx
    @wget -qO- https://releases.hashicorp.com/vault/$(VAULT_RELEASE)/vault_$(VAULT_RELEASE)_linux_amd64.zip | unzip -d /usr/local/bin vault_$(VAULT_RELEASE)_linux_amd64.zip
endif
.PHONY: prerequisites
# custom - the custom target should hold custom deployment instructions
custom: prerequisites
ifdef CI
endif
.PHONY: custom
# Internal Function to Check Defined Makefile Variables
# NOTE: this function is NOT for checking shell variables
check_defined = \
    $(strip $(foreach 1,$1,$(call __check_defined,$1,$(strip $(value 2)))))
__check_defined = \
    $(if $(value $1),,$(error Undefined $1$(if $2, ($2))$(if $(value @), required by target `$@')))
