import { parseArgs } from "util"
import { runClientExample } from "./src/viem_ts_client"
import { runReadContractExample } from "./src/viem_ts_read_contract"

async function main() {
    const { values } = parseArgs({
        args: Bun.argv.slice(2),
        options: {
            example: { type: "string", default: "all" }
        }
    })

    console.log("╔═══════════════════════════════════════╗")
    console.log("║         viem Examples Runner          ║")
    console.log("╚═══════════════════════════════════════╝")

    switch (values.example) {
        case "all":
            await runAll()
            break
        case "client":
            await runClientExample()
            break
        case "read-contract":
            await runReadContractExample()
            break
        default:
            console.error(`Unknown example: ${values.example}`)
            process.exit(1)
    }

    console.log("========================================")
    console.log("All examples completed successfully!")
    console.log("========================================")
}

async function runAll() {
    await runClientExample()
    await runReadContractExample()
}

main()
