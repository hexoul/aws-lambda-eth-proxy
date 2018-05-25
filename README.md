# Ethereum JSON-RPC on AWS Lambda
AWS Lambda code to call and relay JSON-RPC of ethereum written in Golang

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
