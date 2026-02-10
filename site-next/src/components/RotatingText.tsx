"use client";
import { useState, useEffect } from "react";
import { motion, AnimatePresence, useReducedMotion } from "framer-motion";

const phrases = [
  "Viem in Go",
  "20× Faster",
  "Type-Safe",
  "Modular",
];

interface RotatingTextProps {
  disableAnimation?: boolean;
}

const RotatingText = ({ disableAnimation = false }: RotatingTextProps) => {
  const [index, setIndex] = useState(0);
  const shouldReduceMotion = useReducedMotion();
  const noAnimation = disableAnimation || shouldReduceMotion;

  useEffect(() => {
    // Don't rotate at all if animations are disabled
    if (disableAnimation) return;
    
    // Slow down rotation interval for reduced motion users
    const intervalMs = shouldReduceMotion ? 5000 : 3000;
    const interval = setInterval(() => {
      setIndex((prev) => (prev + 1) % phrases.length);
    }, intervalMs);
    return () => clearInterval(interval);
  }, [shouldReduceMotion, disableAnimation]);

  // On mobile with disabled animations, just show static text
  if (disableAnimation) {
    return (
      <span className="inline-block relative text-4xl sm:text-5xl md:text-7xl font-bold">
        <span className="inline-block text-primary">{phrases[0]}</span>
      </span>
    );
  }

  return (
    <span className="inline-block relative text-4xl sm:text-5xl md:text-7xl font-bold">
      <AnimatePresence mode="wait">
        <motion.span
          key={phrases[index]}
          initial={noAnimation ? { opacity: 0 } : { opacity: 0, y: 20 }}
          animate={noAnimation ? { opacity: 1 } : { opacity: 1, y: 0 }}
          exit={noAnimation ? { opacity: 0 } : { opacity: 0, y: -20 }}
          transition={{ duration: noAnimation ? 0.2 : 0.4, ease: "easeInOut" }}
          className="inline-block text-primary"
        >
          {phrases[index]}
        </motion.span>
      </AnimatePresence>
    </span>
  );
};

export default RotatingText;
