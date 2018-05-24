build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/eth-rpc eth-rpc/*.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/test test/*.go

clean:
	rm eth-rpc bin/eth-rpc
