/**
 * ENS Utility Benchmarks (viem TypeScript)
 *
 * These benchmarks mirror the Go benchmarks in ../go/ens_bench_test.go
 * for fair cross-language comparison.
 *
 * Pure CPU-bound operations -- no network involved.
 */

import { bench, describe } from 'vitest'
import { namehash, labelhash, normalize } from 'viem/ens'

const iterations = Number(process.env.BENCH_ITERATIONS ?? '100')

const benchOptions = {
  time: 0,
  warmupTime: 0,
  warmupIterations: 0,
  iterations,
}

describe('Ens', () => {
  bench('viem-ts: ens (namehash)', () => {
    namehash('vitalik.eth')
  }, benchOptions)

  bench('viem-ts: ens (namehash deep)', () => {
    namehash('sub.domain.vitalik.eth')
  }, benchOptions)

  bench('viem-ts: ens (labelhash)', () => {
    labelhash('vitalik')
  }, benchOptions)

  bench('viem-ts: ens (normalize)', () => {
    normalize('Vitalik.ETH')
  }, benchOptions)

  bench('viem-ts: ens (normalize long)', () => {
    normalize('My.Long.SubDomain.Name.vitalik.eth')
  }, benchOptions)
})
