build:
	dep ensure
#	env CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o bin/eth-rpc eth-rpc/main.go
	xgo --deps=https://gmplib.org/download/gmp/gmp-6.0.0a.tar.bz2 \
			--targets=linux/amd64 \
			github.com/hexoul/eth-rpc-on-aws-lambda/eth-rpc
	mkdir -p bin
	mv eth-rpc-linux-amd64 bin/eth-rpc
