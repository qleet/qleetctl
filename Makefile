#help: @ List available tasks
help:
	@clear
	@echo "Usage: make COMMAND"
	@echo "Commands :"
	@grep -E '[a-zA-Z\.\-]+:.*?@ .*$$' $(MAKEFILE_LIST)| tr -d '#' | awk 'BEGIN {FS = ":.*?@ "}; {printf "\033[32m%-19s\033[0m - %s\n", $$1, $$2}'

#build: @ Build workload controller binary
build:
	go generate
	@export GOFLAGS=$(GOFLAGS); export CGO_ENABLED=0; go build -a -o qleetctl main.go

#install: @ Install the qleetctl CLI
install: build
	sudo mv ./qleetctl /usr/local/bin/

