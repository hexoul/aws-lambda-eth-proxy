ifeq ($(GOPATH),)
  export GOPATH=$(HOME)/go
endif
ifeq ($(GOPATH),$(GOROOT))
  export GOPATH=$(GOROOT)/../workspace
endif
export GOBIN=$(GOPATH)/bin

$(shell dep ensure)
$(shell go get github.com/karalabe/xgo)

lambda:	
	xgo --deps=https://gmplib.org/download/gmp/gmp-6.0.0a.tar.bz2 \
			--targets=linux/amd64 -out bin/eth-proxy \
			./
	mv bin/eth-proxy-linux-amd64 bin/eth-proxy

local:
	go build -o bin/eth-proxy