build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/eth-rpc eth-rpc/*.go
