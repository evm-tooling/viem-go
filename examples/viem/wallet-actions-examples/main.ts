/**
 * Wallet Actions Examples — Entry Point
 *
 * Usage:
 *   bun run main.ts                       # Run all examples
 *   bun run main.ts signMessage           # Run a single example
 *   bun run main.ts signMessage signTx    # Run multiple examples
 *
 * Available examples:
 *   signMessage, signTypedData, signTransaction, signAuthorization,
 *   sendTransaction, sendRawTransaction, writeContract, deployContract,
 *   sendCalls
 */

import * as deployContract from './deploy_contract.ts'
import * as sendCalls from './send_calls.ts'
import * as sendRawTransaction from './send_raw_transaction.ts'
import * as sendTransaction from './send_transaction.ts'
import { printAccountInfo, printHeader } from './shared.ts'
import * as signAuthorization from './sign_authorization.ts'
import * as signMessage from './sign_message.ts'
import * as signTransaction from './sign_transaction.ts'
import * as signTypedData from './sign_typed_data.ts'
import * as writeContract from './write_contract.ts'

// ---------------------------------------------------------------------------
// Registry
// ---------------------------------------------------------------------------

const examples: Record<string, { run: () => Promise<void> }> = {
  signMessage,
  signTypedData,
  signTransaction,
  signAuthorization,
  sendTransaction,
  sendRawTransaction,
  writeContract,
  deployContract,
  sendCalls,
}

const exampleNames = Object.keys(examples)

// ---------------------------------------------------------------------------
// CLI
// ---------------------------------------------------------------------------

async function main() {
  const args = process.argv.slice(2)

  // Resolve which examples to run
  let selected: string[]
  if (args.length === 0) {
    selected = exampleNames
  } else {
    // Validate args
    const invalid = args.filter((a) => !examples[a])
    if (invalid.length > 0) {
      console.error(`Unknown example(s): ${invalid.join(', ')}`)
      console.error(`Available: ${exampleNames.join(', ')}`)
      process.exit(1)
    }
    selected = args
  }

  printHeader(`Wallet Actions Examples (viem / bun) — Polygon Mainnet`)
  await printAccountInfo()

  console.log(`\nRunning ${selected.length} example(s): ${selected.join(', ')}`)

  for (const name of selected) {
    try {
      const test = examples[name]
      await test?.run()
    } catch (err) {
      console.error(`\n[${name}] Unhandled error: ${err}`)
    }
  }

  printHeader('Done')
}

main().catch(console.error)
