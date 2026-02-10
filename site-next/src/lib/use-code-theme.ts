"use client";

import { useState, useEffect } from "react";
import type { PrismTheme } from "prism-react-renderer";

/** One Dark inspired — for dark mode */
const darkCodeTheme: PrismTheme = {
  plain: { color: "#abb2bf", backgroundColor: "transparent" },
  styles: [
    { types: ["comment", "prolog", "doctype", "cdata"], style: { color: "#5c6370", fontStyle: "italic" as const } },
    { types: ["keyword", "operator", "tag"], style: { color: "#c678dd" } },
    { types: ["property", "function"], style: { color: "#61afef" } },
    { types: ["string", "attr-value", "template-string"], style: { color: "#d19a66" } },
    { types: ["number", "boolean"], style: { color: "#d19a66" } },
    { types: ["builtin", "class-name", "maybe-class-name"], style: { color: "#e5c07b" } },
    { types: ["punctuation"], style: { color: "#abb2bf" } },
    { types: ["attr-name"], style: { color: "#d19a66" } },
    { types: ["char", "constant", "symbol"], style: { color: "#56b6c2" } },
    { types: ["variable"], style: { color: "#e06c75" } },
    { types: ["parameter"], style: { color: "#e06c75" } },
    { types: ["regex", "important"], style: { color: "#c678dd" } },
    { types: ["plain"], style: { color: "#abb2bf" } },
  ],
};

/** GitHub Light — exact colors from primer/github-syntax-light */
const lightCodeTheme: PrismTheme = {
  plain: { color: "#24292e", backgroundColor: "transparent" },
  styles: [
    { types: ["comment", "prolog", "doctype", "cdata"], style: { color: "#6a737d", fontStyle: "italic" as const } },
    { types: ["keyword", "operator"], style: { color: "#d73a49" } },
    { types: ["tag", "deleted"], style: { color: "#22863a" } },
    { types: ["function", "method"], style: { color: "#6f42c1" } },
    { types: ["string", "attr-value", "template-string", "template-punctuation"], style: { color: "#0450AE" } },
    { types: ["number", "boolean", "constant"], style: { color: "#005cc5" } },
    { types: ["builtin", "class-name", "maybe-class-name"], style: { color: "#6f42c1" } },
    { types: ["punctuation"], style: { color: "#24292e" } },
    { types: ["property"], style: { color: "#005cc5" } },
    { types: ["attr-name"], style: { color: "#e36209" } },
    { types: ["char", "symbol"], style: { color: "#005cc5" } },
    { types: ["variable"], style: { color: "#e36209" } },
    { types: ["regex", "important"], style: { color: "#032f62" } },
    { types: ["inserted"], style: { color: "#22863a" } },
    { types: ["changed"], style: { color: "#e36209" } },
    { types: ["parameter"], style: { color: "#e36209" } },
    { types: ["plain"], style: { color: "#24292e" } },
  ],
};

/**
 * Returns the correct Prism code theme based on the current
 * light/dark class on <html>. Reacts to changes in real-time.
 */
export function useCodeTheme(): PrismTheme {
  const [isLight, setIsLight] = useState(false);

  useEffect(() => {
    const html = document.documentElement;
    setIsLight(html.classList.contains("light"));

    const observer = new MutationObserver(() => {
      setIsLight(html.classList.contains("light"));
    });

    observer.observe(html, { attributes: true, attributeFilter: ["class"] });
    return () => observer.disconnect();
  }, []);

  return isLight ? lightCodeTheme : darkCodeTheme;
}

export { darkCodeTheme, lightCodeTheme };
