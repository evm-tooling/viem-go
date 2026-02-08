/**
 * ABI Encoding/Decoding Benchmarks (viem TypeScript)
 *
 * These benchmarks mirror the Go benchmarks in ../go/abi_bench_test.go
 * for fair cross-language comparison.
 *
 * Pure CPU-bound operations -- no network involved.
 */

import { bench, describe } from 'vitest'
import {
  encodeFunctionData,
  decodeFunctionResult,
  encodePacked,
  parseAbi,
} from 'viem'

const erc20Abi = parseAbi([
  'function name() view returns (string)',
  'function symbol() view returns (string)',
  'function decimals() view returns (uint8)',
  'function totalSupply() view returns (uint256)',
  'function balanceOf(address owner) view returns (uint256)',
  'function allowance(address owner, address spender) view returns (uint256)',
  'function transfer(address to, uint256 amount) returns (bool)',
  'function approve(address spender, uint256 amount) returns (bool)',
  'function transferFrom(address from, address to, uint256 amount) returns (bool)',
])

const VITALIK = '0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045' as const
const ANVIL_0 = '0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266' as const

// Pre-encoded balanceOf return data (uint256 = 1000000)
const balanceOfReturnData = '0x00000000000000000000000000000000000000000000000000000000000f4240' as const

const iterations = Number(process.env.BENCH_ITERATIONS ?? '100')

const benchOptions = {
  time: 0,
  warmupTime: 0,
  warmupIterations: 0,
  iterations,
}

describe('Abi', () => {
  bench('viem-ts: abi (encode simple)', () => {
    encodeFunctionData({
      abi: erc20Abi,
      functionName: 'balanceOf',
      args: [VITALIK],
    })
  }, benchOptions)

  bench('viem-ts: abi (encode complex)', () => {
    encodeFunctionData({
      abi: erc20Abi,
      functionName: 'transfer',
      args: [VITALIK, 1000000n],
    })
  }, benchOptions)

  bench('viem-ts: abi (encode multi-arg)', () => {
    encodeFunctionData({
      abi: erc20Abi,
      functionName: 'transferFrom',
      args: [ANVIL_0, VITALIK, 1000000n],
    })
  }, benchOptions)

  bench('viem-ts: abi (decode result)', () => {
    decodeFunctionResult({
      abi: erc20Abi,
      functionName: 'balanceOf',
      data: balanceOfReturnData,
    })
  }, benchOptions)

  bench('viem-ts: abi (encodePacked)', () => {
    encodePacked(
      ['address', 'uint256'],
      ['0x14dC79964da2C08b23698B3D3cc7Ca32193d9955', 420n],
    )
  }, benchOptions)

  bench('viem-ts: abi (encodePacked multi)', () => {
    encodePacked(
      ['address', 'string', 'uint256', 'bool'],
      ['0x14dC79964da2C08b23698B3D3cc7Ca32193d9955', 'hello world', 420n, true],
    )
  }, benchOptions)
})
