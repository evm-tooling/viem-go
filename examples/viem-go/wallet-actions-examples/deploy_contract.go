package main

import (
	"fmt"

	"github.com/ChefBingbong/viem-go/actions/public"
	"github.com/ChefBingbong/viem-go/actions/wallet"
)

// RunDeployContract demonstrates deploying a contract via walletClient.DeployContract.
func RunDeployContract() {
	PrintSection("deployContract - Simulate then deploy a contract")

	// Minimal SimpleStorage bytecode
	bytecodeHex := "0x6080604052348015600e575f5ffd5b50604051606f380380606f833981016040819052602b916035565b5f55604b565b5f60208284031215604457005b5051919050565b6018806100575f395ff3fe6080604052005ffea164736f6c634300081d000a"
	abiJSON := `[{"type":"constructor","inputs":[{"name":"_value","type":"uint256"}],"stateMutability":"nonpayable"}]`

	fmt.Printf("Bytecode: %s...\n", Truncate(bytecodeHex, 42))

	// Step 1: Simulate deployment via eth_call (no `to` = contract creation)
	fmt.Println("\nSimulating deployment...")
	result, simErr := public.Call(ctx, PublicClient, public.CallParameters{
		Account: &SourceAddr,
		Data:    hexStringToBytes(bytecodeHex),
	})
	if simErr != nil {
		fmt.Printf("Simulation note: %s\n", Truncate(simErr.Error(), 120))
		fmt.Println("(Expected — account is unfunded on Polygon)")
	} else {
		fmt.Printf("Simulation result length: %d bytes\n", len(result.Data))
	}

	// Step 2: Deploy via walletClient.DeployContract (delegates to wallet.DeployContract)
	fmt.Println("\nDeploying via walletClient.DeployContract...")
	hash, err := WalletCl.DeployContract(ctx, wallet.DeployContractParameters{
		ABI:      abiJSON,
		Bytecode: bytecodeHex,
		Args:     []any{uint64(42)},
	})
	if err != nil {
		fmt.Printf("Deploy error: %s\n", Truncate(err.Error(), 120))
		fmt.Println("(Expected on mainnet with unfunded account)")
	} else {
		fmt.Printf("Deploy tx hash: %s\n", hash)
	}
}

func hexStringToBytes(s string) []byte {
	if len(s) >= 2 && (s[:2] == "0x" || s[:2] == "0X") {
		s = s[2:]
	}
	b := make([]byte, len(s)/2)
	for i := 0; i < len(s); i += 2 {
		b[i/2] = hexCharToByte(s[i])<<4 | hexCharToByte(s[i+1])
	}
	return b
}

func hexCharToByte(c byte) byte {
	switch {
	case c >= '0' && c <= '9':
		return c - '0'
	case c >= 'a' && c <= 'f':
		return c - 'a' + 10
	case c >= 'A' && c <= 'F':
		return c - 'A' + 10
	default:
		return 0
	}
}
