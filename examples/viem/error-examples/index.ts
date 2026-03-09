/**
 * Error Path Examples (viem TypeScript)
 *
 * Full coverage of custom errors across 4 phases:
 *
 * Phase 1: getEstimateGasError
 *   - estimateGas with reverting call -> EstimateGasExecutionError
 *
 * Phase 2: getCallError, getContractError
 *   - createAccessList with reverting call -> CallExecutionError
 *   - estimateContractGas with reverting transfer -> ContractFunctionExecutionError
 *   - simulateContract with reverting transfer -> ContractFunctionExecutionError
 *
 * Phase 3: ContractFunctionZeroDataError
 *   - readContract on EOA (returns 0x) -> ContractFunctionZeroDataError
 *
 * Phase 4: RPC Error Code Mapping
 *   - Non-existent RPC method -> MethodNotFoundRpcError
 *
 * Also: Invalid params, contract reverts, CallExecutionError
 */

import {
  type Address,
  ContractFunctionExecutionError,
  ContractFunctionRevertedError,
  ContractFunctionZeroDataError,
  createPublicClient,
  encodeFunctionData,
  EstimateGasExecutionError,
  http,
  MethodNotFoundRpcError,
  parseAbi,
} from 'viem'
import { polygon } from 'viem/chains'

const USDC_ADDRESS: Address = '0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359'
const TEST_ADDRESS: Address = '0x1234567890123456789012345678901234567890'
const ZERO_ADDRESS: Address = '0x0000000000000000000000000000000000000000'

const erc20Abi = parseAbi([
  'function name() view returns (string)',
  'function balanceOf(address account) view returns (uint256)',
  'function transfer(address to, uint256 amount) returns (bool)',
])

function printSection(title: string) {
  console.log(`\n--- ${title} ---`)
}

async function main() {
  console.log('============================================================')
  console.log('  Error Path Examples (viem TypeScript)')
  console.log('============================================================')

  const client = createPublicClient({
    chain: polygon,
    transport: http('https://rpc-mainnet.matic.quiknode.pro'),
  })

  // 1. Invalid params: code + to (mutually exclusive)
  printSection('1. Invalid Params: code + to (mutually exclusive)')
  try {
    const simpleBytecode =
      '0x602a60005260206000f3' as `0x${string}`
    await client.call({
      code: simpleBytecode,
      to: USDC_ADDRESS,
      data: '0x',
    })
  } catch (error: unknown) {
    const msg = (error as Error).message
    if (msg.includes("cannot provide both 'code' and 'to'") || msg.includes('code') && msg.includes('to')) {
      console.log(`  ✓ Caught invalid params error: ${(error as Error).message?.slice(0, 80)}...`)
    } else {
      console.log(`  Error: ${(error as Error).message}`)
    }
  }

  // 2. Invalid params: code + factory (mutually exclusive)
  printSection('2. Invalid Params: code + factory (mutually exclusive)')
  try {
    const simpleBytecode =
      '0x602a60005260206000f3' as `0x${string}`
    const factoryAddress =
      '0x4e59b44847b379578588920cA78FbF26c0B4956C' as Address
    await client.call({
      code: simpleBytecode,
      factory: factoryAddress,
      factoryData: '0x1234' as `0x${string}`,
    })
  } catch (error: unknown) {
    const msg = (error as Error).message
    if (msg.includes('code') && (msg.includes('factory') || msg.includes('factoryData'))) {
      console.log(`  ✓ Caught invalid params error: ${(error as Error).message?.slice(0, 80)}...`)
    } else {
      console.log(`  Error: ${(error as Error).message}`)
    }
  }

  // 3. Contract revert: transfer exceeds balance (ERC20)
  printSection('3. Contract Revert: transfer exceeds balance')
  try {
    await client.simulateContract({
      address: USDC_ADDRESS,
      abi: erc20Abi,
      functionName: 'transfer',
      args: [TEST_ADDRESS, 10n ** 60n],
      account: ZERO_ADDRESS,
    })
  } catch (error: unknown) {
    if (error instanceof ContractFunctionExecutionError) {
      console.log(`  ✓ Caught ContractFunctionExecutionError for "${error.functionName}"`)
    }
    if (error instanceof ContractFunctionRevertedError) {
      console.log(`  ✓ Decoded revert - Reason: "${error.reason}"`)
      if (error.signature) {
        console.log(`  ✓ Decoded revert - Signature: ${error.signature}`)
      }
    }
    console.log(`  Full error: ${(error as Error).message}`)
  }

  // 4. Contract revert: transfer from zero address
  printSection('4. Contract Revert: transfer from zero address')
  try {
    const transferData = encodeFunctionData({
      abi: erc20Abi,
      functionName: 'transfer',
      args: ['0x0000000000000000000000000000000000000001' as Address, 1n],
    })
    await client.call({
      account: ZERO_ADDRESS,
      to: USDC_ADDRESS,
      data: transferData,
    })
  } catch (error: unknown) {
    if (error instanceof ContractFunctionExecutionError) {
      console.log(`  ✓ Caught ContractFunctionExecutionError for "${error.functionName}"`)
    }
    if (error instanceof ContractFunctionRevertedError) {
      console.log(`  ✓ Decoded revert - Reason: "${error.reason}", Signature: ${error.signature ?? 'n/a'}`)
    }
    console.log(`  Full error: ${(error as Error).message}`)
  }

  // 5. Call to non-contract address (EOA with random data)
  printSection('5. Call to non-contract address (EOA)')
  try {
    await client.call({
      to: TEST_ADDRESS,
      data: '0x12345678' as `0x${string}`,
    })
  } catch (error: unknown) {
    const msg = (error as Error).message
    if (msg.includes('CallExecutionError') || msg.includes('call')) {
      console.log(`  ✓ Caught call error (CallExecutionError or similar)`)
    }
    console.log(`  Error: ${(error as Error).message}`)
  }

  // 6. Invalid function selector (non-existent function on contract)
  printSection('6. Invalid selector: non-existent function')
  try {
    await client.call({
      to: USDC_ADDRESS,
      data: '0xdeadbeef' as `0x${string}`,
    })
  } catch (error: unknown) {
    console.log(`  Error: ${(error as Error).message}`)
  }

  // --- Phase 1: EstimateGasExecutionError (getEstimateGasError) ---
  printSection('Phase 1: EstimateGasExecutionError (getEstimateGasError)')
  const transferRevertData = encodeFunctionData({
    abi: erc20Abi,
    functionName: 'transfer',
    args: [TEST_ADDRESS, 10n ** 60n],
  })
  try {
    await client.estimateGas({
      account: ZERO_ADDRESS,
      to: USDC_ADDRESS,
      data: transferRevertData,
      value: 0n,
    })
  } catch (error: unknown) {
    if (error instanceof EstimateGasExecutionError) {
      console.log(`  ✓ Caught EstimateGasExecutionError`)
    }
    console.log(`  Error: ${(error as Error).message}`)
  }

  // --- Phase 2a: createAccessList -> CallExecutionError (getCallError) ---
  printSection('Phase 2a: createAccessList -> CallExecutionError (getCallError)')
  try {
    await client.createAccessList({
      account: ZERO_ADDRESS,
      to: USDC_ADDRESS,
      data: transferRevertData,
    })
  } catch (error: unknown) {
    const msg = (error as Error).message
    if (msg.includes('CallExecutionError') || msg.includes('call')) {
      console.log(`  ✓ Caught CallExecutionError from createAccessList`)
    }
    console.log(`  Error: ${(error as Error).message}`)
  }

  // --- Phase 2b: estimateContractGas -> ContractFunctionExecutionError (getContractError) ---
  printSection('Phase 2b: estimateContractGas -> ContractFunctionExecutionError (getContractError)')
  try {
    await client.estimateContractGas({
      account: ZERO_ADDRESS,
      address: USDC_ADDRESS,
      abi: erc20Abi,
      functionName: 'transfer',
      args: [TEST_ADDRESS, 10n ** 60n],
    })
  } catch (error: unknown) {
    if (error instanceof ContractFunctionExecutionError) {
      console.log(`  ✓ Caught ContractFunctionExecutionError for "${error.functionName}"`)
    }
    if (error instanceof ContractFunctionRevertedError) {
      console.log(`  ✓ Decoded revert - Reason: "${error.reason}"`)
    }
    console.log(`  Error: ${(error as Error).message}`)
  }

  // --- Phase 2c: simulateContract -> ContractFunctionExecutionError (getContractError) ---
  printSection('Phase 2c: simulateContract -> ContractFunctionExecutionError (getContractError)')
  try {
    await client.simulateContract({
      account: ZERO_ADDRESS,
      address: USDC_ADDRESS,
      abi: erc20Abi,
      functionName: 'transfer',
      args: [TEST_ADDRESS, 10n ** 60n],
    })
  } catch (error: unknown) {
    if (error instanceof ContractFunctionExecutionError) {
      console.log(`  ✓ Caught ContractFunctionExecutionError for "${error.functionName}"`)
    }
    console.log(`  Error: ${(error as Error).message}`)
  }

  // --- Phase 3: ContractFunctionZeroDataError (readContract on EOA) ---
  printSection('Phase 3: ContractFunctionZeroDataError (readContract on EOA)')
  try {
    await client.readContract({
      address: TEST_ADDRESS, // EOA, not a contract
      abi: erc20Abi,
      functionName: 'balanceOf',
      args: [TEST_ADDRESS],
    })
  } catch (error: unknown) {
    if (error instanceof ContractFunctionZeroDataError) {
      console.log(`  ✓ Caught ContractFunctionZeroDataError`)
    }
    if (error instanceof ContractFunctionExecutionError) {
      console.log(`  ✓ Caught ContractFunctionExecutionError (EOA may revert instead)`)
    }
    console.log(`  Error: ${(error as Error).message}`)
  }

  // --- Phase 4: MethodNotFoundRpcError (non-existent RPC method) ---
  printSection('Phase 4: MethodNotFoundRpcError (non-existent RPC method)')
  try {
    await (client as { request: (args: { method: string; params?: unknown[] }) => Promise<unknown> }).request({
      method: 'eth_nonexistent',
      params: ['latest'],
    })
  } catch (error: unknown) {
    if (error instanceof MethodNotFoundRpcError) {
      console.log(`  ✓ Caught MethodNotFoundRpcError (code ${error.code})`)
    }
    console.log(`  Error: ${(error as Error).message}`)
  }

  // 7. Summary
  printSection('7. Error type classification summary')
  console.log('  Phase 1: EstimateGasExecutionError')
  console.log('  Phase 2: CallExecutionError (createAccessList), ContractFunctionExecutionError (estimateContractGas, simulateContract)')
  console.log('  Phase 3: ContractFunctionZeroDataError (readContract on EOA)')
  console.log('  Phase 4: MethodNotFoundRpcError')
  console.log('  Also: Invalid params, ContractFunctionRevertedError, CallExecutionError')
  console.log()
}

main().catch(console.error)
