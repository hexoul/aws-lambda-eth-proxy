build:
	dep ensure
	xgo --deps=https://gmplib.org/download/gmp/gmp-6.0.0a.tar.bz2 \
			--targets=linux/amd64 -out bin/eth-rpc \
			./eth-rpc
	mv bin/eth-rpc-linux-amd64 bin/eth-rpc
