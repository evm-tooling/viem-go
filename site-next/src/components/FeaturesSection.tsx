'use client'
import AnimatedSection from "./AnimatedSection";
import { Card } from "@/components/ui/card";
import {
  Workflow,
  FileCode2,
  BookOpenText,
  Users,
  Calculator,
  Braces,
  Coins,
  ShieldCheck,
} from "lucide-react";

const features = [
  { icon: Workflow, title: "JSON-RPC Abstractions", desc: "High-level APIs over the JSON-RPC API to make your life easier", accent: "primary" as const },
  { icon: FileCode2, title: "Smart Contracts", desc: "First-class APIs for interacting with Smart Contracts", accent: "tertiary" as const },
  { icon: BookOpenText, title: "Ethereum Terminology", desc: "Language closely aligned to official Ethereum terminology", accent: "primary" as const },
  { icon: Users, title: "Multiple Accounts", desc: "Private Key, Mnemonic, and HD Wallet account types", accent: "tertiary" as const },
  { icon: Calculator, title: "Safe Numerics", desc: "Go's native big.Int for safe numeric operations", accent: "primary" as const },
  { icon: Braces, title: "ABI Utilities", desc: "Encoding, decoding, and inspection utilities for ABIs", accent: "tertiary" as const },
  { icon: Coins, title: "Token Bindings", desc: "Pre-built bindings for ERC20, ERC721, and ERC1155", accent: "primary" as const },
  { icon: ShieldCheck, title: "Battle Tested", desc: "Test suite running against Anvil for real-world compatibility", accent: "tertiary" as const },
];

export default function FeaturesSection() {
  return (
    <section className="relative w-full py-16 px-8 section-bg-dark overflow-hidden">
      {/* Glow orbs behind cards */}
      <div aria-hidden className="glow-orb-ellipse top-1/2 left-1/2 -translate-x-1/2 -translate-y-[55%] h-[600px] w-[1200px] opacity-30 blur-[130px]" />
      <div aria-hidden className="glow-orb-ellipse top-1/2 left-1/2 -translate-x-1/2 -translate-y-[25%] h-[300px] w-[420px] opacity-20 blur-[80px]" />

      <div className="relative max-w-[1120px] mx-auto mb-12 text-center">
        <AnimatedSection>
          <div className="mb-4">
            <span className="section-badge section-badge-primary">Built-in</span>
          </div>
          <h2 className="mb-4">Features</h2>
          <p className="text-lead max-w-[600px] mx-auto">
            viem-go supports all main features from the original viem typescript
            library. Every feature was built using the same syntax, method-names
            and patterns so that the developer friendly nature still remains.
          </p>
        </AnimatedSection>
      </div>

      <div className="relative grid grid-cols-4 gap-6 max-w-[1120px] mx-auto max-lg:grid-cols-2 max-sm:grid-cols-1">
        {features.map((f, i) => (
          <AnimatedSection key={f.title} delay={i * 0.06}>
            <Card variant="surfaceInteractive" padding="md" className="group relative h-full hover:-translate-y-1 transition-transform duration-200">
              {/* Top accent bar */}
              <div
                className={`absolute top-0 left-4 right-4 h-0.5 rounded-b-full ${
                  f.accent === "tertiary" ? "bg-tertiary" : "bg-primary"
                } opacity-60`}
              />

              <div className={`mb-4 icon-box h-12 w-12 rounded-xl ${
                f.accent === "tertiary" ? "icon-box-tertiary" : "icon-box-primary"
              }`}>
                <f.icon className="h-6 w-6" />
              </div>
              <h4 className="mb-2 text-foreground">{f.title}</h4>
              <p className="text-sm text-foreground-secondary">{f.desc}</p>
            </Card>
          </AnimatedSection>
        ))}
      </div>
    </section>
  );
}
