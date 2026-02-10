"use client";

import * as React from "react";
import useEmblaCarousel from "embla-carousel-react";

interface Slide {
  src: string;
  label: string;
}

interface BenchmarkSliderProps {
  title?: string;
  slides: Slide[];
}

export default function BenchmarkSlider({ title, slides }: BenchmarkSliderProps) {
  const [emblaRef, emblaApi] = useEmblaCarousel({
    loop: false,
    align: "center",
    containScroll: "trimSnaps",
  });

  const [selectedIndex, setSelectedIndex] = React.useState(0);
  const [canScrollPrev, setCanScrollPrev] = React.useState(false);
  const [canScrollNext, setCanScrollNext] = React.useState(false);

  const onSelect = React.useCallback(() => {
    if (!emblaApi) return;
    setSelectedIndex(emblaApi.selectedScrollSnap());
    setCanScrollPrev(emblaApi.canScrollPrev());
    setCanScrollNext(emblaApi.canScrollNext());
  }, [emblaApi]);

  React.useEffect(() => {
    if (!emblaApi) return;
    onSelect();
    emblaApi.on("select", onSelect);
    emblaApi.on("reInit", onSelect);
    return () => {
      emblaApi.off("select", onSelect);
      emblaApi.off("reInit", onSelect);
    };
  }, [emblaApi, onSelect]);

  const scrollPrev = React.useCallback(() => emblaApi?.scrollPrev(), [emblaApi]);
  const scrollNext = React.useCallback(() => emblaApi?.scrollNext(), [emblaApi]);
  const scrollTo = React.useCallback(
    (index: number) => emblaApi?.scrollTo(index),
    [emblaApi]
  );

  // Keyboard navigation
  React.useEffect(() => {
    function handleKey(e: KeyboardEvent) {
      if (!emblaApi) return;
      // Only handle if this component or its children are focused
      if (e.key === "ArrowLeft") {
        e.preventDefault();
        emblaApi.scrollPrev();
      } else if (e.key === "ArrowRight") {
        e.preventDefault();
        emblaApi.scrollNext();
      }
    }

    const node = emblaApi?.rootNode();
    if (!node) return;

    const container = node.closest("[data-benchmark-slider]");
    if (!container) return;

    container.addEventListener("keydown", handleKey as EventListener);
    return () =>
      container.removeEventListener("keydown", handleKey as EventListener);
  }, [emblaApi]);

  return (
    <div
      data-benchmark-slider
      tabIndex={0}
      className="relative my-8 rounded-xl border border-[hsl(var(--border))] bg-[hsl(var(--background-secondary))] overflow-hidden focus:outline-none focus:ring-1 focus:ring-[hsl(var(--ring))]"
    >
      {/* Header */}
      <div className="flex items-center justify-between px-5 pt-4 pb-2">
        <div className="flex items-center gap-3">
          {title && (
            <span className="text-sm font-semibold text-[hsl(var(--foreground))]">
              {title}
            </span>
          )}
          <span className="text-xs text-[hsl(var(--foreground-muted))]">
            {selectedIndex + 1} / {slides.length}
          </span>
        </div>
        <span className="text-sm font-medium text-[hsl(var(--tertiary))]">
          {slides[selectedIndex]?.label}
        </span>
      </div>

      {/* Carousel area */}
      <div className="relative group">
        {/* Prev button */}
        <button
          onClick={scrollPrev}
          disabled={!canScrollPrev}
          aria-label="Previous benchmark"
          className="absolute left-2 top-1/2 -translate-y-1/2 z-10 flex h-9 w-9 items-center justify-center rounded-full bg-[hsl(var(--background-tertiary))] border border-[hsl(var(--border))] text-[hsl(var(--foreground-secondary))] opacity-0 group-hover:opacity-100 transition-opacity disabled:opacity-0 hover:bg-[hsl(var(--background-elevated))] hover:text-[hsl(var(--foreground))]"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
            <polyline points="15 18 9 12 15 6" />
          </svg>
        </button>

        {/* Next button */}
        <button
          onClick={scrollNext}
          disabled={!canScrollNext}
          aria-label="Next benchmark"
          className="absolute right-2 top-1/2 -translate-y-1/2 z-10 flex h-9 w-9 items-center justify-center rounded-full bg-[hsl(var(--background-tertiary))] border border-[hsl(var(--border))] text-[hsl(var(--foreground-secondary))] opacity-0 group-hover:opacity-100 transition-opacity disabled:opacity-0 hover:bg-[hsl(var(--background-elevated))] hover:text-[hsl(var(--foreground))]"
        >
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
            <polyline points="9 18 15 12 9 6" />
          </svg>
        </button>

        {/* Embla viewport */}
        <div ref={emblaRef} className="overflow-hidden px-4">
          <div className="flex">
            {slides.map((slide, i) => (
              <div
                key={i}
                className="flex-[0_0_100%] min-w-0 px-1"
              >
                <img
                  src={slide.src}
                  alt={slide.label}
                  className="w-full h-auto rounded-lg"
                  draggable={false}
                />
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Dot indicators */}
      <div className="flex items-center justify-center gap-1.5 py-3">
        {slides.map((_, i) => (
          <button
            key={i}
            onClick={() => scrollTo(i)}
            aria-label={`Go to ${slides[i].label}`}
            className={`h-2 rounded-full transition-all duration-200 ${
              i === selectedIndex
                ? "w-6 bg-[hsl(var(--tertiary))]"
                : "w-2 bg-[hsl(var(--foreground-muted)/0.35)] hover:bg-[hsl(var(--foreground-muted)/0.6)]"
            }`}
          />
        ))}
      </div>
    </div>
  );
}
