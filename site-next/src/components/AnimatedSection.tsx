"use client";
import { LazyMotion, domAnimation, m } from "framer-motion";

// Use m.div instead of motion.div

import { ReactNode } from "react";

interface AnimatedSectionProps {
  children: ReactNode;
  className?: string;
  delay?: number;
}

const AnimatedSection = ({ children, className = "", delay = 0 }: AnimatedSectionProps) => (
    <LazyMotion features={domAnimation}>
  <m.div
    initial={{ opacity: 0, y: 40 }}
    whileInView={{ opacity: 1, y: 0 }}
    viewport={{ once: true, margin: "-80px" }}
    transition={{ duration: 0.6, delay, ease: "easeOut" }}
    className={className}
  >
    {children}
  </m.div>
  </LazyMotion>
);

export default AnimatedSection;
