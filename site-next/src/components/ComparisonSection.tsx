'use client'
import { motion } from "framer-motion";
import AnimatedSection from "./AnimatedSection";
import { Check, X, Minus } from "lucide-react";
import { Card } from "@/components/ui/card";

type Status = "yes" | "no" | "partial";

const rows: { feature: string; viemGo: Status; goEth: Status }[] = [
  { feature: "High-level API", viemGo: "yes", goEth: "no" },
  { feature: "Familiar to viem users", viemGo: "yes", goEth: "no" },
  { feature: "Type-safe ABI encoding", viemGo: "yes", goEth: "yes" },
  { feature: "Multiple account types", viemGo: "yes", goEth: "partial" },
  { feature: "Transport abstraction", viemGo: "yes", goEth: "no" },
  { feature: "Contract bindings", viemGo: "yes", goEth: "partial" },
  { feature: "Low learning curve", viemGo: "yes", goEth: "no" },
];

const StatusCell = ({ status }: { status: Status }) => {
  if (status === "yes")
    return (
      <div className="flex items-center gap-2">
        <div className="flex h-6 w-6 items-center justify-center rounded-full bg-success/15">
          <Check className="h-3.5 w-3.5 text-success" />
        </div>
        <span className="text-xs text-success font-medium hidden sm:inline">Supported</span>
      </div>
    );
  if (status === "no")
    return (
      <div className="flex items-center gap-2">
        <div className="flex h-6 w-6 items-center justify-center rounded-full bg-destructive/15">
          <X className="h-3.5 w-3.5 text-destructive" />
        </div>
        <span className="text-xs text-destructive font-medium hidden sm:inline">No</span>
      </div>
    );
  return (
    <div className="flex items-center gap-2">
      <div className="flex h-6 w-6 items-center justify-center rounded-full bg-warning/15">
        <Minus className="h-3.5 w-3.5 text-warning" />
      </div>
      <span className="text-xs text-warning font-medium hidden sm:inline">Partial</span>
    </div>
  );
};

const ComparisonSection = () => (
  <section className="w-full py-8 px-8 section-bg-dark">
    <div className="relative mx-auto max-w-4xl px-6">
      <AnimatedSection>
        <div className="mb-4 text-center">
          <span className="section-badge section-badge-tertiary">Head to Head</span>
        </div>
        <h1 className="text-center mb-4">Viem-go vs Go-ethereum</h1>
        <p className="mx-auto max-w-2xl text-center mb-12">
          See how viem-go compares to the standard go-ethereum library.
        </p>
      </AnimatedSection>

      <AnimatedSection delay={0.15}>
        <Card variant="surface" padding="none" className="overflow-hidden shadow-xl shadow-primary/5">
          {/* Table header — uses default-light bg tier */}
          <div className="grid grid-cols-[1fr_140px_140px] sm:grid-cols-[1fr_180px_180px] border-b border-border bg-background-elevated/40">
            <div className="px-6 py-4 text-xs font-bold uppercase tracking-wider text-foreground-muted">Feature</div>
            <div className="px-4 py-4 text-center">
              <span className="inline-flex items-center gap-1.5 rounded-full bg-primary/10 px-3 py-1 text-xs font-bold text-primary">viem-go</span>
            </div>
            <div className="px-4 py-4 text-center">
              <span className="inline-flex items-center gap-1.5 rounded-full bg-background-secondary px-3 py-1 text-xs font-bold text-foreground-secondary">go-ethereum</span>
            </div>
          </div>

          {/* Table rows */}
          {rows.map((row, i) => (
            <motion.div
              key={row.feature}
              initial={{ opacity: 0, x: -20 }}
              whileInView={{ opacity: 1, x: 0 }}
              viewport={{ once: true }}
              transition={{ delay: i * 0.05 }}
              className={`grid grid-cols-[1fr_140px_140px] sm:grid-cols-[1fr_180px_180px] items-center transition-colors hover:bg-background-elevated/40 ${
                i !== rows.length - 1 ? "border-b border-border/50" : ""
              }`}
            >
              <div className="px-6 py-4 text-sm font-medium text-foreground">{row.feature}</div>
              <div className="flex justify-center px-4 py-4"><StatusCell status={row.viemGo} /></div>
              <div className="flex justify-center px-4 py-4"><StatusCell status={row.goEth} /></div>
            </motion.div>
          ))}

          {/* Summary bar */}
          <div className="grid grid-cols-[1fr_140px_140px] sm:grid-cols-[1fr_180px_180px] border-t border-border bg-background-elevated/30">
            <div className="px-6 py-3 text-xs font-bold uppercase tracking-wider text-foreground-muted">Score</div>
            <div className="px-4 py-3 text-center"><span className="text-lg font-black text-primary">7/7</span></div>
            <div className="px-4 py-3 text-center"><span className="text-lg font-black text-foreground-muted">2/7</span></div>
          </div>
        </Card>
      </AnimatedSection>
    </div>
  </section>
);

export default ComparisonSection;
