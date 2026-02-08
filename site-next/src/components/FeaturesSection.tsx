import { Card, CardDescription, CardTitle } from "@/components/ui/card";

const features = [
  {
    title: "JSON-RPC Abstractions",
    description:
      "High-level APIs over the JSON-RPC API to make your life easier",
  },
  {
    title: "Smart Contracts",
    description: "First-class APIs for interacting with Smart Contracts",
  },
  {
    title: "Ethereum Terminology",
    description:
      "Language closely aligned to official Ethereum terminology",
  },
  {
    title: "Multiple Accounts",
    description: "Private Key, Mnemonic, and HD Wallet account types",
  },
  {
    title: "Safe Numerics",
    description: "Go's native big.Int for safe numeric operations",
  },
  {
    title: "ABI Utilities",
    description:
      "Encoding, decoding, and inspection utilities for ABIs",
  },
  {
    title: "Token Bindings",
    description: "Pre-built bindings for ERC20, ERC721, and ERC1155",
  },
  {
    title: "Battle Tested",
    description:
      "Test suite running against Anvil for real-world compatibility",
  },
];

export default function FeaturesSection() {
  return (
    <section className="w-full py-8 px-8 bg-gradient-to-b from-dark-deep/30 to-transparent">
      <div className="max-w-[1120px] mx-auto mb-10 text-center">
        <h2 className="heading-2 mb-2">Features</h2>
        <p className="text-lead max-w-[600px] mx-auto">
          viem-go supports all main features from the original viem typescript
          library. Every feature was built using the same syntax, method-names
          and patterns so that the developer friendly nature still remains.
        </p>
      </div>
      <div className="grid grid-cols-4 gap-4 max-w-[1120px] mx-auto max-lg:grid-cols-2 max-sm:grid-cols-1">
        {features.map((feature) => (
          <Card key={feature.title} variant="interactive">
            <CardTitle>{feature.title}</CardTitle>
            <CardDescription>{feature.description}</CardDescription>
          </Card>
        ))}
      </div>
    </section>
  );
}
