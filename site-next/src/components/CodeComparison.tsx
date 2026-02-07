import { CodeGroup } from "./CodePanel";

export default function CodeComparison() {
  return (
    <section className="max-w-[70ch] mx-auto my-16">
      <h2 className="text-[2.5rem] font-semibold mb-4 text-white">Overview</h2>
      <p className="mb-6 text-[1.1rem] text-gray-2 leading-relaxed">
        viem-go brings the developer experience of{" "}
        <a href="https://viem.sh" className="text-accent hover:underline">
          viem
        </a>{" "}
        to the Go ecosystem. If you&apos;re familiar with viem in TypeScript,
        you&apos;ll feel right at home.
      </p>

      <CodeGroup
        tabs={[
          {
            title: "viem-go",
            language: "go",
            code: `// 1. Import modules.
import (
    "github.com/ChefBingbong/viem-go/client"
    "github.com/ChefBingbong/viem-go/client/transport"
    "github.com/ChefBingbong/viem-go/chain/definitions"
)

// 2. Set up your client with desired chain & transport.
c, _ := client.CreatePublicClient(client.PublicClientConfig{
    Chain:     definitions.Mainnet,
    Transport: transport.HTTP("https://eth.llamarpc.com"),
})
defer c.Close()

// 3. Consume an action!
blockNumber, _ := c.GetBlockNumber(context.Background())`,
          },
          {
            title: "viem",
            language: "ts",
            code: `// 1. Import modules.
import { createPublicClient, http } from 'viem'
import { mainnet } from 'viem/chains'

// 2. Set up your client with desired chain & transport.
const client = createPublicClient({
  chain: mainnet,
  transport: http(),
})

// 3. Consume an action!
const blockNumber = await client.getBlockNumber()`,
          },
        ]}
      />
    </section>
  );
}
