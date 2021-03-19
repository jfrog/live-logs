####################################################################
## JFrog CLI Plugins registry will use `go build` to build plugin binary.
## This file is added for user's benefit to work with code in local.
####################################################################

SHELL := /bin/bash

.DEFAULT_GOAL = help

GOCMD ?= go
TEST_TAGS ?= -tags=test
.DEFAULT_GOAL = build

help:				## Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

# BUILD:

build: clean			## Build live-logs plugin
	$(GOCMD) build

fmt-fix:			## Gofmt fix errors
	gofmt -w -s .

vet:				## GoVet
	$(GOCMD) vet $(TEST_TAGS) ./...

clean:				## Clean from created bins
	rm -f live-logs

run:				## Run the plugin
	$(GOCMD) run main.go

# TEST EXECUTION

test:				## Run all tests
	time $(GOCMD) test ./... $(TEST_TAGS) -count=1

test-list:			## List all tests
	$(GOCMD) list ./...

cover:				## Shows coverage details
	$(GOCMD) test ./... $(TEST_TAGS) -count=1 -coverprofile=coverage


# PLUGIN INSTALLATION

install:			## Install the plugin to jfrog cli
	jfrog plugin install live-logs

uninstall:			## Uninstall the plugin to jfrog cli
	jfrog plugin uninstall live-logs

.PHONY: help build fmt-fix vet clean run test test-list cover install uninstall


