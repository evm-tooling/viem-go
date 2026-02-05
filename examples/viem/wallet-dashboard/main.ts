/**
 * Wallet Dashboard - viem TypeScript Example
 *
 * A simple CLI dashboard demonstrating key viem features:
 * - Public client creation
 * - Fetching network info (block number, chain ID)
 * - Gas price estimation
 * - Address balance lookup
 * - Message signing and verification
 */

import {
  type Address,
  createPublicClient,
  formatEther,
  formatGwei,
  http,
} from 'viem'
import { privateKeyToAccount } from 'viem/accounts'
import { mainnet, polygon } from 'viem/chains'

// Example address to check balance (Vitalik's address)
const EXAMPLE_ADDRESS: Address = '0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045'

// Example private key for signing demo (DO NOT use real keys!)
const DEMO_PRIVATE_KEY =
  '0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80'

function printHeader(title: string) {
  console.log(`\n${'='.repeat(50)}`)
  console.log(`  ${title}`)
  console.log('='.repeat(50))
}

function printSection(title: string) {
  console.log(`\n--- ${title} ---`)
}

function formatDuration(ms: number): string {
  if (ms < 1) {
    return `${(ms * 1000).toFixed(0)}µs`
  }
  if (ms < 1000) {
    return `${ms.toFixed(2)}ms`
  }
  return `${(ms / 1000).toFixed(2)}s`
}

async function main() {
  const totalStart = performance.now()

  console.log('╔══════════════════════════════════════════════════╗')
  console.log('║           Wallet Dashboard (viem)                ║')
  console.log('║        Ethereum Network Information              ║')
  console.log('╚══════════════════════════════════════════════════╝')

  // 1. Create Public Client
  printSection('Creating Public Client')
  const publicClient = createPublicClient({
    chain: polygon,
    transport: http('https://polygon-rpc.com'),
  })
  console.log(`Connected to: ${mainnet.name} (Chain ID: ${mainnet.id})`)

  // 2. Fetch Network Information
  printSection('Network Information')

  const blockNumber = await publicClient.getBlockNumber()
  const chainId = await publicClient.getChainId()

  console.log(`Current Block Number: ${blockNumber.toLocaleString()}`)
  console.log(`Chain ID: ${chainId}`)

  // 3. Get Gas Prices
  printSection('Gas Prices')

  const gasPrice = await publicClient.getGasPrice()
  const feeHistory = await publicClient.getFeeHistory({
    blockCount: 4,
    rewardPercentiles: [25, 50, 75],
  })

  console.log(`Current Gas Price: ${formatGwei(gasPrice)} Gwei`)

  // Calculate average base fee from fee history
  const avgBaseFee =
    feeHistory.baseFeePerGas.reduce((a, b) => a + b, 0n) /
    BigInt(feeHistory.baseFeePerGas.length)
  console.log(`Avg Base Fee (last 4 blocks): ${formatGwei(avgBaseFee)} Gwei`)

  try {
    const maxPriorityFee = await publicClient.estimateMaxPriorityFeePerGas()
    console.log(`Max Priority Fee: ${formatGwei(maxPriorityFee)} Gwei`)
  } catch {
    console.log('Max Priority Fee: Not available on this network')
  }

  // 4. Check Address Balance
  printSection('Address Balance')

  const balance = await publicClient.getBalance({
    address: EXAMPLE_ADDRESS,
  })

  console.log(`Address: ${EXAMPLE_ADDRESS}`)
  console.log(`Balance: ${formatEther(balance)} ETH`)

  // Get transaction count (nonce)
  const txCount = await publicClient.getTransactionCount({
    address: EXAMPLE_ADDRESS,
  })
  console.log(`Transaction Count: ${txCount.toLocaleString()}`)

  // 5. Message Signing Demo
  printSection('Message Signing Demo')

  const account = privateKeyToAccount(DEMO_PRIVATE_KEY)
  console.log(`Demo Account Address: ${account.address}`)

  const message = 'Hello from Wallet Dashboard!'
  console.log(`Message: "${message}"`)

  const signature = await account.signMessage({ message })
  console.log(`Signature: ${signature.slice(0, 42)}...`)

  // Verify the signature
  const isValid = await publicClient.verifyMessage({
    address: account.address,
    message,
    signature,
  })
  console.log(`Signature Valid: ${isValid ? 'Yes' : 'No'}`)

  // 6. Get Latest Block Info
  printSection('Latest Block Info')

  const block = await publicClient.getBlock()
  console.log(`Block Hash: ${block.hash?.slice(0, 42)}...`)
  console.log(
    `Timestamp: ${new Date(Number(block.timestamp) * 1000).toISOString()}`,
  )
  console.log(`Transactions: ${block.transactions.length}`)
  console.log(`Gas Used: ${block.gasUsed.toLocaleString()}`)
  console.log(`Gas Limit: ${block.gasLimit.toLocaleString()}`)

  // Summary
  const totalElapsed = performance.now() - totalStart
  printHeader('Dashboard Summary')
  console.log(`  Network: ${mainnet.name}`)
  console.log(`  Block: #${blockNumber.toLocaleString()}`)
  console.log(`  Gas Price: ${formatGwei(gasPrice)} Gwei`)
  console.log(`  Demo Account: ${account.address.slice(0, 10)}...`)
  console.log(`  Total Runtime: ${formatDuration(totalElapsed)}`)
  console.log('')
}

main().catch(console.error)
