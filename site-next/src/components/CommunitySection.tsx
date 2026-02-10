'use client'
import { motion, useMotionValue, useTransform, animate } from "framer-motion";
import { useEffect, useRef } from "react";
import AnimatedSection from "./AnimatedSection";
import { MessageSquare, Bug, BookOpen, Heart, Star, GitFork, ArrowRight } from "lucide-react";

/* ── Animated counter ── */
const Counter = ({ target, label, suffix = "" }: { target: number; label: string; suffix?: string }) => {
  const count = useMotionValue(0);
  const rounded = useTransform(count, (v) => Math.floor(v));
  const ref = useRef<HTMLSpanElement>(null);

  useEffect(() => {
    const unsubscribe = rounded.on("change", (v) => {
      if (ref.current) ref.current.textContent = `${v}${suffix}`;
    });
    return unsubscribe;
  }, [rounded, suffix]);

  return (
    <AnimatedSection>
      <motion.div
        className="text-center"
        onViewportEnter={() => {
          animate(count, target, { duration: 2, ease: "easeOut" });
        }}
        viewport={{ once: true }}
      >
        <span ref={ref} className="text-4xl font-black text-primary md:text-5xl">
          0{suffix}
        </span>
        <p className="mt-2 text-xs font-semibold uppercase tracking-wider text-foreground-muted">{label}</p>
      </motion.div>
    </AnimatedSection>
  );
};

/* ── Community links ── */
const communityLinks = [
  { icon: MessageSquare, title: "Discussions", desc: "Ask questions and share ideas", href: "#" },
  { icon: Bug, title: "Issues", desc: "Report bugs or request features", href: "#" },
  { icon: BookOpen, title: "viem Docs", desc: "Original viem documentation", href: "#" },
];

const CommunitySection = () => (
  <section className="relative pb-24 pt-6 overflow-hidden">
    {/* Tertiary glow orb */}
    <div className="glow-orb-tertiary top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px]" />

    <div className="relative mx-auto max-w-7xl px-6">
      <AnimatedSection>
        <div className="mb-4 text-center">
          <span className="section-badge section-badge-primary">Join Us</span>
        </div>
        <h2 className="text-center mb-4">Community</h2>
        <p className="mx-auto max-w-lg text-center mb-16">
          Check out the following places for more viem-go content
        </p>
      </AnimatedSection>

      {/* Stats counters */}
      <div className="grid grid-cols-3 gap-8 mb-20 max-w-2xl mx-auto">
        <Counter target={95} label="Test Coverage" suffix="%" />
        <Counter target={2} label="GitHub Stars" />
        <Counter target={12} label="Contributors" />
      </div>

      {/* Links as horizontal cards with arrows */}
      <div className="flex flex-col gap-4 mb-20 max-w-3xl mx-auto">
        {communityLinks.map((link, i) => (
          <AnimatedSection key={link.title} delay={i * 0.1}>
            <motion.a
              href={link.href}
              whileHover={{ x: 8 }}
              className="group flex items-center gap-6 rounded-xl border border-card-border bg-card p-5 transition-all duration-300 hover:border-primary/30 hover:shadow-lg hover:shadow-primary/5"
            >
              <div className="icon-box icon-box-primary h-12 w-12 shrink-0 rounded-xl transition-colors group-hover:bg-primary/20">
                <link.icon className="h-6 w-6" />
              </div>
              <div className="flex-1">
                <h4 className="text-base font-semibold text-foreground group-hover:text-primary transition-colors">
                  {link.title}
                </h4>
                <p className="text-sm text-foreground-secondary">{link.desc}</p>
              </div>
              <ArrowRight className="h-5 w-5 text-foreground-muted group-hover:text-primary transition-all opacity-0 group-hover:opacity-100 -translate-x-2 group-hover:translate-x-0" />
            </motion.a>
          </AnimatedSection>
        ))}
      </div>

      {/* Support CTA — split layout */}
      <AnimatedSection>
        <div className="relative rounded-2xl border border-card-border overflow-hidden">
          <div className="grid md:grid-cols-2">
            {/* Left: text */}
            <div className="relative p-10 md:p-12 bg-gradient-to-br from-card to-card/80">
              <div className="absolute -left-16 -bottom-16 h-48 w-48 rounded-full border border-tertiary/10" />
              <h3 className="relative mb-3">Support the Project</h3>
              <p className="relative mb-6 max-w-sm">
                Help make viem-go a sustainable open-source project
              </p>
              <div className="relative flex flex-wrap gap-3">
                <motion.button
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.97 }}
                  className="inline-flex items-center gap-2 rounded-lg bg-tertiary px-6 py-2.5 text-sm font-bold text-tertiary-foreground shadow-lg shadow-tertiary/20 transition hover:shadow-xl"
                >
                  <Heart className="h-4 w-4" />
                  Sponsor
                </motion.button>
                <motion.button
                  whileHover={{ scale: 1.05 }}
                  whileTap={{ scale: 0.97 }}
                  className="inline-flex items-center gap-2 rounded-lg border border-border bg-background-secondary px-6 py-2.5 text-sm font-bold text-foreground transition hover:bg-background-tertiary"
                >
                  <Star className="h-4 w-4" />
                  Star
                </motion.button>
              </div>
            </div>
            {/* Right: decorative metrics */}
            <div className="relative bg-background-tertiary/30 p-10 md:p-12 flex flex-col justify-center gap-6 border-t md:border-t-0 md:border-l border-border">
              <div className="absolute -right-20 -top-20 h-60 w-60 rounded-full border border-primary/10" />
              {[
                { icon: Star, label: "Stars on GitHub", value: "2" },
                { icon: GitFork, label: "Forks", value: "0" },
                { icon: Heart, label: "Sponsors", value: "0" },
              ].map((stat) => (
                <div key={stat.label} className="relative flex items-center gap-4">
                  <stat.icon className="h-5 w-5 text-primary/60" />
                  <div>
                    <span className="text-lg font-bold text-foreground">{stat.value}</span>
                    <span className="ml-2 text-sm text-foreground-muted">{stat.label}</span>
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      </AnimatedSection>
    </div>
  </section>
);

export default CommunitySection;
