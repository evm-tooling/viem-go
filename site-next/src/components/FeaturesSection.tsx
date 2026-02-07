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
        <h2 className="text-[2.5rem] font-semibold text-white mb-2">
          Features
        </h2>
        <p className="text-[1.1rem] text-gray-2 max-w-[600px] mx-auto leading-relaxed">
          viem-go supports all main features from the original viem typescript
          library. Every feature was built using the same syntax, method-names
          and patterns so that the developer friendly nature still remains.
        </p>
      </div>
      <div className="grid grid-cols-4 gap-4 max-w-[1120px] mx-auto max-lg:grid-cols-2 max-sm:grid-cols-1">
        {features.map((feature) => (
          <div
            key={feature.title}
            className="bg-gray-6/50 border border-accent/20 rounded-xl p-6 transition-all duration-200 hover:border-accent/40 hover:bg-gray-6/70 hover:-translate-y-0.5"
          >
            <h3 className="text-[1.2rem] font-semibold text-white mb-1.5">
              {feature.title}
            </h3>
            <p className="text-[1rem] text-gray-2 leading-relaxed">
              {feature.description}
            </p>
          </div>
        ))}
      </div>
    </section>
  );
}
