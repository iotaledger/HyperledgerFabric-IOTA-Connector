package iota

import (
	"log"

	"github.com/iotaledger/iota.go/api"
	"github.com/iotaledger/iota.go/mam/v1"
)

func Fetch(root string, mode string, sideKey string) []string {
  	currentRoot := root

	cm, err := mam.ParseChannelMode(mode)
	if err != nil {
		log.Fatal(err)
	}

	api, err := api.ComposeAPI(api.HTTPClientSettings{
		URI: endpoint,
	})
	if err != nil {
		log.Fatal(err)
	}

	receiver := mam.NewReceiver(api)
	if err := receiver.SetMode(cm, sideKey); err != nil {
		log.Fatal(err)
	}

	var channelMessages []string

	loop:
		nextRoot, messages, err := receiver.Receive(currentRoot)
		if err != nil {
			log.Fatal(err)
		}

		if len(messages) > 0 {
			currentRoot = nextRoot
			channelMessages = append(channelMessages, messages...)
      		goto loop
		}

	return channelMessages
}
