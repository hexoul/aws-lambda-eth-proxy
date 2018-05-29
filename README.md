# Ethereum JSON-RPC on AWS Lambda
AWS Lambda code to call and relay JSON-RPC of ethereum written in Golang.

In addition, this project try porting web3 to Golang.

Furthermore IPFS will be applied to this project to maximize service utility by supporting token development.

# Usage
1. Build
```shell
cd $GOPATH/src/{repo}/eth-rpc
go get
make
```
2. Setting Lambda on AWS with binary file in $GOPATH/src/{repo}/eth-rpc/bin
3. Setting API Gateway as proxy on AWS
4. Link between Lambda and API Gateway

# Test
1. Move each module directory such as json, rpc and so on
2. Run testunit
```shell
go test -v
```

# Reference
[1] https://github.com/aws/aws-lambda-go

[2] https://github.com/ethereum/go-ethereum

[3] https://github.com/ipfs/go-ipfs-api

# License
MIT
