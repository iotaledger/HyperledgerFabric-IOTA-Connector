package iota

import (
	"log"

	"github.com/iotaledger/iota.go/mam/v1"
)

func PublishAndReturnState(message string, useTransmitter bool, seedFromStorage string, mamStateFromStorage string, mode string, sideKey string) (string, string, string) {
	var t *mam.Transmitter = nil
	if useTransmitter == true {
		mamState := StringToMamState(mamStateFromStorage)
		t = ReconstructTransmitter(seedFromStorage, mamState)
	}

	transmitter, seed, root := Publish(message, t, mode, sideKey)
	channel := transmitter.Channel()

	return MamStateToString(channel), root, seed
}

func Publish(message string, t *mam.Transmitter, mode string, sideKey string) (*mam.Transmitter, string, string) {
	transmitter, seed := GetTransmitter(t, mode, sideKey)

	root, err := transmitter.Transmit(message, transactionTag)
	if err != nil {
		log.Fatal(err)
	}

	return transmitter, seed, root
}
