const features = [
  {
    title: "Familiar API",
    description:
      "Same Client/Transport and Actions patterns as viem for TypeScript developers",
  },
  {
    title: "Idiomatic Go",
    description:
      "Built with Go conventions: explicit errors, context, and interfaces",
  },
  {
    title: "Type Safe",
    description:
      "Go's static typing for contract ABIs, transactions, and RPC calls",
  },
  {
    title: "go-ethereum",
    description: "Built on proven go-ethereum cryptographic primitives",
  },
];

export default function FeaturesGrid() {
  return (
    <div className="grid grid-cols-4 gap-2 mt-12 max-lg:grid-cols-2 max-sm:grid-cols-1">
      {features.map((feature) => (
        <div
          key={feature.title}
          className="relative h-[168px] rounded-lg border border-accent/15 p-5 mt-2 flex flex-col justify-start bg-gray-6/50 backdrop-blur-sm transition-all duration-200 hover:border-accent/35 hover:bg-gray-6/70 max-lg:h-[142px]"
        >
          <h3 className="text-[1.25rem] font-semibold m-0 text-white">
            {feature.title}
          </h3>
          <p className="text-[1rem] text-gray-2 m-0 leading-relaxed mt-2">
            {feature.description}
          </p>
        </div>
      ))}
    </div>
  );
}
