package main

import (
	"fmt"
	"github.com/stein-f/algo-scripts/config"
)

func main() {
	cfg := config.LoadConfig()
	fmt.Println("address: ", cfg.Account.Address.String())
}
