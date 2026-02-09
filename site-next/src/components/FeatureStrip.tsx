'use client'
import { motion } from "framer-motion";
import { Code2, Layers, Shield, Cpu, Zap, GitBranch, Lock, Database } from "lucide-react";
import { Card } from "@/components/ui/card";

const features = [
  { icon: Layers, title: "Familiar API", desc: "Same Client/Transport and Actions patterns as viem for TypeScript developers" },
  { icon: Code2, title: "Idiomatic Go", desc: "Built with Go conventions: explicit errors, context, and interfaces" },
  { icon: Shield, title: "Type Safe", desc: "Go's static typing for contract ABIs, transactions, and RPC calls" },
  { icon: Cpu, title: "go-ethereum", desc: "Built on proven go-ethereum cryptographic primitives" },
  { icon: Zap, title: "High Performance", desc: "Leverage Go's concurrency model for parallel RPC operations" },
  { icon: GitBranch, title: "Composable", desc: "Modular architecture lets you import only what you need" },
  { icon: Lock, title: "Battle Tested", desc: "Comprehensive test suite running against Anvil" },
  { icon: Database, title: "ABI Utilities", desc: "Encoding, decoding, and inspection utilities for ABIs" },
];

const doubled = [...features, ...features];

const FeaturesStrip = () => (
  <section className="relative overflow-hidden py-10">
    <div className="absolute top-0 left-0 right-0 h-px bg-gradient-to-r from-transparent via-primary/50 to-transparent" />

    <div className="relative flex overflow-hidden [mask-image:linear-gradient(to_right,transparent,black_50%,black_80%,transparent)]">
      <motion.div
        className="flex shrink-0 gap-5"
        animate={{ x: ["0%", "-50%"] }}
        transition={{ duration: 45, repeat: Infinity, ease: "linear" }}
      >
        {doubled.map((f, i) => (
          <Card
            key={i}
            variant="surfaceInteractive"
            className="group relative shrink-0 w-[280px] border-card-border/60 bg-card/40 backdrop-blur-md p-6 shadow-lg shadow-primary/5 duration-300 hover:border-primary/40 hover:shadow-xl hover:shadow-primary/10"
          >
            {/* Hover glow */}
            <div className="absolute -inset-px rounded-xl bg-gradient-to-b from-primary/5 to-transparent opacity-0 group-hover:opacity-100 transition-opacity -z-10" />

            <div className="mb-4 flex h-10 w-10 items-center justify-center rounded-lg bg-primary/10 text-primary transition-colors group-hover:bg-primary/20">
              <f.icon className="h-5 w-5" />
            </div>
            <h4 className="mb-2 text-foreground group-hover:text-primary transition-colors">{f.title}</h4>
            <p className="text-sm text-foreground-secondary">{f.desc}</p>
          </Card>
        ))}
      </motion.div>
    </div>
  </section>
);

export default FeaturesStrip;
