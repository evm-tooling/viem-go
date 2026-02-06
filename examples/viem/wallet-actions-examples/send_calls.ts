/**
 * sendCalls — EIP-5792 batch call sending.
 *
 * Simulates each call individually via eth_call before attempting
 * to send the batch via wallet_sendCalls.
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
  printSection('sendCalls - EIP-5792 batch calls')

  const calls = [
    { to: TARGET, value: parseEther('0.0001') },
    { to: TARGET, value: parseEther('0.0002') },
    { to: TARGET, value: parseEther('0.0003') },
  ]

  // Step 1: Simulate each call individually
  for (let i = 0; i < calls.length; i++) {
    const call = calls[i]
    if (!call) continue

    const { to, value } = call
    console.log(
      `Simulating call ${i + 1}: ${account.address} -> ${to} (${formatEther(value)} POL)`,
    )
    try {
      await publicClient.call({
        account: account.address,
        to,
        value,
      })
      console.log(`  Simulation passed.`)
    } catch (err) {
      console.log(`  Simulation note: ${err as Error}`)
    }
  }

  // Step 2: Send batch via wallet_sendCalls (EIP-5792)
  console.log('\nSending batch via wallet_sendCalls...')
  try {
    const result = await walletClient.sendCalls({
      account,
      calls: calls.map((c) => ({ to: c.to, value: c.value })),
    })
    console.log(`Batch ID: ${result}`)
  } catch (err) {
    console.log(`sendCalls error: ${(err as Error).message.slice(0, 120)}`)
    console.log('(Expected — most public RPCs do not support wallet_sendCalls)')
  }
}
