package main

import (
	"context"
	"fmt"
	"log"

	"geth-test/ierc20"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {

	// Connect to a geth node (when using Infura, you need to use your own API key)
	//client, err := ethclient.Dial("https://cloudflare-eth.com")
	client, err := ethclient.Dial("https://subnets.avax.network/swimmer/mainnet/rpc")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	fmt.Println("we have a connection")
	_ = client // we'll use this in the upcoming sections

	//acc3
	account := common.HexToAddress("0x23784e22B4Cf9662EA21AdAdA29652d6495eb8AE")
	balance, err := client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(balance) // returns TUS holdings - TUS is native to Swimmer

	fmt.Println("test custom functions")
	account3 := getAccount("0x23784e22B4Cf9662EA21AdAdA29652d6495eb8AE")
	acc3_balance := getBalance(client, account3)
	fmt.Println(account3)
	fmt.Println(acc3_balance)

	// Instantiate IERC20 interface to contract address - CRA
	tokenAddress := common.HexToAddress("0xC1a1F40D558a3E82C3981189f61EF21e17d6EB48") //CRA
	token, err := ierc20.NewIerc20(tokenAddress, client)
	if err != nil {
		log.Fatalf("Failed to instantiate a IERC20 Interface: %v", err)
	}

	// Use interface to access token methods
	bal, err := token.BalanceOf(&bind.CallOpts{}, account3)
	if err != nil {
		log.Fatalf("Failed to retrieve balanceOf: %v", err)
	}

	fmt.Printf("cra: %s\n", bal)
	fmt.Println(scaleDecimals(bal)) // "cra: 5.0625"

	//// Retrieve a block by number
	ctx := context.Background()

	// Bind to an already deployed contract -CRA
	//ctr, err := contract.NewContract(tokenAddress, client)

	// Watch for a Deposited event
	watchOpts := &bind.WatchOpts{Context: ctx, Start: nil}
	// Setup a channel for results
	transfers := make(chan *ierc20.Ierc20Transfer)

	// Start a goroutine which watches new events
	sub, err := token.WatchTransfer(watchOpts, transfers, nil, nil)
	//defer sub.Unsubscribe()

	if err != nil {
		fmt.Println("Listining error")
		log.Fatal(err)
	}

	// Receive events from the channel
	//event := <-channel

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case t := <-transfers:
			fmt.Printf("%s -> %s : ", t.From, t.To)
		}
	}
}
