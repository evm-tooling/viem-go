import { NextResponse } from "next/server";

export const dynamic = "force-dynamic";

export async function POST(req: Request) {
  const runnerUrl = process.env.DOCS_RUNNER_URL || "http://localhost:3000";

  try {
    const body = (await req.json()) as {
      rpcUrl?: unknown;
      tokenAddress?: unknown;
      userAddress?: unknown;
    };

    const base =
      runnerUrl.endsWith("/") ? runnerUrl.slice(0, Math.max(0, runnerUrl.length - 1)) : runnerUrl;

    const res = await fetch(`${base}/run/read-contract`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({
        rpcUrl: typeof body.rpcUrl === "string" ? body.rpcUrl : undefined,
        tokenAddress: typeof body.tokenAddress === "string" ? body.tokenAddress : undefined,
        userAddress: typeof body.userAddress === "string" ? body.userAddress : undefined,
      }),
      cache: "no-store",
    });

    const text = await res.text();
    return new NextResponse(text, {
      status: res.status,
      headers: {
        "content-type": res.headers.get("content-type") ?? "application/json",
        "cache-control": "no-store",
      },
    });
  } catch (e) {
    return NextResponse.json(
      { error: e instanceof Error ? e.message : "Bad request." },
      { status: 400, headers: { "cache-control": "no-store" } },
    );
  }
}

