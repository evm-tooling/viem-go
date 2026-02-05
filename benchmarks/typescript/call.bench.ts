/**
 * Call Action Benchmarks (viem TypeScript)
 *
 * These benchmarks mirror the Go benchmarks in ../go/call_bench_test.go
 * for fair cross-language comparison.
 *
 * Both benchmark suites should be run against the same Anvil instance
 * using the scripts/anvil.sh management script.
 */

import { bench, describe } from 'vitest'
import {
  createPublicClient,
  http,
  encodeFunctionData,
  type Address,
  type Hex,
} from 'viem'
import { mainnet } from 'viem/chains'

// Get RPC URL from environment or use default
const rpcUrl = process.env.ANVIL_RPC_URL || 'http://127.0.0.1:8545'

// Create the public client
const client = createPublicClient({
  chain: mainnet,
  transport: http(rpcUrl),
})

// Test addresses (same as Go benchmarks)
const USDC_ADDRESS: Address = '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48'
const VITALIK_ADDRESS: Address = '0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045'
const ANVIL_ACCOUNT_0: Address = '0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266'

// ERC20 function selectors (same as Go benchmarks)
const NAME_SELECTOR: Hex = '0x06fdde03'
const DECIMALS_SELECTOR: Hex = '0x313ce567'
const SYMBOL_SELECTOR: Hex = '0x95d89b41'

// Simple ERC20 ABI for encoding
const erc20Abi = [
  {
    name: 'balanceOf',
    type: 'function',
    stateMutability: 'view',
    inputs: [{ name: 'account', type: 'address' }],
    outputs: [{ type: 'uint256' }],
  },
] as const

// Pre-encoded calldata for balanceOf(vitalikAddress)
const balanceOfVitalikData = encodeFunctionData({
  abi: erc20Abi,
  functionName: 'balanceOf',
  args: [VITALIK_ADDRESS],
})

// Log connection info
console.log(`\n[viem-ts] RPC URL: ${rpcUrl}`)

// Verify connection before running benchmarks - do a warmup call
// This ensures Anvil is running and helps stabilize benchmark results
const warmup = async () => {
  try {
    const blockNumber = await client.getBlockNumber()
    console.log(`[viem-ts] Connected to Anvil, block number: ${blockNumber}`)
    // Do a warmup call
    await client.call({
      to: USDC_ADDRESS,
      data: NAME_SELECTOR,
    })
    console.log('[viem-ts] Warmup complete, starting benchmarks...')
  } catch (error) {
    console.error(`[viem-ts] FATAL: Failed to connect to Anvil at ${rpcUrl}`)
    console.error(error)
    process.exit(1)
  }
}

// Run warmup synchronously before benchmarks
await warmup()

describe('Call', () => {
  /**
   * BenchmarkCall_Basic - Simple contract call reading the token name.
   *
   * Equivalent to Go:
   *   public.Call(ctx, client, public.CallParameters{
   *     To:   &usdcAddress,
   *     Data: nameSelector,
   *   })
   */
  bench('viem-ts: call (basic)', async () => {
    await client.call({
      to: USDC_ADDRESS,
      data: NAME_SELECTOR,
    })
  })

  /**
   * BenchmarkCall_WithData - Call with encoded function parameters.
   *
   * Equivalent to Go:
   *   public.Call(ctx, client, public.CallParameters{
   *     To:   &usdcAddress,
   *     Data: balanceOfVitalikData,
   *   })
   */
  bench('viem-ts: call (with data)', async () => {
    await client.call({
      to: USDC_ADDRESS,
      data: balanceOfVitalikData,
    })
  })

  /**
   * BenchmarkCall_WithAccount - Call with a specified sender address.
   *
   * Equivalent to Go:
   *   public.Call(ctx, client, public.CallParameters{
   *     Account: &anvilAccount0,
   *     To:      &usdcAddress,
   *     Data:    nameSelector,
   *   })
   */
  bench('viem-ts: call (with account)', async () => {
    await client.call({
      account: ANVIL_ACCOUNT_0,
      to: USDC_ADDRESS,
      data: NAME_SELECTOR,
    })
  })

  /**
   * BenchmarkCall_Decimals - Reading the decimals of a token.
   */
  bench('viem-ts: call (decimals)', async () => {
    await client.call({
      to: USDC_ADDRESS,
      data: DECIMALS_SELECTOR,
    })
  })

  /**
   * BenchmarkCall_Symbol - Reading the symbol of a token.
   */
  bench('viem-ts: call (symbol)', async () => {
    await client.call({
      to: USDC_ADDRESS,
      data: SYMBOL_SELECTOR,
    })
  })

  /**
   * BenchmarkCall_BalanceOfMultiple - Multiple balanceOf calls with different addresses.
   */
  const addresses = [VITALIK_ADDRESS, ANVIL_ACCOUNT_0, USDC_ADDRESS]
  let callIndex = 0

  bench('viem-ts: call (balanceOf multiple)', async () => {
    const addr = addresses[callIndex % addresses.length]
    callIndex++

    await client.call({
      to: USDC_ADDRESS,
      data: encodeFunctionData({
        abi: erc20Abi,
        functionName: 'balanceOf',
        args: [addr],
      }),
    })
  })
})
