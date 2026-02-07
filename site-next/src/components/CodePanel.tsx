"use client";

import * as React from "react";
import { Highlight, type PrismTheme } from "prism-react-renderer";

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
  bash: "Terminal",
  shell: "Terminal",
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

export function CodeGroup({ tabs: tabsInput, title }: CodeGroupProps) {
  const [selectedIndex, setSelectedIndex] = React.useState(0);

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

  if (tabs.length === 1) {
    const tab = tabs[0]!;
    const codeStr = (tab.code || "").trim();
    return (
      <div className="my-6 rounded-lg overflow-hidden border border-accent/20 bg-transparent">
        <div className="flex items-center justify-between px-3 h-11 bg-gray-6 border-b border-accent/10">
          <span className="font-mono text-[0.8125rem] font-medium text-white leading-none">
            {title || tab.title}
          </span>
          <CopyButton text={codeStr} />
        </div>
        <Highlight
          theme={codeTheme}
          code={codeStr}
          language={tab.language}
        >
          {({ tokens, getLineProps, getTokenProps }) => (
            <pre className="!m-0 !py-2 !px-3 !border-0 overflow-auto text-[0.8125rem] leading-relaxed bg-transparent">
              {tokens.map((line, i) => (
                <div
                  key={i}
                  {...getLineProps({ line })}
                  className="table-row m-0"
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
      </div>
    );
  }

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
            <CopyButton text={(tabs[selectedIndex]?.code || "").trim()} />
          </div>
        </div>

        <TabPanels className="!mt-0 !pt-0 !pb-0">
          {tabs.map((tab, index) => {
            const codeStr = (tab.code || "").trim();
            return (
              <TabPanel key={index}>
                <Highlight
                  theme={codeTheme}
                  code={codeStr}
                  language={tab.language}
                >
                  {({ tokens, getLineProps, getTokenProps }) => (
                    <pre className="!m-0 !pb-4  !pt-3 !px-4 !border-0 overflow-auto text-[0.8125rem] leading-relaxed bg-transparent">
                      {tokens.map((line, i) => (
                        <div
                          key={i}
                          {...getLineProps({ line })}
                          className="table-row m-0"
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
