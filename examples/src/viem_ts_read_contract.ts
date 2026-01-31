import 'dotenv/config'
import { createPublicClient, erc20Abi, formatUnits, http } from 'viem'
import { mainnet } from 'viem/chains'

const URL = process.env.TENDERLY_RPC_URL
const ERC20_ADDRESS = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"

export async function runReadContractExample() {
    const publicClient = createPublicClient({
        chain: mainnet,
        transport: http(URL)
    })

    // Setup token
    const token = {
        address: ERC20_ADDRESS,
        abi: erc20Abi,
    } as const

    // Call contract functions
    const name = await publicClient.readContract({ ...token, functionName: 'name' })
    const symbol = await publicClient.readContract({ ...token, functionName: 'symbol' })
    const decimals = await publicClient.readContract({ ...token, functionName: 'decimals' })
    const totalSupply = await publicClient.readContract({ ...token, functionName: 'totalSupply' })

    // Log token details
    console.log(`\nReading from ${ERC20_ADDRESS}\n`)
    console.log(`Name: ${name}`)
    console.log(`Symbol: ${symbol}`)
    console.log(`Decimals: ${decimals}`)
    console.log(`Total Supply: ${totalSupply}\n`)

    // Get ERC20 balance
    const USER = "0x830690922a56f31Cb96851951587D8A2f45C0EBA"

    const balance = await publicClient.readContract({
        ...token,
        functionName: 'balanceOf',
        args: [USER]
    })

    // Log balances
    console.log(`Balance Returned: ${balance}`)
    console.log(`Balance Formatted: ${formatUnits(balance, Number(decimals))}\n`)
}

// Run if called directly
if (import.meta.main) {
    runReadContractExample()
}
