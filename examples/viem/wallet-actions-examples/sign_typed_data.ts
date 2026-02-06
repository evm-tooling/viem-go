/**
 * signTypedData — EIP-712 structured typed data signing.
 *
 * Signs typed data locally. Demonstrates the classic "Ether Mail" example.
 */

import { account, printSection, walletClient } from './shared.ts'

export async function run() {
  printSection('signTypedData - Sign EIP-712 typed data')

  const sig = await walletClient.signTypedData({
    account,
    domain: {
      name: 'Ether Mail',
      version: '1',
      chainId: 137, // Polygon
      verifyingContract: '0xCcCCccccCCCCcCCCCCCcCcCccCcCCCcCcccccccC',
    },
    types: {
      Person: [
        { name: 'name', type: 'string' },
        { name: 'wallet', type: 'address' },
      ],
      Mail: [
        { name: 'from', type: 'Person' },
        { name: 'to', type: 'Person' },
        { name: 'contents', type: 'string' },
      ],
    },
    primaryType: 'Mail',
    message: {
      from: {
        name: 'Cow',
        wallet: '0xCD2a3d9F938E13CD947Ec05AbC7FE734Df8DD826',
      },
      to: {
        name: 'Bob',
        wallet: '0xbBbBBBBbbBBBbbbBbbBbbbbBBbBbbbbBbBbbBBbB',
      },
      contents: 'Hello, Bob!',
    },
  })

  console.log(`Typed data signature: ${sig.slice(0, 42)}...`)
}
