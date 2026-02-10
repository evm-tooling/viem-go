"use client";
import { useState, useEffect } from "react";
import { motion, AnimatePresence, useReducedMotion } from "framer-motion";

const phrases = [
  "Viem in Go",
  "20× Faster",
  "Type-Safe",
  "Modular",
];

const RotatingText = () => {
  const [index, setIndex] = useState(0);
  const shouldReduceMotion = useReducedMotion();

  useEffect(() => {
    // Slow down rotation interval for reduced motion users (or skip animation)
    const intervalMs = shouldReduceMotion ? 5000 : 3000;
    const interval = setInterval(() => {
      setIndex((prev) => (prev + 1) % phrases.length);
    }, intervalMs);
    return () => clearInterval(interval);
  }, [shouldReduceMotion]);

  return (
    <span className="inline-block relative text-4xl sm:text-5xl md:text-7xl font-bold">
      <AnimatePresence mode="wait">
        <motion.span
          key={phrases[index]}
          initial={shouldReduceMotion ? { opacity: 0 } : { opacity: 0, y: 20 }}
          animate={shouldReduceMotion ? { opacity: 1 } : { opacity: 1, y: 0 }}
          exit={shouldReduceMotion ? { opacity: 0 } : { opacity: 0, y: -20 }}
          transition={{ duration: shouldReduceMotion ? 0.2 : 0.4, ease: "easeInOut" }}
          className="inline-block text-primary"
        >
          {phrases[index]}
        </motion.span>
      </AnimatePresence>
    </span>
  );
};

export default RotatingText;
