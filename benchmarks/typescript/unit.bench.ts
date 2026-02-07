/**
 * Unit Parsing/Formatting Benchmarks (viem TypeScript)
 *
 * These benchmarks mirror the Go benchmarks in ../go/unit_bench_test.go
 * for fair cross-language comparison.
 *
 * Pure CPU-bound operations -- no network involved.
 */

import { bench, describe } from 'vitest'
import {
  parseEther,
  formatEther,
  parseUnits,
  parseGwei,
  formatUnits,
} from 'viem'

const benchOptions = {
  time: 2000,
  warmupTime: 0,
  warmupIterations: 0,
}

describe('Unit', () => {
  bench('viem-ts: unit (parseEther)', () => {
    parseEther('1.5')
  }, benchOptions)

  bench('viem-ts: unit (parseEther large)', () => {
    parseEther('123456789.123456789012345678')
  }, benchOptions)

  bench('viem-ts: unit (formatEther)', () => {
    formatEther(1500000000000000000n)
  }, benchOptions)

  bench('viem-ts: unit (parseUnits 6)', () => {
    parseUnits('100.50', 6)
  }, benchOptions)

  bench('viem-ts: unit (parseGwei)', () => {
    parseGwei('20.5')
  }, benchOptions)

  bench('viem-ts: unit (formatUnits)', () => {
    formatUnits(100500000n, 6)
  }, benchOptions)
})
