.DEFAULT_GOAL := help
CURRENTTAG:=$(shell git describe --tags --abbrev=0)
NEWTAG ?= $(shell bash -c 'read -p "Please provide a new tag (currnet tag - ${CURRENTTAG}): " newtag; echo $$newtag')
GOFLAGS=-mod=mod
GOPRIVATE=github.com/threeport/*,github.com/qleet/*

#help: @ List available tasks
help:
	@clear
	@echo "Usage: make COMMAND"
	@echo "Commands :"
	@grep -E '[a-zA-Z\.\-]+:.*?@ .*$$' $(MAKEFILE_LIST)| tr -d '#' | awk 'BEGIN {FS = ":.*?@ "}; {printf "\033[32m%-19s\033[0m - %s\n", $$1, $$2}'

#test: @ Run tests
test:
	@export GOPRIVATE=$(GOPRIVATE); export GOFLAGS=$(GOFLAGS); go test $(go list ./... | grep -v /internal/setup)

#build: @ Build workload controller binary
build:
	@go generate
	@export GOPRIVATE=$(GOPRIVATE); export GOFLAGS=$(GOFLAGS); export CGO_ENABLED=0; go build -a -o qleetctl main.go

#get: @ Download and install dependency packages
get:
	@export GOPRIVATE=$(GOPRIVATE); export GOFLAGS=$(GOFLAGS); go get . ; go mod tidy

#install: @ Install the qleetctl CLI
install: build
	sudo mv ./qleetctl /usr/local/bin/

#release: @ Create and push a new tag
release: build
	$(eval NT=$(NEWTAG))
	@echo -n "Are you sure to create and push ${NT} tag? [y/N] " && read ans && [ $${ans:-N} = y ]
	@git commit -a -s -m "Cut ${NT} release"
	@git tag ${NT}
	@git push origin ${NT}
	@git push
	@echo "Done."

#update: @ Update dependencies to latest versions
update:
	@export GOPRIVATE=$(GOPRIVATE); export GOFLAGS=$(GOFLAGS); go get -u; go mod tidy

#version: @ Print current version(tag)
version:
	@echo $(shell git describe --tags --abbrev=0)
