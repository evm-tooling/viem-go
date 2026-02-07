/**
 * Signature Benchmarks (viem TypeScript)
 *
 * These benchmarks mirror the Go benchmarks in ../go/signature_bench_test.go
 * for fair cross-language comparison.
 *
 * Pure CPU-bound operations -- no network involved.
 */

import { bench, describe } from 'vitest'
import {
  hashMessage,
  recoverMessageAddress,
  verifyMessage,
  parseSignature,
} from 'viem'

// Signature from Anvil account 0 signing "hello world"
const testSignature = '0x6e100a352ec6ad1b70802290e18aeed190704973570f3b8ed42cb9808e2ea6bf4a90a229a244495b41890987806fcbd2d5d23fc0dbe5f5256c2613c039d76db81c' as const
const testAddress = '0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266' as const

const longMessage = 'The quick brown fox jumps over the lazy dog. ' +
  'This is a much longer message that simulates real-world signing scenarios ' +
  'where users might sign terms of service, governance proposals, or other text content.'

const benchOptions = {
  time: 2000,
  warmupTime: 0,
  warmupIterations: 0,
}

describe('Signature', () => {
  bench('viem-ts: signature (hashMessage)', () => {
    hashMessage('hello world')
  }, benchOptions)

  bench('viem-ts: signature (hashMessage long)', () => {
    hashMessage(longMessage)
  }, benchOptions)

  bench('viem-ts: signature (recoverAddress)', async () => {
    await recoverMessageAddress({
      message: 'hello world',
      signature: testSignature,
    })
  }, benchOptions)

  bench('viem-ts: signature (verifyMessage)', async () => {
    await verifyMessage({
      address: testAddress,
      message: 'hello world',
      signature: testSignature,
    })
  }, benchOptions)

  bench('viem-ts: signature (parseSignature)', () => {
    parseSignature(testSignature)
  }, benchOptions)
})
