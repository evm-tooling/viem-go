package examples

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/chain/definitions"
	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/contracts/erc20"
	"github.com/ChefBingbong/viem-go/utils/unit"
)

const ERC20_ADDRESS = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"

func RunReadContractExample() error {
	ctx := context.Background()

	publicClient, _ := client.CreatePublicClient(client.PublicClientConfig{
		Chain:     &definitions.Mainnet,
		Transport: transport.HTTP("https://eth.merkle.io"),
	})

	token := client.ReadContractOptions{
		Address: common.HexToAddress(ERC20_ADDRESS),
		ABI:     erc20.ContractABI,
	}

	// Call contract functions
	name, _ := publicClient.ReadContract(ctx, token.WithFunction("name"))
	symbol, _ := publicClient.ReadContract(ctx, token.WithFunction("symbol"))
	decimals, _ := publicClient.ReadContract(ctx, token.WithFunction("decimals"))
	totalSupply, _ := publicClient.ReadContract(ctx, token.WithFunction("totalSupply"))

	// Log token details
	fmt.Printf("\nReading from %s\n\n", ERC20_ADDRESS)
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Symbol: %s\n", symbol)
	fmt.Printf("Decimals: %d\n", decimals)
	fmt.Printf("Total Supply: %s\n\n", totalSupply)

	// Get ERC20 balance
	USER := common.HexToAddress("0x830690922a56f31Cb96851951587D8A2f45C0EBA")

	balance, _ := publicClient.ReadContract(ctx, token.WithFunction("balanceOf", USER))

	// Log balances
	fmt.Printf("Balance Returned: %s\n", balance)
	fmt.Printf("Balance Formatted: %s\n\n", unit.FormatUnits(balance.(*big.Int), int(decimals.(uint8))))

	publicClient.Close()
	return nil
}
