/**
 * Shared setup for all wallet action examples.
 * Uses Polygon mainnet with a local private key for signing demos.
 */

import {
  type Address,
  createPublicClient,
  createWalletClient,
  formatEther,
  type Hex,
  http,
} from 'viem'
import { privateKeyToAccount } from 'viem/accounts'
import { polygon } from 'viem/chains'

// ---------------------------------------------------------------------------
// Constants
// ---------------------------------------------------------------------------

const POLYGON_RPC = 'https://polygon-rpc.com'

/**
 * Private key loaded from the PRIVATE_KEY environment variable.
 * Falls back to the Anvil account #0 key for offline signing demos.
 * Set PRIVATE_KEY in a .env file or export it in your shell.
 */
const DEFAULT_PRIVATE_KEY =
  '0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80'

export const PRIVATE_KEY = (process.env.PRIVATE_KEY ??
  DEFAULT_PRIVATE_KEY) as Hex

export const account = privateKeyToAccount(PRIVATE_KEY)

/** A well-known address used as a transfer target in examples. */
export const TARGET: Address = '0x5d9339C29f1582e08F2b69bfa5740D11E0898D1F'

/** USDC on Polygon */
export const USDC_ADDRESS: Address =
  '0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359'

// ---------------------------------------------------------------------------
// Clients
// ---------------------------------------------------------------------------

export const publicClient = createPublicClient({
  chain: polygon,
  transport: http(POLYGON_RPC),
})

export const walletClient = createWalletClient({
  account,
  chain: polygon,
  transport: http(POLYGON_RPC),
})

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

export function printHeader(title: string) {
  console.log()
  console.log('='.repeat(70))
  console.log(`  ${title}`)
  console.log('='.repeat(70))
}

export function printSection(title: string) {
  console.log(`\n--- ${title} ---`)
}

export async function printAccountInfo() {
  console.log(`Account:  ${account.address}`)
  console.log(`Chain:    Polygon (${polygon.id})`)
  try {
    const bal = await publicClient.getBalance({
      address: account.address,
    })
    console.log(`Balance:  ${formatEther(bal)} POL`)
  } catch {
    console.log('Balance:  (could not fetch)')
  }
}
