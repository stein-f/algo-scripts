package main

import (
	"context"
	"fmt"
	"github.com/stein-f/algo-scripts/config"
	"log"
	"strings"
)

const (
	unitNamePrefix = "MNGO"
	account        = "STEINCMH2IQXMU37WR7SJH4WXXUGC2TB35WVMQRH3S5TOZE3VQRZEFJE5E"
)

func main() {
	conf := config.LoadConfig()

	nextToken := ""
	for {
		accountAssetsRes, err := conf.IndexerClient.LookupAccountAssets(account).
			Limit(10).
			Next(nextToken).
			Do(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		for _, asset := range accountAssetsRes.Assets {
			if asset.Amount == 1 {
				_, assetInfo, err := conf.IndexerClient.LookupAssetByID(asset.AssetId).Do(context.Background())
				if err != nil {
					log.Fatal(err)
				}
				if strings.HasPrefix(assetInfo.Params.Creator, unitNamePrefix) {
					fmt.Println(assetInfo.Params.Name, asset.AssetId)
				}
			}
		}

		if accountAssetsRes.NextToken == "" {
			break
		}
		nextToken = accountAssetsRes.NextToken
	}
}
