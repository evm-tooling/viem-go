interface ComparisonItem {
  label: string;
  status: "yes" | "no" | "partial";
}

const viemGoItems: ComparisonItem[] = [
  { label: "High-level API", status: "yes" },
  { label: "Familiar to viem users", status: "yes" },
  { label: "Type-safe ABI encoding", status: "yes" },
  { label: "Multiple account types", status: "yes" },
  { label: "Transport abstraction", status: "yes" },
  { label: "Contract bindings", status: "yes" },
  { label: "Low learning curve", status: "yes" },
];

const goEthereumItems: ComparisonItem[] = [
  { label: "High-level API", status: "no" },
  { label: "Familiar to viem users", status: "no" },
  { label: "Type-safe ABI encoding", status: "yes" },
  { label: "Multiple account types", status: "partial" },
  { label: "Transport abstraction", status: "no" },
  { label: "Contract bindings (codegen)", status: "partial" },
  { label: "Low learning curve", status: "no" },
];

function StatusIcon({ status }: { status: ComparisonItem["status"] }) {
  if (status === "yes") {
    return (
      <span className="w-5 h-5 flex items-center justify-center font-semibold text-sm rounded-full text-[#4ade80] bg-[rgba(74,222,128,0.15)]">
        &#10003;
      </span>
    );
  }
  if (status === "no") {
    return (
      <span className="w-5 h-5 flex items-center justify-center font-semibold text-sm rounded-full text-[#f87171] bg-[rgba(248,113,113,0.15)]">
        &#10007;
      </span>
    );
  }
  return (
    <span className="w-5 h-5 flex items-center justify-center font-semibold text-sm rounded-full text-[#fbbf24] bg-[rgba(251,191,36,0.15)]">
      &#9675;
    </span>
  );
}

function ComparisonCard({
  title,
  items,
  highlighted,
}: {
  title: string;
  items: ComparisonItem[];
  highlighted?: boolean;
}) {
  return (
    <div
      className={`bg-gray-6/50 border rounded-xl overflow-hidden ${
        highlighted ? "border-accent/40" : "border-accent/20"
      }`}
    >
      <div
        className={`px-6 py-3 text-[1.1rem] font-semibold text-white border-b border-accent/15 ${
          highlighted ? "bg-accent/15" : "bg-dark-deep/60"
        }`}
      >
        {title}
      </div>
      <div className="px-6 py-4">
        {items.map((item, idx) => (
          <div
            key={idx}
            className={`flex items-center gap-3 py-2 text-[0.9rem] text-gray-2 ${
              idx < items.length - 1
                ? "border-b border-accent/[0.08]"
                : ""
            }`}
          >
            <StatusIcon status={item.status} />
            {item.label}
          </div>
        ))}
      </div>
    </div>
  );
}

export default function ComparisonSection() {
  return (
    <section className="w-full py-8 px-8 bg-gradient-to-b from-dark-deep/30 to-transparent">
      <div className="max-w-[1120px] mx-auto mb-10 text-center">
        <h2 className="text-[2.5rem] font-semibold text-white mb-2">
          viem-go vs go-ethereum
        </h2>
        <p className="text-[1.1rem] text-gray-2 max-w-[600px] mx-auto leading-relaxed">
          Viem-go aims to hit that middleground of low level control and high
          level abstraction resulting in the best developer experience. See how
          viem-go compares to the standard go-ethereum library
        </p>
      </div>
      <div className="grid grid-cols-2 gap-6 max-w-[900px] mx-auto max-md:grid-cols-1">
        <ComparisonCard title="viem-go" items={viemGoItems} highlighted />
        <ComparisonCard title="go-ethereum" items={goEthereumItems} />
      </div>
    </section>
  );
}
