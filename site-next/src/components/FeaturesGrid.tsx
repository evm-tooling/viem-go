'use client'
import { useRef } from "react";
import { motion, useReducedMotion, useInView } from "framer-motion";
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

const FeaturesGrid = () => {
  const shouldReduceMotion = useReducedMotion();
  const ref = useRef(null);
  const isInView = useInView(ref, { once: false, margin: "100px" });

  return (
    <section ref={ref} className="relative overflow-hidden pb-10">
      <div className="relative flex overflow-hidden [mask-image:linear-gradient(to_right,transparent,black_5%,black_95%,transparent)]">
        <motion.div
          className="flex shrink-0 gap-5"
          animate={shouldReduceMotion || !isInView ? {} : { x: ["0%", "-50%"] }}
          transition={shouldReduceMotion ? {} : { duration: 45, repeat: Infinity, ease: "linear" }}
        >
          {doubled.map((f, i) => (
            <Card
              key={i}
              variant="surfaceInteractive"
              className="group relative shrink-0 w-[280px] border-card-border/60 bg-card/40 sm:backdrop-blur-md p-6 shadow-lg shadow-primary/5"
            >
              <div className="mb-4 icon-box icon-box-primary h-10 w-10 rounded-lg">
                <f.icon className="h-5 w-5" />
              </div>
              <h4 className="mb-2 text-foreground">{f.title}</h4>
              <p className="text-sm text-foreground-secondary">{f.desc}</p>
            </Card>
          ))}
        </motion.div>
      </div>
    </section>
  );
};

export default FeaturesGrid;
