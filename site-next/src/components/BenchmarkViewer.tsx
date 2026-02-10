"use client";

import * as React from "react";
import { AnimatePresence, motion } from "framer-motion";
import { Highlight, type PrismTheme } from "prism-react-renderer";

// ============================================================================
// Code theme (matches CodePanel)
// ============================================================================

const codeTheme: PrismTheme = {
  plain: { color: "#abb2bf", backgroundColor: "transparent" },
  styles: [
    { types: ["comment", "prolog", "doctype", "cdata"], style: { color: "#5c6370", fontStyle: "italic" as const } },
    { types: ["keyword", "operator", "tag"], style: { color: "#c678dd" } },
    { types: ["property", "function"], style: { color: "#61afef" } },
    { types: ["string", "attr-value", "template-string"], style: { color: "#d19a66" } },
    { types: ["number", "boolean"], style: { color: "#d19a66" } },
    { types: ["builtin", "class-name", "maybe-class-name"], style: { color: "#e5c07b" } },
    { types: ["punctuation"], style: { color: "#abb2bf" } },
    { types: ["char", "constant", "symbol"], style: { color: "#56b6c2" } },
    { types: ["variable"], style: { color: "#e06c75" } },
  ],
};

const monoFont =
  'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace';

// ============================================================================
// Types
// ============================================================================

interface Slide {
  src: string;
  label: string;
  summary: string;
  code?: string;
}

interface Suite {
  id: string;
  label: string;
  badge?: string;
  slides: Slide[];
}

interface BenchmarkViewerProps {
  suites: Suite[];
}

// ============================================================================
// Inline SVG loader
// ============================================================================

function InlineSvg({ src, alt }: { src: string; alt: string }) {
  const [svg, setSvg] = React.useState<string | null>(null);

  React.useEffect(() => {
    let cancelled = false;
    fetch(src)
      .then((r) => r.text())
      .then((text) => {
        if (cancelled) return;
        const patched = text.replace(
          /(<rect\s+width="[^"]*"\s+height="[^"]*"\s+fill=")#1E1E20(")/,
          "$1transparent$2"
        );
        setSvg(patched);
      })
      .catch(() => {});
    return () => { cancelled = true; };
  }, [src]);

  if (!svg) {
    return (
      <div className="w-full aspect-[1200/440] rounded-lg bg-[hsl(var(--code-bg))] animate-pulse" />
    );
  }

  return (
    <div
      className="w-full rounded-lg bg-[hsl(var(--code-bg))] [&>svg]:w-full [&>svg]:h-auto [&>svg]:block"
      role="img"
      aria-label={alt}
      dangerouslySetInnerHTML={{ __html: svg }}
    />
  );
}

// ============================================================================
// Icons
// ============================================================================

function ChevronLeft() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
      <polyline points="15 18 9 12 15 6" />
    </svg>
  );
}

function ChevronRight() {
  return (
    <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
      <polyline points="9 18 15 12 9 6" />
    </svg>
  );
}

function ChevronDown() {
  return (
    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
      <polyline points="6 9 12 15 18 9" />
    </svg>
  );
}

// ============================================================================
// Slide animation variants
// ============================================================================

const slideVariants = {
  enter: (dir: number) => ({
    x: dir > 0 ? 250 : -250,
    opacity: 0,
  }),
  center: {
    x: 0,
    opacity: 1,
  },
  exit: (dir: number) => ({
    x: dir > 0 ? -250 : 250,
    opacity: 0,
  }),
};

const summaryVariants = {
  enter: { y: 12, opacity: 0 },
  center: { y: 0, opacity: 1 },
  exit: { y: -8, opacity: 0 },
};

// ============================================================================
// Main component
// ============================================================================

export default function BenchmarkViewer({ suites }: BenchmarkViewerProps) {
  const [activeSuiteIdx, setActiveSuiteIdx] = React.useState(0);
  const [activeSlideIdx, setActiveSlideIdx] = React.useState(0);
  const [direction, setDirection] = React.useState(0); // +1 = right, -1 = left
  const [dropdownOpen, setDropdownOpen] = React.useState(false);
  const dropdownRef = React.useRef<HTMLDivElement>(null);

  const suite = suites[activeSuiteIdx];
  const slide = suite.slides[activeSlideIdx];
  const canPrev = activeSlideIdx > 0;
  const canNext = activeSlideIdx < suite.slides.length - 1;

  // Close dropdown on outside click
  React.useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target as Node)) {
        setDropdownOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClick);
    return () => document.removeEventListener("mousedown", handleClick);
  }, []);

  function goTo(idx: number) {
    setDirection(idx > activeSlideIdx ? 1 : -1);
    setActiveSlideIdx(idx);
  }

  function goPrev() {
    if (!canPrev) return;
    setDirection(-1);
    setActiveSlideIdx((i) => i - 1);
  }

  function goNext() {
    if (!canNext) return;
    setDirection(1);
    setActiveSlideIdx((i) => i + 1);
  }

  function switchSuite(idx: number) {
    setActiveSuiteIdx(idx);
    setActiveSlideIdx(0);
    setDirection(1);
    setDropdownOpen(false);
  }

  // Keyboard navigation
  React.useEffect(() => {
    function handleKey(e: KeyboardEvent) {
      const viewer = dropdownRef.current?.closest("[data-benchmark-viewer]");
      if (!viewer) return;
      if (!viewer.contains(document.activeElement) && document.activeElement !== viewer) return;
      if (e.key === "ArrowLeft") { e.preventDefault(); goPrev(); }
      if (e.key === "ArrowRight") { e.preventDefault(); goNext(); }
    }
    document.addEventListener("keydown", handleKey);
    return () => document.removeEventListener("keydown", handleKey);
  });

  return (
    <div
      data-benchmark-viewer
      tabIndex={0}
      className="my-8 rounded-xl border border-[hsl(var(--border))] bg-[hsl(var(--background-secondary))] overflow-hidden focus:outline-none focus:ring-1 focus:ring-[hsl(var(--ring))]"
    >
      {/* ── Header bar ── */}
      <div className="flex items-center justify-between px-5 pt-4 pb-3 border-b border-[hsl(var(--border)/0.5)]">
        {/* Suite dropdown */}
        <div className="relative" ref={dropdownRef}>
          <button
            onClick={() => setDropdownOpen(!dropdownOpen)}
            className="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-[hsl(var(--background-tertiary))] border border-[hsl(var(--border))] text-sm font-semibold text-[hsl(var(--foreground))] hover:bg-[hsl(var(--background-elevated))] transition-colors"
          >
            <span>{suite.label}</span>
            {suite.badge && (
              <span className="text-xs font-medium text-[hsl(var(--tertiary))] opacity-80">
                {suite.badge}
              </span>
            )}
            <ChevronDown />
          </button>

          {dropdownOpen && (
            <div className="absolute left-0 top-full mt-1 z-20 min-w-[220px] rounded-lg border border-[hsl(var(--border))] bg-[hsl(var(--background-tertiary))] shadow-xl shadow-black/30 py-1">
              {suites.map((s, i) => (
                <button
                  key={s.id}
                  onClick={() => switchSuite(i)}
                  className={`w-full text-left px-3 py-2 text-sm flex items-center justify-between gap-3 transition-colors ${
                    i === activeSuiteIdx
                      ? "bg-[hsl(var(--primary)/0.12)] text-[hsl(var(--primary))]"
                      : "text-[hsl(var(--foreground-secondary))] hover:bg-[hsl(var(--background-elevated))] hover:text-[hsl(var(--foreground))]"
                  }`}
                >
                  <span className="font-medium">{s.label}</span>
                  {s.badge && (
                    <span className="text-xs text-[hsl(var(--foreground-muted))]">
                      {s.badge}
                    </span>
                  )}
                </button>
              ))}
            </div>
          )}
        </div>

        {/* Current slide info */}
        <div className="flex items-center gap-3">
          <span className="text-xs text-[hsl(var(--foreground-muted))]">
            {activeSlideIdx + 1} / {suite.slides.length}
          </span>
          <span className="text-sm font-medium text-[hsl(var(--tertiary))]">
            {slide?.label}
          </span>
        </div>
      </div>

      {/* ── Chart area ── */}
      <div className="relative group px-4 pt-2 pb-1">
        {/* Prev */}
        <button
          onClick={goPrev}
          disabled={!canPrev}
          aria-label="Previous benchmark"
          className="absolute left-2 top-1/2 -translate-y-1/2 z-10 flex h-9 w-9 items-center justify-center rounded-full bg-[hsl(var(--background-tertiary))] border border-[hsl(var(--border))] text-[hsl(var(--foreground-secondary))] opacity-0 group-hover:opacity-100 transition-opacity disabled:opacity-0 hover:bg-[hsl(var(--background-elevated))] hover:text-[hsl(var(--foreground))]"
        >
          <ChevronLeft />
        </button>

        {/* Next */}
        <button
          onClick={goNext}
          disabled={!canNext}
          aria-label="Next benchmark"
          className="absolute right-2 top-1/2 -translate-y-1/2 z-10 flex h-9 w-9 items-center justify-center rounded-full bg-[hsl(var(--background-tertiary))] border border-[hsl(var(--border))] text-[hsl(var(--foreground-secondary))] opacity-0 group-hover:opacity-100 transition-opacity disabled:opacity-0 hover:bg-[hsl(var(--background-elevated))] hover:text-[hsl(var(--foreground))]"
        >
          <ChevronRight />
        </button>

        {/* Animated slide */}
        <div className="overflow-hidden">
          <AnimatePresence mode="wait" custom={direction}>
            <motion.div
              key={`${activeSuiteIdx}-${activeSlideIdx}`}
              custom={direction}
              variants={slideVariants}
              initial="enter"
              animate="center"
              exit="exit"
              transition={{ duration: 0.2, ease: "easeOut" }}
            >
              <InlineSvg src={slide.src} alt={slide.label} />
            </motion.div>
          </AnimatePresence>
        </div>
      </div>

      {/* ── Dots ── */}
      <div className="flex items-center justify-center gap-1.5 pb-2">
        {suite.slides.map((_, i) => (
          <button
            key={i}
            onClick={() => goTo(i)}
            aria-label={`Go to ${suite.slides[i].label}`}
            className={`h-2 rounded-full transition-all duration-200 ${
              i === activeSlideIdx
                ? "w-6 bg-[hsl(var(--tertiary))]"
                : "w-2 bg-[hsl(var(--foreground-muted)/0.35)] hover:bg-[hsl(var(--foreground-muted)/0.6)]"
            }`}
          />
        ))}
      </div>

      {/* ── Summary box ── */}
      {slide?.summary && (
          <div
            key={`summary-${activeSuiteIdx}-${activeSlideIdx}`}
        
            className="mx-5 mb-4 mt-1 rounded-lg border border-[hsl(var(--border)/0.5)] bg-[hsl(var(--background-tertiary)/0.5)] px-4 py-3"
          >
            <div className="flex items-start gap-2">
              <svg
                className="mt-0.5 shrink-0 text-[hsl(var(--foreground-muted))]"
                width="14"
                height="14"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <circle cx="12" cy="12" r="10" />
                <line x1="12" y1="16" x2="12" y2="12" />
                <line x1="12" y1="8" x2="12.01" y2="8" />
              </svg>
              <p className="text-sm text-[hsl(var(--foreground-secondary))] leading-relaxed m-0">
                {slide.summary}
              </p>
            </div>

            {slide.code && (
              <div className="mt-3 rounded-md border border-[hsl(var(--code-border))] bg-[hsl(var(--code-bg))] overflow-hidden">
                <div className="flex items-center h-8 px-3 border-b border-[hsl(var(--code-border))] bg-[hsl(var(--code-bg)/0.5)]">
                  <span className="text-[0.6875rem] font-medium text-[hsl(var(--foreground-muted))]">Go</span>
                </div>
                <Highlight theme={codeTheme} code={slide.code.trim()} language="go">
                  {({ tokens, getLineProps, getTokenProps }) => (
                    <pre
                      style={{ fontFamily: monoFont }}
                      className="!m-0 !bg-transparent !py-2.5 !px-3 !border-0 overflow-auto text-[0.75rem] leading-relaxed"
                    >
                      {tokens.map((line, i) => (
                        <div key={i} {...getLineProps({ line })} className="table-row">
                          <span className="table-cell">
                            {line.map((token, key) => (
                              <span key={key} {...getTokenProps({ token })} />
                            ))}
                          </span>
                        </div>
                      ))}
                    </pre>
                  )}
                </Highlight>
              </div>
            )}
          </div>
      )}
    </div>
  );
}
