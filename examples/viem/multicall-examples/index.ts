/**
 * Multicall Examples (viem TypeScript)
 *
 * Comprehensive examples demonstrating the multicall action features:
 * - Basic multicall batching
 * - Multiple contract calls in parallel
 * - Error handling with allowFailure
 * - Deployless multicall
 * - Cross-contract queries
 */

import {
  type Address,
  createPublicClient,
  formatUnits,
  http,
  parseAbi,
} from 'viem'
import { polygon } from 'viem/chains'

// Example addresses on Polygon
const USDC_ADDRESS: Address = '0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359'
const WETH_ADDRESS: Address = '0x7ceB23fD6bC0adD59E62ac25578270cFf1b9f619'
const WMATIC_ADDRESS: Address = '0x0d500B1d8E8eF31E21C99d1Db9A6444d3ADf1270'
const VITALIK_ADDRESS: Address = '0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045'
const INVALID_ADDRESS: Address = '0x0000000000000000000000000000000000000001'

// ERC20 ABI
const erc20Abi = parseAbi([
  'function name() view returns (string)',
  'function symbol() view returns (string)',
  'function decimals() view returns (uint8)',
  'function totalSupply() view returns (uint256)',
  'function balanceOf(address account) view returns (uint256)',
])

// Helper functions
function printHeader(title: string) {
  console.log()
  console.log('='.repeat(60))
  console.log(`  ${title}`)
  console.log('='.repeat(60))
}

function printSection(title: string) {
  console.log(`\n--- ${title} ---`)
}

function truncateAddress(addr: string): string {
  return `${addr.slice(0, 10)}...${addr.slice(-4)}`
}

async function main() {
  printHeader('Multicall Action Examples (viem TypeScript)')

  // Create Public Client
  printSection('1. Creating Public Client')
  const client = createPublicClient({
    chain: polygon,
    transport: http('https://polygon-rpc.com'),
  })
  console.log('Connected to Polygon Mainnet')

  // Example 1: Basic Multicall - Read multiple token names
  printSection('2. Basic Multicall - Read Multiple Token Names')
  try {
    const results = await client.multicall({
      contracts: [
        {
          address: USDC_ADDRESS,
          abi: erc20Abi,
          functionName: 'name',
        },
        {
          address: WETH_ADDRESS,
          abi: erc20Abi,
          functionName: 'name',
        },
        {
          address: WMATIC_ADDRESS,
          abi: erc20Abi,
          functionName: 'name',
        },
      ],
    })
    console.log('Token names retrieved in single RPC call:')
    const tokens = ['USDC', 'WETH', 'WMATIC']
    results.forEach((result, i) => {
      if (result.status === 'success') {
        console.log(`  ${tokens[i]}: ${result.result}`)
      } else {
        console.log(`  ${tokens[i]}: failed - ${result.error}`)
      }
    })
  } catch (error) {
    console.log(`Error: ${error}`)
  }

  // Example 2: Multicall for Token Metadata
  printSection('3. Multicall - Complete Token Metadata')
  try {
    const results = await client.multicall({
      contracts: [
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'name' },
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'symbol' },
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'decimals' },
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'totalSupply' },
      ],
    })
    console.log('USDC Token Metadata:')
    const fields = ['Name', 'Symbol', 'Decimals', 'Total Supply']
    results.forEach((result, i) => {
      if (result.status === 'success') {
        let value: string | bigint | number = result.result as
          | string
          | bigint
          | number
        if (fields[i] === 'Total Supply' && typeof value === 'bigint') {
          value = `${formatUnits(value, 6)} USDC`
        }
        console.log(`  ${fields[i]}: ${value}`)
      } else {
        console.log(`  ${fields[i]}: failed - ${result.error}`)
      }
    })
  } catch (error) {
    console.log(`Error: ${error}`)
  }

  // Example 3: Multicall for Multiple Balances
  printSection('4. Multicall - Multiple Balance Queries')
  try {
    const results = await client.multicall({
      contracts: [
        {
          address: USDC_ADDRESS,
          abi: erc20Abi,
          functionName: 'balanceOf',
          args: [VITALIK_ADDRESS],
        },
        {
          address: WETH_ADDRESS,
          abi: erc20Abi,
          functionName: 'balanceOf',
          args: [VITALIK_ADDRESS],
        },
        {
          address: WMATIC_ADDRESS,
          abi: erc20Abi,
          functionName: 'balanceOf',
          args: [VITALIK_ADDRESS],
        },
      ],
    })
    console.log(`Vitalik's Polygon balances (${truncateAddress(VITALIK_ADDRESS)}):`)
    const tokens = [
      { name: 'USDC', decimals: 6 },
      { name: 'WETH', decimals: 18 },
      { name: 'WMATIC', decimals: 18 },
    ]
    results.forEach((result, i) => {
      if (result.status === 'success') {
        const balance = result.result as bigint
        console.log(`  ${tokens[i].name}: ${formatUnits(balance, tokens[i].decimals)}`)
      } else {
        console.log(`  ${tokens[i].name}: failed - ${result.error}`)
      }
    })
  } catch (error) {
    console.log(`Error: ${error}`)
  }

  // Example 4: Multicall with allowFailure
  printSection('5. Multicall with allowFailure=true (default)')
  try {
    const results = await client.multicall({
      contracts: [
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'name' },
        { address: INVALID_ADDRESS, abi: erc20Abi, functionName: 'name' }, // Will fail
        { address: WETH_ADDRESS, abi: erc20Abi, functionName: 'name' },
      ],
    })
    console.log('Results (with one failing call):')
    results.forEach((result, i) => {
      if (result.status === 'success') {
        console.log(`  Call ${i + 1}: success - ${result.result}`)
      } else {
        console.log(`  Call ${i + 1}: failure - ${result.error?.message?.slice(0, 50)}...`)
      }
    })
  } catch (error) {
    console.log(`Error: ${error}`)
  }

  // Example 5: Multicall with allowFailure=false
  printSection('6. Multicall with allowFailure=false')
  try {
    await client.multicall({
      allowFailure: false,
      contracts: [
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'name' },
        { address: INVALID_ADDRESS, abi: erc20Abi, functionName: 'name' }, // Will cause error
      ],
    })
    console.log('Unexpected: no error')
  } catch (error) {
    const errMsg = (error as Error).message?.slice(0, 80) || 'Unknown error'
    console.log(`Expected error with allowFailure=false: ${errMsg}...`)
  }

  // Example 6: Large Multicall with Batching
  printSection('7. Large Multicall with Automatic Batching')
  try {
    const addresses = [USDC_ADDRESS, WETH_ADDRESS, WMATIC_ADDRESS]
    const contracts = Array.from({ length: 30 }, (_, i) => ({
      address: addresses[i % 3],
      abi: erc20Abi,
      functionName: 'balanceOf' as const,
      args: [VITALIK_ADDRESS] as const,
    }))

    const results = await client.multicall({
      contracts,
      batchSize: 512, // Small batch size to force multiple chunks
    })

    const successCount = results.filter((r) => r.status === 'success').length
    console.log(`Executed ${successCount} calls successfully`)
    console.log('BatchSize: 512 bytes (forces multiple batches)')
  } catch (error) {
    console.log(`Error: ${error}`)
  }

  // Example 7: Multicall at Specific Block
  printSection('8. Multicall at Specific Block Number')
  try {
    const results = await client.multicall({
      blockNumber: 52000000n, // Historical Polygon block
      contracts: [
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'totalSupply' },
      ],
    })
    if (results[0].status === 'success') {
      const supply = results[0].result as bigint
      console.log(`USDC Total Supply at block 52000000: ${formatUnits(supply, 6)}`)
    }
  } catch (error) {
    console.log(`Error: ${error}`)
  }

  // Example 8: Deployless Multicall
  printSection('9. Deployless Multicall')
  try {
    const results = await client.multicall({
      deployless: true,
      contracts: [
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'symbol' },
        { address: WETH_ADDRESS, abi: erc20Abi, functionName: 'symbol' },
      ],
    })
    console.log('Symbols via deployless multicall:')
    results.forEach((result, i) => {
      if (result.status === 'success') {
        console.log(`  Token ${i + 1}: ${result.result}`)
      }
    })
  } catch (error) {
    console.log(`Error: ${error}`)
  }

  // Example 9: Cross-contract Queries
  printSection('10. Cross-Contract Queries')
  try {
    const results = await client.multicall({
      contracts: [
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'name' },
        { address: USDC_ADDRESS, abi: erc20Abi, functionName: 'decimals' },
        { address: WETH_ADDRESS, abi: erc20Abi, functionName: 'name' },
        { address: WETH_ADDRESS, abi: erc20Abi, functionName: 'decimals' },
        { address: WMATIC_ADDRESS, abi: erc20Abi, functionName: 'name' },
        { address: WMATIC_ADDRESS, abi: erc20Abi, functionName: 'decimals' },
      ],
    })
    console.log('Multi-contract metadata in single RPC call:')
    for (let i = 0; i < 6; i += 2) {
      const name =
        results[i].status === 'success' ? results[i].result : 'unknown'
      const decimals =
        results[i + 1].status === 'success' ? results[i + 1].result : '?'
      console.log(`  ${name}: ${decimals} decimals`)
    }
  } catch (error) {
    console.log(`Error: ${error}`)
  }

  // Summary
  printHeader('Examples Complete')
  console.log('Demonstrated Multicall features:')
  console.log('  - Basic multicall batching')
  console.log('  - Token metadata retrieval')
  console.log('  - Multiple balance queries')
  console.log('  - Error handling with allowFailure')
  console.log('  - Large multicalls with batching')
  console.log('  - Historical block queries')
  console.log('  - Deployless multicall')
  console.log('  - Cross-contract queries')
  console.log()
}

main().catch(console.error)
