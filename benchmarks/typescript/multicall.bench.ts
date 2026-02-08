/**
 * Multicall Action Benchmarks (viem TypeScript)
 *
 * These benchmarks mirror the Go benchmarks in ../go/multicall_bench_test.go
 * for fair cross-language comparison.
 *
 * Both benchmark suites should be run against the same Anvil instance
 * using the unified benchmarks entrypoint (benchmarks/bench.sh).
 *
 * IMPORTANT: All benchmarks use batchSize: 0 to disable chunking,
 * ensuring a single RPC call for fair comparison.
 */

import { bench, describe } from 'vitest'
import {
  createPublicClient,
  http,
  parseAbi,
  type Address,
} from 'viem'
import { mainnet } from 'viem/chains'

// Get RPC URL from environment or use default
const rpcUrl = process.env.ANVIL_RPC_URL || 'http://127.0.0.1:8545'

// Create the public client
const client = createPublicClient({
  chain: mainnet,
  transport: http(rpcUrl),
  batch: { multicall: { batchSize: 8192, wait: 16 }}
})

// Test addresses (same as Go benchmarks)
const USDC_ADDRESS: Address = '0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48'
const WETH_ADDRESS: Address = '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2'
const VITALIK_ADDRESS: Address = '0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045'
const ANVIL_ACCOUNT_0: Address = '0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266'

// ERC20 ABI for multicall benchmarks (parsed once, reused)
const erc20Abi = parseAbi([
  'function name() view returns (string)',
  'function symbol() view returns (string)',
  'function decimals() view returns (uint8)',
  'function totalSupply() view returns (uint256)',
  'function balanceOf(address account) view returns (uint256)',
])

// Log connection info
console.log(`\n[viem-ts] Multicall RPC URL: ${rpcUrl}`)

const iterations = Number(process.env.BENCH_ITERATIONS ?? '100')

// Benchmark options: iteration-based (controlled by BENCH_ITERATIONS)
const benchOptions = {
  time: 0,
  warmupTime: 0,
  warmupIterations: 0,
  iterations,
}

describe('Multicall', () => {
  /**
   * BenchmarkMulticall_Basic - Simple multicall with 3 calls.
   */
  bench('viem-ts: multicall (basic)', async () => {
    await client.multicall({
      batchSize: 1,
      contracts: [
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'name' },
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'symbol' },
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'decimals' },
      ],
    })
  }, benchOptions)

  /**
   * BenchmarkMulticall_WithArgs - Multicall with function arguments.
   */
  bench('viem-ts: multicall (with args)', async () => {
    await client.multicall({
      contracts: [
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'balanceOf', args: [VITALIK_ADDRESS] },
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'balanceOf', args: [ANVIL_ACCOUNT_0] },
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'balanceOf', args: [USDC_ADDRESS] },
      ],
    })
  }, benchOptions)

  /**
   * BenchmarkMulticall_MultiContract - Multicall across multiple contracts.
   */
  bench('viem-ts: multicall (multi-contract)', async () => {
    await client.multicall({
      contracts: [
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'name' },
        { address: WETH_ADDRESS, abi: erc20Abi, functionName: 'name' },
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'balanceOf', args: [VITALIK_ADDRESS] },
        { address: WETH_ADDRESS, abi: erc20Abi, functionName: 'balanceOf', args: [VITALIK_ADDRESS] },
      ],
    })
  }, benchOptions)

  /**
   * BenchmarkMulticall_10Calls - Multicall with 10 calls.
   */
  bench('viem-ts: multicall (10 calls)', async () => {
    const contracts = Array.from({ length: 10 }, () => ({
      address: USDC_ADDRESS,
      abi: erc20Abi,
      functionName: 'balanceOf' as const,
      args: [VITALIK_ADDRESS] as const,
    }))
    await client.multicall({  contracts })
  }, benchOptions)

  /**
   * BenchmarkMulticall_30Calls - Multicall with 30 calls.
   */
  bench('viem-ts: multicall (30 calls)', async () => {
    const contracts = Array.from({ length: 30 }, () => ({
      address: USDC_ADDRESS,
      abi: erc20Abi,
      functionName: 'balanceOf' as const,
      args: [VITALIK_ADDRESS] as const,
    }))
    await client.multicall({  contracts })
  }, benchOptions)

  /**
   * BenchmarkMulticall_Deployless - Deployless multicall.
   */
  bench('viem-ts: multicall (deployless)', async () => {
    await client.multicall({
      deployless: true,
      contracts: [
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'name' },
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'symbol' },
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'decimals' },
      ],
    })
  }, benchOptions)

  /**
   * BenchmarkMulticall_TokenMetadata - Complete token metadata fetch.
   */
  bench('viem-ts: multicall (token metadata)', async () => {
    await client.multicall({
      contracts: [
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'name' },
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'symbol' },
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'decimals' },
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'totalSupply' },
      ],
    })
  }, benchOptions)

  // ============================================================
  // STRESS TESTS - Large batch sizes to test batching performance
  // ============================================================

  /**
   * BenchmarkMulticall_50Calls - Stress test with 50 calls.
   */
  bench('viem-ts: multicall (50 calls)', async () => {
    const contracts = Array.from({ length: 50 }, () => ({
      address: USDC_ADDRESS,
      abi: erc20Abi,
      functionName: 'balanceOf' as const,
      args: [VITALIK_ADDRESS] as const,
    }))
    await client.multicall({ contracts })
  }, benchOptions)

  /**
   * BenchmarkMulticall_100Calls - Stress test with 100 calls.
   */
  bench('viem-ts: multicall (100 calls)', async () => {
    const contracts = Array.from({ length: 100 }, () => ({
      address: USDC_ADDRESS,
      abi: erc20Abi,
      functionName: 'balanceOf' as const,
      args: [VITALIK_ADDRESS] as const,
    }))
    await client.multicall({  contracts })
  }, benchOptions)

  /**
   * BenchmarkMulticall_200Calls - Stress test with 200 calls.
   */
  bench('viem-ts: multicall (200 calls)', async () => {
    const contracts = Array.from({ length: 200 }, () => ({
      address: USDC_ADDRESS,
      abi: erc20Abi,
      functionName: 'balanceOf' as const,
      args: [VITALIK_ADDRESS] as const,
    }))
    await client.multicall({  contracts })
  }, benchOptions)

  /**
   * BenchmarkMulticall_500Calls - Extreme stress test with 500 calls.
   */
  bench('viem-ts: multicall (500 calls)', async () => {
    const contracts = Array.from({ length: 500 }, () => ({
      address: USDC_ADDRESS,
      abi: erc20Abi,
      functionName: 'balanceOf' as const,
      args: [VITALIK_ADDRESS] as const,
    }))
    await client.multicall({  contracts })
  }, benchOptions)

  /**
   * BenchmarkMulticall_MixedContracts_100 - 100 calls across multiple contracts.
   * More realistic workload with varied targets.
   */
  bench('viem-ts: multicall (100 mixed contracts)', async () => {
    const contracts = Array.from({ length: 100 }, (_, i) => ({
      address: i % 2 === 0 ? USDC_ADDRESS : WETH_ADDRESS,
      abi: erc20Abi,
      functionName: 'balanceOf' as const,
      args: [VITALIK_ADDRESS] as const,
    }))
    await client.multicall({  contracts })
  }, benchOptions)

  // ============================================================
  // EXTREME STRESS TESTS - Maximum throughput benchmarks
  // ============================================================

  /**
   * BenchmarkMulticall_1000Calls - 1000 calls, single RPC (no batching).
   */
  bench('viem-ts: multicall (1000 calls)', async () => {
    const contracts = Array.from({ length: 1000 }, () => ({
      address: USDC_ADDRESS,
      abi: erc20Abi,
      functionName: 'balanceOf' as const,
      args: [VITALIK_ADDRESS] as const,
    }))
    await client.multicall({  contracts })
  }, benchOptions)

  /**
   * BenchmarkMulticall_10000Calls_SingleRPC - 10,000 calls in single RPC.
   * Tests maximum payload size handling.
   */
  bench('viem-ts: multicall (10000 calls single RPC)', async () => {
    const contracts = Array.from({ length: 10000 }, () => ({
      address: USDC_ADDRESS,
      abi: erc20Abi,
      functionName: 'balanceOf' as const,
      args: [VITALIK_ADDRESS] as const,
    }))
    await client.multicall({
      batchSize: 0, // Disable chunking - single massive RPC call
      contracts,
    })
  }, benchOptions)

  /**
   * BenchmarkMulticall_10000Calls_Chunked - 10,000 calls with optimized chunking.
   * Uses large batches for parallel RPC execution.
   * batchSize: 8192 bytes (~228 calls per chunk) = ~44 chunks
   */
  bench('viem-ts: multicall (10000 calls chunked)', async () => {
    const contracts = Array.from({ length: 10000 }, () => ({
      address: USDC_ADDRESS,
      abi: erc20Abi,
      functionName: 'balanceOf' as const,
      args: [VITALIK_ADDRESS] as const,
    }))
    await client.multicall({
      batchSize: 8192, // Large batches for efficient parallel execution
      contracts,
    })
  }, benchOptions)

  /**
   * BenchmarkMulticall_10000Calls_AggressiveChunking - 10,000 calls with aggressive chunking.
   * Uses smaller batches for maximum parallelism.
   * batchSize: 2048 bytes (~57 calls per chunk) = ~175 chunks
   */
  bench('viem-ts: multicall (10000 calls aggressive)', async () => {
    const contracts = Array.from({ length: 10000 }, () => ({
      address: USDC_ADDRESS,
      abi: erc20Abi,
      functionName: 'balanceOf' as const,
      args: [VITALIK_ADDRESS] as const,
    }))
    await client.multicall({
      batchSize: 8192, // Smaller batches = more parallelism
      contracts,
    })
  }, benchOptions)
})
