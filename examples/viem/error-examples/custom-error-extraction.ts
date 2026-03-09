/**
 * Custom Error Extraction (viem TypeScript)
 *
 * Showcases the err.walk() pattern for extracting ContractFunctionRevertedError
 * and accessing errorName / args for custom error handling.
 *
 * Pattern from viem docs:
 *   err.walk(err => err instanceof ContractFunctionRevertedError)
 *   revertError.data?.errorName
 */

import {
  type Address,
  BaseError,
  ContractFunctionRevertedError,
  createPublicClient,
  http,
  parseAbi,
} from 'viem'
import { polygon } from 'viem/chains'

const USDC_ADDRESS: Address = '0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359'
const TEST_ADDRESS: Address = '0x1234567890123456789012345678901234567890'
const ZERO_ADDRESS: Address = '0x0000000000000000000000000000000000000000'

const erc20Abi = parseAbi([
  'function balanceOf(address account) view returns (uint256)',
  'function transfer(address to, uint256 amount) returns (bool)',
])

async function main() {
  console.log('============================================================')
  console.log('  Custom Error Extraction (viem TypeScript)')
  console.log('  Pattern: err.walk() -> ContractFunctionRevertedError')
  console.log('============================================================\n')

  const client = createPublicClient({
    chain: polygon,
    transport: http('https://rpc-mainnet.matic.quiknode.pro'),
  })

  // Simulate a reverting transfer (from zero address -> Error(string))
  try {
    await client.simulateContract({
      address: USDC_ADDRESS,
      abi: erc20Abi,
      functionName: 'transfer',
      args: [TEST_ADDRESS, 10n ** 60n],
      account: ZERO_ADDRESS,
    })
  } catch (err) {
    if (err instanceof BaseError) {
      const revertError = err.walk(
        (e) => e instanceof ContractFunctionRevertedError,
      )
      if (revertError instanceof ContractFunctionRevertedError) {
        const errorName = revertError.data?.errorName ?? ''
        const args = revertError.data?.args

        console.log('✓ Extracted ContractFunctionRevertedError via err.walk()')
        console.log(`  errorName: "${errorName}"`)
        console.log(`  reason:   "${revertError.reason ?? 'n/a'}"`)
        if (args && args.length > 0) {
          console.log(`  args:     ${JSON.stringify(args)}`)
        }
        console.log()

        // Example: branch on error type
        switch (errorName) {
          case 'Error':
            console.log(
              '  → Standard Error(string) - handle generic revert message',
            )
            break
          case 'Panic':
            console.log(
              '  → Standard Panic(uint256) - handle assertion failure',
            )
            break
          case 'ERC20InsufficientBalance':
          case 'ERC20InsufficientAllowance':
            console.log(
              '  → Custom ERC20 error - handle insufficient balance/allowance',
            )
            break
          default:
            console.log(`  → Custom error "${errorName}" - handle as needed`)
        }
      } else {
        console.log('  No ContractFunctionRevertedError in chain')
      }
    }
    console.log('\n  Full error:', (err as Error).message?.split('\n')[0])
  }

  console.log('\n--- Summary ---')
  console.log(
    '  viem supports: err.walk(fn) to find ContractFunctionRevertedError,',
  )
  console.log(
    '  then access revertError.data?.errorName and revertError.data?.args',
  )
  console.log()
}

main().catch(console.error)
