/**
 * deployContract — Simulate then deploy a contract.
 *
 * Demonstrates the deploy pattern: bytecode + ABI-encoded constructor args.
 * Since the account is unfunded on Polygon, the actual deploy will fail,
 * but the simulation and encoding pattern is correct.
 */

import { encodeDeployData, type Hex } from 'viem'
import { account, printSection, publicClient, walletClient } from './shared.ts'

const simpleStorageAbi = [
  {
    type: 'constructor',
    inputs: [{ name: '_value', type: 'uint256' }],
    stateMutability: 'nonpayable',
  },
  {
    name: 'storedData',
    type: 'function',
    stateMutability: 'view',
    inputs: [],
    outputs: [{ name: '', type: 'uint256' }],
  },
  {
    name: 'set',
    type: 'function',
    stateMutability: 'nonpayable',
    inputs: [{ name: 'x', type: 'uint256' }],
    outputs: [],
  },
] as const

// Minimal SimpleStorage bytecode (stores a uint256 in constructor)
const bytecode: Hex =
  '0x6080604052348015600e575f5ffd5b50604051606f380380606f833981016040819052602b916035565b5f55604b565b5f60208284031215604457005b5051919050565b6018806100575f395ff3fe6080604052005ffea164736f6c634300081d000a'

export async function run() {
  printSection('deployContract - Simulate then deploy a contract')

  // Encode the deploy data (bytecode + constructor args)
  const deployData = encodeDeployData({
    abi: simpleStorageAbi,
    bytecode,
    args: [42n],
  })
  console.log(`Deploy data (bytecode + constructor(42)):`)
  console.log(`  Length: ${deployData.length} hex chars`)
  console.log(`  Prefix: ${deployData.slice(0, 42)}...`)

  // Step 1: Simulate the deployment via eth_call (no `to` = contract creation)
  console.log('\nSimulating deployment...')
  try {
    const result = await publicClient.call({
      account: account.address,
      data: deployData,
    })
    console.log(`Simulation result length: ${result.data?.length ?? 0} chars`)
  } catch (err) {
    console.log(`Simulation note: ${(err as Error).message.slice(0, 120)}`)
    console.log('(Expected — account is unfunded on Polygon)')
  }

  // Step 2: Deploy (would work if account had POL for gas)
  console.log('\nDeploying contract...')
  try {
    const hash = await walletClient.deployContract({
      account,
      abi: simpleStorageAbi,
      bytecode,
      args: [42n],
    })
    console.log(`Deploy tx hash: ${hash}`)

    const receipt = await publicClient.waitForTransactionReceipt({ hash })
    console.log(
      `Status: ${receipt.status} | Contract: ${receipt.contractAddress}`,
    )
  } catch (err) {
    console.log(`Deploy error: ${(err as Error).message.slice(0, 120)}`)
    console.log('(Expected on mainnet with unfunded account)')
  }
}
