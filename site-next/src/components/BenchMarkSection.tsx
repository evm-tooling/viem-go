'use client'
import { useState, useEffect } from "react";
import { motion, useReducedMotion } from "framer-motion";
import AnimatedSection from "./AnimatedSection";

const metrics = [
  { label: "ABI Encode", value: 19, max: 32 },
  { label: "Hash Functions", value: 32, max: 32 },
  { label: "RLP Encoding", value: 24, max: 32 },
  { label: "Sig Recovery", value: 9, max: 32 },
];

// Hook to detect mobile screen size
function useIsMobile(breakpoint = 768) {
  const [isMobile, setIsMobile] = useState(false);
  
  useEffect(() => {
    const checkMobile = () => setIsMobile(window.innerWidth < breakpoint);
    checkMobile();
    window.addEventListener("resize", checkMobile);
    return () => window.removeEventListener("resize", checkMobile);
  }, [breakpoint]);
  
  return isMobile;
}

const RadialGauge = ({ label, value, max, delay, disableAnimation }: { label: string; value: number; max: number; delay: number; disableAnimation: boolean }) => {
  const percentage = (value / max) * 100;
  const circumference = 2 * Math.PI * 40;
  const strokeDashoffset = circumference - (percentage / 100) * circumference;

  // If animations are disabled, render static version
  if (disableAnimation) {
    return (
      <div className="flex flex-col items-center gap-3">
        <div className="relative h-28 w-28">
          <svg className="h-full w-full -rotate-90" viewBox="0 0 96 96">
            <circle
              cx="48" cy="48" r="40"
              fill="none"
              stroke="hsl(var(--background-tertiary))"
              strokeWidth="6"
            />
            <circle
              cx="48" cy="48" r="40"
              fill="none"
              stroke="url(#gaugeGradient)"
              strokeWidth="6"
              strokeLinecap="round"
              strokeDasharray={circumference}
              strokeDashoffset={strokeDashoffset}
            />
            <defs>
              <linearGradient id="gaugeGradient" x1="0%" y1="0%" x2="100%" y2="0%">
                <stop offset="0%" stopColor="hsl(var(--primary))" />
                <stop offset="100%" stopColor="hsl(var(--tertiary))" />
              </linearGradient>
            </defs>
          </svg>
          <div className="absolute inset-0 flex items-center justify-center">
            <span className="text-xl font-black text-foreground">{value}×</span>
          </div>
        </div>
        <span className="text-xs font-medium text-foreground-secondary text-center">{label}</span>
      </div>
    );
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 30 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true }}
      transition={{ delay, duration: 0.5, ease: "easeOut" }}
      className="flex flex-col items-center gap-3"
    >
      <div className="relative h-28 w-28">
        <svg className="h-full w-full -rotate-90" viewBox="0 0 96 96">
          {/* Background circle */}
          <circle
            cx="48" cy="48" r="40"
            fill="none"
            stroke="hsl(var(--background-tertiary))"
            strokeWidth="6"
          />
          {/* Animated progress */}
          <motion.circle
            cx="48" cy="48" r="40"
            fill="none"
            stroke="url(#gaugeGradient)"
            strokeWidth="6"
            strokeLinecap="round"
            strokeDasharray={circumference}
            initial={{ strokeDashoffset: circumference }}
            whileInView={{ strokeDashoffset }}
            viewport={{ once: true }}
            transition={{ delay: delay + 0.2, duration: 1.2, ease: "easeOut" }}
          />
          <defs>
            <linearGradient id="gaugeGradient" x1="0%" y1="0%" x2="100%" y2="0%">
              <stop offset="0%" stopColor="hsl(var(--primary))" />
              <stop offset="100%" stopColor="hsl(var(--tertiary))" />
            </linearGradient>
          </defs>
        </svg>
        <div className="absolute inset-0 flex items-center justify-center">
          <span className="text-xl font-black text-foreground">{value}×</span>
        </div>
      </div>
      <span className="text-xs font-medium text-foreground-secondary text-center">{label}</span>
    </motion.div>
  );
};

const BenchmarkBigNumber = () => {
  const shouldReduceMotion = useReducedMotion();
  const isMobile = useIsMobile();
  const disableAnimations = shouldReduceMotion || isMobile;

  return (
    <section className="relative py-12 overflow-hidden bg-background">
      <div className="absolute top-0 left-0 right-0 h-px bg-gradient-to-r from-transparent via-primary/30 to-transparent" />

      {/* Subtle radial glow behind the big number */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[600px] rounded-full opacity-10"
        style={{ background: `radial-gradient(circle, hsl(var(--primary)) 0%, transparent 70%)` }}
      />

      <div className="relative mx-auto max-w-4xl px-6">
        <AnimatedSection>
          <div className="mb-4 text-center">
            <span className="inline-block rounded-full border border-primary/30 bg-primary/10 px-4 py-1 text-xs font-semibold text-primary uppercase tracking-wider">
              Performance
            </span>
          </div>
        </AnimatedSection>

        {/* Big hero number */}
        <AnimatedSection>
          <div className="text-center mb-4">
            <span className="text-[8rem] sm:text-[10rem] font-black leading-none text-primary tracking-tighter">
              20×
            </span>
          </div>
        </AnimatedSection>

        <AnimatedSection delay={0.2}>
          <h2 className="text-center mb-3 text-2xl sm:text-3xl">
            Faster than viem. On average.
          </h2>
          <p className="mx-auto max-w-lg text-center mb-12">
            Viem-go is blazingly fast. along with careful optimisations, Go's compiled runtime crushes JavaScript across every benchmark. No JIT warm-up, no GC pauses.
          </p>
        </AnimatedSection>

        {/* Radial gauges */}
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-8 justify-items-center">
          {metrics.map((m, i) => (
            <RadialGauge key={m.label} {...m} delay={i * 0.12} disableAnimation={disableAnimations} />
          ))}
        </div>

        <AnimatedSection delay={0.7}>
          <p className="text-center text-xs text-foreground-muted mt-10">
            * Benchmarks measured on Apple M2, Go 1.22, Node.js 20. Results may vary.
          </p>
        </AnimatedSection>
      </div>
    </section>
  );
};

export default BenchmarkBigNumber;
