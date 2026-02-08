import { NextResponse } from "next/server";

export const dynamic = "force-dynamic";

type PlaygroundEvent = {
  Message: string;
  Kind: "stdout" | "stderr";
  Delay?: number;
};

type PlaygroundCompileResponse = {
  Errors?: string;
  VetErrors?: string;
  Events?: PlaygroundEvent[];
};

export async function POST(req: Request) {
  try {
    const body = (await req.json()) as { code?: unknown };
    const code = typeof body.code === "string" ? body.code : "";

    if (!code.trim()) {
      return NextResponse.json(
        { errors: "No code provided.", vetErrors: "", events: [] },
        { status: 400, headers: { "cache-control": "no-store" } },
      );
    }

    // Keep this endpoint "cosmetic" and safe-ish.
    if (code.length > 50_000) {
      return NextResponse.json(
        { errors: "Code is too large.", vetErrors: "", events: [] },
        { status: 413, headers: { "cache-control": "no-store" } },
      );
    }

    const params = new URLSearchParams();
    params.set("version", "2");
    params.set("body", code);

    const upstream = await fetch("https://go.dev/_/compile", {
      method: "POST",
      headers: { "content-type": "application/x-www-form-urlencoded" },
      body: params.toString(),
      // Avoid caching user code at the edge.
      cache: "no-store",
    });

    if (!upstream.ok) {
      const text = await upstream.text();
      return NextResponse.json(
        { errors: text || `Upstream error (${upstream.status})`, vetErrors: "", events: [] },
        { status: 502, headers: { "cache-control": "no-store" } },
      );
    }

    const json = (await upstream.json()) as PlaygroundCompileResponse;

    return NextResponse.json(
      {
        errors: json.Errors ?? "",
        vetErrors: json.VetErrors ?? "",
        events: json.Events ?? [],
      },
      { headers: { "cache-control": "no-store" } },
    );
  } catch (e) {
    return NextResponse.json(
      { errors: e instanceof Error ? e.message : "Bad request.", vetErrors: "", events: [] },
      { status: 400, headers: { "cache-control": "no-store" } },
    );
  }
}

