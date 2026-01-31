package examples

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/chain/definitions"
	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/utils/unit"
)

// RunClientExample demonstrates basic public client usage.
// This mirrors the TypeScript _viem_ts_client.ts example.
func RunClientExample() error {
	ctx := context.Background()

	publicClient, _ := client.CreatePublicClient(client.PublicClientConfig{
		Chain:     &definitions.Mainnet,
		Transport: transport.HTTP("https://eth.merkle.io"),
	})

	blockNumber, _ := publicClient.GetBlockNumber(ctx)
	fmt.Printf("\nCurrent Block Number: %d\n", blockNumber)

	chainID, _ := publicClient.GetChainID(ctx)
	fmt.Printf("Chain ID: %d\n", chainID)

	address := common.HexToAddress("0x73BCEb1Cd57C711feaC4224D062b0F6ff338501e")
	balance, _ := publicClient.GetBalance(ctx, address)

	ethBalance := unit.FormatUnits(balance, 18)
	fmt.Printf("\nETH Balance of %s: %s ETH\n\n", address.Hex(), ethBalance)
	publicClient.Close()

	return nil
}
