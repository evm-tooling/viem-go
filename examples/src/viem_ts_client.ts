import { createPublicClient, formatUnits, http } from 'viem'
import { mainnet } from 'viem/chains'

export async function runClientExample() {
    const publicClient = createPublicClient({
        chain: mainnet,
        transport: http(),
    })

    // Fetch current block number
    const blockNumber = await publicClient.getBlockNumber()
    console.log(`\nCurrent Block Number: ${blockNumber}\n`)

    // Get balance of an address
    const address = '0x73BCEb1Cd57C711feaC4224D062b0F6ff338501e'
    const balance = await publicClient.getBalance({ address })

    console.log(`ETH Balance of ${address}: ${formatUnits(balance, 18)} ETH\n`)
}

// Run if called directly
if (import.meta.main) {
    runClientExample()
}
