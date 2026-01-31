package examples

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ChefBingbong/viem-go/client"
	"github.com/ChefBingbong/viem-go/client/transport"
	"github.com/ChefBingbong/viem-go/contract"
	"github.com/ChefBingbong/viem-go/contracts/erc20"
)

// USDC on Ethereum mainnet
var usdcAddress = common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")

func main() {
	// Create a public client
	c, err := client.CreatePublicClient(client.PublicClientConfig{
		Transport: transport.HTTP("https://eth.llamarpc.com"),
	})
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	owner := common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045") // vitalik.eth

	fmt.Println("=== Tier 1: Generic ReadContract[T] ===")
	tier1Example(c, ctx, owner)

	fmt.Println("\n=== Tier 2: Typed Method Descriptors ===")
	tier2Example(c, ctx, owner)

	fmt.Println("\n=== Tier 3: Classic Typed Bindings ===")
	tier3Example(c, ctx, owner)

	fmt.Println("\n=== Type Safety Demonstration ===")
	typeSafetyDemo()
}

// =============================================================================
// Tier 1: Generic ReadContract[T]
// =============================================================================
//
// Use this for one-off calls where you specify the return type explicitly.
// Type safety: Return type is enforced at compile time.

func tier1Example(c *client.PublicClient, ctx context.Context, owner common.Address) {
	// The type parameter [string] tells Go what type to return
	// If the ABI function doesn't return a string, you'll get a runtime error
	name, err := contract.ReadContract[string](c, contract.ReadContractParams{
		Address:      usdcAddress,
		ABI:          erc20.ContractABI,
		FunctionName: "name",
	})
	if err != nil {
		log.Printf("Error reading name: %v", err)
		return
	}
	fmt.Printf("Token name: %s\n", name)

	// For balanceOf, we specify [*big.Int] as the return type
	balance, err := contract.ReadContract[*big.Int](c, contract.ReadContractParams{
		Address:      usdcAddress,
		ABI:          erc20.ContractABI,
		FunctionName: "balanceOf",
		Args:         []any{owner},
	})
	if err != nil {
		log.Printf("Error reading balance: %v", err)
		return
	}
	fmt.Printf("Balance: %s\n", balance.String())

	// Decimals returns uint8
	decimals, err := contract.ReadContract[uint8](c, contract.ReadContractParams{
		Address:      usdcAddress,
		ABI:          erc20.ContractABI,
		FunctionName: "decimals",
	})
	if err != nil {
		log.Printf("Error reading decimals: %v", err)
		return
	}
	fmt.Printf("Decimals: %d\n", decimals)
}

// =============================================================================
// Tier 2: Typed Method Descriptors
// =============================================================================
//
// Use this when you want compile-time type checking for both arguments AND return types.
// This is the most type-safe approach for reusable contract interactions.

func tier2Example(c *client.PublicClient, ctx context.Context, owner common.Address) {
	// Bind the contract once
	token, err := contract.Bind(usdcAddress, []byte(erc20.ContractABI), c)
	if err != nil {
		log.Fatal(err)
	}

	// erc20.Methods.Name is typed as Fn[string]
	// contract.Call returns string - the compiler knows this!
	name, err := contract.Call(token, ctx, erc20.Methods.Name)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Token name: %s\n", name)

	// erc20.Methods.BalanceOf is typed as Fn1[common.Address, *big.Int]
	// - First type param (common.Address) = argument type
	// - Second type param (*big.Int) = return type
	// The compiler enforces that 'owner' must be common.Address
	balance, err := contract.Call1(token, ctx, erc20.Methods.BalanceOf, owner)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Balance: %s\n", balance.String())

	// erc20.Methods.Allowance is typed as Fn2[common.Address, common.Address, *big.Int]
	// Both arguments must be common.Address, returns *big.Int
	spender := common.HexToAddress("0x1111111254EEB25477B68fb85Ed929f73A960582") // 1inch router
	allowance, err := contract.Call2(token, ctx, erc20.Methods.Allowance, owner, spender)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Allowance: %s\n", allowance.String())
}

// =============================================================================
// Tier 3: Classic Typed Bindings
// =============================================================================
//
// Use this for the most ergonomic API with fully typed methods.
// This is similar to go-ethereum's abigen output.

func tier3Example(c *client.PublicClient, ctx context.Context, owner common.Address) {
	// Create a typed ERC20 instance
	token, err := erc20.New(usdcAddress, c)
	if err != nil {
		log.Fatal(err)
	}

	// All methods are fully typed - no generics needed at call site
	name, err := token.Name(ctx)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Token name: %s\n", name)

	// BalanceOf takes common.Address, returns *big.Int
	balance, err := token.BalanceOf(ctx, owner)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Balance: %s\n", balance.String())

	// You can also access the Bound() method to use Tier 2 style
	decimals, err := contract.Call(token.Bound(), ctx, erc20.Methods.Decimals)
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	fmt.Printf("Decimals: %d\n", decimals)
}

// =============================================================================
// Type Safety Demonstration
// =============================================================================
//
// This section shows what the compiler catches vs what it doesn't.

func typeSafetyDemo() {
	fmt.Println("Type safety examples (these are compile-time checks):")
	fmt.Println()

	// -------------------------------------------------------------------------
	// TIER 2 TYPE SAFETY - Compile-time argument checking
	// -------------------------------------------------------------------------

	fmt.Println("Tier 2 (Fn descriptors) provides compile-time argument type checking:")
	fmt.Println()

	// This is the method descriptor - it encodes both input and output types
	fmt.Println("  var BalanceOf = contract.Fn1[common.Address, *big.Int]{Name: \"balanceOf\"}")
	fmt.Println()

	fmt.Println("  // ✅ COMPILES: argument is common.Address")
	fmt.Println("  balance, err := contract.Call1(token, ctx, BalanceOf, owner)")
	fmt.Println()

	fmt.Println("  // ❌ WON'T COMPILE: argument is string, not common.Address")
	fmt.Println("  balance, err := contract.Call1(token, ctx, BalanceOf, \"0x123...\")")
	fmt.Println("  // Error: cannot use \"0x123...\" (type string) as common.Address")
	fmt.Println()

	fmt.Println("  // ❌ WON'T COMPILE: wrong number of arguments")
	fmt.Println("  balance, err := contract.Call1(token, ctx, BalanceOf)")
	fmt.Println("  // Error: not enough arguments")
	fmt.Println()

	fmt.Println("  // ❌ WON'T COMPILE: assigning to wrong type")
	fmt.Println("  var name string = contract.Call1(token, ctx, BalanceOf, owner)")
	fmt.Println("  // Error: cannot use *big.Int as string")
	fmt.Println()

	// -------------------------------------------------------------------------
	// COMPARISON WITH UNTYPED API
	// -------------------------------------------------------------------------

	fmt.Println("Compare with the untyped API (runtime errors only):")
	fmt.Println()

	fmt.Println("  // ⚠️ COMPILES but fails at runtime: wrong argument type")
	fmt.Println("  result, err := contract.Read(ctx, \"balanceOf\", \"not-an-address\")")
	fmt.Println()

	fmt.Println("  // ⚠️ COMPILES but fails at runtime: wrong number of arguments")
	fmt.Println("  result, err := contract.Read(ctx, \"balanceOf\")")
	fmt.Println()

	fmt.Println("  // ⚠️ COMPILES but panics: type assertion fails")
	fmt.Println("  balance := result[0].(string)  // runtime panic!")
	fmt.Println()

	// -------------------------------------------------------------------------
	// DEFINE YOUR OWN TYPED METHODS
	// -------------------------------------------------------------------------

	fmt.Println("You can define typed methods for any contract:")
	fmt.Println()
	fmt.Println("  // Uniswap V2 Pair example")
	fmt.Println("  var UniswapPair = struct {")
	fmt.Println("      GetReserves contract.Fn[PairReserves]")
	fmt.Println("      Token0      contract.Fn[common.Address]")
	fmt.Println("      Token1      contract.Fn[common.Address]")
	fmt.Println("      Swap        contract.FnWrite4[*big.Int, *big.Int, common.Address, []byte]")
	fmt.Println("  }{")
	fmt.Println("      GetReserves: contract.Fn[PairReserves]{Name: \"getReserves\"},")
	fmt.Println("      Token0:      contract.Fn[common.Address]{Name: \"token0\"},")
	fmt.Println("      Token1:      contract.Fn[common.Address]{Name: \"token1\"},")
	fmt.Println("      Swap:        contract.FnWrite4[...]{Name: \"swap\"},")
	fmt.Println("  }")
	fmt.Println()
	fmt.Println("  // Usage - fully type-safe!")
	fmt.Println("  reserves, err := contract.Call(pair, ctx, UniswapPair.GetReserves)")
	fmt.Println("  // reserves.Reserve0, reserves.Reserve1, reserves.BlockTimestamp")
}

// PairReserves is an example struct for multi-return values
type PairReserves struct {
	Reserve0           *big.Int
	Reserve1           *big.Int
	BlockTimestampLast uint32
}
