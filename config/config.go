package config

import (
	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/client/v2/indexer"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type Config struct {
	Account       crypto.Account
	AlgodClient   *algod.Client
	IndexerClient *indexer.Client
}

func LoadConfig() Config {
	yamlFile, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Fatal(err)
	}

	var cfg config
	if err = yaml.Unmarshal(yamlFile, &cfg); err != nil {
		log.Fatal(err)
	}

	config := Config{}

	if cfg.Passphrase != "" {
		privateKey, err := mnemonic.ToPrivateKey(cfg.Passphrase)
		if err != nil {
			log.Fatal(err)
		}
		account, err := crypto.AccountFromPrivateKey(privateKey)
		if err != nil {
			log.Fatal(err)
		}
		config.Account = account
	}

	config.AlgodClient, err = algod.MakeClient(cfg.AlgodURL, "")
	if err != nil {
		log.Fatal(err)
	}

	config.IndexerClient, err = indexer.MakeClient(cfg.IndexerURL, "")
	if err != nil {
		log.Fatal(err)
	}

	return config
}

type config struct {
	Passphrase string `yaml:"passphrase"`
	AlgodURL   string `yaml:"algodUrl"`
	IndexerURL string `yaml:"indexerUrl"`
}
