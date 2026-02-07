/**
 * Address Utility Benchmarks (viem TypeScript)
 *
 * These benchmarks mirror the Go benchmarks in ../go/address_bench_test.go
 * for fair cross-language comparison.
 *
 * Pure CPU-bound operations -- no network involved.
 */

import { bench, describe } from 'vitest'
import {
  isAddress,
  getAddress,
  getContractAddress,
} from 'viem'

const benchOptions = {
  time: 2000,
  warmupTime: 0,
  warmupIterations: 0,
}

describe('Address', () => {
  bench('viem-ts: address (isAddress)', () => {
    isAddress('0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48')
  }, benchOptions)

  bench('viem-ts: address (isAddress lower)', () => {
    isAddress('0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48')
  }, benchOptions)

  bench('viem-ts: address (checksum)', () => {
    getAddress('0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48')
  }, benchOptions)

  bench('viem-ts: address (create)', () => {
    getContractAddress({
      from: '0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266',
      nonce: 1n,
      opcode: 'CREATE',
    })
  }, benchOptions)

  bench('viem-ts: address (create2)', () => {
    getContractAddress({
      from: '0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266',
      salt: '0x0000000000000000000000000000000000000000000000000000000000000001',
      bytecode: '0x6080604052348015',
      opcode: 'CREATE2',
    })
  }, benchOptions)
})
