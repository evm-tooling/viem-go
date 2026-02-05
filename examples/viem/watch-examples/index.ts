/**
 * Watch Action Examples (viem TypeScript)
 *
 * Demonstrates the viem watch actions for real-time blockchain monitoring:
 * - watchBlockNumber: Watch for new block numbers
 * - watchBlocks: Watch for new blocks with full data
 * - watchPendingTransactions: Watch for pending transactions
 * - watchEvent: Watch for generic event logs
 * - watchContractEvent: Watch for contract events with ABI decoding
 *
 * This is the TypeScript equivalent of the Go watch-examples/main.go
 */

import {
  type Address,
  createPublicClient,
  formatEther,
  formatUnits,
  http,
  parseAbiItem,
  webSocket,
} from 'viem'
import { polygon } from 'viem/chains'

// USDC contract address on Polygon
const USDC_ADDRESS: Address = '0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359'

// ERC20 Transfer event ABI
const transferEventAbi = parseAbiItem(
  'event Transfer(address indexed from, address indexed to, uint256 value)',
)

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

function truncateHash(hash: string): string {
  return `${hash.slice(0, 18)}...${hash.slice(-4)}`
}

// Utility to run with timeout
function runWithTimeout<T>(
  fn: (signal: AbortSignal) => Promise<T>,
  timeoutMs: number,
): Promise<T | void> {
  return new Promise((resolve) => {
    const controller = new AbortController()
    const timeout = setTimeout(() => {
      controller.abort()
      resolve()
    }, timeoutMs)

    fn(controller.signal)
      .then(resolve)
      .catch(() => resolve())
      .finally(() => clearTimeout(timeout))
  })
}

/**
 * Example 1: Watch Block Number
 * Monitors incoming block numbers with polling
 */
async function watchBlockNumberExample(signal: AbortSignal) {
  printSection('Watch Block Number')
  console.log('Watching for new block numbers...')
  console.log('(Will stop after 5 blocks)')

  const client = createPublicClient({
    chain: polygon,
    transport: http('https://rough-purple-market.matic.quiknode.pro/c1a568726a34041d3c5d58603f5981951e6a8503'),
  })

  let count = 0
  let prevBlockNumber: bigint | undefined

  const unwatch = client.watchBlockNumber({
    emitOnBegin: true,
    emitMissed: true,
    // pollingInterval: 2_000, // 2 seconds
    onBlockNumber: (blockNumber) => {
      const prev = prevBlockNumber !== undefined ? prevBlockNumber.toString() : 'nil'
      console.log(`Block: ${blockNumber} (prev: ${prev})`)
      prevBlockNumber = blockNumber

      count++
      if (count >= 5) {
        console.log('Received 5 block numbers, stopping...')
        unwatch()
      }
    },
    onError: (error) => {
      console.log(`Error: ${error.message}`)
    },
  })

  // Handle abort signal
  signal.addEventListener('abort', () => unwatch())

  // Keep running until stopped
  await new Promise<void>((resolve) => {
    const checkInterval = setInterval(() => {
      if (signal.aborted || count >= 5) {
        clearInterval(checkInterval)
        resolve()
      }
    }, 500)
  })
}

/**
 * Example 2: Watch Blocks
 * Monitors incoming blocks with full block data
 */
async function watchBlocksExample(signal: AbortSignal) {
  printSection('Watch Blocks')
  console.log('Watching for new blocks with full data...')
  console.log('(Will stop after 3 blocks)')

  const client = createPublicClient({
    chain: polygon,
    transport: http('https://polygon-rpc.com'),
  })

  let count = 0

  const unwatch = client.watchBlocks({
    emitOnBegin: true,
    emitMissed: true,
    includeTransactions: false,
    pollingInterval: 2_000,
    onBlock: (block) => {
      console.log(`Block ${block.number}:`)
      console.log(`  Hash:         ${truncateHash(block.hash)}`)
      console.log(`  Timestamp:    ${block.timestamp}`)
      console.log(`  Gas Used:     ${block.gasUsed}`)
      console.log(`  Transactions: ${block.transactions.length}`)

      count++
      if (count >= 3) {
        console.log('Received 3 blocks, stopping...')
        unwatch()
      }
    },
    onError: (error) => {
      console.log(`Error: ${error.message}`)
    },
  })

  signal.addEventListener('abort', () => unwatch())

  await new Promise<void>((resolve) => {
    const checkInterval = setInterval(() => {
      if (signal.aborted || count >= 3) {
        clearInterval(checkInterval)
        resolve()
      }
    }, 500)
  })
}

/**
 * Example 3: Watch Pending Transactions
 * Monitors pending transaction hashes
 * Note: Many RPC providers don't support this
 */
async function watchPendingTransactionsExample(signal: AbortSignal) {
  printSection('Watch Pending Transactions')
  console.log('Watching for pending transactions...')
  console.log('Note: Some RPC providers do not support pending transaction filters')
  console.log('(Will stop after 20 transactions or error)')

  const client = createPublicClient({
    chain: polygon,
    transport: http('https://polygon-rpc.com'),
  })

  let totalTx = 0

  const unwatch = client.watchPendingTransactions({
    batch: true,
    pollingInterval: 2_000,
    onTransactions: (hashes) => {
      console.log(`Received ${hashes.length} pending transaction(s):`)
      for (let i = 0; i < Math.min(5, hashes.length); i++) {
        console.log(`  - ${truncateHash(hashes[i])}`)
      }
      if (hashes.length > 5) {
        console.log(`  ... and ${hashes.length - 5} more`)
      }

      totalTx += hashes.length
      if (totalTx >= 20) {
        console.log(`Received ${totalTx} total pending transactions, stopping...`)
        unwatch()
      }
    },
    onError: (error) => {
      console.log(`Error: ${error.message}`)
      // Many providers don't support pending tx filters
      unwatch()
    },
  })

  signal.addEventListener('abort', () => unwatch())

  await new Promise<void>((resolve) => {
    const checkInterval = setInterval(() => {
      if (signal.aborted || totalTx >= 20) {
        clearInterval(checkInterval)
        resolve()
      }
    }, 500)
  })
}

/**
 * Example 4: Watch Event
 * Monitors generic event logs (Transfer events from USDC)
 */
async function watchEventExample(signal: AbortSignal) {
  printSection('Watch Event')
  console.log('Watching for Transfer events from USDC contract...')
  console.log('(Will stop after 10 events)')

  const client = createPublicClient({
    chain: polygon,
    transport: http('https://polygon-rpc.com'),
  })

  let count = 0

  const unwatch = client.watchEvent({
    address: USDC_ADDRESS,
    event: transferEventAbi,
    batch: true,
    // pollingInterval: 2_000,
    onLogs: (logs) => {
      console.log(`Received ${logs.length} Transfer event(s):`)
      for (let i = 0; i < Math.min(3, logs.length); i++) {
        const log = logs[i]
        console.log(`  Block ${log.blockNumber}, Tx: ${truncateHash(log.transactionHash ?? '')}`)
        if (log.args) {
          console.log(`    From: ${truncateAddress(log.args.from as string)}`)
          console.log(`    To:   ${truncateAddress(log.args.to as string)}`)
        }
      }
      if (logs.length > 3) {
        console.log(`  ... and ${logs.length - 3} more`)
      }

      count += logs.length
      if (count >= 100) {
        console.log(`Received ${count} total events, stopping...`)
        unwatch()
      }
    },
    onError: (error) => {
      console.log(`Error: ${error.message}`)
    },
  })

  signal.addEventListener('abort', () => unwatch())

  await new Promise<void>((resolve) => {
    const checkInterval = setInterval(() => {
      if (signal.aborted || count >= 10) {
        clearInterval(checkInterval)
        resolve()
      }
    }, 500)
  })
}

/**
 * Example 5: Watch Contract Event
 * Monitors contract events with full ABI decoding
 */
async function watchContractEventExample(signal: AbortSignal) {
  printSection('Watch Contract Event')
  console.log('Watching for Transfer events with ABI decoding...')
  console.log('(Will stop after 10 events)')

  const client = createPublicClient({
    chain: polygon,
    transport: http('https://polygon-rpc.com'),
  })

  // Full ERC20 ABI for contract event watching
  const erc20Abi = [
    {
      anonymous: false,
      inputs: [
        { indexed: true, name: 'from', type: 'address' },
        { indexed: true, name: 'to', type: 'address' },
        { indexed: false, name: 'value', type: 'uint256' },
      ],
      name: 'Transfer',
      type: 'event',
    },
  ] as const

  let count = 0

  const unwatch = client.watchContractEvent({
    address: USDC_ADDRESS,
    abi: erc20Abi,
    eventName: 'Transfer',
    batch: true,
    pollingInterval: 2_000,
    onLogs: (logs) => {
      console.log(`Received ${logs.length} decoded Transfer event(s):`)
      for (let i = 0; i < Math.min(3, logs.length); i++) {
        const log = logs[i]
        console.log(`  Block ${log.blockNumber}:`)
        console.log(`    Event: ${log.eventName}`)
        if (log.args) {
          console.log(`    From:  ${(log.args.from)}`)
          console.log(`    To:    ${(log.args.to)}`)
          // USDC has 6 decimals
          console.log(`    Value: ${formatUnits(log.args.value ?? 0n, 6)} USDC`)
        }
      }
      if (logs.length > 3) {
        console.log(`  ... and ${logs.length - 3} more`)
      }

      count += logs.length
      if (count >= 10) {
        console.log(`Received ${count} total events, stopping...`)
        unwatch()
      }
    },
    onError: (error) => {
      console.log(`Error: ${error.message}`)
    },
  })

  signal.addEventListener('abort', () => unwatch())

  await new Promise<void>((resolve) => {
    const checkInterval = setInterval(() => {
      if (signal.aborted || count >= 10) {
        clearInterval(checkInterval)
        resolve()
      }
    }, 500)
  })
}

/**
 * Main function - runs selected example or all examples
 */
async function main() {
  const args = process.argv.slice(2)
  const example = args[0]

  if (!example) {
    console.log('Usage: npx ts-node index.ts <example>')
    console.log('Examples:')
    console.log('  watch-block-number    - Watch for new block numbers')
    console.log('  watch-blocks          - Watch for new blocks with full data')
    console.log('  watch-pending-tx      - Watch for pending transactions')
    console.log('  watch-event           - Watch for generic events')
    console.log('  watch-contract-event  - Watch for ERC20 Transfer events')
    console.log('  all                   - Run all examples')
    process.exit(1)
  }

  printHeader('Watch Action Examples (viem TypeScript)')

  // Handle graceful shutdown
  const controller = new AbortController()
  process.on('SIGINT', () => {
    console.log('\nShutting down...')
    controller.abort()
    process.exit(0)
  })
  process.on('SIGTERM', () => {
    controller.abort()
    process.exit(0)
  })

  switch (example) {
    case 'watch-block-number':
      await watchBlockNumberExample(controller.signal)
      break
    case 'watch-blocks':
      await watchBlocksExample(controller.signal)
      break
    case 'watch-pending-tx':
      await watchPendingTransactionsExample(controller.signal)
      break
    case 'watch-event':
      await watchEventExample(controller.signal)
      break
    case 'watch-contract-event':
      await watchContractEventExample(controller.signal)
      break
    case 'all':
      console.log('Running all examples in sequence...')

      console.log('\n=== Watch Block Number ===')
      await runWithTimeout(watchBlockNumberExample, 15_000)

      console.log('\n=== Watch Blocks ===')
      await runWithTimeout(watchBlocksExample, 15_000)

      console.log('\n=== Watch Pending Transactions ===')
      await runWithTimeout(watchPendingTransactionsExample, 10_000)

      console.log('\n=== Watch Event ===')
      await runWithTimeout(watchEventExample, 15_000)

      console.log('\n=== Watch Contract Event ===')
      await runWithTimeout(watchContractEventExample, 15_000)

      console.log('\nAll examples completed!')
      break
    default:
      console.log(`Unknown example: ${example}`)
      process.exit(1)
  }

  printHeader('Examples Complete')
  console.log('Demonstrated watch features:')
  console.log('  - watchBlockNumber: Real-time block number monitoring')
  console.log('  - watchBlocks: Full block data streaming')
  console.log('  - watchPendingTransactions: Pending tx pool monitoring')
  console.log('  - watchEvent: Generic event log filtering')
  console.log('  - watchContractEvent: ABI-decoded contract events')
  console.log()
  console.log('Key differences from Go implementation:')
  console.log('  - Callbacks instead of channels (JavaScript event model)')
  console.log('  - unwatch() function instead of context cancellation')
  console.log('  - Native async/await patterns')
  console.log()
}

main().catch(console.error)
