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

var assetsToOptIn = []uint64{
	871152580,
	870720719,
	386956116,
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	conf := config.LoadConfig()

	txParams, err := conf.AlgodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		panic(err)
	}

	for _, id := range assetsToOptIn {
		optInTxn, err := future.MakeAssetTransferTxn(
			conf.Account.Address.String(),
			conf.Account.Address.String(),
			0,
			nil,
			txParams,
			"",
			id,
		)
		if err != nil {
			panic(err)
		}

		txID, signedTxn, err := crypto.SignTransaction(conf.Account.PrivateKey, optInTxn)
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

		fmt.Printf("completed opt-in for %d. tx: %s\n", id, sendResponseTxID)
	}
}
