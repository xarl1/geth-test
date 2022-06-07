package main

import (
	"context"
	"fmt"
	"geth-test/dai"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/quan8/go-ethereum/accounts/abi"
)

func main() {

	// Connect to a geth node (when using Infura, you need to use your own API key)
	//client, err := ethclient.Dial("https://cloudflare-eth.com")
	client, err := ethclient.Dial("wss://eth-mainnet.alchemyapi.io/v2/l4hAoiQSbj_QfT6N763bygpPEtF9vvgM")
	if err != nil {#
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	fmt.Println("we have a connection")
	defer client.Close()

	// Instantiate IERC20 interface to contract address - UniswapV2Router02
	//uniswapAddress := common.HexToAddress("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D")
	daiAddress := common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F") //DAI
	query := ethereum.FilterQuery{
		Addresses: []common.Address{daiAddress},
	}

	logs := make(chan types.Log)

	sub, err := client.SubscribeFilterLogs(context.Background(), query, logs)
	if err != nil {
		log.Fatal(err)
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(dai.DaiABI)))
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case vLog := <-logs:
			fmt.Println("receive log: ", vLog.BlockNumber, vLog.Index)
			//fmt.Println(vLog) // pointer to event log
			//fmt.Println("...")
			//fmt.Println("BlockHash: ", vLog.BlockHash.Hex())
			//fmt.Println("...")
			fmt.Println("Tx Hash: ", vLog.TxHash.Hex())
			fmt.Println("...")
			fmt.Println("Topics: ", vLog.Topics)
			fmt.Println("...")
			fmt.Println("Address: ", vLog.Address)
			//fmt.Println("...")
			//fmt.Println("Data: ", vLog.Data)

			event := struct {
				From common.Address
				To   common.Address
				Wad  *big.Int
			}{}

			//Unpack() will not parse indexed event types because those are stored under topics - From and TO are indexed.
			// is that why from and to are 0x00000...
			err := contractAbi.Unpack(&event, "Transfer", vLog.Data)
			if err != nil {
				fmt.Println(err) //                log.Println("Failed to unpack")
				log.Fatal(err)
			}
			fmt.Println("DECODED")
			//fmt.Println(string(event.From[:])) // foo
			//fmt.Println(string(event.To[:]))   // bar
			//fmt.Println("From:", event.From.Hex())
			//fmt.Println("To:", event.To)
			event.From = common.HexToAddress(vLog.Topics[1].Hex())
			event.To = common.HexToAddress(vLog.Topics[2].Hex())

			fmt.Println("From:", event.From.Hex())
			fmt.Println("To:", event.To.Hex())
			fmt.Println("Amount:", event.Wad)
			fmt.Println("DECODED END")

		}
	}

}
