/**
 * sendRawTransaction — Sign locally, then broadcast a pre-signed transaction.
 *
 * Pattern: prepareTransactionRequest -> signTransaction -> sendRawTransaction.
 * Simulates first via eth_call.
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
  printSection('sendRawTransaction - Local sign + broadcast')

  const value = parseEther('0.0005')

  // Step 1: Simulate
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
  }

  // Step 2: Sign locally
  const tx = await walletClient.prepareTransactionRequest({
    account,
    to: TARGET,
    value,
  })
  console.log('\nSigning transaction locally...')
  const signedTx = await walletClient.signTransaction(tx)
  console.log(`Signed tx: ${signedTx.slice(0, 42)}...`)

  // Step 3: Broadcast
  console.log('\nBroadcasting raw transaction...')
  try {
    const hash = await walletClient.sendRawTransaction({
      serializedTransaction: signedTx,
    })
    console.log(`Tx hash: ${hash}`)

    const receipt = await publicClient.waitForTransactionReceipt({ hash })
    console.log(`Status: ${receipt.status} | Gas used: ${receipt.gasUsed}`)
  } catch (err) {
    console.log(`Broadcast error: ${err as Error}`)
    console.log('(Expected on mainnet with unfunded account)')
  }
}
