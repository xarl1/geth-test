package main

import (
	"context"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func getAccount(address string) common.Address {
	return common.HexToAddress(address)
}

func getBalance(client *ethclient.Client, account common.Address) *big.Float {
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	return ethValue
}

//scale back 10**18
func scaleDecimals(token *big.Int) *big.Float {
	fbalance := new(big.Float)
	fbalance.SetString(token.String())
	tokenValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	return tokenValue
}
