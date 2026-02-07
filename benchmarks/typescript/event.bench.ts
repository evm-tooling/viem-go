/**
 * Event Log Decoding Benchmarks (viem TypeScript)
 *
 * These benchmarks mirror the Go benchmarks in ../go/event_bench_test.go
 * for fair cross-language comparison.
 *
 * Pure CPU-bound operations -- no network involved.
 */

import { bench, describe } from 'vitest'
import { decodeEventLog, Hex, parseAbi } from 'viem'

const eventAbi = parseAbi([
  'event Transfer(address indexed from, address indexed to, uint256 value)',
  'event Approval(address indexed owner, address indexed spender, uint256 value)',
])

// Transfer event log data
const transferTopics = [
  '0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef', // Transfer signature
  '0x000000000000000000000000f39Fd6e51aad88F6F4ce6aB8827279cffFb92266', // from
  '0x000000000000000000000000d8dA6BF26964aF9D7eEd9e03E53415D37aA96045', // to
] as [Hex, Hex, Hex]

// ABI-encoded uint256 value: 1000000
const transferData = '0x00000000000000000000000000000000000000000000000000000000000f4240' as Hex

const benchOptions = {
  time: 2000,
  warmupTime: 0,
  warmupIterations: 0,
}

describe('Event', () => {
  bench('viem-ts: event (decode transfer)', () => {
    decodeEventLog({
      abi: eventAbi,
      topics: transferTopics,
      data: transferData,
    })
  }, benchOptions)

  bench('viem-ts: event (decode batch 10)', () => {
    for (let i = 0; i < 10; i++) {
      decodeEventLog({
        abi: eventAbi,
        topics: transferTopics,
        data: transferData,
      })
    }
  }, benchOptions)

  bench('viem-ts: event (decode batch 100)', () => {
    for (let i = 0; i < 100; i++) {
      decodeEventLog({
        abi: eventAbi,
        topics: transferTopics,
        data: transferData,
      })
    }
  }, benchOptions)
})
