# Ethereum JSON-RPC on AWS Lambda
AWS Lambda code to call and relay JSON-RPC of ethereum written in Golang.

In addition, this project try porting web3 to Golang.

Furthermore IPFS will be applied to this project to maximize service utility by supporting token development.

# Build
```shell
cd $GOPATH/src/{repo}/eth-rpc
go get
make
```

# Test
1. Move each module directory such as json, rpc and so on
2. Run testunit
```shell
go test -v
```

# Deploy
1. Set Lambda on AWS
  - Function package: compressed binary file in $GOPATH/src/{repo}/eth-rpc/bin
  - Handler: eth-rpc (binary file name, it is optional)
  - Runtime: Go 1.x
2. Set API Gateway as proxy on AWS
3. Add API Gateway as Lambda trigger
4. Add CloudWatch Logs
5. Check logs at CloudWatch console

# Usage
1. JSON-RPC relay
2. Ecrecover
3. ...

# Reference
[1] https://github.com/aws/aws-lambda-go

[2] https://github.com/ethereum/go-ethereum

[3] https://ipfs.io/

[4] https://github.com/ipfs/go-ipfs-api

# License
MIT
