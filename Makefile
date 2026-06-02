BINARY := lalachaind
BUILD_DIR := ./build
GOMOD := github.com/lala-protocol/lalachain/chain
CHAIN_ID := lalachain-testnet-1

LDFLAGS := -X github.com/cosmos/cosmos-sdk/version.Name=lalachain \
           -X github.com/cosmos/cosmos-sdk/version.AppName=$(BINARY) \
           -X github.com/cosmos/cosmos-sdk/version.Version=0.1.0 \
           -X github.com/cosmos/cosmos-sdk/version.Commit=$(shell git rev-parse --short HEAD 2>/dev/null || echo "dev")

.PHONY: build install test lint proto clean docker testnet

build:
	cd chain && go build -ldflags '$(LDFLAGS)' -o ../$(BUILD_DIR)/$(BINARY) ./cmd/lalachaind

install:
	cd chain && go install -ldflags '$(LDFLAGS)' ./cmd/lalachaind

test:
	cd chain && go test ./...

lint:
	cd chain && golangci-lint run ./...

proto:
	cd chain/proto && buf generate

clean:
	rm -rf $(BUILD_DIR) testnet

docker:
	docker build -t lalachain:latest .

docker-testnet:
	docker-compose up --build -d

testnet:
	bash scripts/init-testnet.sh

proto-lint:
	cd chain/proto && buf lint

proto-format:
	cd chain/proto && buf format -w

tidy:
	cd chain && go mod tidy

simulation:
	cd simulation && python -X utf8 simulation.py
