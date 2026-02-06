/**
 * signTransaction — Sign a transaction without broadcasting.
 *
 * Demonstrates EIP-1559, EIP-2930, and legacy transaction signing.
 * All signing happens locally; no RPC call is needed.
 */

import { parseEther, parseGwei } from 'viem'
import { account, printSection, TARGET, walletClient } from './shared.ts'

export async function run() {
  printSection('signTransaction - Sign without broadcasting')

  // EIP-1559
  const signedEIP1559 = await walletClient.signTransaction({
    account,
    to: TARGET,
    value: parseEther('0.01'),
    gas: 21_000n,
    maxFeePerGas: parseGwei('50'),
    maxPriorityFeePerGas: parseGwei('2'),
    type: 'eip1559',
  })
  console.log(`Signed tx (EIP-1559): ${signedEIP1559.slice(0, 42)}...`)

  // Legacy
  const signedLegacy = await walletClient.signTransaction({
    account,
    to: TARGET,
    value: parseEther('0.01'),
    gas: 21_000n,
    gasPrice: parseGwei('50'),
    type: 'legacy',
  })
  console.log(`Signed tx (legacy):   ${signedLegacy.slice(0, 42)}...`)

  // EIP-2930
  const signedEIP2930 = await walletClient.signTransaction({
    account,
    to: TARGET,
    value: parseEther('0.01'),
    gas: 21_000n,
    gasPrice: parseGwei('50'),
    type: 'eip2930',
    accessList: [
      {
        address: TARGET,
        storageKeys: [
          '0x0000000000000000000000000000000000000000000000000000000000000001',
        ],
      },
    ],
  })
  console.log(`Signed tx (EIP-2930): ${signedEIP2930.slice(0, 42)}...`)
}
