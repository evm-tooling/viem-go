/**
 * signAuthorization / prepareAuthorization — EIP-7702 authorization signing.
 *
 * Signs an authorization object locally that allows an EOA to delegate
 * to a contract implementation.
 */

import { account, printSection, TARGET, walletClient } from './shared.ts'

export async function run() {
  printSection('signAuthorization - EIP-7702 authorization signing')

  try {
    const authorization = await walletClient.signAuthorization({
      account,
      contractAddress: TARGET,
    })
    console.log(`Authorization address: ${authorization.address}`)
    console.log(`Authorization chainId: ${authorization.chainId}`)
    console.log(`Authorization nonce:   ${authorization.nonce}`)
    console.log(`Signature r: ${authorization.r.slice(0, 22)}...`)
    console.log(`Signature s: ${authorization.s.slice(0, 22)}...`)
  } catch (err) {
    // Some RPC providers may not support the nonce lookup for EIP-7702
    console.log(
      `signAuthorization error: ${(err as Error).message.slice(0, 120)}`,
    )
    console.log(
      '(This is expected if the RPC does not support pending nonce lookups for EIP-7702)',
    )
  }
}
