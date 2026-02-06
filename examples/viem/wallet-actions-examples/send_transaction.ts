/**
 * sendTransaction — Simulate then send an ETH transfer.
 *
 * Uses Polygon mainnet. The transaction is simulated via `publicClient.call()`
 * before the actual send. Since the example account is not funded on Polygon,
 * the simulation will reflect insufficient balance, but the pattern is correct.
 */

import { formatEther, parseEther } from 'viem'
import {
  account,
  printSection,
  publicClient,
  TARGET,
  walletClient,
} from './shared.ts'

export async function run() {
  printSection('sendTransaction - Simulate then send ETH transfer')

  const value = parseEther('0.0001')

  // Step 1: Simulate the transaction via eth_call
  console.log(
    `Simulating: ${account.address} -> ${TARGET} (${formatEther(value)} POL)`,
  )
  try {
    await publicClient.call({
      account: account.address,
      to: TARGET,
      value,
    })
    console.log('Simulation passed.')
  } catch (err) {
    console.log(`Simulation note: ${(err as Error).message.slice(0, 120)}`)
    console.log('(Expected on mainnet with unfunded account)')
  }

  // Step 2: Send the real transaction (local account -> prepare + sign + sendRaw)
  console.log('\nSending transaction...')
  try {
    const hash = await walletClient.sendTransaction({
      account,
      to: TARGET,
      value,
    })
    console.log(`Tx hash: ${hash}`)

    const receipt = await publicClient.waitForTransactionReceipt({ hash })
    console.log(`Status: ${receipt.status} | Gas used: ${receipt.gasUsed}`)
  } catch (err) {
    console.log(`Send error: ${(err as Error).message.slice(0, 120)}`)
    console.log('(Expected on mainnet with unfunded account)')
  }
}
