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

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	conf := config.LoadConfig()

	txParams, err := conf.AlgodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		panic(err)
	}

	createTxn, err := future.MakeAssetCreateTxn(
		conf.Account.Address.String(),
		nil,
		txParams,
		10000,
		0,
		false,
		conf.Account.Address.String(),
		"",
		conf.Account.Address.String(),
		conf.Account.Address.String(),
		"TOK",
		"Test Token",
		"template-ipfs://{ipfscid:1:dag-pb:reserve:sha2-256}",
		"",
	)
	if err != nil {
		panic(err)
	}

	txID, signedTxn, err := crypto.SignTransaction(conf.Account.PrivateKey, createTxn)
	if err != nil {
		log.Error().Msgf("Failed to sign transaction: %s", err)
		panic(err)
	}

	sendResponseTxID, err := conf.AlgodClient.SendRawTransaction(signedTxn).Do(context.Background())
	if err != nil {
		log.Error().Msgf("failed to send transaction: %s", err)
		panic(err)
	}

	res, err := future.WaitForConfirmation(conf.AlgodClient, txID, 4, context.Background())
	if err != nil {
		log.Error().Msgf("Error waiting for confirmation on txID: %s", txID)
		panic(err)
	}

	fmt.Printf("created asa opt-in for %d. tx: %s\n", res.AssetIndex, sendResponseTxID)
}
