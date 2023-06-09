package main

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stein-f/algo-scripts/config"
	"os"
)

var assetsToSend = []uint64{
	1092064840,
}

const recipientAccount = "UCURP277N3KJVSKVWVJX7I4UUQ4DNAQ6ITP3JOPJ77VWJIMJZYVKKXK6WQ"

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	conf := config.LoadConfig()

	txParams, err := conf.AlgodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		panic(err)
	}

	for _, id := range assetsToSend {
		sendTx, err := future.MakeAssetTransferTxn(
			conf.Account.Address.String(),
			recipientAccount,
			1000000,
			nil,
			txParams,
			"",
			id,
		)
		if err != nil {
			panic(err)
		}

		txID, signedTxn, err := crypto.SignTransaction(conf.Account.PrivateKey, sendTx)
		if err != nil {
			log.Error().Msgf("Failed to sign transaction: %s", err)
			panic(err)
		}

		sendResponseTxID, err := conf.AlgodClient.SendRawTransaction(signedTxn).Do(context.Background())
		if err != nil {
			log.Error().Msgf("failed to send transaction: %s", err)
			panic(err)
		}

		_, err = future.WaitForConfirmation(conf.AlgodClient, txID, 4, context.Background())
		if err != nil {
			log.Error().Msgf("Error waiting for confirmation on txID: %s", txID)
			panic(err)
		}

		fmt.Printf("completed send for %d. tx: %s\n", id, sendResponseTxID)
	}
}
