# API Reference

## Exported functions

### iota.GenerateRandomSeedString()

The `GenerateRandomSeedString()` method generates a random string of the given length using the following alphabet: "ABCDEFGHIJKLMNOPQRSTUVWXYZ9".

#### Syntax
```
iota.GenerateRandomSeedString(length)
```
#### Parameters
* length **int**  
Desired length of the seed string.

#### Return values
* seed **string**  
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
* shortString **string**  
Short string that needs to be extended.

#### Return values
* extendedString **string**  
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
* message **string**  
Message to be published on the Tangle.
* transmitter *mam.Transmitter | **nil**  
Transmitter object, will be created if nil is provided.
* mode **string**  
MAM channel mode. Possible values: "public", "private", "restricted".
* sideKey **string**  
Encryption key for MAM message. Can be an empty string. 

#### Return values
* transmitter *mam.Transmitter  
Transmitter object, holds the MAM channel state for eventual publishing of the further messages.
* seed **string**  
Randomly generated seed string of the MAM channel. 
* root **string**  
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
* message **string**  
Message to be published on the Tangle.
* useTransmitter **bool**  
Flag indicating that existing channel should be used to publish a nem message into. In this case the Transmitter object will be reconstructed from the provided seed and mamState string values. If false, new MAM channel will be created. 
* seed **string**  
Seed string of the MAM channel. Empty string if useTransmitter flag is false. 
* mamState **string**  
Stringified MAM channel state. Empty string if useTransmitter flag is false. 
* mode **string**  
MAM channel mode. Possible values: "public", "private", "restricted".
* sideKey **string**  
Encryption key for MAM message. Can be an empty string. 

#### Return values
* mamState **string**  
Stringified MAM channel state for eventual further messages.
* root **string**  
Root address of the transaction. Used to fetch MAM message.
* seed **string**  
Seed string (randomly generated for a new channel) of the MAM channel. 

#### Example usage

New channel
```
mamState, root, seed := iota.PublishAndReturnState(string(payloadAsBytes), false, "", "", mode, sideKey)

iotaPayload := IotaPayload{Seed: seed, MamState: mamState, Root: root, Mode: mode, SideKey: sideKey}
iotaPayloadAsBytes, _ := json.Marshal(iotaPayload)
APIstub.PutState("IOTA", iotaPayloadAsBytes)
```
Appending to an existing channel
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
* root **string**  
Root address of the initial MAM channel transaction.
* mode **string**  
MAM channel mode. Possible values: "public", "private", "restricted".
* sideKey **string**  
Decryption key for MAM messages. Can be an empty string. 

#### Return values
* messages **[]string**  
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
* walletAddress **string**  
Root address of the wallet. Used to receive IOTA tokens and balance check.
* seed **string**  
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
* seed **string**  
Seed string of the own wallet.
* keyIndex **uint64**  
Index of the current unspent address. Starts with 0 for a new wallet, increases by 1 after each outgoing transaction.
* recipientAddress **string**  
Wallet address of the recipient. 

#### Return values
* keyIndex **uint64**  
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


How to Integrate the IOTA Connectors
We will start with a simple chaincode file suitable for a supply chain project and will add a few IOTA connectors one by one with a short usage description.
The plain Hyperledger Fabric chaincode written in Go has a namespace, which can be defined is the docker-compose.yml file under volumes:  

![image](https://user-images.githubusercontent.com/18738760/69976583-82c4cf00-1529-11ea-82ca-defbaf1ec381.png)


In our case, the chaincode folder is mapped to the github.com namespace. All external packages, including IOTA connector, need to be imported using the same namespace.
We start by adding the IOTA connector code into the chaincode folder.  
![image](https://user-images.githubusercontent.com/18738760/69976604-8ce6cd80-1529-11ea-8bae-1d5b24487a76.png)


In addition, you will need to add 3 other complementary packages.
You can do this either by adding an instruction to the shell script that installs and starts the Hyperledger instance or by manually downloading and copying the projects into the chaincode folder.
To add instructions to the shell script, add the following three lines before the “chaincode install” command  
![image](https://user-images.githubusercontent.com/18738760/69976659-a425bb00-1529-11ea-8ed7-e6006395e432.png)


This will ensure that required packages are downloaded and placed into the chaincode folder, mapped as github.com
Depending on your version of Go, the last command that downloads the iota.go library might fail. In this case, you can download it manually and place it into the chaincode folder.  
![image](https://user-images.githubusercontent.com/18738760/69976675-aab43280-1529-11ea-8b48-36ca1d150424.png)


Next, add import of the IOTA connector library to the list of imports of your chaincode.go file. Please note, gibhub.com namespace in from of the iota package name. If your organization uses a different namespace, please adjust accordingly.  
![image](https://user-images.githubusercontent.com/18738760/69976686-b0117d00-1529-11ea-8897-1d144c1d4120.png)


Digital Twin
Hyperledger chaincode contains a function called initLedger, where you can define an initial structure of the ledger and store data as key/value pairs. For complex types, you would typically use structures to describe object fields.
In our example, the initial ledger structure consists of a number of container definitions, where each container includes the Holder field among many others  
![image](https://user-images.githubusercontent.com/18738760/69976699-b56ec780-1529-11ea-8683-d9abf9053906.png)


Once the structure is defined, it is stored on a ledger one by one  
![image](https://user-images.githubusercontent.com/18738760/69976716-bbfd3f00-1529-11ea-97ab-c6c1e45eee60.png)


This is one possible way to store data on a ledger. We do not claim to develop a best possible option.
We will add the IOTA connector here, to ensure that the digital twin is stored on the Tangle in MAM message streams. We will create an individual message stream for each container, where all changes to the container data, location, holder etc. will be tracked.
First, define the structure of the IOTA MAM stream. It should contain randomly generated Seed, MamState object, Root address to fetch data from, stream mode and sideKey (encryption key).  
![image](https://user-images.githubusercontent.com/18738760/69976727-c1f32000-1529-11ea-9634-0352555962b4.png)


Then, define mode and encryption key for every asset, that will get a digital twin on the Tangle.
You can use one of the following modes: “public”, “private”, “restricted”.
SideKey (encryption key) is only required for the “restricted” mode. Otherwise, it can remain an empty string.
We provide two helper functions for the sideKey. Both functions are accessible from the iota namespace:
iota.GenerateRandomSeedString(length) will generate a random string of the given length. It can be used as seed or encryption key.
iota.PadSideKey() will automatically adjust the length of the short key to 81 characters.
If you do not want to define your own values for mode and sideKey, you can use default values iota.MamMode and iota.MamSideKey, which you can inspect and modify under chaincode/iota/config.go
After that, you can call a function iota.PublishAndReturnState(), which will publish the message payload as a new MAM channel. This function returns the mamState, root and seed values.
MamState and seed are needed for further appends to the same channel. Root value is used to read data from the channel.
These values need to be stored on the ledger and communicated to each peer of the organization.
If you do not want to append new messages to the same channel in the future, you can call iota.Publish() function instead, which won’t return the MamState.
iota.PublishAndReturnState() requires the following parameters:
Message payload (string)
Existing MAM stream flag (bool)
Existing seed value (string)
Existing MamState value (string)
Mode (string)
sideKey (string)
If you create a new MAM stream, set values for existing MAM stream, seed and mamState to false, “”, “”
Once the message was published, it is time to persist values on the ledger. Create a new object of type IotaPayload and put it on the ledger, similarly as you created records for container assets.
Please note that we recommend to add a prefix like “IOTA_” in front of the asset ID.
APIstub.PutState(“IOTA_” + strconv.Itoa(i+1), iotaPayloadAsBytes)  
![image](https://user-images.githubusercontent.com/18738760/69976758-cd464b80-1529-11ea-9e43-ec52e8205dcb.png)


If your smart contract contains a function to add new records to the ledger, please update this function by adding the IOTA connector code to it.
Please note that the ID of the new asset is used for the new IOTA object
APIstub.PutState(“IOTA_” + args[0], iotaPayloadAsBytes)  
![image](https://user-images.githubusercontent.com/18738760/69976772-d46d5980-1529-11ea-8bbe-2ba3933a807e.png)


Similarly, if your smart contract contains a function to modify existing records, please update this function by adding the IOTA connector code to it as well.
Please note that in this case, we are appending the modified asset data to the existing MAM stream of this asset. Therefore we retrieve the iotaPayload for specific assets from the ledger, and communicate this information to the iota.PublishAndReturnState() function, along with the existing MAM stream flag set to true.
Also note, at this point, the mode and encryption key of the existing MAM stream can not be changed. Please use existing values stored on the ledger.
Once the new message was added to the existing stream, the new MamState should replace the previous state on the ledger. All other values should remain the same as before.  
![image](https://user-images.githubusercontent.com/18738760/69976785-dd5e2b00-1529-11ea-9b28-08357cbf4ca4.png)


If your smart contract contains a function to query a specific record from the ledger, you might want to also fetch and return the state stored on the Tangle. To perform the query from the Tangle, you can use the following function available from the IOTA connector code:
iota.Fetch(root, mode, sideKey)
Please note that the ID of the asset is used to retrieve the corresponding IOTA object from the ledger.
iotaPayloadAsBytes, _ := APIstub.GetState(“IOTA_” + args[0])  
![image](https://user-images.githubusercontent.com/18738760/69976798-e2bb7580-1529-11ea-8c6a-3fb2564f90d7.png)


Fetched messages are returned as an array of strings. If you want to join them together into one string and add to the output object, you will also need to import the “strings” Go package.
In addition, if you want to output the MamState values, to be able to perform other actions on the UI, like confirm and compare data from the ledger with data stored on the Tangle, you can add the MamState object to the output  
![image](https://user-images.githubusercontent.com/18738760/69976811-e818c000-1529-11ea-9f7b-dcece5f14ff8.png)


Sending Payments
Supply chain projects might benefit from the built-in payment solution, where payments for certain services and goods can be sent between supply chain participants.
One possible use case could be when retailers or end consumers send payments to the producers, logistics and fulfillment providers for the ordered assets.
Since none of the Hyperledger projects support cryptocurrency or any other type of payments, IOTA connector can be used to perform fee-less payments between participants at the moment where a smart contract confirms successful transaction.
To send payment using the IOTA wallet, you will need to store wallet seed and keyIndex on the ledger. Seed is used to initiate a transaction, and keyIndex is specific for IOTA implementation and represents the index of the current wallet address, which holds the tokens.
After every outgoing payment transaction, the remaining tokens are transferred to the next address in order to prevent double-spending. The index of the new address (called remainderAddress) should be stored on the ledger and used for the next outgoing payment. Incoming payments do not trigger address or index change.
In the example, we will maintain only one wallet for outgoing payments. This wallet will be assigned to the retailer, who is the end consumer of the asset in this supply chain project.
The payment will be sent to the previous asset holder each time the holder is changed, which indicates asset movement towards the end consumer. In other words, once a producer prepares a container for shipment and transfer it over to a freight forwarder, the retailer will pay to the producer in IOTA tokens. Then, once a freight forwarder delivers the container to the next destination, the retailer will transfer IOTA tokens to pay for this service.
All possible participants in this sample of a supply chain project are defined upon initialisation of the ledger. For simplicity, we assume that there is only one participant with the role of “Producer”, “Shipper” and so on.  
![image](https://user-images.githubusercontent.com/18738760/69976831-f1099180-1529-11ea-89c5-981cf0f36ff9.png)


We will start with the definition of the structure for the wallet object  
![image](https://user-images.githubusercontent.com/18738760/69976842-f535af00-1529-11ea-84b9-7d5f117d17da.png)


This structure contains seed and keyIndex as described above. In addition, it also contains the actual address where tokens are currently stored. You can perform balance check to ensure sufficient balance of the wallet. Enter your wallet address on this page to check the current balance.
Next, we will extend the existing Participant structure by adding the IotaWallet part into it  
![image](https://user-images.githubusercontent.com/18738760/69976856-fd8dea00-1529-11ea-92e3-6406f4cb8982.png)


And then we will generate a new empty wallet and add wallet information to every participant record.
To generate a new wallet, you can use a function from the IOTA connector:
walletAddress, walletSeed := iota.CreateWallet()  
![image](https://user-images.githubusercontent.com/18738760/69976864-02529e00-152a-11ea-8bde-1d5ebfb44f22.png)


Since we generate a new wallet for every record, we set the keyIndex value to 0. If you are about to use existing wallets, please adjust keyIndex values accordingly.
The generated wallets are empty and currently can only receive IOTA tokens.
In order to send tokens, you need to maintain at least one wallet funded with IOTA tokens.
Wallet data of this wallet should be stored on the ledger.
You can provide wallet data upon ledger initialisation, or you can modify values in the configuration file under chaincode/iota/config.go  
![image](https://user-images.githubusercontent.com/18738760/69976875-08e11580-152a-11ea-8db0-9437927f24ff.png)


To store a wallet on the ledger, please update the initLedger() function by adding the following code. You can replace values for iota.DefaultWalletSeed and iota.DefaultWalletKeyIndex with respective values of your wallet.  
![image](https://user-images.githubusercontent.com/18738760/69976897-10a0ba00-152a-11ea-8151-fdb810242605.png)


Once the wallets are configured, we can add functionality to perform payments. This consists of 3 simple steps:
Identify function in your smart contract which should trigger payment.
Identify the payment recipient. Retrieve wallet address of the recipient.
Identify the payment sender. Retrieve wallet seed and keyIndex of the sender. Perform token transfer, then update keyIndex to the new value and store it on the ledger.
In our example, we will perform payments once the asset holder was changed. So, the function that triggers payments called changeContainerHolder(). Payment recipient is the previous container holder. So, we need to preserve the holder value before changing, to be able to retrieve the corresponding wallet data from the ledger.
On the screenshot below you see the required updates to the function. Container data is retrieved based on the provided ID. Before the holder is reassigned, we store the original holder value. Later we query the ledger in order to get the wallet address of the original container holder.  
![image](https://user-images.githubusercontent.com/18738760/69976913-16969b00-152a-11ea-9fe0-02324fbdc043.png)


For step 3 we will request the IOTA wallet data of the retailer. Then we will trigger the following function and submit seed and keyIndex values of the sender and address value of the recipient.
iota.TransferTokens(seed, keyIndex, address)
Once this is done, we just need to update keyIndex to the new value and store it on the ledger.  
![image](https://user-images.githubusercontent.com/18738760/69976930-1c8c7c00-152a-11ea-9597-0046423f8dfe.png)


Token transfer usually requires a few seconds in order to be confirmed. Please do not attempt to trigger multiple transfers from the same wallet within a very short timeframe (less than 10 seconds), as it will result into “Invalid balance” error, and wallet keyIndex should be reset to previous value manually.
As always, you can check the status of the token transfer on this page by entering the wallet address.
Following the open sourcing of the IOTA Connector, we will be seeding the capabilities as a Hyperledger Bridge to the Linux Foundation. At the current time we’re working with the Linux Foundation to ensure the code is being seeding into the project that it will be most relevant for and look to finalize that decision in the coming weeks.
Resources
The source code of the IOTA Connector can be downloaded from this open-source GitHub repository.
The supplementary iota.go library, that need to be located within the same chaincode folder, can be downloaded from this open-source GitHub repository.
The demo project used to showcase the IOTA Connector integration is open-source and can be found in this GitHub repository.
