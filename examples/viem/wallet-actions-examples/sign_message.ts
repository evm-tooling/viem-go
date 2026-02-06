/**
 * signMessage — EIP-191 personal message signing.
 *
 * Signs a message locally using the private key account.
 * Demonstrates string messages, raw hex messages, and raw byte messages.
 */

import { account, printSection, walletClient } from './shared.ts'

export async function run() {
  printSection('signMessage - Sign a plain text message (EIP-191)')

  // String message
  const sig = await walletClient.signMessage({
    account,
    message: 'hello world',
  })
  console.log(`Message:   "hello world"`)
  console.log(`Signature: ${sig.slice(0, 42)}...`)

  // Raw hex message (same content: "hello world" = 0x68656c6c6f20776f726c64)
  const rawHexSig = await walletClient.signMessage({
    account,
    message: { raw: '0x68656c6c6f20776f726c64' },
  })
  console.log(`Raw hex:   ${rawHexSig.slice(0, 42)}... (same content)`)

  // Raw bytes message
  const rawBytesSig = await walletClient.signMessage({
    account,
    message: {
      raw: Uint8Array.from([
        104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100,
      ]),
    },
  })
  console.log(`Raw bytes: ${rawBytesSig.slice(0, 42)}... (same content)`)

  // All three should produce the same signature
  console.log(`\nAll match: ${sig === rawHexSig && sig === rawBytesSig}`)
}
