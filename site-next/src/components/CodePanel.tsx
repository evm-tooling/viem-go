"use client";

import * as React from "react";
import { Highlight, type PrismTheme } from "prism-react-renderer";

const viemMonoFontFamily =
  'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace';

type CodeLineMeta = {
  /** Focused lines stay at full opacity; others are dimmed when any focus exists. */
  focused?: boolean;
  /** Highlighted lines get a subtle background. */
  highlighted?: boolean;
  /** Diff markers (optional). */
  added?: boolean;
  removed?: boolean;
};

type ParsedCode = {
  /** Code with all `[!code ...]` directives removed. */
  cleanCode: string;
  /** Per-line metadata (0-based index). */
  metaByLine: CodeLineMeta[];
  /** Whether any focus directives were found. */
  hasFocus: boolean;
};

/** Custom muted dark theme (similar to One Dark / viem docs) */
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
    { types: ["attr-name"], style: { color: "#d19a66" } },
    { types: ["char", "constant", "symbol"], style: { color: "#56b6c2" } },
    { types: ["variable"], style: { color: "#e06c75" } },
    { types: ["regex", "important"], style: { color: "#c678dd" } },
    { types: ["plain"], style: { color: "#abb2bf" } },
  ],
};
import {
  Tab,
  TabGroup,
  TabList,
  TabPanel,
  TabPanels,
} from "@headlessui/react";
import CopyButton from "./CopyButton";

const languageNames: Record<string, string> = {
  js: "JavaScript",
  jsx: "JavaScript",
  ts: "TypeScript",
  tsx: "TypeScript",
  javascript: "JavaScript",
  typescript: "TypeScript",
  go: "Go",
  bash: "TypeScript",
  shell: "TypeScript",
  json: "JSON",
  yaml: "YAML",
  css: "CSS",
  html: "HTML",
};

interface CodeTab {
  code: string;
  language?: string;
  title?: string;
  showLineNumbers?: boolean;
}

interface CodeGroupProps {
  tabs: CodeTab[];
  title?: string;
}

function parseCodeDirectives(code: string): ParsedCode {
  const lines = code.split("\n");
  const metaByLine: CodeLineMeta[] = Array.from({ length: lines.length }, () => ({}));
  const cleanLines: string[] = [];

  let hasFocus = false;

  for (let i = 0; i < lines.length; i++) {
    const rawLine = lines[i] ?? "";
    let line = rawLine;

    // Collect all directives on the line (supports multiple per line).
    // Matches patterns like:
    //   // [!code focus]
    //   // [!code focus:10]
    //   // [!code ++]
    //   #  [!code focus]
    const directiveRegex = /(?:\/\/|#|--)\s*\[!code\s+([^\]]+)\]\s*/g;
    const directives: string[] = [];
    let m: RegExpExecArray | null;
    while ((m = directiveRegex.exec(rawLine)) !== null) {
      directives.push((m[1] ?? "").trim());
    }

    if (directives.length > 0) {
      // Remove all directive fragments from the visible code.
      line = line.replace(directiveRegex, "").replace(/[ \t]+$/g, "");

      for (const directive of directives) {
        // focus or focus:n
        if (directive.startsWith("focus")) {
          hasFocus = true;
          const [, countStr] = directive.split(":");
          const count = countStr ? Math.max(1, Number.parseInt(countStr, 10) || 1) : 1;
          for (let j = i; j < Math.min(lines.length, i + count); j++) {
            metaByLine[j]!.focused = true;
          }
          continue;
        }

        if (directive === "highlight") {
          metaByLine[i]!.highlighted = true;
          continue;
        }

        // Diff-style markers (optional, but matches viem docs notation)
        if (directive === "++") {
          metaByLine[i]!.added = true;
          continue;
        }
        if (directive === "--") {
          metaByLine[i]!.removed = true;
          continue;
        }
      }
    }

    cleanLines.push(line);
  }

  return {
    cleanCode: cleanLines.join("\n"),
    metaByLine,
    hasFocus,
  };
}

function lineClassName(meta: CodeLineMeta | undefined, hasFocus: boolean, active: boolean) {
  const classes: string[] = ["table-row", "m-0"];

  if (hasFocus && !meta?.focused) classes.push("opacity-35");
  if (active ) classes.push("opacity-100");

  if (hasFocus && meta?.focused) classes.push("opacity-100");

  if (meta?.focused) classes.push("font-semibold");
  if (meta?.highlighted) classes.push("font-semibold");

  if (meta?.added) classes.push("bg-emerald-500/10");
  if (meta?.removed) classes.push("bg-red-500/10");

  return classes.join(" ");
}

export function CodeGroup({ tabs: tabsInput, title }: CodeGroupProps) {
  const [selectedIndex, setSelectedIndex] = React.useState(0);
  const [active, setActive] = React.useState(false);


  const tabs = tabsInput.map((tab) => ({
    title:
      tab.title ||
      languageNames[tab.language || ""] ||
      tab.language ||
      "Code",
    language: tab.language || "typescript",
    code: tab.code || "",
    showLineNumbers: tab.showLineNumbers ?? true,
  }));

  if (tabs.length === 0) return null;

  return (
    <div className="my-6 rounded-lg overflow-hidden border border-accent/20 bg-transparent ">
      <TabGroup selectedIndex={selectedIndex} onChange={setSelectedIndex}>
        <div className="flex items-center justify-between h-11 bg-[rgb(23,26,38,0.7)] border-b border-accent/10">
          <TabList className="flex h-full items-stretch">
            {tabs.map((tab, index) => (
              <Tab
                key={index}
                className={`flex items-center justify-center px-3.5 text-[0.8125rem] font-medium cursor-pointer transition-all duration-150 h-11 border-b-2 outline-none ${
                  selectedIndex === index
                    ? "text-white bg-white/5 border-accent"
                    : "text-gray-3 bg-transparent border-transparent hover:text-white"
                }`}
              >
                {tab.title}
              </Tab>
            ))}
          </TabList>
          <div className="flex items-center pr-2">
            <CopyButton
              text={parseCodeDirectives((tabs[selectedIndex]?.code || "").trim()).cleanCode.trim()}
            />
          </div>
        </div>

        <TabPanels className="!mt-0 !pt-0 !pb-0">
          {tabs.map((tab, index) => {
            const parsed = parseCodeDirectives((tab.code || "").trim());
            const codeStr = parsed.cleanCode.trim();
            return (
              <TabPanel key={index}>
                <Highlight
                  theme={codeTheme}
                  code={codeStr}
                  language={tab.language}
                >
                  {({ tokens, getLineProps, getTokenProps }) => (
                    <pre
                      style={{ fontFamily: viemMonoFontFamily }}
                      className="!m-0 !pb-4  !pt-3 !px-4 !border-0 overflow-auto text-[0.8125rem] leading-relaxed bg-transparent" 
                      onMouseEnter={() => {
                        if (active) return
                        setActive(true)
                      }} 
                      onMouseLeave={() =>  {
                        if (!active) return 
                        setActive(false)
                      }}>
                      {tokens.map((line, i) => (
                        <div
                          key={i}
                          {...getLineProps({ line })}
                          className={`transition-all duration-300 ${lineClassName(parsed.metaByLine[i], parsed.hasFocus, active)}`}
                        >
                          {tab.showLineNumbers && (
                            <span className="table-cell pr-3 text-right text-gray-4 select-none min-w-6">
                              {i + 1}
                            </span>
                          )}
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
              </TabPanel>
            );
          })}
        </TabPanels>
      </TabGroup>
    </div>
  );
}
