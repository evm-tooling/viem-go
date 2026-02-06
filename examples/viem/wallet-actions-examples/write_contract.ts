/**
 * writeContract — Simulate then call a contract write function.
 *
 * Demonstrates the simulate-before-send pattern using a real ERC-20 transfer
 * on Polygon mainnet. The simulation uses eth_call, and if the account had
 * funds the actual write would follow.
 */

import { encodeFunctionData, formatUnits, parseUnits } from 'viem'
import {
  account,
  printSection,
  publicClient,
  TARGET,
  USDC_ADDRESS,
  walletClient,
} from './shared.ts'

const erc20Abi = [
  {
    name: 'transfer',
    type: 'function',
    stateMutability: 'nonpayable',
    inputs: [
      { name: 'to', type: 'address' },
      { name: 'amount', type: 'uint256' },
    ],
    outputs: [{ name: '', type: 'bool' }],
  },
  {
    name: 'balanceOf',
    type: 'function',
    stateMutability: 'view',
    inputs: [{ name: 'account', type: 'address' }],
    outputs: [{ name: '', type: 'uint256' }],
  },
] as const

export async function run() {
  printSection('writeContract - Simulate then call contract write')

  const amount = parseUnits('10', 6) // 10 USDC

  // Encode the calldata
  const data = encodeFunctionData({
    abi: erc20Abi,
    functionName: 'transfer',
    args: [TARGET, amount],
  })
  console.log(`Encoded transfer(${TARGET}, ${formatUnits(amount, 6)} USDC)`)
  console.log(`Calldata: ${data.slice(0, 42)}...`)

  // Step 1: Simulate the write via eth_call
  console.log('\nSimulating contract write...')
  try {
    const result = await publicClient.call({
      account: account.address,
      to: USDC_ADDRESS,
      data,
    })
    console.log(`Simulation result: ${result.data}`)
  } catch (err) {
    console.log(`Simulation note: ${(err as Error).message.slice(0, 120)}`)
    console.log('(Expected — account has no USDC on Polygon)')
  }

  // Step 2: Send the actual write (would work if account had USDC)
  console.log('\nSending writeContract...')
  try {
    const hash = await walletClient.writeContract({
      account,
      address: USDC_ADDRESS,
      abi: erc20Abi,
      functionName: 'transfer',
      args: [TARGET, amount],
    })
    console.log(`Tx hash: ${hash}`)

    const receipt = await publicClient.waitForTransactionReceipt({ hash })
    console.log(`Status: ${receipt.status} | Gas used: ${receipt.gasUsed}`)
  } catch (err) {
    console.log(`Write error: ${(err as Error).message.slice(0, 120)}`)
    console.log('(Expected on mainnet with unfunded account)')
  }
}
