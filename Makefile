build:
	go build -o bin/bitcoin-rate-app ./cmd/bitcoin-rate-app

run:
	go run ./cmd/bitcoin-rate-app --config configs/config.yaml

lint:
	golangci-lint run
