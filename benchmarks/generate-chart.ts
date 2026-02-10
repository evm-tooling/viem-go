/**
 * Benchmark Comparison Chart Generator
 *
 * Generates individual dark-themed SVG comparison charts for each Go vs TypeScript
 * benchmark test, inspired by the site's bench-isaddress.svg style.
 *
 * Each benchmark gets its own chart showing Go vs TS total time (iterations * ns/op).
 *
 * Usage:
 *   bun run generate-chart.ts --suite abi --iter 5000
 *   bun run generate-chart.ts --suite address --iter 500 --run 20260208-020510
 *   bun run generate-chart.ts --all --iter 5000
 */

import { readFileSync, existsSync, writeFileSync, mkdirSync, readdirSync, statSync } from 'fs'
import { join } from 'path'
import { parseArgs } from 'util'

// ============================================================================
// Types
// ============================================================================

interface GoBenchmark {
  name: string
  iterations: number
  nsPerOp: number
  bytesPerOp: number
  allocsPerOp: number
}

interface TSBenchmark {
  name: string
  hz: number
  nsPerOp: number
  samples: number
}

interface MatchedBenchmark {
  label: string
  slug: string
  goNsPerOp: number
  tsNsPerOp: number
  goIterations: number
  tsIterations: number
  goTotalNs: number
  tsTotalNs: number
  winner: 'go' | 'ts' | 'tie'
  speedup: number
}

// ============================================================================
// CLI
// ============================================================================

const { values: args } = parseArgs({
  options: {
    suite: { type: 'string' },
    iter: { type: 'string', default: '5000' },
    run: { type: 'string' },
    all: { type: 'boolean', default: false },
    help: { type: 'boolean', short: 'h' },
  },
  allowPositionals: true,
})

if (args.help) {
  console.log(`
Benchmark Comparison Chart Generator

Usage:
  bun run generate-chart.ts --suite <name> --iter <number>
  bun run generate-chart.ts --all --iter <number>

Options:
  --suite <name>    Test suite name (e.g. abi, address, hash, ens, event, signature, unit, call, multicall)
  --iter <number>   Iteration category (e.g. 1, 5, 10, 50, 500, 5000)
  --run <timestamp> Specific run timestamp (default: latest)
  --all             Generate charts for all suites in the iteration folder
  -h, --help        Show this help
`)
  process.exit(0)
}

if (!args.all && !args.suite) {
  console.error('Error: --suite <name> or --all is required')
  process.exit(1)
}

const iterCount = parseInt(args.iter || '5000', 10)

// ============================================================================
// Directory Resolution
// ============================================================================

const RESULTS_DIR = join(import.meta.dir, 'results')
const SINGLE_RUN_DIR = join(RESULTS_DIR, 'single-run')

function findLatestRun(): string {
  if (!existsSync(SINGLE_RUN_DIR)) {
    console.error(`No single-run results directory found at ${SINGLE_RUN_DIR}`)
    process.exit(1)
  }
  const runs = readdirSync(SINGLE_RUN_DIR)
    .filter(d => d.startsWith('run-') && statSync(join(SINGLE_RUN_DIR, d)).isDirectory())
    .sort()
  if (runs.length === 0) {
    console.error('No run directories found')
    process.exit(1)
  }
  return runs[runs.length - 1]
}

function resolveIterDir(runDir: string): string {
  const iterDir = join(runDir, `iter-${iterCount}`)
  if (!existsSync(iterDir)) {
    console.error(`Iteration directory not found: ${iterDir}`)
    process.exit(1)
  }
  return iterDir
}

function discoverSuites(iterDir: string): string[] {
  return readdirSync(iterDir)
    .filter(d => d !== '_overall' && statSync(join(iterDir, d)).isDirectory())
    .sort()
}

// ============================================================================
// Parsing
// ============================================================================

function parseGoMd(content: string): GoBenchmark[] {
  const results: GoBenchmark[] = []
  for (const line of content.split('\n')) {
    const match = line.match(
      /^(Benchmark\w+)(?:-\d+)?\s+(\d+)\s+([\d.]+)\s+ns\/op(?:\s+([\d.]+)\s+B\/op)?(?:\s+([\d.]+)\s+allocs\/op)?/
    )
    if (match) {
      results.push({
        name: match[1],
        iterations: parseInt(match[2], 10),
        nsPerOp: parseFloat(match[3]),
        bytesPerOp: match[4] ? parseFloat(match[4]) : 0,
        allocsPerOp: match[5] ? parseFloat(match[5]) : 0,
      })
    }
  }
  return results
}

function parseTsMd(content: string): TSBenchmark[] {
  const results: TSBenchmark[] = []
  for (const line of content.split('\n')) {
    const match = line.match(
      /[·✓]\s*(viem-ts:\s*\w+\s*\([^)]+\))\s+([\d,]+\.?\d*)\s+([\d.]+)\s+([\d.]+)\s+([\d.]+)/
    )
    if (match) {
      const name = match[1].trim()
      const hz = parseFloat(match[2].replace(/,/g, ''))

      // Extract samples count from end of line
      const samplesMatch = line.match(/(\d+)\s*(?:fastest|slowest)?\s*$/)
      const samples = samplesMatch ? parseInt(samplesMatch[1], 10) : 0

      results.push({
        name,
        hz,
        nsPerOp: hz > 0 ? 1_000_000_000 / hz : Infinity,
        samples,
      })
    }
  }
  return results
}

// ============================================================================
// Benchmark Name Mapping (from compare.ts)
// ============================================================================

const STATIC_MAPPINGS: Record<string, string> = {
  BenchmarkCall_Basic: 'viem-ts: call (basic)',
  BenchmarkCall_WithData: 'viem-ts: call (with data)',
  BenchmarkCall_WithAccount: 'viem-ts: call (with account)',
  BenchmarkCall_Decimals: 'viem-ts: call (decimals)',
  BenchmarkCall_Symbol: 'viem-ts: call (symbol)',
  BenchmarkCall_BalanceOfMultiple: 'viem-ts: call (balanceOf multiple)',
  BenchmarkMulticall_Basic: 'viem-ts: multicall (basic)',
  BenchmarkMulticall_WithArgs: 'viem-ts: multicall (with args)',
  BenchmarkMulticall_MultiContract: 'viem-ts: multicall (multi-contract)',
  BenchmarkMulticall_10Calls: 'viem-ts: multicall (10 calls)',
  BenchmarkMulticall_30Calls: 'viem-ts: multicall (30 calls)',
  BenchmarkMulticall_ChunkedParallel: 'viem-ts: multicall (chunked parallel)',
  BenchmarkMulticall_Deployless: 'viem-ts: multicall (deployless)',
  BenchmarkMulticall_TokenMetadata: 'viem-ts: multicall (token metadata)',
  BenchmarkMulticall_100Calls: 'viem-ts: multicall (100 calls)',
  BenchmarkMulticall_50Calls: 'viem-ts: multicall (50 calls)',
  BenchmarkMulticall_100Calls: 'viem-ts: multicall (100 calls)',
  BenchmarkMulticall_200Calls: 'viem-ts: multicall (200 calls)',
  BenchmarkMulticall_500Calls: 'viem-ts: multicall (500 calls)',
  BenchmarkMulticall_MixedContracts_100: 'viem-ts: multicall (100 mixed contracts)',
  BenchmarkMulticall_1000Calls: 'viem-ts: multicall (1000 calls)',
  BenchmarkMulticall_10000Calls_SingleRPC: 'viem-ts: multicall (10000 calls single RPC)',
  BenchmarkMulticall_10000Calls_Chunked: 'viem-ts: multicall (10000 calls chunked)',
  BenchmarkMulticall_10000Calls_AggressiveChunking: 'viem-ts: multicall (10000 calls aggressive)',
  BenchmarkAbi_EncodeSimple: 'viem-ts: abi (encode simple)',
  BenchmarkAbi_EncodeComplex: 'viem-ts: abi (encode complex)',
  BenchmarkAbi_EncodeMultiArg: 'viem-ts: abi (encode multi-arg)',
  BenchmarkAbi_DecodeResult: 'viem-ts: abi (decode result)',
  BenchmarkAbi_EncodePacked: 'viem-ts: abi (encodePacked)',
  BenchmarkAbi_EncodePackedMulti: 'viem-ts: abi (encodePacked multi)',
  BenchmarkHash_Keccak256Short: 'viem-ts: hash (keccak256 short)',
  BenchmarkHash_Keccak256Long: 'viem-ts: hash (keccak256 long)',
  BenchmarkHash_Keccak256Hex: 'viem-ts: hash (keccak256 hex)',
  BenchmarkHash_Sha256Short: 'viem-ts: hash (sha256 short)',
  BenchmarkHash_Sha256Long: 'viem-ts: hash (sha256 long)',
  BenchmarkHash_FunctionSelector: 'viem-ts: hash (function selector)',
  BenchmarkHash_EventSelector: 'viem-ts: hash (event selector)',
  BenchmarkSignature_HashMessage: 'viem-ts: signature (hashMessage)',
  BenchmarkSignature_HashMessageLong: 'viem-ts: signature (hashMessage long)',
  BenchmarkSignature_RecoverAddress: 'viem-ts: signature (recoverAddress)',
  BenchmarkSignature_VerifyMessage: 'viem-ts: signature (verifyMessage)',
  BenchmarkSignature_ParseSignature: 'viem-ts: signature (parseSignature)',
  BenchmarkUnit_ParseEther: 'viem-ts: unit (parseEther)',
  BenchmarkUnit_ParseEtherLarge: 'viem-ts: unit (parseEther large)',
  BenchmarkUnit_FormatEther: 'viem-ts: unit (formatEther)',
  BenchmarkUnit_ParseUnits6: 'viem-ts: unit (parseUnits 6)',
  BenchmarkUnit_ParseGwei: 'viem-ts: unit (parseGwei)',
  BenchmarkUnit_FormatUnits: 'viem-ts: unit (formatUnits)',
  BenchmarkAddress_IsAddress: 'viem-ts: address (isAddress)',
  BenchmarkAddress_IsAddressLower: 'viem-ts: address (isAddress lower)',
  BenchmarkAddress_Checksum: 'viem-ts: address (checksum)',
  BenchmarkAddress_Create: 'viem-ts: address (create)',
  BenchmarkAddress_Create2: 'viem-ts: address (create2)',
  BenchmarkEvent_DecodeTransfer: 'viem-ts: event (decode transfer)',
  BenchmarkEvent_DecodeBatch10: 'viem-ts: event (decode batch 10)',
  BenchmarkEvent_DecodeBatch100: 'viem-ts: event (decode batch 100)',
  BenchmarkEns_Namehash: 'viem-ts: ens (namehash)',
  BenchmarkEns_NamehashDeep: 'viem-ts: ens (namehash deep)',
  BenchmarkEns_Labelhash: 'viem-ts: ens (labelhash)',
  BenchmarkEns_Normalize: 'viem-ts: ens (normalize)',
  BenchmarkEns_NormalizeLong: 'viem-ts: ens (normalize long)',
}

function matchBenchmarks(goResults: GoBenchmark[], tsResults: TSBenchmark[]): MatchedBenchmark[] {
  const matched: MatchedBenchmark[] = []

  for (const goBench of goResults) {
    let tsName = STATIC_MAPPINGS[goBench.name]

    // Auto-match fallback
    if (!tsName) {
      const m = goBench.name.match(/^Benchmark(\w+)_(\w+)$/)
      if (m) {
        const suite = m[1].toLowerCase()
        const variant = m[2].replace(/([A-Z])/g, ' $1').trim().toLowerCase()
        const tsBench = tsResults.find(
          ts => ts.name.toLowerCase().includes(suite) && ts.name.toLowerCase().includes(variant.split(' ')[0])
        )
        if (tsBench) tsName = tsBench.name
      }
    }

    if (!tsName) continue
    const tsBench = tsResults.find(ts => ts.name === tsName)
    if (!tsBench) continue

    // Total time = iterations * ns/op
    const goTotalNs = goBench.iterations * goBench.nsPerOp
    const tsTotalNs = tsBench.samples * tsBench.nsPerOp

    const ratio = goTotalNs / tsTotalNs
    let winner: 'go' | 'ts' | 'tie'
    if (Math.abs(ratio - 1) < 0.05) winner = 'tie'
    else winner = ratio > 1 ? 'ts' : 'go'

    const speedup = ratio > 1 ? ratio : 1 / ratio

    const label = formatLabel(goBench.name)
    const slug = label
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, '-')
      .replace(/^-|-$/g, '')

    matched.push({
      label,
      slug,
      goNsPerOp: goBench.nsPerOp,
      tsNsPerOp: tsBench.nsPerOp,
      goIterations: goBench.iterations,
      tsIterations: tsBench.samples,
      goTotalNs,
      tsTotalNs,
      winner,
      speedup,
    })
  }

  return matched
}

// ============================================================================
// Label Formatting (from compare.ts)
// ============================================================================

function formatLabel(benchmark: string): string {
  const firstUnderscore = benchmark.indexOf('_')
  if (firstUnderscore === -1) return benchmark.replace('Benchmark', '')
  return benchmark
    .slice(firstUnderscore + 1)
    .replace(/10000Calls_?/g, '10K ')
    .replace(/1000Calls/g, '1K Calls')
    .replace(/(\d+)Calls/g, '$1 Calls')
    .replace(/MixedContracts_(\d+)/g, '$1 Mixed')
    .replace(/MultiContract/g, 'Multi-Contract')
    .replace(/BalanceOfMultiple/g, 'BalanceOf Multi')
    .replace(/WithData/g, 'With Data')
    .replace(/WithAccount/g, 'With Account')
    .replace(/WithArgs/g, 'With Args')
    .replace(/TokenMetadata/g, 'Token Metadata')
    .replace(/AggressiveChunking/g, 'Aggressive')
    .replace(/SingleRPC/g, 'Single RPC')
    .replace(/ChunkedParallel/g, 'Chunked')
    .replace(/Keccak256/g, 'Keccak256 ')
    .replace(/Sha256/g, 'SHA-256 ')
    .replace(/EncodePacked/g, 'EncodePacked ')
    .replace(/EncodeSimple/g, 'Encode Simple')
    .replace(/EncodeComplex/g, 'Encode Complex')
    .replace(/EncodeMultiArg/g, 'Encode 3-Arg')
    .replace(/DecodeResult/g, 'Decode Result')
    .replace(/HashMessage/g, 'HashMessage ')
    .replace(/RecoverAddress/g, 'Recover Address')
    .replace(/VerifyMessage/g, 'Verify Message')
    .replace(/ParseSignature/g, 'Parse Signature')
    .replace(/ParseEther/g, 'ParseEther ')
    .replace(/ParseUnits6/g, 'ParseUnits(6)')
    .replace(/ParseGwei/g, 'ParseGwei')
    .replace(/FormatEther/g, 'FormatEther')
    .replace(/FormatUnits/g, 'FormatUnits')
    .replace(/IsAddress/g, 'IsAddress ')
    .replace(/Checksum/g, 'Checksum')
    .replace(/Create2/g, 'CREATE2')
    .replace(/Create(?!2)/g, 'CREATE')
    .replace(/DecodeBatch/g, 'Decode Batch ')
    .replace(/DecodeTransfer/g, 'Decode Transfer')
    .replace(/Namehash/g, 'Namehash ')
    .replace(/NamehashDeep/g, 'Namehash Deep')
    .replace(/Labelhash/g, 'Labelhash')
    .replace(/Normalize/g, 'Normalize ')
    .replace(/NormalizeLong/g, 'Normalize Long')
    .replace(/FunctionSelector/g, 'Fn Selector')
    .replace(/EventSelector/g, 'Event Selector')
    .replace(/_/g, ' ')
    .replace(/\s+/g, ' ')
    .trim()
}

// ============================================================================
// Helpers
// ============================================================================

function escapeXml(s: string): string {
  return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;')
}

function formatTime(ns: number): string {
  if (ns >= 1_000_000_000) return `${(ns / 1_000_000_000).toFixed(2)}s`
  if (ns >= 1_000_000) return `${(ns / 1_000_000).toFixed(2)}ms`
  if (ns >= 1_000) return `${(ns / 1_000).toFixed(2)}µs`
  return `${ns.toFixed(1)}ns`
}

// ============================================================================
// SVG Chart Generation — one chart per benchmark (dark theme, bench-isaddress style)
// ============================================================================

function generateSingleChart(suite: string, bench: MatchedBenchmark): string {
  const BG = '#1E1E20'
  const GO_COLOR = '#00ADD8'
  const TS_COLOR = '#3178C6'
  const TEXT_COLOR = '#FFFFFF'
  const DIM_TEXT = 'rgba(255,255,255,0.55)'
  const BORDER_STROKE = 'white'
  const BORDER_OPACITY = '0.23'
  const BAR_BORDER = '#787878'
  const ORANGE_START = '#FF8A00'
  const ORANGE_END = '#FFBA07'

  // Layout (inspired by bench-isaddress.svg proportions)
  const svgW = 1200
  const rowH = 110        // height per bar row
  const rowGap = 20       // gap between rows
  const marginTop = 130   // space for title area
  const marginBottom = 70 // space for footer
  const barStartX = 260   // where bars begin (after accent bar + label area)
  const barMaxW = 780     // max bar width
  const rowRx = 14        // bar corner radius
  const barH = 80         // inner bar height
  const accentBarW = 16
  const accentBarX = 80

  const contentH = 2 * rowH + rowGap
  const svgH = marginTop + contentH + marginBottom

  // The winner bar fills to max; loser is proportional
  const maxTotal = Math.max(bench.goTotalNs, bench.tsTotalNs)
  const goBarW = Math.max((bench.goTotalNs / maxTotal) * barMaxW, 6)
  const tsBarW = Math.max((bench.tsTotalNs / maxTotal) * barMaxW, 6)

  const suiteTitle = suite.charAt(0).toUpperCase() + suite.slice(1)
  const speedupLabel = bench.winner === 'tie'
    ? 'Similar performance'
    : `${bench.winner === 'go' ? 'Go' : 'TypeScript'} ${bench.speedup >= 10 ? bench.speedup.toFixed(0) : bench.speedup.toFixed(1)}x faster`

  const goRow1Y = marginTop
  const tsRow1Y = marginTop + rowH + rowGap

  const out: string[] = []
  const w = (s: string) => out.push(s)

  w(`<svg xmlns="http://www.w3.org/2000/svg" width="${svgW}" height="${svgH}" viewBox="0 0 ${svgW} ${svgH}" fill="none">`)

  // Background
  w(`<rect width="${svgW}" height="${svgH}" fill="${BG}"/>`)

  // Defs
  w(`<defs>`)
  // Orange accent gradient
  w(`  <radialGradient id="accentGrad" cx="0.5" cy="0.5" r="0.8">`)
  w(`    <stop stop-color="${ORANGE_START}"/>`)
  w(`    <stop offset="1" stop-color="${ORANGE_END}"/>`)
  w(`  </radialGradient>`)
  // Go bar gradient
  w(`  <linearGradient id="goFill" x1="0" y1="0" x2="1" y2="0">`)
  w(`    <stop offset="0%" stop-color="${GO_COLOR}" stop-opacity="0.9"/>`)
  w(`    <stop offset="100%" stop-color="#0891B2" stop-opacity="0.7"/>`)
  w(`  </linearGradient>`)
  // TS bar gradient
  w(`  <linearGradient id="tsFill" x1="0" y1="0" x2="1" y2="0">`)
  w(`    <stop offset="0%" stop-color="${TS_COLOR}" stop-opacity="0.9"/>`)
  w(`    <stop offset="100%" stop-color="#1D4ED8" stop-opacity="0.7"/>`)
  w(`  </linearGradient>`)
  // Accent shadow filter
  w(`  <filter id="accentShadow" x="-80%" y="-20%" width="260%" height="150%">`)
  w(`    <feDropShadow dx="0" dy="4" stdDeviation="5.5" flood-color="${ORANGE_START}" flood-opacity="0.98"/>`)
  w(`  </filter>`)
  w(`</defs>`)

  // Title: benchmark label
  w(`<text x="${svgW / 2}" y="50" text-anchor="middle" font-family="system-ui, -apple-system, sans-serif" font-size="34" font-weight="700" fill="${TEXT_COLOR}">${escapeXml(bench.label)}</text>`)

  // Subtitle: suite + speedup
  w(`<text x="${svgW / 2}" y="82" text-anchor="middle" font-family="system-ui, -apple-system, sans-serif" font-size="17" fill="${DIM_TEXT}">${escapeXml(suiteTitle)} · ${escapeXml(speedupLabel)} · ${bench.goIterations.toLocaleString()} iterations</text>`)

  // Orange accent bar (spans both rows)
  w(`<rect x="${accentBarX}" y="${marginTop - 8}" width="${accentBarW}" height="${contentH + 16}" rx="8" fill="url(#accentGrad)" fill-opacity="0.96" filter="url(#accentShadow)"/>`)

  // --- Go Row ---
  // Outer border
  w(`<rect x="${barStartX - 10}" y="${goRow1Y}" width="${barMaxW + 80}" height="${rowH}" rx="30" stroke="${BORDER_STROKE}" stroke-opacity="${BORDER_OPACITY}" stroke-width="3" fill="none"/>`)
  // Bar fill
  w(`<rect x="${barStartX}" y="${goRow1Y + (rowH - barH) / 2}" width="${goBarW.toFixed(1)}" height="${barH}" rx="${rowRx}" fill="url(#goFill)" stroke="${BAR_BORDER}" stroke-width="2"/>`)
  // "Go" label (inside left of row)
  w(`<text x="${barStartX - 28}" y="${goRow1Y + rowH / 2 + 7}" text-anchor="end" font-family="system-ui, -apple-system, sans-serif" font-size="22" font-weight="700" fill="${GO_COLOR}">Go</text>`)
  // Time value (right of bar)
  w(`<text x="${(barStartX + goBarW + 16).toFixed(1)}" y="${goRow1Y + rowH / 2 + 7}" font-family="system-ui, -apple-system, sans-serif" font-size="20" font-weight="600" fill="${TEXT_COLOR}">${formatTime(bench.goTotalNs)}</text>`)

  // --- TS Row ---
  // Outer border
  w(`<rect x="${barStartX - 10}" y="${tsRow1Y}" width="${barMaxW + 80}" height="${rowH}" rx="30" stroke="${BORDER_STROKE}" stroke-opacity="${BORDER_OPACITY}" stroke-width="3" fill="none"/>`)
  // Bar fill
  w(`<rect x="${barStartX}" y="${tsRow1Y + (rowH - barH) / 2}" width="${tsBarW.toFixed(1)}" height="${barH}" rx="${rowRx}" fill="url(#tsFill)" stroke="${BAR_BORDER}" stroke-width="2"/>`)
  // "TS" label
  w(`<text x="${barStartX - 28}" y="${tsRow1Y + rowH / 2 + 7}" text-anchor="end" font-family="system-ui, -apple-system, sans-serif" font-size="22" font-weight="700" fill="${TS_COLOR}">TS</text>`)
  // Time value
  w(`<text x="${(barStartX + tsBarW + 16).toFixed(1)}" y="${tsRow1Y + rowH / 2 + 7}" font-family="system-ui, -apple-system, sans-serif" font-size="20" font-weight="600" fill="${TEXT_COLOR}">${formatTime(bench.tsTotalNs)}</text>`)

  // Footer: legend + note
  const footerY = svgH - 25
  w(`<text x="${svgW / 2}" y="${footerY}" text-anchor="middle" font-family="system-ui, -apple-system, sans-serif" font-size="13" fill="${DIM_TEXT}">total time over ${bench.goIterations.toLocaleString()} iterations · lower is better</text>`)

  w(`</svg>`)
  return out.join('\n')
}

// ============================================================================
// Main
// ============================================================================

function generateForSuite(suite: string, iterDir: string) {
  const suiteDir = join(iterDir, suite)
  const goFile = join(suiteDir, 'go.md')
  const tsFile = join(suiteDir, 'ts.md')

  if (!existsSync(goFile)) {
    console.warn(`  Skipping ${suite}: go.md not found at ${goFile}`)
    return
  }
  if (!existsSync(tsFile)) {
    console.warn(`  Skipping ${suite}: ts.md not found at ${tsFile}`)
    return
  }

  const goContent = readFileSync(goFile, 'utf-8')
  const tsContent = readFileSync(tsFile, 'utf-8')

  const goResults = parseGoMd(goContent)
  const tsResults = parseTsMd(tsContent)

  if (goResults.length === 0) {
    console.warn(`  Skipping ${suite}: no Go benchmarks parsed`)
    return
  }
  if (tsResults.length === 0) {
    console.warn(`  Skipping ${suite}: no TS benchmarks parsed`)
    return
  }

  const matched = matchBenchmarks(goResults, tsResults)
  if (matched.length === 0) {
    console.warn(`  Skipping ${suite}: no matching benchmarks found between Go and TS`)
    return
  }

  const outDir = join(RESULTS_DIR, 'bench-comparison-charts', `bench-${suite}`)
  mkdirSync(outDir, { recursive: true })

  let count = 0
  for (const bench of matched) {
    const svg = generateSingleChart(suite, bench)
    const outPath = join(outDir, `${bench.slug}.svg`)
    writeFileSync(outPath, svg, 'utf-8')
    count++
  }

  console.log(`  ✓ ${suite}: ${count} charts → ${outDir}/`)
}

// Resolve run directory
const runName = args.run ? `run-${args.run}` : findLatestRun()
const runDir = join(SINGLE_RUN_DIR, runName)
if (!existsSync(runDir)) {
  console.error(`Run directory not found: ${runDir}`)
  process.exit(1)
}

const iterDir = resolveIterDir(runDir)
console.log(`Using: ${runName} / iter-${iterCount}`)

if (args.all) {
  const suites = discoverSuites(iterDir)
  console.log(`Generating charts for ${suites.length} suites: ${suites.join(', ')}`)
  for (const suite of suites) {
    generateForSuite(suite, iterDir)
  }
} else {
  generateForSuite(args.suite!, iterDir)
}

console.log('Done.')
