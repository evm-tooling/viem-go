'use client'
import { motion } from "framer-motion";
import AnimatedSection from "./AnimatedSection";
import { Server, Terminal, Network, Globe } from "lucide-react";
import { Card } from "@/components/ui/card";

const reasons = [
  {
    icon: Server,
    title: "Backend Services",
    desc: "Build high-performance blockchain indexers, APIs, and microservices",
    accent: "primary",
  },
  {
    icon: Terminal,
    title: "CLI Tools",
    desc: "Create command-line utilities for wallet management and contract interactions",
    accent: "tertiary",
  },
  {
    icon: Network,
    title: "Infrastructure",
    desc: "Power validators, relayers, and other blockchain infrastructure",
    accent: "primary",
  },
  {
    icon: Globe,
    title: "Cross-platform",
    desc: "Compile to a single binary for any OS without dependencies",
    accent: "tertiary",
  },
];

const WhyGoSection = () => (
  <section className="relative overflow-hidden  bg-gradient-to-b from-dark-deep/30 to-transparent py-24">
    {/* Diagonal accent line */}
    <div className="absolute top-0 left-0 right-0 h-px  bg-gradient-to-b from-dark-deep/30 to-transparent" />
    <svg className="absolute top-0 right-0 w-96 h-96 opacity-5" viewBox="0 0 400 400">
      <circle cx="350" cy="50" r="300" fill="none" stroke="currentColor" strokeWidth="0.5" className="text-primary" />
      <circle cx="350" cy="50" r="200" fill="none" stroke="currentColor" strokeWidth="0.5" className="text-primary" />
      <circle cx="350" cy="50" r="100" fill="none" stroke="currentColor" strokeWidth="0.5" className="text-primary" />
    </svg>

    {/* Bright centered glow behind cards */}
    <div
      aria-hidden
      className="pointer-events-none absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-[50%] h-[400px] w-[1200px] rounded-full opacity-40 blur-[130px]"
      style={{ background: "radial-gradient(ellipse, hsl(215 90% 55% / 0.5), transparent 70%)" }}
    />
    <div
      aria-hidden
      className="pointer-events-none absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-[30%] h-[280px] w-[400px] rounded-full opacity-25 blur-[80px]"
      style={{ background: "radial-gradient(ellipse, hsl(215 90% 65% / 0.55), transparent 65%)" }}
    />

    <div className="relative mx-auto max-w-7xl px-6">
      <AnimatedSection>
        <div className="mb-4 text-center">
          <span className="inline-block rounded-full border border-primary/30 bg-primary/10 px-4 py-1 text-xs font-semibold text-primary uppercase tracking-wider">
            Use Cases
          </span>
        </div>
        <h2 className="text-center mb-4">Why Go?</h2>
        <p className="mx-auto max-w-2xl text-center mb-16">
          viem-go is designed for teams and projects that need Ethereum tooling in the Go ecosystem and have familiarity
          with the viem typescript library.
        </p>
      </AnimatedSection>

      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        {reasons.map((r, i) => (
          <AnimatedSection key={r.title} delay={i * 0.1}>
            <Card asChild variant="surfaceInteractive" padding="md" className="group relative h-full duration-300">
              <motion.div whileHover={{ y: -6 }}>
              {/* Top accent bar */}
              <div
                className={`absolute top-0 left-4 right-4 h-0.5 rounded-b-full ${r.accent === "tertiary" ? "bg-tertiary" : "bg-primary"} opacity-60`}
              />

              <div
                className={`mb-4 flex h-12 w-12 items-center justify-center rounded-xl ${
                  r.accent === "tertiary" ? "bg-tertiary/10 text-tertiary" : "bg-primary/10 text-primary"
                } transition-transform group-hover:scale-110`}
              >
                <r.icon className="h-6 w-6" />
              </div>
              <h4 className="mb-2 text-foreground">{r.title}</h4>
              <p className="text-sm text-foreground-secondary">{r.desc}</p>
              </motion.div>
            </Card>
          </AnimatedSection>
        ))}
      </div>
    </div>
  </section>
);

export default WhyGoSection;
