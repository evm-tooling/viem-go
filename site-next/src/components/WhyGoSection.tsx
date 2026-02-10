'use client'
import AnimatedSection from "./AnimatedSection";
import { Server, Terminal, Network, Globe } from "lucide-react";
import { Card } from "@/components/ui/card";

const reasons = [
  { icon: Server, title: "Backend Services", desc: "Build high-performance blockchain indexers, APIs, and microservices", accent: "primary" },
  { icon: Terminal, title: "CLI Tools", desc: "Create command-line utilities for wallet management and contract interactions", accent: "tertiary" },
  { icon: Network, title: "Infrastructure", desc: "Power validators, relayers, and other blockchain infrastructure", accent: "primary" },
  { icon: Globe, title: "Cross-platform", desc: "Compile to a single binary for any OS without dependencies", accent: "tertiary" },
];

const WhyGoSection = () => (
  <section className="relative overflow-hidden section-bg-dark py-24">
    {/* Top divider */}
    <div className="absolute top-0 left-0 right-0 h-px bg-gradient-to-r from-transparent via-tertiary/40 to-transparent" />

    {/* Decorative circles */}
    <svg className="absolute top-0 right-0 w-96 h-96 opacity-5" viewBox="0 0 400 400">
      <circle cx="350" cy="50" r="300" fill="none" stroke="currentColor" strokeWidth="0.5" className="text-primary" />
      <circle cx="350" cy="50" r="200" fill="none" stroke="currentColor" strokeWidth="0.5" className="text-primary" />
      <circle cx="350" cy="50" r="100" fill="none" stroke="currentColor" strokeWidth="0.5" className="text-primary" />
    </svg>

    {/* Glow orbs */}
    <div aria-hidden className="glow-orb-ellipse top-1/2 left-1/2 -translate-x-1/2 -translate-y-[50%] h-[400px] w-[1200px] opacity-40 blur-[130px]" />
    <div aria-hidden className="glow-orb-ellipse top-1/2 left-1/2 -translate-x-1/2 -translate-y-[30%] h-[280px] w-[400px] opacity-25 blur-[80px]" />

    <div className="relative mx-auto max-w-7xl px-6">
      <AnimatedSection>
        <div className="mb-4 text-center">
          <span className="section-badge section-badge-primary">Use Cases</span>
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
            <Card variant="surfaceInteractive" padding="md" className="group relative h-full hover:-translate-y-1 transition-transform duration-200">
              <div
                className={`absolute top-0 left-4 right-4 h-0.5 rounded-b-full ${r.accent === "tertiary" ? "bg-tertiary" : "bg-primary"} opacity-60`}
              />
              <div className={`mb-4 icon-box h-12 w-12 rounded-xl ${
                r.accent === "tertiary" ? "icon-box-tertiary" : "icon-box-primary"
              }`}>
                <r.icon className="h-6 w-6" />
              </div>
              <h4 className="mb-2 text-foreground">{r.title}</h4>
              <p className="text-sm text-foreground-secondary">{r.desc}</p>
            </Card>
          </AnimatedSection>
        ))}
      </div>
    </div>
  </section>
);

export default WhyGoSection;
