package bench

import (
	"testing"

	"github.com/ChefBingbong/viem-go/actions/public"
)

// BenchmarkCall_Basic benchmarks a simple contract call reading the token name.
// This is the most basic call scenario with minimal parameters.
//
// Equivalent to TypeScript:
//
//	await client.call({ to: USDC, data: '0x06fdde03' })
func BenchmarkCall_Basic(b *testing.B) {
	params := public.CallParameters{
		To:   &usdcAddress,
		Data: nameSelector,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := public.Call(benchCtx, benchClient, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCall_WithData benchmarks a call with encoded function parameters.
// This tests the balanceOf(address) call which requires parameter encoding.
//
// Equivalent to TypeScript:
//
//	await client.call({
//	  to: USDC,
//	  data: encodeFunctionData({ abi: erc20Abi, functionName: 'balanceOf', args: [address] })
//	})
func BenchmarkCall_WithData(b *testing.B) {
	params := public.CallParameters{
		To:   &usdcAddress,
		Data: balanceOfVitalikData,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := public.Call(benchCtx, benchClient, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCall_WithAccount benchmarks a call with a specified sender address.
// This simulates calling from a specific account (msg.sender).
//
// Equivalent to TypeScript:
//
//	await client.call({
//	  account: anvilAccount0,
//	  to: USDC,
//	  data: '0x06fdde03'
//	})
func BenchmarkCall_WithAccount(b *testing.B) {
	params := public.CallParameters{
		Account: &anvilAccount0,
		To:      &usdcAddress,
		Data:    nameSelector,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := public.Call(benchCtx, benchClient, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCall_Decimals benchmarks reading the decimals of a token.
// This is another simple view function call.
func BenchmarkCall_Decimals(b *testing.B) {
	params := public.CallParameters{
		To:   &usdcAddress,
		Data: decimalsSelector,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := public.Call(benchCtx, benchClient, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCall_Symbol benchmarks reading the symbol of a token.
func BenchmarkCall_Symbol(b *testing.B) {
	params := public.CallParameters{
		To:   &usdcAddress,
		Data: symbolSelector,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := public.Call(benchCtx, benchClient, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkCall_BalanceOfMultiple benchmarks multiple balanceOf calls in sequence.
// This tests repeated calls with different parameters.
func BenchmarkCall_BalanceOfMultiple(b *testing.B) {
	addresses := []struct {
		name string
		data []byte
	}{
		{"vitalik", encodeBalanceOf(vitalikAddress)},
		{"anvil0", encodeBalanceOf(anvilAccount0)},
		{"usdc", encodeBalanceOf(usdcAddress)},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		addr := addresses[i%len(addresses)]
		params := public.CallParameters{
			To:   &usdcAddress,
			Data: addr.data,
		}
		_, err := public.Call(benchCtx, benchClient, params)
		if err != nil {
			b.Fatal(err)
		}
	}
}
