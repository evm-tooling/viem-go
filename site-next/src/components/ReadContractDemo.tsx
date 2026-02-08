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
          <CopyButton text={defaultCode} />
        </div>
      </div>

      <div className="border-b border-accent/10">
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
          }}
          theme="viem-dark"
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
          }}
        />
      </div>

      <div className="bg-dark-deep/70">
        <div className="flex items-center justify-between px-3 py-2 border-b border-accent/10">
          <div className="flex items-center gap-2 text-[0.75rem] text-gray-4 min-w-0">
            <span className="font-medium text-gray-2">Terminal</span>
            <span className="text-gray-4">•</span>
            <span className="text-gray-4">docs runner</span>
            <span className="text-gray-4">•</span>
            <span className="text-gray-4">Command</span>
            <code className="text-gray-2 bg-gray-6/40 border border-gray-5/40 rounded px-1.5 py-0.5 font-normal truncate">
              viem-go read-contract
            </code>
          </div>
          <button
            type="button"
            className="h-[30px] px-3 rounded-md text-[0.75rem] font-medium bg-accent/30 hover:bg-accent/40 text-primary-foreground border border-accent/40 disabled:opacity-60 active:scale-[0.99] transition"
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
                  ? "text-emerald-300"
                  : l.kind === "command"
                    ? "text-gray-2"
                  : l.kind === "stderr"
                    ? "text-red-300"
                    : l.kind === "meta"
                      ? "text-gray-4"
                      : "text-gray-2";
              return (
                <div key={idx} className={cls}>
                  {l.text}
                </div>
              );
            })}
            {running ? (
              <div className="text-gray-2">
                <span className="animate-cursor-blink">▍</span>
              </div>
            ) : null}
          </div>
          {!running && exitCode != null ? (
            <div className="mt-2 text-gray-4 whitespace-pre">
              {exitCode === 0 ? "exit 0" : "exit 1"}
            </div>
          ) : null}
        </div>
      </div>
    </div>
  );
}

