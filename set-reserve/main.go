package main

import (
	"context"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stein-f/algo-scripts/config"
	"os"
)

const (
	reserve = "A3OWTJUKUWDRQ54UTEW4N6U7ALQTZEWG2XGVGAPKL22WKJYBJHFX2SRT4M"
)

var assetIDs = []uint64{
	1092064840,
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	conf := config.LoadConfig()

	txParams, err := conf.AlgodClient.SuggestedParams().Do(context.Background())
	if err != nil {
		panic(err)
	}

	for _, assetID := range assetIDs {
		log.Info().Msgf("Setting reserve. address=%s assetId=%d", conf.Account.Address.String(), assetID)

		transaction, err := future.MakeAssetConfigTxn(
			conf.Account.Address.String(),
			nil,
			txParams,
			assetID,
			conf.Account.Address.String(),
			reserve,
			conf.Account.Address.String(),
			conf.Account.Address.String(),
			false,
		)
		if err != nil {
			panic(err)
		}

		txID, signedTxn, err := crypto.SignTransaction(conf.Account.PrivateKey, transaction)
		if err != nil {
			panic(err)
		}

		_, err = conf.AlgodClient.SendRawTransaction(signedTxn).Do(context.Background())
		if err != nil {
			panic(err)
		}

		_, err = future.WaitForConfirmation(conf.AlgodClient, txID, 4, context.Background())
		if err != nil {
			panic(err)
		}

		log.Info().Msgf("updated reserve %d", assetID)
	}
}
