"use client";

import * as React from "react";
import { Highlight, themes } from "prism-react-renderer";
import { Tab, TabGroup, TabList, TabPanel, TabPanels } from "@headlessui/react";
import { ClipboardIcon, CheckIcon } from "@heroicons/react/24/outline";

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

function CopyButton({ code }: { code: string }) {
  const [copied, setCopied] = React.useState(false);

  const handleCopy = async () => {
    await navigator.clipboard.writeText(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <button
      onClick={handleCopy}
      aria-label={copied ? "Copied!" : "Copy code"}
      style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        gap: '0.25rem',
        padding: '0.25rem 0.5rem',
        borderRadius: '0.25rem',
        fontSize: '0.75rem',
        color: copied ? '#4ade80' : '#5a6778',
        background: 'transparent',
        border: 'none',
        cursor: 'pointer',
        transition: 'all 0.2s',
        height: '100%',
      }}
    >
      {copied ? (
        <>
          <CheckIcon style={{ width: '0.875rem', height: '0.875rem' }} />
          <span>Copied!</span>
        </>
      ) : (
        <>
          <ClipboardIcon style={{ width: '0.875rem', height: '0.875rem' }} />
          <span>Copy</span>
        </>
      )}
    </button>
  );
}

interface CodeBlockProps {
  code: string;
  language?: string;
  title?: string;
  showLineNumbers?: boolean;
}

export function CodeBlock({ 
  code, 
  language = "typescript", 
  title,
  showLineNumbers = true 
}: CodeBlockProps) {
  const trimmedCode = code.trim();
  
  return (
    <div style={{
      marginBottom: '1.5rem',
      borderRadius: '0.5rem',
      overflow: 'hidden',
      border: '1px solid rgba(57, 145, 205, 0.2)',
      background: 'rgba(37, 45, 58, 0.8)',
    }}>
      {/* Header */}
      <div style={{
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'space-between',
        padding: '0 0.75rem',
        height: '44px',
        background: 'rgba(23, 28, 36, 0.6)',
      }}>
        <span style={{
          fontFamily: 'var(--sl-font-mono)',
          fontSize: '0.8125rem',
          fontWeight: 500,
          color: '#ffffff',
          lineHeight: 1,
        }}>
          {title || languageNames[language] || language}
        </span>
        <CopyButton code={trimmedCode} />
      </div>
      
      {/* Code */}
      <Highlight theme={themes.nightOwl} code={trimmedCode} language={language}>
        {({ tokens, getLineProps, getTokenProps }) => (
          <pre style={{
            margin: 0,
            padding: '0.375rem 0.75rem 0.5rem 0.75rem',
            overflow: 'auto',
            fontSize: '0.8125rem',
            lineHeight: 1.6,
            background: 'transparent',
          }}>
            {tokens.map((line, i) => (
              <div key={i} {...getLineProps({ line })} style={{ display: 'table-row', margin: "0px" }}>
                {showLineNumbers && (
                  <span style={{
                    display: 'table-cell',
                    paddingRight: '0.75rem',
                    textAlign: 'right',
                    color: '#5a6778',
                    userSelect: 'none',
                    minWidth: '1.5rem',
                  }}>
                    {i + 1}
                  </span>
                )}
                <span style={{ display: 'table-cell' }}>
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
  
  // Process tabs
  const tabs = tabsInput.map((tab) => ({
    title: tab.title || languageNames[tab.language || ""] || tab.language || "Code",
    language: tab.language || "typescript",
    code: tab.code || "",
    showLineNumbers: tab.showLineNumbers ?? true,
  }));

  if (tabs.length === 0) return null;
  
  const containerStyle: React.CSSProperties = {
    margin: '1.5rem 0',
    borderRadius: '0.5rem',
    overflow: 'hidden',
    border: '1px solid rgba(57, 145, 205, 0.2)',
    background: 'rgba(37, 45, 58, 0.8)',
  };

  const headerStyle: React.CSSProperties = {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    height: '44px',
    background: 'rgba(23, 28, 36, 0.9)',
  };

  const preStyle: React.CSSProperties = {
    margin: 0,
    padding: '0.375rem 0.75rem 0.5rem 0.75rem',
    overflow: 'auto',
    fontSize: '0.8125rem',
    lineHeight: 1.6,
    background: 'transparent',
  };

  // Single code block - no tabs needed
  if (tabs.length === 1) {
    const tab = tabs[0]!;
    const codeStr = (tab.code || "").trim();
    return (
      <div style={containerStyle}>
        <div style={{ ...headerStyle, padding: '0 0.75rem' }}>
          <span style={{
            fontFamily: 'var(--sl-font-mono)',
            fontSize: '0.8125rem',
            fontWeight: 500,
            color: '#ffffff',
            lineHeight: 1,
          }}>
            {title || tab.title}
          </span>
          <CopyButton code={codeStr} />
        </div>
        <Highlight theme={themes.nightOwl} code={codeStr} language={tab.language}>
          {({ tokens, getLineProps, getTokenProps }) => (
            <pre style={preStyle}>
              {tokens.map((line, i) => (
                <div key={i} {...getLineProps({ line })} style={{ display: 'table-row', margin: "0px" }}>
                  {tab.showLineNumbers && (
                    <span style={{
                      display: 'table-cell',
                      paddingRight: '0.75rem',
                      textAlign: 'right',
                      color: '#5a6778',
                      userSelect: 'none',
                      minWidth: '1.5rem',
                    }}>
                      {i + 1}
                    </span>
                  )}
                  <span style={{ display: 'table-cell' }}>
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

  // Multiple code blocks - show tabs
  return (
    <div style={containerStyle}>
      <TabGroup selectedIndex={selectedIndex} onChange={setSelectedIndex}>
        {/* Header with tabs */}
        <div style={headerStyle}>
          <TabList style={{ display: 'flex', height: '100%', alignItems: 'stretch' }}>
            {tabs.map((tab, index) => (
              <Tab
                key={index}
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  justifyContent: 'center',
                  padding: '0 0.875rem',
                  fontSize: '0.8125rem',
                  fontWeight: 500,
                  color: selectedIndex === index ? '#ffffff' : '#8494a8',
                  background: selectedIndex === index ? 'rgba(255, 255, 255, 0.05)' : 'transparent',
                  border: 'none',
                  borderBottom: selectedIndex === index ? '2px solid #3991cd' : '2px solid transparent',
                  cursor: 'pointer',
                  transition: 'all 0.15s',
                  height: '44px',
                  marginTop: "0px",
                  boxSizing: 'border-box',
                }}
              >
                {tab.title}
              </Tab>
            ))}
          </TabList>
          <div style={{ display: 'flex', alignItems: 'center', margin: '0px', paddingRight: '0.5rem' }}>
            <CopyButton code={(tabs[selectedIndex]?.code || "").trim()} />
          </div>
        </div>
        
        {/* Code panels */}
        <TabPanels style={{ marginTop: "0px"}}>
          {tabs.map((tab, index) => {
            const codeStr = (tab.code || "").trim();
            return (
              <TabPanel key={index} >
                <Highlight theme={themes.nightOwl} code={codeStr} language={tab.language}>
                  {({ tokens, getLineProps, getTokenProps }) => (
                    <pre style={preStyle}>
                      {tokens.map((line, i) => (
                        <div key={i} {...getLineProps({ line })} style={{ display: 'table-row', margin: "0px" }}>
                          {tab.showLineNumbers && (
                            <span style={{
                              display: 'table-cell',
                              paddingRight: '0.75rem',
                              textAlign: 'right',
                              color: '#5a6778',
                              userSelect: 'none',
                              minWidth: '1.5rem',
                            }}>
                              {i + 1}
                            </span>
                          )}
                          <span style={{ display: 'table-cell' }}>
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

// For backwards compatibility
export const Code = CodeBlock;
export const Pre = CodeBlock;
