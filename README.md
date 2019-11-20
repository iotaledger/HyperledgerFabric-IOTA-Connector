## Hyperledger Fabric Sample Application

This application demonstrates the creation and transfer of container shipments between actors leveraging Hyperledger Fabric in the supply chain. In this demo app we will set up the minimum number of nodes required to develop chaincode. It has a single peer and a single organization.

if getting error about running ./startFabric.sh permission

chmod a+x startFabric.sh

This code is based on code written by the Hyperledger Fabric community. Source code can be found here: (https://github.com/hyperledger/fabric-samples).


## Cleanup

Stop existing HL Fabric Instance

```
cd basic-network

./stop.sh

./teardown.sh
```


Remove any pre-existing containers, as it may conflict with commands in this project

```
docker rm -f $(docker ps -aq)
```

Remove key store contents

```
cd ~ && rm -rf .hfc-key-store/
```

## Starting the Application

1. Install the required libraries from the package.json file

```
yarn
```

2. Set GOPATH variable

```
export GOPATH=$GOPATH:~/go:~/go/src
```

3. Install Go dependencies

```
go get github.com/cespare/xxhash
go get github.com/pkg/errors
go get github.com/iotaledger/iota.go/address
go get github.com/iotaledger/iota.go/api
go get github.com/iotaledger/iota.go/bundle
go get github.com/iotaledger/iota.go/consts
go get github.com/iotaledger/iota.go/pow
```

4. Start the Hyperledger Fabric network

```
./startFabric.sh
```

5. Register the Admin and User components of our network

```
node registerAdmin.js

node registerUser.js
```

6. Start the client application

```
yarn dev
```

Load the client simply by opening localhost:3000 in any browser window of your choice, and you should see the user interface for our simple application at this URL.
