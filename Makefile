build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/eth-rpc eth-rpc/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/test test/main.go

clean:
	rm eth-rpc bin/eth-rpc
