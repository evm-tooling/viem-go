"use client";

import * as React from "react";
import { Highlight } from "prism-react-renderer";
import { useCodeTheme } from "@/lib/use-code-theme";

const viemMonoFontFamily =
  'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace';

type CodeLineMeta = {
  focused?: boolean;
  highlighted?: boolean;
  added?: boolean;
  removed?: boolean;
};

type ParsedCode = {
  cleanCode: string;
  metaByLine: CodeLineMeta[];
  hasFocus: boolean;
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

    const directiveRegex = /(?:\/\/|#|--)\s*\[!code\s+([^\]]+)\]\s*/g;
    const directives: string[] = [];
    let m: RegExpExecArray | null;
    while ((m = directiveRegex.exec(rawLine)) !== null) {
      directives.push((m[1] ?? "").trim());
    }

    if (directives.length > 0) {
      line = line.replace(directiveRegex, "").replace(/[ \t]+$/g, "");

      for (const directive of directives) {
        if (directive.startsWith("focus")) {
          hasFocus = true;
          const [, countStr] = directive.split(":");
          const count = countStr ? Math.max(1, Number.parseInt(countStr, 10) || 1) : 1;
          for (let j = i; j < Math.min(lines.length, i + count); j++) {
            metaByLine[j]!.focused = true;
          }
          continue;
        }
        if (directive === "highlight") { metaByLine[i]!.highlighted = true; continue; }
        if (directive === "++") { metaByLine[i]!.added = true; continue; }
        if (directive === "--") { metaByLine[i]!.removed = true; continue; }
      }
    }

    cleanLines.push(line);
  }

  return { cleanCode: cleanLines.join("\n"), metaByLine, hasFocus };
}

function lineClassName(meta: CodeLineMeta | undefined, hasFocus: boolean, active: boolean) {
  const classes: string[] = ["table-row", "m-0"];

  if (hasFocus && !meta?.focused) classes.push("opacity-35 bg-code-bg-deep");
  if (active) classes.push("opacity-100 bg-code-bg-deep/40");
  if (hasFocus && meta?.focused) classes.push("opacity-100 bg-code-bg-deep/40");
  // if (meta?.focused) classes.push("font-semibold");
  // if (meta?.highlighted) classes.push("font-semibold");
  if (meta?.added) classes.push("diff-added");
  if (meta?.removed) classes.push("diff-removed");

  return classes.join(" ");
}

export function CodeGroup({ tabs: tabsInput, title }: CodeGroupProps) {
  const [selectedIndex, setSelectedIndex] = React.useState(0);
  const [active, setActive] = React.useState(false);
  const codeTheme = useCodeTheme();

  const tabs = tabsInput.map((tab) => ({
    title: tab.title || languageNames[tab.language || ""] || tab.language || "Code",
    language: tab.language || "typescript",
    code: tab.code || "",
    showLineNumbers: tab.showLineNumbers ?? true,
  }));

  if (tabs.length === 0) return null;

  return (
    <div className="group/code my-6 rounded-lg border-1 border-code-border transition-all duration-300 !hover:border-tertiary/30 hover:shadow-[0_0_20px_-1px_hsl(var(--tertiary)/0.15)]">
    <div className="rounded-lg overflow-hidden bg-code-bg/80">
      <TabGroup selectedIndex={selectedIndex} onChange={setSelectedIndex}>
        <div className="flex items-center justify-between h-11 bg-code-bg/50 border-b border-code-border">
          <TabList className="flex h-full items-stretch">
            {tabs.map((tab, index) => (
              <Tab
                key={index}
                className={`flex items-center justify-center px-3.5 text-[0.8125rem] font-medium cursor-pointer transition-all duration-150 h-11 border-b-2 outline-none ${
                  selectedIndex === index
                    ? "text-primary bg-code-bg !border-primary"
                    : "text-foreground-muted bg-transparent border-transparent hover:text-foreground-secondary hover:bg-code-bg/30"
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
                <Highlight theme={codeTheme} code={codeStr} language={tab.language}>
                  {({ tokens, getLineProps, getTokenProps }) => (
                    <pre
                      style={{ fontFamily: viemMonoFontFamily }}
                      className="!m-0 !bg-code-bg-deep/20 !pb-4 !pt-3 !px-0 !border-0 overflow-auto text-[0.8125rem] leading-relaxed"
                      onMouseEnter={() => { if (!active) setActive(true); }}
                      onMouseLeave={() => { if (active) setActive(false); }}
                    >
                      <code className="table w-full">
                        {tokens.map((line, i) => (
                          <div
                            key={i}
                            {...getLineProps({ line })}
                            className={`transition-all duration-300 ${lineClassName(parsed.metaByLine[i], parsed.hasFocus, active)}`}
                          >
                            {tab.showLineNumbers && (
                              <span className="table-cell pl-4 pr-3 text-right text-foreground-muted select-none min-w-6">
                                {i + 1}
                              </span>
                            )}
                            <span className={`table-cell w-full pr-4${tab.showLineNumbers ? "" : " pl-4"}`}>
                              {line.map((token, key) => (
                                <span key={key} {...getTokenProps({ token })} />
                              ))}
                            </span>
                          </div>
                        ))}
                      </code>
                    </pre>
                  )}
                </Highlight>
              </TabPanel>
            );
          })}
        </TabPanels>
      </TabGroup>
    </div>
    </div>
  );
}
