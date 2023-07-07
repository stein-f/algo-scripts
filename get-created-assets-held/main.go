package main

import (
	"context"
	"fmt"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/stein-f/algo-scripts/config"
	"log"
)

const (
	account = "PIGMPSL6PU5RA2EJEBPZELJZWQ7R6V4LVDXGPL4NU5QDBCXWFGRJG5M5HY"
)

func main() {
	conf := config.LoadConfig()

	createdAssets := lookupCreatedAssets(conf)
	heldAssets := lookupHeldAssets(conf)

	for _, createdAsset := range createdAssets {
		if holdsAsset(heldAssets, createdAsset) {
			fmt.Printf("id: %v, name: %s\n", createdAsset.Index, createdAsset.Params.Name)
		}
	}
}

func holdsAsset(heldAssets []models.AssetHolding, createdAsset models.Asset) bool {
	for _, heldAsset := range heldAssets {
		if heldAsset.AssetId == createdAsset.Index && heldAsset.Amount > 0 {
			return true
		}
	}
	return false
}

func lookupCreatedAssets(conf config.Config) []models.Asset {
	nextToken := ""
	var createdAssets []models.Asset
	for {
		accountAssetsRes, err := conf.IndexerClient.LookupAccountCreatedAssets(account).
			Limit(1000).
			Next(nextToken).
			Do(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		for _, asset := range accountAssetsRes.Assets {
			createdAssets = append(createdAssets, asset)
		}

		if accountAssetsRes.NextToken == "" {
			break
		}
		nextToken = accountAssetsRes.NextToken
	}
	return createdAssets
}

func lookupHeldAssets(conf config.Config) []models.AssetHolding {
	nextToken := ""
	var assetHoldings []models.AssetHolding
	for {
		accountAssetsRes, err := conf.IndexerClient.LookupAccountAssets(account).
			Limit(1000).
			Next(nextToken).
			Do(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		for _, asset := range accountAssetsRes.Assets {
			assetHoldings = append(assetHoldings, asset)
		}

		if accountAssetsRes.NextToken == "" {
			break
		}
		nextToken = accountAssetsRes.NextToken
	}
	return assetHoldings
}
