# API Reference

## Exported functions

### iota.GenerateRandomSeedString()

The `GenerateRandomSeedString()` method generates a random string of the given length using the following alphabet: "ABCDEFGHIJKLMNOPQRSTUVWXYZ9".

#### Syntax
```
iota.GenerateRandomSeedString(length)
```
#### Parameters
length **int**
Desired length of the seed string.

#### Return values
seed **string** 
Randomly generated string of the given length.

#### Example usage
```
sideKey := iota.GenerateRandomSeedString(50)
```

### iota.PadSideKey()

The `PadSideKey()` method extends short strings with "9" up until the length of 81.

#### Syntax
```
iota.PadSideKey()
```
#### Parameters
shortString **string**
Short string that needs to be extended.

#### Return values
extendedString **string** 
String extended with "9" with the length of 81.

#### Example usage
```
sideKey := iota.PadSideKey(shortKey)
```

### iota.Publish()

The `Publish()` method publishes a given message into a new or existing MAM channel with the given mode and encryption key. Channel state can be reconstructed from the returned transmitter and seed values.

#### Syntax
```
iota.Publish(message, transmitter, mode, sideKey)
```
#### Parameters
message **string**
Message to be published on the Tangle.
transmitter *mam.Transmitter | **nil**
Transmitter object, will be created if nil is provided.
mode **string**
MAM channel mode. Possible values: "public", "private", "restricted".
sideKey **string**
Encryption key for MAM message. Can be an empty string. 

#### Return values
transmitter *mam.Transmitter 
Transmitter object, holds the MAM channel state for eventual publishing of the further messages.
seed **string** 
Randomly generated seed string of the MAM channel. 
root **string** 
Root address of the transaction. Used to fetch MAM message.

#### Example usage
```
transmitter, seed, root := iota.Publish(message, transmitter, mode, sideKey)
channelStateObject := transmitter.Channel()
channelState := iota.MamStateToString(channelStateObject)
```

### iota.PublishAndReturnState()

The `PublishAndReturnState()` method publishes a given message into a new or existing MAM channel with the given mode and encryption key. Channel state is returned to simplify publishing of further messages.

#### Syntax
```
iota.PublishAndReturnState(message, useTransmitter, seed, mamState, mode, sideKey)
```
#### Parameters
message **string**
Message to be published on the Tangle.
useTransmitter **bool**
Flag indicating that existing channel should be used to publish a nem message into. In this case the Transmitter object will be reconstructed from the provided seed and mamState string values. If false, new MAM channel will be created. 
seed **string** 
Seed string of the MAM channel. Empty string if useTransmitter flag is false. 
mamState **string**
Stringified MAM channel state. Empty string if useTransmitter flag is false. 
mode **string**
MAM channel mode. Possible values: "public", "private", "restricted".
sideKey **string**
Encryption key for MAM message. Can be an empty string. 

#### Return values
mamState **string**
Stringified MAM channel state for eventual further messages.
root **string**
Root address of the transaction. Used to fetch MAM message.
seed **string**
Seed string (randomly generated for a new channel) of the MAM channel. 

#### Example usage

**New channel**
```
mamState, root, seed := iota.PublishAndReturnState(string(payloadAsBytes), false, "", "", mode, sideKey)

iotaPayload := IotaPayload{Seed: seed, MamState: mamState, Root: root, Mode: mode, SideKey: sideKey}
iotaPayloadAsBytes, _ := json.Marshal(iotaPayload)
APIstub.PutState("IOTA", iotaPayloadAsBytes)
```
**Appending to an existing channel**
```
iotaPayloadAsBytes, _ := APIstub.GetState("IOTA")
iotaPayload := IotaPayload{}
json.Unmarshal(iotaPayloadAsBytes, &iotaPayload)

mamState, _, _ := iota.PublishAndReturnState(string(payloadAsBytes), true, iotaPayload.Seed, iotaPayload.MamState, iotaPayload.Mode, iotaPayload.SideKey)

iotaPayloadNew := IotaPayload{Seed: iotaPayload.Seed, MamState: mamState, Root: iotaPayload.Root, Mode: iotaPayload.Mode, SideKey: iotaPayload.SideKey}
iotaPayloadNewAsBytes, _ := json.Marshal(iotaPayloadNew)
APIstub.PutState("IOTA", iotaPayloadNewAsBytes)
```

### iota.Fetch()

The `Fetch()` method retrieves messages from the MAM channel. Messages of a restricted channel will be decrypted using provided decryption key.

#### Syntax
```
iota.Fetch(root, mode, sideKey)
```
#### Parameters
root **string** 
Root address of the initial MAM channel transaction.
mode **string**
MAM channel mode. Possible values: "public", "private", "restricted".
sideKey **string**
Decryption key for MAM messages. Can be an empty string. 

#### Return values
messages **[]string**
Array of decoded and decrypted messages from the MAM channel.

#### Example usage
```
iotaPayloadAsBytes, _ := APIstub.GetState("IOTA")
iotaPayload := IotaPayload{}
json.Unmarshal(iotaPayloadAsBytes, &iotaPayload)

messages := iota.Fetch(iotaPayload.Root, iotaPayload.Mode, iotaPayload.SideKey)
```

### iota.CreateWallet()

The `CreateWallet()` method generates and returns a new seed and address of an empty IOTA wallet.

#### Syntax
```
iota.CreateWallet()
```
#### Return values
walletAddress **string**
Root address of the wallet. Used to receive IOTA tokens and balance check.
seed **string** 
Seed string (randomly generated) of the wallet. Used to send IOTA tokens.

#### Example usage
```
walletAddress, walletSeed := iota.CreateWallet()
```

### iota.TransferTokens()

The `TransferTokens()` method transfers IOTA tokens from the existing wallet to a provided address. Transferred amount can be changed in iota.DefaultAmount

#### Syntax
```
iota.TransferTokens(seed, keyIndex, recipientAddress)
```
#### Parameters
seed **string**
Seed string of the own wallet.
keyIndex **uint64**
Index of the current unspent address. Starts with 0 for a new wallet, increases by 1 after each outgoing transaction.
recipientAddress **string**
Wallet address of the recipient. 

#### Return values
keyIndex **uint64**
Index of the next unspent address. 

#### Example usage
```
iotaWalletAsBytes, _ := APIstub.GetState("IOTA_WALLET")
iotaWallet := IotaWallet{}
json.Unmarshal(iotaWalletAsBytes, &iotaWallet)

newKeyIndex := iota.TransferTokens(iotaWallet.Seed, iotaWallet.KeyIndex, participant.Address)

iotaWallet.KeyIndex = newKeyIndex
iotaWalletAsBytes, _ = json.Marshal(iotaWallet)
err = APIstub.PutState("IOTA_WALLET", iotaWalletAsBytes)
```

## Exported constants

All values can be changed in iota/config.go file

### iota.MamMode
Mode of the MAM channel. Possible values: "public", "private", "restricted".
Default value: "public"

### iota.MamSideKey
Encryption key for MAM message. Used to encrypt messages in "restricted" mode. Can be an empty string.
Default value: ""

### iota.DefaultWalletSeed
Seed string of the default wallet to transfer IOTA tokens from. 

### iota.DefaultWalletKeyIndex
Index of the current unspent address of the default wallet to transfer IOTA tokens from. 

### iota.DefaultAmount
Amount that will be sent by default to any recipient.


## Non exported constants

All values can be changed in iota/config.go file

### endpoint
IOTA provider URL. All IOTA-related requests will be directed to this endpoint. Can point to a Mainnet node, Devnet node or a private Tangle node.

### mwm
Difficulty of Proof-of-Work required to attach transaction to the Tangle.
Minimum value on mainnet & spamnet is `14`, `9` on Devnet and private nets.
Default value: 9

### depth
Depth or how far to go for tip selection entry point. How many milestones back to start the random walk from.
Default value: 3

### transactionTag
Default tag added to every data transaction (MAM channel message).
Default value: "HYPERLEDGER"

# Testing and Deployment

### Setting up a Demo Hyperledger Fabric  Application

The demo project used to showcase the IOTA Connector integration is open-source and can be found in this GitHub repository. 

This application demonstrates the creation and transfer of container shipments between actors leveraging Hyperledger Fabric in the supply chain. In this demo app we will set up the minimum number of nodes required to develop chaincode. It has a single peer and a single organization.
This code is based on code written by the Hyperledger Fabric community. Source code can be found here: (https://github.com/hyperledger/fabric-samples).
 
### Starting the Application
1. Install the required libraries from the `package.json` file
```
yarn
```

2. Set `GOPATH` variable
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
go get github.com/iotaledger/iota.go/mam/v1
go get github.com/iotaledger/iota.go/pow
go get github.com/iotaledger/iota.go/trinary
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

7. Load the client by opening localhost:3000 in any browser window of your choice, and you should see the user interface for our simple application at this URL.


### Testing

You can test the chaincode by invoking any of its functions from the inside the Docker container or from the outside of it.

To invoke the chaincode function from the outside, run the following command and provide parameters to the function as an array of strings:
```
docker exec cli peer chaincode invoke -C mychannel -n fabric-demo-app -c '{"Args":["changeContainerHolder", "1", "Thomas"]}'
```

To invoke the chaincode inside the Docker container, do the following:

1. Run 
```
docker exec -it cli bash
```

2. Call a chaincode function
```
peer chaincode invoke -C mychannel -n fabric-demo-app -c '{"Args":["changeContainerHolder", "1", "Thomas"]}'
```

### Logging

The logs for chaincodes are in their respective Docker containers. To inspect logs for a chaincode called fabric-demo-app at version 1.0 on peer0 of org and see the output of any fmt.Print*s, run the following command:
```
docker logs -f dev-peer0.org1.example.com-fabric-demo-app-1.0
```

### Updating the chaincode

1. Stop chaincode Docker container
2. Remove chaincode Docker container
3. Remove chaincode Docker container image
4. Run `./startFabric.sh` to create and deploy new chaincode

### Cleanup
1. Stop existing HL Fabric Instance
```
cd basic-network
./stop.sh
./teardown.sh
```

2. Remove any pre-existing containers and images, as it may conflict with commands in this project
```
docker image rm -f XXXXXXXXXXXX
docker rm -f $(docker ps -aq)
```

3. Remove key store contents
```
cd ~ && rm -rf .hfc-key-store/
```

### Troubleshooting

In case of permission error while running ./startFabric.sh execute the following:
```
chmod a+x startFabric.sh
```
