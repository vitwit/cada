#!/usr/bin/make -f

DOCKER := $(shell which docker)

export GO111MODULE = on

###############################################################################
###                                     e2e                                 ###
###############################################################################

ictest-basic:
	@echo "Running basic integration tests"
	@cd interchaintest && go test -race -v -run TestBasicChain .


###############################################################################
###                                  Docker                                 ###
###############################################################################

get-heighliner:
	git clone https://github.com/strangelove-ventures/heighliner.git
	cd heighliner && go install

local-image:
ifeq (,$(shell which heighliner))
	echo 'heighliner' binary not found. Consider running `make get-heighliner`
else
	heighliner build -c rollchain --local -f chains.yaml --go-version 1.22.1
endif

###################
###  Protobuf  ####
###################

protoVer=0.13.2
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating protobuf files..."
	@$(protoImage) sh ./scripts/protocgen.sh
	@go mod tidy

proto-format:
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint

.PHONY: proto-all proto-gen proto-format proto-lint

##################
###  Linting  ####
##################

golangci_lint_cmd=golangci-lint
golangci_version=v1.55.2

lint:
	@echo "--> Running linter"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
	@$(golangci_lint_cmd) run ./... --timeout 15m

lint-fix:
	@echo "--> Running linter and fixing issues"
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@$(golangci_version)
	@$(golangci_lint_cmd) run ./... --fix --timeout 15m

.PHONY: lint lint-fix

GO := go
TARGET := availd
BINDIR ?= $(GOPATH)/bin

.PHONY: all build install clean

all: build

build:
	$(GO) build -o build/$(TARGET)

install: build
	@echo "Installing build/$(TARGET) to $(BINDIR)"
	@cp build/$(TARGET) $(BINDIR)

clean:
	@echo "Cleaning up"
	rm -f $(TARGET)



###############################################################################
###                                    testnet                              ###
###############################################################################

# Run this before testnet keys are added
# chainid-1 is used in the testnet.json
set-testnet-configs:
	availd config set client chain-id chainid-1
	availd config set client keyring-backend test
	availd config set client output text