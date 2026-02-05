/**
 * Call Action Examples (viem TypeScript)
 *
 * Comprehensive examples demonstrating the call action features:
 * - Basic contract calls
 * - Call with various parameters (gas, value, fees)
 * - State and block overrides
 * - Deployless calls
 * - Access lists
 * - Error handling
 */

import {
  type Address,
  createPublicClient,
  decodeFunctionResult,
  encodeFunctionData,
  formatEther,
  formatUnits,
  type Hex,
  http,
  parseEther,
  parseGwei,
} from 'viem'
import { polygon } from 'viem/chains'

// Example addresses
const USDC_ADDRESS: Address = '0x3c499c542cEF5E3811e1192ce70d8cC03d5c3359'
const VITALIK_ADDRESS: Address = '0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045'
const TEST_ADDRESS: Address = '0x1234567890123456789012345678901234567890'

// Simple ERC20 ABI for examples
const erc20Abi = [
  {
    name: 'name',
    type: 'function',
    stateMutability: 'view',
    inputs: [],
    outputs: [{ type: 'string' }],
  },
  {
    name: 'symbol',
    type: 'function',
    stateMutability: 'view',
    inputs: [],
    outputs: [{ type: 'string' }],
  },
  {
    name: 'decimals',
    type: 'function',
    stateMutability: 'view',
    inputs: [],
    outputs: [{ type: 'uint8' }],
  },
  {
    name: 'totalSupply',
    type: 'function',
    stateMutability: 'view',
    inputs: [],
    outputs: [{ type: 'uint256' }],
  },
  {
    name: 'balanceOf',
    type: 'function',
    stateMutability: 'view',
    inputs: [{ name: 'account', type: 'address' }],
    outputs: [{ type: 'uint256' }],
  },
] as const

// Helper functions
function printHeader(title: string) {
  console.log()
  console.log('='.repeat(60))
  console.log(`  ${title}`)
  console.log('='.repeat(60))
}

function printSection(title: string) {
  console.log(`\n--- ${title} ---`)
}

function truncateAddress(addr: string): string {
  return `${addr.slice(0, 10)}...${addr.slice(-4)}`
}

async function main() {
  printHeader('Call Action Examples (viem TypeScript)')

  // Create Public Client
  printSection('1. Creating Public Client')
  const client = createPublicClient({
    chain: polygon,
    transport: http('https://polygon-rpc.com'),
  })
  console.log('Connected to Polygon Mainnet')

  // Example 1: Basic Call - Read contract name
  printSection('2. Basic Call - Read ERC20 Name')
  try {
    const nameData = encodeFunctionData({
      abi: erc20Abi,
      functionName: 'name',
    })
    const result = await client.call({
      to: USDC_ADDRESS,
      data: nameData,
    })
    if (result.data) {
      const name = decodeFunctionResult({
        abi: erc20Abi,
        functionName: 'name',
        data: result.data,
      })
      console.log(`Contract Name: ${name}`)
    }
  } catch (error) {
    console.log(`Error: ${error}`)
  }

  // Example 2: Call with Address Parameter - balanceOf
  printSection('3. Call with Parameter - balanceOf(address)')
  try {
    const balanceOfData = encodeFunctionData({
      abi: erc20Abi,
      functionName: 'balanceOf',
      args: [VITALIK_ADDRESS],
    })
    const result = await client.call({
      to: USDC_ADDRESS,
      data: balanceOfData,
    })
    if (result.data) {
      const balance = decodeFunctionResult({
        abi: erc20Abi,
        functionName: 'balanceOf',
        data: result.data,
      })
      // USDC has 6 decimals
      console.log(`Vitalik's USDC Balance: ${formatUnits(balance, 6)} USDC`)
    }
  } catch (error) {
    console.log(`Error: ${error}`)
  }

  // Example 3: Call with From Address (Account)
  printSection('4. Call with Account (from address)')
  try {
    const decimalsData = encodeFunctionData({
      abi: erc20Abi,
      functionName: 'decimals',
    })
    const result = await client.call({
      account: VITALIK_ADDRESS,
      to: USDC_ADDRESS,
      data: decimalsData,
    })
    if (result.data) {
      const decimals = decodeFunctionResult({
        abi: erc20Abi,
        functionName: 'decimals',
        data: result.data,
      })
      console.log(
        `USDC Decimals: ${decimals} (called from ${truncateAddress(VITALIK_ADDRESS)})`,
      )
    }
  } catch (error) {
    console.log(`Error: ${error}`)
  }

  // Example 4: Call with Block Number
  printSection('5. Call at Specific Block Number')
  try {
    const totalSupplyData = encodeFunctionData({
      abi: erc20Abi,
      functionName: 'totalSupply',
    })
    const result = await client.call({
      to: USDC_ADDRESS,
      data: totalSupplyData,
      blockNumber: 82563588n, // Historical block
    })
    if (result.data) {
      const supply = decodeFunctionResult({
        abi: erc20Abi,
        functionName: 'totalSupply',
        data: result.data,
      })
      console.log(
        `USDC Total Supply at block 18000000: ${formatUnits(supply, 6)} USDC`,
      )
    }
  } catch (error) {
    console.log(`Error: ${error}`)
  }

  // Example 5: Call with Block Tag
  printSection('6. Call with Block Tag (pending)')
  try {
    const totalSupplyData = encodeFunctionData({
      abi: erc20Abi,
      functionName: 'totalSupply',
    })
    const result = await client.call({
      to: USDC_ADDRESS,
      data: totalSupplyData,
      blockTag: 'pending',
    })
    if (result.data) {
      const supply = decodeFunctionResult({
        abi: erc20Abi,
        functionName: 'totalSupply',
        data: result.data,
      })
      console.log(`USDC Total Supply (pending): ${formatUnits(supply, 6)} USDC`)
    }
  } catch (error) {
    console.log(`Error: ${error}`)
  }

  // Example 6: Call with Gas Parameters
  printSection('7. Call with Gas Parameters')
  try {
    const nameData = encodeFunctionData({
      abi: erc20Abi,
      functionName: 'name',
    })
    const result = await client.call({
      to: USDC_ADDRESS,
      data: nameData,
      gas: 100000n,
      gasPrice: parseGwei('500'),
    })
    console.log(`Call succeeded with gas=100000, gasPrice=500 gwei`)
    console.log(`Result data length: ${result.data?.length ?? 0} chars`)
  } catch (error) {
    console.log(`Error: ${error}`)
  }

  // Example 7: Call with EIP-1559 Fees
  printSection('8. Call with EIP-1559 Fees')
  try {
    const decimalsData = encodeFunctionData({
      abi: erc20Abi,
      functionName: 'decimals',
    })
    const result = await client.call({
      to: USDC_ADDRESS,
      data: decimalsData,
      maxFeePerGas: parseGwei('500'),
      maxPriorityFeePerGas: parseGwei('2'),
    })
    console.log(
      `Call succeeded with maxFeePerGas=50 gwei, maxPriorityFeePerGas=2 gwei ${result.data}`,
    )
  } catch (error) {
    console.log(`Error: ${error}`)
  }

  // Example 8: Call with Value (simulating ETH transfer check)
  printSection('9. Call with Value')
  try {
    const value = parseEther('0.1')
    const result = await client.call({
      account: VITALIK_ADDRESS,
      to: TEST_ADDRESS,
      value,
      data: '0x', // Empty call with value
    })
    console.log(
      `Simulated transfer of ${formatEther(value)} ETH would succeed ${result.data}`,
    )
  } catch (error) {
    console.log(`Simulated transfer error: ${error}`)
  }

  // Example 9: Call with State Override
  printSection('10. Call with State Override')
  try {
    const result = await client.call({
      account: TEST_ADDRESS,
      to: VITALIK_ADDRESS,
      value: parseEther('100'), // Transfer 100 ETH
      stateOverride: [
        {
          address: TEST_ADDRESS,
          balance: parseEther('1000'), // Override to 1000 ETH
        },
      ],
    })
    console.log(
      `State override successful! Test address had balance overridden to 1000 ETH ${result.data}`,
    )
    console.log('Simulated 100 ETH transfer succeeded')
  } catch (error) {
    console.log(`Error with state override: ${error}`)
  }

  // Example 10: Call with Block Override
  printSection('11. Call with Block Override')
  try {
    const totalSupplyData = encodeFunctionData({
      abi: erc20Abi,
      functionName: 'totalSupply',
    })
    const result = await client.call({
      to: USDC_ADDRESS,
      data: totalSupplyData,
      blockOverrides: {
        gasLimit: 50000000n,
        baseFeePerGas: parseGwei('1'),
        time: 1700000000n,
      },
    })
    console.log(`Block override successful! ${result.data}`)
    console.log('  Simulated gasLimit: 50000000')
    console.log('  Simulated baseFee: 1 gwei')
  } catch (error) {
    console.log(`Error with block override: ${error}`)
  }

  // Example 11: Call with Access List
  printSection('12. Call with Access List (EIP-2930)')
  try {
    const nameData = encodeFunctionData({
      abi: erc20Abi,
      functionName: 'name',
    })
    const result = await client.call({
      to: USDC_ADDRESS,
      data: nameData,
      accessList: [
        {
          address: USDC_ADDRESS,
          storageKeys: [
            '0x0000000000000000000000000000000000000000000000000000000000000000',
          ],
        },
      ],
    })
    console.log(`Call with access list succeeded ${result.data}`)
    console.log(`  Pre-warmed contract: ${truncateAddress(USDC_ADDRESS)}`)
  } catch (error) {
    console.log(`Error: ${error}`)
  }

  // Example 12: Deployless Call (Code parameter)
  printSection('13. Deployless Call (execute bytecode without deployment)')
  try {
    // Simple bytecode that returns 42
    // PUSH1 0x2a PUSH1 0x00 MSTORE PUSH1 0x20 PUSH1 0x00 RETURN
    const simpleBytecode: Hex = '0x602a60005260206000f3'
    const result = await client.call({
      code: simpleBytecode,
      data: '0x',
    })
    if (result.data) {
      const returnValue = BigInt(result.data)
      console.log(`Deployless call returned: ${returnValue}`)
    }
  } catch (error) {
    console.log(`Deployless call error: ${error}`)
  }

  // Example 13: Combined State + Block Override
  printSection('14. Combined State and Block Overrides')
  try {
    const balanceOfData = encodeFunctionData({
      abi: erc20Abi,
      functionName: 'balanceOf',
      args: [TEST_ADDRESS],
    })
    const result = await client.call({
      account: TEST_ADDRESS,
      to: USDC_ADDRESS,
      data: balanceOfData,
      stateOverride: [
        {
          address: TEST_ADDRESS,
          balance: parseEther('10000'),
          nonce: 100,
        },
      ],
      blockOverrides: {
        gasLimit: 100000000n,
      },
    })
    console.log(`Combined state + block override successful! ${result.data}`)
  } catch (error) {
    console.log(`Combined override error: ${error}`)
  }

  // Example 14: State Override with Storage
  printSection('15. State Override with Storage Slots')
  try {
    const nameData = encodeFunctionData({
      abi: erc20Abi,
      functionName: 'name',
    })
    const result = await client.call({
      to: USDC_ADDRESS,
      data: nameData,
      stateOverride: [
        {
          address: USDC_ADDRESS,
          // Override specific storage slots
          stateDiff: [
            {
              slot: '0x0000000000000000000000000000000000000000000000000000000000000000',
              value:
                '0x0000000000000000000000000000000000000000000000000000000000000001',
            },
          ],
        },
      ],
    })
    console.log(`State override with storage slots succeeded ${result.data}`)
  } catch (error) {
    console.log(`Storage override error: ${error}`)
  }

  // Example 15: Deployless Call via Factory
  printSection('16. Deployless Call via Factory Pattern')
  try {
    // This demonstrates the factory pattern for deployless calls
    // In practice, you'd use a real factory address and data
    const factoryAddress: Address = '0x4e59b44847b379578588920cA78FbF26c0B4956C' // CREATE2 deployer
    const result = await client.call({
      to: TEST_ADDRESS,
      data: '0x',
      factory: factoryAddress,
      factoryData: '0x1234', // Factory deployment data
    })
    console.log(`Deployless factory call completed ${result.data}`)
  } catch (error) {
    // This will likely fail without proper factory setup, but demonstrates the API
    console.log(
      `Factory call (expected to need proper setup): ${(error as Error).message?.slice(0, 50)}...`,
    )
  }

  // Example 16: Error Handling
  printSection('17. Error Handling Examples')

  // Test calling non-existent function
  console.log('\nTest: Call to non-contract address...')
  try {
    const result = await client.call({
      to: TEST_ADDRESS,
      data: '0x12345678', // Random selector
    })
    console.log(`  Call returned (empty contract) ${result.data}`)
  } catch (error) {
    console.log(`  Error: ${(error as Error).message?.slice(0, 60)}...`)
  }

  // Summary
  printHeader('Examples Complete')
  console.log('Demonstrated call features:')
  console.log('  - Basic contract calls')
  console.log('  - Calls with parameters (account, gas, value)')
  console.log('  - Block number and block tag queries')
  console.log('  - EIP-1559 fee parameters')
  console.log('  - State overrides (modify account state)')
  console.log('  - Block overrides (modify block context)')
  console.log('  - Access lists (EIP-2930)')
  console.log('  - Deployless calls (execute bytecode)')
  console.log('  - Factory pattern for deployless calls')
  console.log('  - Storage slot overrides')
  console.log('  - Error handling')
  console.log()
}

main().catch(console.error)
