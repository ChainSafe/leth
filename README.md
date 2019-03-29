# leth
golang tools for compiling, deploying, testing, and interacting with smart contracts

work in progress ~

currently compiles + deploys; testing is under development.

# dependencies

solc
```
sudo add-apt-repository ppa:ethereum/ethereum
sudo apt-get update
sudo apt-get install solc
```

go-ethereum

`go get github.com/ethereum/go-ethereum`

abigen
```
cd $GOPATH/src/github.com/ethereum/go-ethereum
godep go install ./cmd/abigen
```

# get 

`go get github.com/ChainSafe/leth`

# usage

```
cd $GOPATH/src/github.com/ChainSafe/leth
go build
go install
```

keystore setup: if you have an existing geth keystore, copy the keystore/ directory into this directory.

`leth init` to initialize setup

`leth compile` to compile all contracts in contracts/ directory

`leth deploy` to deploy contracts to network

`leth test` to run the tests as specified in the project main.go

the default network to connect to is `default`. if you wish to connect to another network (as specified in config.json), use `leth deploy --network network_name`

see config.json for example network configurations.

# example 

for an example please see [this repo](https://github.com/noot/token) using leth to write tests for a simple token contract.
