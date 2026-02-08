"use client";

import * as React from "react";
import Editor from "@monaco-editor/react";
import CopyButton from "@/components/CopyButton";

type PlaygroundEvent = {
  Message: string;
  Kind: "stdout" | "stderr";
  Delay?: number;
};

type CompileResponse = {
  errors?: string;
  vetErrors?: string;
  events?: PlaygroundEvent[];
};

const viemMonoFontFamily =
  'var(--font-jetbrains-mono), ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace';

function buildOutput(res: CompileResponse): string {
  const parts: string[] = [];

  if (res.errors?.trim()) parts.push(res.errors.trimEnd());
  if (res.vetErrors?.trim()) parts.push(res.vetErrors.trimEnd());

  const events = res.events ?? [];
  if (events.length > 0) {
    parts.push(
      events
        .map((e) => {
          if (!e?.Message) return "";
          return e.Message;
        })
        .join(""),
    );
  }

  return parts.filter(Boolean).join("\n");
}

export default function GoPlayground({
  title = "Try it",
  fileName = "main.go",
  initialCode,
  editorHeight = 360,
  terminalHeight = 250,
}: {
  title?: string;
  fileName?: string;
  initialCode: string;
  editorHeight?: number;
  terminalHeight?: number;
}) {
  const [code, setCode] = React.useState(initialCode);
  const [output, setOutput] = React.useState<string>("");
  const [running, setRunning] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);

  const run = React.useCallback(async () => {
    setRunning(true);
    setError(null);
    setOutput("");

    try {
      const res = await fetch("/api/go/compile", {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({ code }),
      });

      if (!res.ok) {
        const text = await res.text();
        throw new Error(text || `Compile failed (${res.status})`);
      }

      const json = (await res.json()) as CompileResponse;
      setOutput(buildOutput(json));
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to run code");
    } finally {
      setRunning(false);
    }
  }, [code]);

  const reset = React.useCallback(() => {
    setCode(initialCode);
    setOutput("");
    setError(null);
  }, [initialCode]);

  return (
    <div className="my-6 rounded-lg overflow-hidden border border-accent/20 bg-gray-6/50">
      <div className="flex items-center justify-between h-10 bg-dark-deep/60 border-b border-accent/10 px-2">
        <div className="flex items-stretch min-w-0">
          <div className="flex items-center gap-2 px-3 border-r border-accent/10 text-[0.75rem] text-gray-4">
            <span className="font-medium text-gray-2 truncate">{title}</span>
          </div>
          <div className="flex items-center px-3 text-[0.8125rem] font-medium text-gray-1 bg-gray-6/40 border-r border-accent/10">
            <span className="text-[0.75rem] text-gray-4 mr-2">Go</span>
            <span className="truncate">{fileName}</span>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <CopyButton text={code} />

          <button
            type="button"
            className="h-[30px] px-2.5 rounded-md text-[0.75rem] font-medium bg-gray-6/40 hover:bg-gray-6/60 text-gray-1 border border-gray-5/40 active:scale-[0.99] transition"
            onClick={reset}
            disabled={running}
          >
            Reset
          </button>

          <button
            type="button"
            className="h-[30px] px-3 rounded-md text-[0.75rem] font-medium bg-accent/30 hover:bg-accent/40 text-primary-foreground border border-accent/40 disabled:opacity-60 active:scale-[0.99] transition"
            onClick={run}
            disabled={running}
          >
            {running ? "Running…" : "Run"}
          </button>
        </div>
      </div>

      <div className="flex flex-col">
        <div className="border-b border-accent/10">
          <Editor
            height={editorHeight}
            language="go"
            value={code}
            onChange={(v) => setCode(v ?? "")}
            theme="viem-dark"
            beforeMount={(monaco) => {
              // Define a theme closer to the StackBlitz dark aesthetic.
              monaco.editor.defineTheme("viem-dark", {
                base: "vs-dark",
                inherit: true,
                rules: [],
                colors: {
                  "editor.background": "#0d1117",
                  "editor.foreground": "#e8eef4",
                  "editorLineNumber.foreground": "#3a4555",
                  "editorLineNumber.activeForeground": "#7bbbe6",
                  "editorGutter.background": "#0d1117",
                  "editorCursor.foreground": "#54b2f0",
                  "editor.selectionBackground": "#1a2942",
                  "editor.inactiveSelectionBackground": "#1a2942",
                  "editorIndentGuide.background1": "#252d3a",
                  "editorIndentGuide.activeBackground1": "#3a4555",
                  "editorWhitespace.foreground": "#252d3a",
                  "scrollbarSlider.background": "#252d3a66",
                  "scrollbarSlider.hoverBackground": "#3a455599",
                  "scrollbarSlider.activeBackground": "#3a4555cc",
                },
              });
            }}
            options={{
              minimap: { enabled: false },
              fontFamily: viemMonoFontFamily,
              fontSize: 13,
              lineHeight: 20,
              padding: { top: 10, bottom: 10 },
              renderLineHighlight: "none",
              scrollBeyondLastLine: false,
              smoothScrolling: true,
              wordWrap: "off",
              tabSize: 2,
              automaticLayout: true,
              bracketPairColorization: { enabled: true },
              guides: { indentation: true },
            }}
          />
        </div>

        <div className="bg-dark-deep/70">
          <div className="flex items-center justify-between px-3 py-2 border-b border-accent/10">
            <div className="flex items-center gap-2 text-[0.75rem] text-gray-4">
              <span className="font-medium text-gray-2">Terminal</span>
              <span className="text-gray-4">•</span>
              <span className="text-gray-4">Go playground</span>
            </div>
          </div>
          <pre
            className="m-0 px-4 py-3 overflow-auto text-[0.8125rem] leading-relaxed text-gray-2"
            style={{ fontFamily: viemMonoFontFamily, height: terminalHeight }}
          >
            {error ? `Error: ${error}` : output || " "}
          </pre>
        </div>
      </div>
    </div>
  );
}

