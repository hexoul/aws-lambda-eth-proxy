# Ethereum JSON-RPC on AWS Lambda
AWS Lambda code to call and relay JSON-RPC of ethereum written in Golang.

In addition, this project try porting web3 to Golang.

Furthermore IPFS will be applied to this project to maximize service utility by supporting token development.

# Prerequisite
1. Docker
  - Install docker (https://docs.docker.com/install/)
2. xgo
  - because of C compile in go-ethereum, we need improved cross-compiler
  ```shell
  docker pull karalabe/xgo-latest
  go get github.com/karalabe/xgo
  ```

# Build
```shell
cd $GOPATH/src/{repo start with github.com}
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
  - (Optional) Include DynamoDB execution role to Lambda execution role  
2. Set API Gateway as proxy on AWS
3. Add API Gateway as Lambda trigger
4. Add CloudWatch Logs

# Usage
1. JSON-RPC relay with Ethereum node
2. Ecrecover
3. Sign with encrypted private key on DynamoDB
4. IPFS
5. fromWei, toWei written in Golang

# Reference
[1] https://github.com/aws/aws-lambda-go

[2] https://github.com/ethereum/go-ethereum

[3] https://ipfs.io/

[4] https://github.com/ipfs/go-ipfs-api

[5] https://github.com/ethereum/go-ethereum/wiki/Cross-compiling-Ethereum

[6] https://github.com/karalabe/xgo

# License
MIT
