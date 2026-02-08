/**
 * Hashing Benchmarks (viem TypeScript)
 *
 * These benchmarks mirror the Go benchmarks in ../go/hash_bench_test.go
 * for fair cross-language comparison.
 *
 * Pure CPU-bound operations -- no network involved.
 */

import { bench, describe } from 'vitest'
import {
  keccak256,
  sha256,
  toFunctionSelector,
  toEventSelector,
  toBytes,
  toHex,
} from 'viem'

// Pre-build data
const shortData = toHex(new TextEncoder().encode('hello world'))

// 1KB of data
const longArray = new Uint8Array(1024)
for (let i = 0; i < 1024; i++) longArray[i] = i % 256
const longData = toHex(longArray)

const iterations = Number(process.env.BENCH_ITERATIONS ?? '100')

const benchOptions = {
  time: 0,
  warmupTime: 0,
  warmupIterations: 0,
  iterations,
}

describe('Hash', () => {
  bench('viem-ts: hash (keccak256 short)', () => {
    keccak256(shortData)
  }, benchOptions)

  bench('viem-ts: hash (keccak256 long)', () => {
    keccak256(longData)
  }, benchOptions)

  bench('viem-ts: hash (keccak256 hex)', () => {
    keccak256('0x68656c6c6f20776f726c64')
  }, benchOptions)

  bench('viem-ts: hash (sha256 short)', () => {
    sha256(shortData)
  }, benchOptions)

  bench('viem-ts: hash (sha256 long)', () => {
    sha256(longData)
  }, benchOptions)

  bench('viem-ts: hash (function selector)', () => {
    toFunctionSelector('function transfer(address to, uint256 amount)')
  }, benchOptions)

  bench('viem-ts: hash (event selector)', () => {
    toEventSelector('event Transfer(address indexed from, address indexed to, uint256 amount)')
  }, benchOptions)
})
