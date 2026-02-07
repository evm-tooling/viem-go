"use client";

import { useEffect, useRef, useState } from "react";
import type { TocEntry } from "@/lib/mdx";

export default function TableOfContents({
  headings,
}: {
  headings: TocEntry[];
}) {
  const [activeId, setActiveId] = useState<string>("");
  const observerRef = useRef<IntersectionObserver | null>(null);

  useEffect(() => {
    // Find all heading elements on the page
    const elements = headings
      .map((h) => document.getElementById(h.id))
      .filter(Boolean) as HTMLElement[];

    if (elements.length === 0) return;

    // Use IntersectionObserver to track which heading is in view
    const scrollRoot = document.querySelector("main") || null;
    observerRef.current = new IntersectionObserver(
      (entries) => {
        // Find the first heading that is intersecting
        for (const entry of entries) {
          if (entry.isIntersecting) {
            setActiveId(entry.target.id);
            break;
          }
        }
      },
      {
        root: scrollRoot,
        rootMargin: "-80px 0px -70% 0px",
        threshold: 0,
      }
    );

    for (const el of elements) {
      observerRef.current.observe(el);
    }

    return () => observerRef.current?.disconnect();
  }, [headings]);

  // Also track scroll position for more accurate active state
  // The scroll container is the <main> element (overflow-y-auto), not the window
  useEffect(() => {
    const scrollContainer =
      document.querySelector("main") || window;

    function onScroll() {
      const scrollTop =
        scrollContainer instanceof HTMLElement
          ? scrollContainer.scrollTop
          : window.scrollY;
      const offset = scrollTop + 120;
      let current = "";
      for (const h of headings) {
        const el = document.getElementById(h.id);
        if (el && el.offsetTop <= offset) {
          current = h.id;
        }
      }
      if (current) setActiveId(current);
    }
    scrollContainer.addEventListener("scroll", onScroll, { passive: true });
    // Set initial
    onScroll();
    return () => scrollContainer.removeEventListener("scroll", onScroll);
  }, [headings]);

  if (headings.length === 0) return null;

  return (
    <aside className="hidden xl:block w-[220px] shrink-0 sticky top-0 self-start h-[calc(100vh-3rem)] border-l border-accent/10">
      <div className="py-8 pl-4 pr-2 overflow-y-auto h-full">
        <p className="text-[0.8125rem] font-semibold text-gray-3 uppercase tracking-wider mb-2">
          On this page
        </p>
        <nav className="flex flex-col gap-0">
          {headings.map((heading) => {
            const isActive = activeId === heading.id;
            return (
              <a
                key={heading.id}
                href={`#${heading.id}`}
                onClick={(e) => {
                  e.preventDefault();
                  const el = document.getElementById(heading.id);
                  if (el) {
                    const mainEl = document.querySelector("main");
                    if (mainEl) {
                      mainEl.scrollTo({
                        top: el.offsetTop - 32,
                        behavior: "smooth",
                      });
                    } else {
                      el.scrollIntoView({ behavior: "smooth" });
                    }
                    setActiveId(heading.id);
                    history.pushState(null, "", `#${heading.id}`);
                  }
                }}
                className={`block text-[0.875rem] leading-snug no-underline py-1.5 transition-colors border-l-2 ${
                  heading.depth === 3 ? "pl-5" : "pl-3"
                } ${
                  isActive
                    ? "text-accent border-accent"
                    : "text-gray-4 border-transparent hover:text-gray-2 hover:border-gray-5"
                }`}
              >
                {heading.text}
              </a>
            );
          })}
        </nav>
      </div>
    </aside>
  );
}
