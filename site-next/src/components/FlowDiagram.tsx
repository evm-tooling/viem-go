interface FlowNode {
  label: string;
  sublabel?: string;
  items?: string[];
  variant?: "default" | "accent" | "muted";
}

interface FlowBranch {
  condition: string;
  label: string;
}

interface FlowDiagramProps {
  nodes: FlowNode[];
  /** Optional branching at the end (e.g. HTTP vs WebSocket) */
  branches?: FlowBranch[];
}

function nodeColors(variant: FlowNode["variant"] = "default") {
  switch (variant) {
    case "accent":
      return "border-tertiary/40 bg-tertiary/8 text-tertiary";
    case "muted":
      return "border-border/60 bg-background-secondary/50 text-foreground-muted";
    default:
      return "border-primary/30 bg-primary/6 text-foreground";
  }
}

function Arrow() {
  return (
    <div className="flex flex-col items-center py-1">
      <div className="w-px h-5 bg-gradient-to-b from-primary/50 to-primary/20" />
      <svg
        width="10"
        height="8"
        viewBox="0 0 10 8"
        className="text-primary/40 -mt-px"
      >
        <path d="M5 8L0 0h10z" fill="currentColor" />
      </svg>
    </div>
  );
}

export default function FlowDiagram({ nodes, branches }: FlowDiagramProps) {
  return (
    <div className="my-6 flex flex-col items-center">
      {nodes.map((node, i) => (
        <div key={i} className="flex flex-col items-center w-full max-w-[480px]">
          {i > 0 && <Arrow />}
          <div
            className={`w-full rounded-lg border px-4 py-3 transition-colors ${nodeColors(node.variant)}`}
          >
            <div className="text-sm font-semibold leading-snug">
              {node.label}
            </div>
            {node.sublabel && (
              <div className="text-xs text-foreground-muted mt-0.5 leading-snug">
                {node.sublabel}
              </div>
            )}
            {node.items && node.items.length > 0 && (
              <ul className="mt-2 space-y-0.5">
                {node.items.map((item, j) => (
                  <li
                    key={j}
                    className="text-xs text-foreground-secondary leading-relaxed flex items-start gap-1.5"
                  >
                    <span className="text-tertiary mt-px shrink-0">-</span>
                    <span className="font-mono">{item}</span>
                  </li>
                ))}
              </ul>
            )}
          </div>
        </div>
      ))}

      {branches && branches.length > 0 && (
        <>
          <div className="flex flex-col items-center py-1">
            <div className="w-px h-4 bg-gradient-to-b from-primary/50 to-primary/20" />
          </div>
          <div className="flex gap-4 w-full max-w-[480px]">
            {branches.map((branch, i) => (
              <div key={i} className="flex-1 flex flex-col items-center">
                <div className="text-[10px] uppercase tracking-wider text-foreground-muted font-semibold mb-1.5">
                  {branch.condition}
                </div>
                <div className="w-full rounded-lg border border-border/60 bg-background-secondary/50 px-3 py-2.5 text-center">
                  <span className="text-xs text-foreground-secondary font-mono leading-relaxed">
                    {branch.label}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </>
      )}
    </div>
  );
}
