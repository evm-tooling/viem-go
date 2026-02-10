"use client";

import * as React from "react";
import Editor from "@monaco-editor/react";
import CopyButton from "@/components/CopyButton";

const viemMonoFontFamily =
  'var(--font-jetbrains-mono), ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace';

const defaultCode = `package main

import (
  "context"
  "fmt"
  "math/big"

  "github.com/ethereum/go-ethereum/common"

  "github.com/ChefBingbong/viem-go/chain/definitions"
  "github.com/ChefBingbong/viem-go/client"
  "github.com/ChefBingbong/viem-go/client/transport"
  "github.com/ChefBingbong/viem-go/contracts/erc20"
  "github.com/ChefBingbong/viem-go/utils/unit"
)

const ERC20_ADDRESS = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48" // USDC
const USER_ADDRESS  = "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045"

func main() {
  ctx := context.Background()

  publicClient, err := client.CreatePublicClient(client.PublicClientConfig{
    Chain:     &definitions.Mainnet,
    Transport: transport.HTTP("https://eth.merkle.io"),
  })
  if err != nil { panic(err) }
  defer publicClient.Close()

  balanceAny, err := publicClient.ReadContract(ctx, client.ReadContractOptions{
    Address:      common.HexToAddress(ERC20_ADDRESS),
    ABI:          erc20.ContractABI,
    FunctionName: "balanceOf",
    Args:         []any{common.HexToAddress(USER_ADDRESS)},
  })
  if err != nil { panic(err) }

  balance := balanceAny.(*big.Int)
  fmt.Println("Balance returned: ", balance.String())
  fmt.Println("Balance formatted:", unit.FormatUnits(balance, 6)) // USDC has 6 decimals
}`;

type RunnerResponse =
  | { rawBalance: string; formattedBalance: string; decimals: number }
  | { error: string };

type TerminalLine = {
  kind: "prompt" | "command" | "stdout" | "stderr" | "meta";
  text: string;
};

function buildCommand(rpcUrl: string, tokenAddress: string, userAddress: string) {
  return [
    "viem-go read-contract \\",
    `  --rpc ${rpcUrl} \\`,
    `  --token ${tokenAddress} \\`,
    `  --user ${userAddress}`,
  ].join("\n");
}

function shellPrompt() {
  return "viem-go@docs:~$";
}

export default function ReadContractDemo({
  title = "Try viem-go (readContract)",
  fileName = "main.go",
  initialTokenAddress = "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48",
  initialUserAddress = "0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
  initialRpcUrl = "https://eth.merkle.io",
  editorHeight = 360,
  terminalHeight = 150,
}: {
  title?: string;
  fileName?: string;
  initialTokenAddress?: string;
  initialUserAddress?: string;
  initialRpcUrl?: string;
  editorHeight?: number;
  terminalHeight?: number;
}) {
  const [tokenAddress, setTokenAddress] = React.useState(initialTokenAddress);
  const [userAddress, setUserAddress] = React.useState(initialUserAddress);
  const [rpcUrl, setRpcUrl] = React.useState(initialRpcUrl);
  const [lines, setLines] = React.useState<TerminalLine[]>([
    { kind: "meta", text: "Click Run to execute a real readContract against RPC." },
  ]);
  const [running, setRunning] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);
  const [exitCode, setExitCode] = React.useState<number | null>(null);

  // Track light/dark mode (same pattern as useCodeTheme)
  const [isLight, setIsLight] = React.useState(false);
  React.useEffect(() => {
    const html = document.documentElement;
    setIsLight(html.classList.contains("light"));
    const observer = new MutationObserver(() => {
      setIsLight(html.classList.contains("light"));
    });
    observer.observe(html, { attributes: true, attributeFilter: ["class"] });
    return () => observer.disconnect();
  }, []);

  const run = React.useCallback(async () => {
    setRunning(true);
    setError(null);
    const startedAt = performance.now();
    setExitCode(null);
    setLines([
      { kind: "prompt", text: shellPrompt() },
      { kind: "command", text: buildCommand(rpcUrl, tokenAddress, userAddress) },
      { kind: "meta", text: "running…" },
    ]);

    try {
      const res = await fetch("/api/demos/read-contract", {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({ rpcUrl, tokenAddress, userAddress }),
      });
      const json = (await res.json()) as RunnerResponse;
      if (!res.ok) {
        const msg = "error" in json ? json.error : `Runner failed (${res.status})`;
        throw new Error(msg);
      }

      if ("error" in json) throw new Error(json.error);

      const durationMs = Math.max(0, Math.round(performance.now() - startedAt));
      setExitCode(0);
      setLines([
        { kind: "prompt", text: shellPrompt() },
        { kind: "command", text: buildCommand(rpcUrl, tokenAddress, userAddress) },
        { kind: "stdout", text: `Balance returned:  ${json.rawBalance}` },
        { kind: "stdout", text: `Balance formatted: ${json.formattedBalance} (decimals=${json.decimals})` },
        { kind: "meta", text: `✓ done in ${durationMs}ms` },
      ]);
    } catch (e) {
      const durationMs = Math.max(0, Math.round(performance.now() - startedAt));
      const msg = e instanceof Error ? e.message : "Failed to run";
      setError(msg);
      setExitCode(1);
      setLines([
        { kind: "prompt", text: shellPrompt() },
        { kind: "command", text: buildCommand(rpcUrl, tokenAddress, userAddress) },
        { kind: "stderr", text: `✗ error: ${msg}` },
        { kind: "meta", text: `failed in ${durationMs}ms` },
      ]);
    } finally {
      setRunning(false);
    }
  }, [rpcUrl, tokenAddress, userAddress]);

  return (
    <div className="group/code my-6 rounded-lg overflow-hidden border-1 border-code-border bg-code-bg/80 transition-all duration-300 !hover:border-tertiary/30 hover:shadow-[0_0_20px_-1px_hsl(var(--tertiary)/0.15)]">
      <div className="flex items-center justify-between h-10 bg-code-bg/50 border-b border-code-border px-2">
        <div className="flex items-stretch min-w-0">
          <div className="flex items-center gap-2 px-3 border-r border-code-border text-[0.75rem] text-foreground-muted">
            <span className="font-medium text-foreground-secondary truncate">{title}</span>
          </div>
          <div className="flex items-center px-3 text-[0.8125rem] font-medium text-foreground bg-code-bg-deep/20 border-r border-code-border">
            <span className="text-[0.75rem] text-foreground-muted mr-2">Go</span>
            <span className="truncate">{fileName}</span>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <CopyButton text={defaultCode} />
        </div>
      </div>

      <div className="border-b border-code-border">
        <Editor
          height={editorHeight}
          language="go"
          value={defaultCode}
          options={{
            readOnly: true,
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
            scrollbar: { vertical: "hidden", horizontal: "hidden", handleMouseWheel: false },
            overviewRulerLanes: 0,
          }}
          theme={isLight ? "viem-light" : "viem-dark"}
          beforeMount={(monaco) => {
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
              },
            });
            monaco.editor.defineTheme("viem-light", {
              base: "vs",
              inherit: true,
              rules: [],
              colors: {
                "editor.background": "#eef1f6",
                "editor.foreground": "#24292e",
                "editorLineNumber.foreground": "#8b949e",
                "editorLineNumber.activeForeground": "#24292e",
                "editorGutter.background": "#eef1f6",
                "editorCursor.foreground": "#2563eb",
                "editor.selectionBackground": "#c8d7ea",
                "editor.inactiveSelectionBackground": "#dce5f0",
                "editorIndentGuide.background1": "#d0d7de",
                "editorIndentGuide.activeBackground1": "#b8c0c8",
                "editorWhitespace.foreground": "#d0d7de",
              },
            });
          }}
        />
      </div>

      <div className="bg-code-bg-deep/40">
        <div className="flex items-center justify-between px-3 py-2 border-b border-code-border">
          <div className="flex items-center gap-2 text-[0.75rem] text-foreground-muted min-w-0">
            <span className="font-medium text-foreground-secondary">Terminal</span>
            <span className="text-foreground-muted">•</span>
            <span className="text-foreground-muted">docs runner</span>
            <span className="text-foreground-muted">•</span>
            <span className="text-foreground-muted">Command</span>
            <code className="text-foreground-secondary bg-code-bg-deep/30 border border-code-border rounded px-1.5 py-0.5 font-normal truncate">
              viem-go read-contract
            </code>
          </div>
          <button
            type="button"
            className="h-[30px] px-3 rounded-md text-[0.75rem] font-medium bg-primary/15 hover:bg-primary/25 text-primary border border-primary/30 disabled:opacity-60 active:scale-[0.99] transition"
            onClick={run}
            disabled={running}
          >
            {running ? "Running…" : "Run"}
          </button>
        </div>
        <div
          className="m-0 px-4 py-3 overflow-auto text-[0.8125rem] leading-relaxed"
          style={{ fontFamily: viemMonoFontFamily, maxHeight: terminalHeight }}
        >
          <div className="whitespace-pre">
            {lines.map((l, idx) => {
              const cls =
                l.kind === "prompt"
                  ? "text-terminal-output"
                  : l.kind === "command"
                    ? "text-foreground-secondary"
                  : l.kind === "stderr"
                    ? "text-destructive"
                    : l.kind === "meta"
                      ? "text-foreground-muted"
                      : "text-foreground-secondary";
              return (
                <div key={idx} className={cls}>
                  {l.text}
                </div>
              );
            })}
            {running ? (
              <div className="text-foreground-secondary">
                <span className="animate-cursor-blink">▍</span>
              </div>
            ) : null}
          </div>
          {!running && exitCode != null ? (
            <div className="mt-2 text-foreground-muted whitespace-pre">
              {exitCode === 0 ? "exit 0" : "exit 1"}
            </div>
          ) : null}
        </div>
      </div>
    </div>
  );
}

