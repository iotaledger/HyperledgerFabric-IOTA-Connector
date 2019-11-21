package iota

import (
	"fmt"
	"github.com/iotaledger/iota.go/address"
	. "github.com/iotaledger/iota.go/api"
	"github.com/iotaledger/iota.go/bundle"
	. "github.com/iotaledger/iota.go/consts"
	"github.com/iotaledger/iota.go/pow"
	"github.com/iotaledger/iota.go/trinary"
)


func TransferTokens(seed string, keyIndex uint64, recipientAddress string) uint64 {

	// get the best available PoW implementation
	_, proofOfWorkFunc := pow.GetFastestProofOfWorkImpl()

	// create a new API instance
	api, err := ComposeAPI(HTTPClientSettings{
		URI: endpoint,
		// (!) if no PoWFunc is supplied, then the connected node is requested to do PoW for us
		// via the AttachToTangle() API call.
		LocalProofOfWorkFunc: proofOfWorkFunc,
	})
	must(err)

	// create a transfer to the given recipient address
	// optionally define a message and tag
	transfers := bundle.Transfers{
		{
			// must be 90 trytes long (include the checksum)
			Address: recipientAddress,
			Value:   DefaultAmount,
		},
	}

	// create inputs for the transfer
	walletAddress, err := address.GenerateAddress(seed, keyIndex, SecurityLevelMedium, true)
	must(err)

	balances, err := api.GetBalances(trinary.Hashes{walletAddress}, 100)
	must(err)
	walletBalance := balances.Balances[0]

	inputs := []Input{
		{
			Address:  walletAddress,
			Security: SecurityLevelMedium,
			KeyIndex: keyIndex,
			Balance:  walletBalance,
		},
	}

	// create an address for the remainder.
	// in this case we will have 20 iotas as the remainder, since we spend 100 from our input
	// address and only send 80 to the recipient.
	remainderAddress, err := address.GenerateAddress(seed, keyIndex + 1, SecurityLevelMedium, true)
	must(err)

	// we don't need to set the security level or timestamp in the options because we supply
	// the input and remainder addresses.
	prepTransferOpts := PrepareTransfersOptions{Inputs: inputs, RemainderAddress: &remainderAddress}

	// prepare the transfer by creating a bundle with the given transfers and inputs.
	// the result are trytes ready for PoW.
	trytes, err := api.PrepareTransfers(seed, transfers, prepTransferOpts)
	must(err)

	// you can decrease your chance of sending to a spent address by checking the address before
	// broadcasting your bundle.
	spent, err := api.WereAddressesSpentFrom(transfers[0].Address)
	must(err)

	if spent[0] {
		fmt.Println("recipient address is spent from, aborting transfer")
		return keyIndex
	}

	// at this point the bundle trytes are signed.
	// now we need to:
	// 1. select two tips
	// 2. do proof-of-work
	// 3. broadcast the bundle
	// 4. store the bundle
	// SendTrytes() conveniently does the steps above for us.
	bndl, err := api.SendTrytes(trytes, depth, mwm)
	must(err)

	fmt.Println("broadcasted bundle with tail tx hash: ", bundle.TailTransactionHash(bndl))
	fmt.Println("remainder address: ", remainderAddress)
	fmt.Println("recipient address: ", recipientAddress)
	fmt.Println("new key index: ", keyIndex + 1)

	return keyIndex + 1
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
