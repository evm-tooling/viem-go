/**
 * Benchmark Comparison Script
 *
 * Parses Go and TypeScript benchmark results and generates comparison reports.
 * Supports both single benchmark and full comparison modes.
 *
 * Usage:
 *   bun run compare.ts                                    # Default: compare all
 *   bun run compare.ts --mode full                        # Full comparison with detailed report
 *   bun run compare.ts --bench call --go-results X --ts-results Y  # Single benchmark mode
 */

import { readFileSync, existsSync, writeFileSync, readdirSync } from 'fs'
import { join, basename } from 'path'
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
  suite: string // e.g., "call", "multicall"
}

interface TSBenchmark {
  name: string
  hz: number // operations per second
  mean: number // mean time in ms
  samples: number
  suite: string // e.g., "call", "multicall"
}

interface ComparisonResult {
  benchmark: string
  suite: string
  category: string
  goNsPerOp: number
  goOpsPerSec: number
  goBytesPerOp: number
  goAllocsPerOp: number
  tsNsPerOp: number
  tsOpsPerSec: number
  winner: 'go' | 'ts' | 'tie'
  ratio: number // go/ts ratio (>1 means TS is faster)
  speedup: number // how much faster the winner is (always >= 1)
  speedupStr: string // "Go 1.5x faster" or "TS 1.2x faster"
}

interface SuiteStats {
  suite: string
  totalBenchmarks: number
  goWins: number
  tsWins: number
  ties: number
  avgGoNsPerOp: number
  avgTsNsPerOp: number
  avgRatio: number
  avgSpeedup: number
  winner: 'go' | 'ts' | 'tie'
  summary: string
}

interface OverallStats {
  totalBenchmarks: number
  totalSuites: number
  goWins: number
  tsWins: number
  ties: number
  avgGoNsPerOp: number
  avgTsNsPerOp: number
  avgRatio: number
  overallWinner: 'go' | 'ts' | 'tie'
  overallSpeedup: number
  overallSummary: string
  suiteStats: SuiteStats[]
}

// ============================================================================
// CLI Argument Parsing
// ============================================================================

const { values: args } = parseArgs({
  options: {
    mode: { type: 'string', default: 'default' }, // 'default', 'full', 'single'
    bench: { type: 'string' }, // benchmark name for single mode
    'go-results': { type: 'string' }, // custom go results path
    'ts-results': { type: 'string' }, // custom ts results path
    help: { type: 'boolean', short: 'h' },
  },
  allowPositionals: true,
})

if (args.help) {
  console.log(`
Benchmark Comparison Tool

Usage:
  bun run compare.ts [options]

Options:
  --mode <mode>           Comparison mode: 'default', 'full', or 'single' (default: 'default')
  --bench <name>          Benchmark name for single mode (e.g., 'call', 'multicall')
  --go-results <path>     Path to Go results file (relative to benchmarks/)
  --ts-results <path>     Path to TypeScript results file (relative to benchmarks/)
  -h, --help              Show this help message

Examples:
  bun run compare.ts                          # Standard comparison
  bun run compare.ts --mode full              # Full detailed report
  bun run compare.ts --bench call \\
    --go-results results/call-go-results.txt \\
    --ts-results results/call-ts-results.txt  # Single benchmark comparison
`)
  process.exit(0)
}

// Determine mode from --bench flag
const mode = args.bench ? 'single' : (args.mode || 'default')

// ============================================================================
// Parsing Functions
// ============================================================================

// Extract suite name from Go benchmark name
function extractSuiteFromGoBench(name: string): string {
  // BenchmarkCall_Basic -> call, BenchmarkMulticall_Basic -> multicall
  const match = name.match(/^Benchmark([A-Z][a-z]+)/)
  return match ? match[1].toLowerCase() : 'unknown'
}

// Extract suite name from TS benchmark name  
function extractSuiteFromTSBench(name: string): string {
  // viem-ts: call (basic) -> call, viem-ts: multicall (basic) -> multicall
  const match = name.match(/viem-ts:\s*(\w+)/)
  return match ? match[1].toLowerCase() : 'unknown'
}

// Parse Go benchmark output
function parseGoResults(content: string): GoBenchmark[] {
  const results: GoBenchmark[] = []
  const lines = content.split('\n')

  for (const line of lines) {
    // Match: BenchmarkCall_Basic-10    1234    5678 ns/op    123 B/op    4 allocs/op
    const match = line.match(
      /^(Benchmark\w+)(?:-\d+)?\s+(\d+)\s+([\d.]+)\s+ns\/op(?:\s+([\d.]+)\s+B\/op)?(?:\s+([\d.]+)\s+allocs\/op)?/
    )

    if (match) {
      const name = match[1]
      results.push({
        name,
        iterations: parseInt(match[2], 10),
        nsPerOp: parseFloat(match[3]),
        bytesPerOp: match[4] ? parseFloat(match[4]) : 0,
        allocsPerOp: match[5] ? parseFloat(match[5]) : 0,
        suite: extractSuiteFromGoBench(name),
      })
    }
  }

  return results
}

// Parse TypeScript vitest benchmark output
function parseTSResults(content: string): TSBenchmark[] {
  const results: TSBenchmark[] = []
  const lines = content.split('\n')

  for (const line of lines) {
    // Match vitest bench table output with various formats
    // Format: · name    hz    min    max    mean    p75    p99    p995    p999    rme    samples
    // Examples:
    //   · viem-ts: call (basic)               3,592.85  0.1888  6.0242  ...
    //   · viem-ts: multicall (basic)          1,234.56  0.5123  2.1234  ...
    const match = line.match(
      /[·✓]\s*(viem-ts:\s*\w+\s*\([^)]+\))\s+([\d,]+\.?\d*)\s+([\d.]+)\s+([\d.]+)\s+([\d.]+)/
    )

    if (match) {
      const name = match[1].trim()
      // Remove commas from hz value before parsing
      const hz = parseFloat(match[2].replace(/,/g, ''))
      
      // Extract samples count from end of line if available
      const samplesMatch = line.match(/(\d+)\s*(?:fastest|slowest)?\s*$/)
      const samples = samplesMatch ? parseInt(samplesMatch[1], 10) : 0

      results.push({
        name,
        hz,
        mean: hz > 0 ? 1000 / hz : 0, // Convert to ms
        samples,
        suite: extractSuiteFromTSBench(name),
      })
      
      if (hz === 0) {
        console.warn(`Warning: Benchmark "${name}" has 0 hz (failed or no samples)`)
      }
    }
  }

  return results
}

// ============================================================================
// Benchmark Name Mapping
// ============================================================================

// Build dynamic mapping from Go to TypeScript benchmark names
function buildBenchmarkMapping(goResults: GoBenchmark[], tsResults: TSBenchmark[]): Map<string, string> {
  const mapping = new Map<string, string>()
  
  // Static mappings for known benchmarks
  const staticMappings: Record<string, string> = {
    // Call benchmarks
    BenchmarkCall_Basic: 'viem-ts: call (basic)',
    BenchmarkCall_WithData: 'viem-ts: call (with data)',
    BenchmarkCall_WithAccount: 'viem-ts: call (with account)',
    BenchmarkCall_Decimals: 'viem-ts: call (decimals)',
    BenchmarkCall_Symbol: 'viem-ts: call (symbol)',
    BenchmarkCall_BalanceOfMultiple: 'viem-ts: call (balanceOf multiple)',
    // Multicall benchmarks
    BenchmarkMulticall_Basic: 'viem-ts: multicall (basic)',
    BenchmarkMulticall_WithArgs: 'viem-ts: multicall (with args)',
    BenchmarkMulticall_MultiContract: 'viem-ts: multicall (multi-contract)',
    BenchmarkMulticall_10Calls: 'viem-ts: multicall (10 calls)',
    BenchmarkMulticall_30Calls: 'viem-ts: multicall (30 calls)',
    BenchmarkMulticall_ChunkedParallel: 'viem-ts: multicall (chunked parallel)',
    BenchmarkMulticall_Deployless: 'viem-ts: multicall (deployless)',
    BenchmarkMulticall_TokenMetadata: 'viem-ts: multicall (token metadata)',
    // Stress test benchmarks
    BenchmarkMulticall_100Calls: 'viem-ts: multicall (100 calls)',
    BenchmarkMulticall_500Calls: 'viem-ts: multicall (500 calls)',
    BenchmarkMulticall_MixedContracts_100: 'viem-ts: multicall (100 mixed contracts)',
    BenchmarkMulticall_1000Calls: 'viem-ts: multicall (1000 calls)',
    BenchmarkMulticall_10000Calls_SingleRPC: 'viem-ts: multicall (10000 calls single RPC)',
    BenchmarkMulticall_10000Calls_Chunked: 'viem-ts: multicall (10000 calls chunked)',
    BenchmarkMulticall_10000Calls_AggressiveChunking: 'viem-ts: multicall (10000 calls aggressive)',
  }
  
  // Apply static mappings
  for (const [goName, tsName] of Object.entries(staticMappings)) {
    mapping.set(goName, tsName)
  }
  
  // Try to auto-match remaining benchmarks by pattern
  for (const goBench of goResults) {
    if (!mapping.has(goBench.name)) {
      // Try to find a matching TS benchmark
      // BenchmarkFoo_Bar -> viem-ts: foo (bar)
      const match = goBench.name.match(/^Benchmark(\w+)_(\w+)$/)
      if (match) {
        const suite = match[1].toLowerCase()
        const variant = match[2]
          .replace(/([A-Z])/g, ' $1')
          .trim()
          .toLowerCase()
        const tsPattern = `viem-ts: ${suite} (${variant})`
        
        const tsBench = tsResults.find(ts => 
          ts.name.toLowerCase().includes(suite) && 
          ts.name.toLowerCase().includes(variant.split(' ')[0])
        )
        
        if (tsBench) {
          mapping.set(goBench.name, tsBench.name)
        }
      }
    }
  }
  
  return mapping
}

// ============================================================================
// Comparison Logic
// ============================================================================

// Categorize benchmark by name
function categorizeBenchmark(name: string): string {
  const lower = name.toLowerCase()
  if (lower.includes('basic')) return 'Basic Operations'
  if (lower.includes('data') || lower.includes('witharg')) return 'With Parameters'
  if (lower.includes('account')) return 'With Account'
  // Stress tests (10000, 1000, 500, 100 calls)
  if (lower.includes('10000') || lower.includes('1000call') || lower.includes('500call')) return 'Extreme Stress Tests'
  if (lower.includes('aggressive') || lower.includes('singlerpc')) return 'Extreme Stress Tests'
  if (lower.includes('multiple') || lower.includes('10call') || lower.includes('30call') || lower.includes('100call')) return 'Batch Operations'
  if (lower.includes('mixedcontract') || lower.includes('mixed')) return 'Multi-Contract'
  if (lower.includes('multicontract') || lower.includes('multi-contract')) return 'Multi-Contract'
  if (lower.includes('chunk') || lower.includes('parallel')) return 'Parallel Execution'
  if (lower.includes('deployless')) return 'Deployless'
  if (lower.includes('metadata')) return 'Metadata Queries'
  if (lower.includes('decimal') || lower.includes('symbol')) return 'Simple Reads'
  return 'Other'
}

// Compare results
function compareResults(
  goResults: GoBenchmark[],
  tsResults: TSBenchmark[]
): ComparisonResult[] {
  const mapping = buildBenchmarkMapping(goResults, tsResults)
  const comparisons: ComparisonResult[] = []

  for (const goBench of goResults) {
    const tsName = mapping.get(goBench.name)
    if (!tsName) continue
    
    const tsBench = tsResults.find((ts) => ts.name === tsName)
    if (!tsBench) continue

    const goOpsPerSec = 1_000_000_000 / goBench.nsPerOp
    const tsOpsPerSec = tsBench.hz
    const tsNsPerOp = tsOpsPerSec > 0 ? 1_000_000_000 / tsOpsPerSec : Infinity

    const ratio = goBench.nsPerOp / tsNsPerOp
    let winner: 'go' | 'ts' | 'tie'
    if (Math.abs(ratio - 1) < 0.05) {
      winner = 'tie'
    } else {
      winner = ratio > 1 ? 'ts' : 'go'
    }

    const speedup = ratio > 1 ? ratio : 1 / ratio
    let speedupStr: string
    if (winner === 'tie') {
      speedupStr = 'Similar'
    } else if (winner === 'go') {
      speedupStr = `Go ${speedup.toFixed(2)}x faster`
    } else {
      speedupStr = `TS ${speedup.toFixed(2)}x faster`
    }

    const benchName = goBench.name.replace('Benchmark', '')
    comparisons.push({
      benchmark: benchName,
      suite: goBench.suite,
      category: categorizeBenchmark(benchName),
      goNsPerOp: goBench.nsPerOp,
      goOpsPerSec,
      goBytesPerOp: goBench.bytesPerOp,
      goAllocsPerOp: goBench.allocsPerOp,
      tsNsPerOp,
      tsOpsPerSec,
      winner,
      ratio,
      speedup,
      speedupStr,
    })
  }

  return comparisons
}

// Calculate suite-level statistics
function calculateSuiteStats(comparisons: ComparisonResult[]): SuiteStats[] {
  const suites = [...new Set(comparisons.map(c => c.suite))]
  
  return suites.map(suite => {
    const suiteComparisons = comparisons.filter(c => c.suite === suite)
    const goWins = suiteComparisons.filter(c => c.winner === 'go').length
    const tsWins = suiteComparisons.filter(c => c.winner === 'ts').length
    const ties = suiteComparisons.filter(c => c.winner === 'tie').length
    
    const avgGoNsPerOp = suiteComparisons.reduce((sum, c) => sum + c.goNsPerOp, 0) / suiteComparisons.length
    const avgTsNsPerOp = suiteComparisons.reduce((sum, c) => sum + c.tsNsPerOp, 0) / suiteComparisons.length
    const avgRatio = avgGoNsPerOp / avgTsNsPerOp
    const avgSpeedup = avgRatio > 1 ? avgRatio : 1 / avgRatio
    
    let winner: 'go' | 'ts' | 'tie'
    if (Math.abs(avgRatio - 1) < 0.05) {
      winner = 'tie'
    } else {
      winner = avgRatio > 1 ? 'ts' : 'go'
    }
    
    let summary: string
    if (winner === 'tie') {
      summary = 'Similar performance'
    } else if (winner === 'go') {
      summary = `Go ${avgSpeedup.toFixed(2)}x faster`
    } else {
      summary = `TS ${avgSpeedup.toFixed(2)}x faster`
    }
    
    return {
      suite,
      totalBenchmarks: suiteComparisons.length,
      goWins,
      tsWins,
      ties,
      avgGoNsPerOp,
      avgTsNsPerOp,
      avgRatio,
      avgSpeedup,
      winner,
      summary,
    }
  })
}

// Calculate overall statistics
function calculateOverallStats(comparisons: ComparisonResult[]): OverallStats {
  const goWins = comparisons.filter((c) => c.winner === 'go').length
  const tsWins = comparisons.filter((c) => c.winner === 'ts').length
  const ties = comparisons.filter((c) => c.winner === 'tie').length

  const avgGoNsPerOp = comparisons.reduce((sum, c) => sum + c.goNsPerOp, 0) / comparisons.length
  const avgTsNsPerOp = comparisons.reduce((sum, c) => sum + c.tsNsPerOp, 0) / comparisons.length
  const avgRatio = avgGoNsPerOp / avgTsNsPerOp

  let overallWinner: 'go' | 'ts' | 'tie'
  if (Math.abs(avgRatio - 1) < 0.05) {
    overallWinner = 'tie'
  } else {
    overallWinner = avgRatio > 1 ? 'ts' : 'go'
  }

  const overallSpeedup = avgRatio > 1 ? avgRatio : 1 / avgRatio

  let overallSummary: string
  if (overallWinner === 'tie') {
    overallSummary = '🤝 Performance is similar between Go and TypeScript'
  } else if (overallWinner === 'go') {
    overallSummary = `🏆 Go is ${overallSpeedup.toFixed(2)}x faster overall`
  } else {
    overallSummary = `🏆 TypeScript is ${overallSpeedup.toFixed(2)}x faster overall`
  }

  const suiteStats = calculateSuiteStats(comparisons)

  return {
    totalBenchmarks: comparisons.length,
    totalSuites: suiteStats.length,
    goWins,
    tsWins,
    ties,
    avgGoNsPerOp,
    avgTsNsPerOp,
    avgRatio,
    overallWinner,
    overallSpeedup,
    overallSummary,
    suiteStats,
  }
}

// ============================================================================
// Formatting Functions
// ============================================================================

function formatNumber(n: number, decimals = 0): string {
  if (!isFinite(n)) return 'N/A'
  return n.toLocaleString('en-US', {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  })
}

function formatDuration(ns: number): string {
  if (ns < 1000) return `${formatNumber(ns, 1)} ns`
  if (ns < 1_000_000) return `${formatNumber(ns / 1000, 2)} µs`
  if (ns < 1_000_000_000) return `${formatNumber(ns / 1_000_000, 2)} ms`
  return `${formatNumber(ns / 1_000_000_000, 3)} s`
}

// ============================================================================
// Console Report Generation
// ============================================================================

function generateConsoleReport(comparisons: ComparisonResult[], stats: OverallStats): void {
  console.log('\n' + '='.repeat(100))
  console.log('  BENCHMARK COMPARISON: viem-go vs viem TypeScript')
  console.log('='.repeat(100) + '\n')

  // Overall Summary Box
  console.log('┌' + '─'.repeat(98) + '┐')
  console.log('│' + stats.overallSummary.padStart(50 + stats.overallSummary.length / 2).padEnd(98) + '│')
  console.log('└' + '─'.repeat(98) + '┘')
  console.log('')

  // Suite Summary (if multiple suites)
  if (stats.suiteStats.length > 1) {
    console.log('BY SUITE')
    console.log('─'.repeat(60))
    for (const suite of stats.suiteStats) {
      const icon = suite.winner === 'go' ? '🟢' : suite.winner === 'ts' ? '🔵' : '⚪'
      console.log(`  ${icon} ${suite.suite.padEnd(15)} ${suite.summary.padEnd(25)} (${suite.totalBenchmarks} benchmarks)`)
    }
    console.log('')
  }

  // Detailed Results Table
  console.log('DETAILED RESULTS')
  console.log('─'.repeat(100))
  console.log(
    '| Benchmark'.padEnd(32) +
      '| Go (ns/op)'.padEnd(14) +
      '| TS (ns/op)'.padEnd(14) +
      '| Go (ops/s)'.padEnd(12) +
      '| TS (ops/s)'.padEnd(12) +
      '| Result'.padEnd(20) + '|'
  )
  console.log('|' + '-'.repeat(31) + '|' + '-'.repeat(13) + '|' + '-'.repeat(13) + '|' + '-'.repeat(11) + '|' + '-'.repeat(11) + '|' + '-'.repeat(19) + '|')

  // Group by suite
  const suites = [...new Set(comparisons.map(c => c.suite))]
  for (const suite of suites) {
    const suiteComparisons = comparisons.filter(c => c.suite === suite)
    for (const c of suiteComparisons) {
      const icon = c.winner === 'go' ? '🟢' : c.winner === 'ts' ? '🔵' : '⚪'
      console.log(
        '| ' +
          c.benchmark.padEnd(30) +
          '| ' +
          formatNumber(c.goNsPerOp, 0).padEnd(12) +
          '| ' +
          formatNumber(c.tsNsPerOp, 0).padEnd(12) +
          '| ' +
          formatNumber(c.goOpsPerSec, 0).padEnd(10) +
          '| ' +
          formatNumber(c.tsOpsPerSec, 0).padEnd(10) +
          '| ' +
          `${icon} ${c.speedupStr}`.padEnd(18) +
          '|'
      )
    }
  }

  console.log('')

  // Summary Statistics
  console.log('SUMMARY')
  console.log('─'.repeat(60))
  console.log(`  Total benchmarks: ${stats.totalBenchmarks}`)
  console.log(`  Total suites:     ${stats.totalSuites}`)
  console.log(`  Go wins:          ${stats.goWins} (${((stats.goWins / stats.totalBenchmarks) * 100).toFixed(0)}%)`)
  console.log(`  TS wins:          ${stats.tsWins} (${((stats.tsWins / stats.totalBenchmarks) * 100).toFixed(0)}%)`)
  console.log(`  Ties:             ${stats.ties} (${((stats.ties / stats.totalBenchmarks) * 100).toFixed(0)}%)`)
  console.log('')
  console.log(`  Avg Go:           ${formatNumber(stats.avgGoNsPerOp, 0)} ns/op (${formatNumber(1_000_000_000 / stats.avgGoNsPerOp, 0)} ops/s)`)
  console.log(`  Avg TS:           ${formatNumber(stats.avgTsNsPerOp, 0)} ns/op (${formatNumber(1_000_000_000 / stats.avgTsNsPerOp, 0)} ops/s)`)
  console.log('')

  // Legend
  console.log('LEGEND')
  console.log('─'.repeat(60))
  console.log('  🟢 Go faster  |  🔵 TS faster  |  ⚪ Similar (within 5%)')
  console.log('  ns/op = nanoseconds per operation (lower is better)')
  console.log('  ops/s = operations per second (higher is better)')
  console.log('')
}

// ============================================================================
// Markdown Report Generation
// ============================================================================

function generateMarkdownReport(comparisons: ComparisonResult[], stats: OverallStats): string {
  let md = '# Benchmark Comparison: viem-go vs viem TypeScript\n\n'
  md += `Generated: ${new Date().toISOString()}\n\n`

  // Overall Summary
  md += '## Overall Summary\n\n'
  if (stats.overallWinner === 'go') {
    md += `**🏆 Go is ${stats.overallSpeedup.toFixed(2)}x faster overall**\n\n`
  } else if (stats.overallWinner === 'ts') {
    md += `**🏆 TypeScript is ${stats.overallSpeedup.toFixed(2)}x faster overall**\n\n`
  } else {
    md += `**🤝 Performance is similar between Go and TypeScript**\n\n`
  }

  md += `| Metric | Go | TypeScript |\n`
  md += `|--------|----|-----------|\n`
  md += `| Avg ns/op | ${formatNumber(stats.avgGoNsPerOp, 0)} | ${formatNumber(stats.avgTsNsPerOp, 0)} |\n`
  md += `| Avg ops/s | ${formatNumber(1_000_000_000 / stats.avgGoNsPerOp, 0)} | ${formatNumber(1_000_000_000 / stats.avgTsNsPerOp, 0)} |\n`
  md += `| Wins | ${stats.goWins}/${stats.totalBenchmarks} | ${stats.tsWins}/${stats.totalBenchmarks} |\n\n`

  // Suite Summary (if multiple)
  if (stats.suiteStats.length > 1) {
    md += '## By Suite\n\n'
    md += '| Suite | Benchmarks | Go Wins | TS Wins | Ties | Winner |\n'
    md += '|-------|------------|---------|---------|------|--------|\n'
    for (const suite of stats.suiteStats) {
      const icon = suite.winner === 'go' ? '🟢' : suite.winner === 'ts' ? '🔵' : '⚪'
      md += `| ${suite.suite} | ${suite.totalBenchmarks} | ${suite.goWins} | ${suite.tsWins} | ${suite.ties} | ${icon} ${suite.summary} |\n`
    }
    md += '\n'
  }

  // Detailed Results
  md += '## Detailed Results\n\n'
  md += '| Benchmark | Go (ns/op) | TS (ns/op) | Go (ops/s) | TS (ops/s) | Result |\n'
  md += '|-----------|------------|------------|------------|------------|--------|\n'

  for (const c of comparisons) {
    const resultIcon = c.winner === 'go' ? '🟢' : c.winner === 'ts' ? '🔵' : '⚪'
    md += `| ${c.benchmark} | ${formatNumber(c.goNsPerOp, 0)} | ${formatNumber(c.tsNsPerOp, 0)} | ${formatNumber(c.goOpsPerSec, 0)} | ${formatNumber(c.tsOpsPerSec, 0)} | ${resultIcon} ${c.speedupStr} |\n`
  }

  // Win Summary
  md += '\n## Win Summary\n\n'
  md += `- 🟢 Go wins: ${stats.goWins} (${((stats.goWins / stats.totalBenchmarks) * 100).toFixed(0)}%)\n`
  md += `- 🔵 TS wins: ${stats.tsWins} (${((stats.tsWins / stats.totalBenchmarks) * 100).toFixed(0)}%)\n`
  md += `- ⚪ Ties: ${stats.ties} (${((stats.ties / stats.totalBenchmarks) * 100).toFixed(0)}%)\n`

  md += '\n## Notes\n\n'
  md += '- Benchmarks run against the same Anvil instance (mainnet fork) for fair comparison\n'
  md += '- ns/op = nanoseconds per operation (lower is better)\n'
  md += '- ops/s = operations per second (higher is better)\n'
  md += '- 🟢 = Go faster, 🔵 = TS faster, ⚪ = Similar (within 5%)\n'

  return md
}

// ============================================================================
// Full Report Generation (Enhanced)
// ============================================================================

function generateFullReport(comparisons: ComparisonResult[], stats: OverallStats): string {
  let md = '# Full Benchmark Report: viem-go vs viem TypeScript\n\n'
  md += `Generated: ${new Date().toISOString()}\n\n`
  md += '---\n\n'

  // Executive Summary
  md += '## Executive Summary\n\n'
  md += `This report compares **${stats.totalBenchmarks}** benchmarks across **${stats.totalSuites}** test suites.\n\n`
  
  if (stats.overallWinner === 'go') {
    md += `### 🏆 Winner: Go (viem-go)\n\n`
    md += `Go is **${stats.overallSpeedup.toFixed(2)}x faster** on average across all benchmarks.\n\n`
  } else if (stats.overallWinner === 'ts') {
    md += `### 🏆 Winner: TypeScript (viem)\n\n`
    md += `TypeScript is **${stats.overallSpeedup.toFixed(2)}x faster** on average across all benchmarks.\n\n`
  } else {
    md += `### 🤝 Result: Tie\n\n`
    md += `Performance is similar between Go and TypeScript (within 5% margin).\n\n`
  }

  // Quick Stats
  md += '### Quick Stats\n\n'
  md += '| Metric | Value |\n'
  md += '|--------|-------|\n'
  md += `| Total Benchmarks | ${stats.totalBenchmarks} |\n`
  md += `| Test Suites | ${stats.totalSuites} |\n`
  md += `| Go Wins | ${stats.goWins} (${((stats.goWins / stats.totalBenchmarks) * 100).toFixed(1)}%) |\n`
  md += `| TypeScript Wins | ${stats.tsWins} (${((stats.tsWins / stats.totalBenchmarks) * 100).toFixed(1)}%) |\n`
  md += `| Ties | ${stats.ties} (${((stats.ties / stats.totalBenchmarks) * 100).toFixed(1)}%) |\n`
  md += `| Avg Go Latency | ${formatDuration(stats.avgGoNsPerOp)} |\n`
  md += `| Avg TS Latency | ${formatDuration(stats.avgTsNsPerOp)} |\n`
  md += `| Go Throughput | ${formatNumber(1_000_000_000 / stats.avgGoNsPerOp, 0)} ops/s |\n`
  md += `| TS Throughput | ${formatNumber(1_000_000_000 / stats.avgTsNsPerOp, 0)} ops/s |\n\n`

  // Suite-by-Suite Analysis
  md += '---\n\n'
  md += '## Suite-by-Suite Analysis\n\n'
  
  for (const suite of stats.suiteStats) {
    const suiteComparisons = comparisons.filter(c => c.suite === suite.suite)
    const icon = suite.winner === 'go' ? '🟢' : suite.winner === 'ts' ? '🔵' : '⚪'
    
    md += `### ${suite.suite.charAt(0).toUpperCase() + suite.suite.slice(1)} Suite\n\n`
    md += `**Result:** ${icon} ${suite.summary}\n\n`
    
    md += '| Benchmark | Go | TS | Diff | Winner |\n'
    md += '|-----------|----|----|------|--------|\n'
    
    for (const c of suiteComparisons) {
      const resultIcon = c.winner === 'go' ? '🟢' : c.winner === 'ts' ? '🔵' : '⚪'
      const variant = c.benchmark.replace(`${suite.suite}_`, '').replace(suite.suite.charAt(0).toUpperCase() + suite.suite.slice(1) + '_', '')
      md += `| ${variant} | ${formatDuration(c.goNsPerOp)} | ${formatDuration(c.tsNsPerOp)} | ${c.speedup.toFixed(2)}x | ${resultIcon} |\n`
    }
    
    md += '\n'
    
    // Suite statistics
    md += `**Suite Statistics:**\n`
    md += `- Benchmarks: ${suite.totalBenchmarks}\n`
    md += `- Go wins: ${suite.goWins}, TS wins: ${suite.tsWins}, Ties: ${suite.ties}\n`
    md += `- Avg Go: ${formatDuration(suite.avgGoNsPerOp)} | Avg TS: ${formatDuration(suite.avgTsNsPerOp)}\n\n`
  }

  // Category Analysis
  md += '---\n\n'
  md += '## Category Analysis\n\n'
  
  const categories = [...new Set(comparisons.map(c => c.category))]
  for (const category of categories) {
    const catComparisons = comparisons.filter(c => c.category === category)
    if (catComparisons.length === 0) continue
    
    const catGoWins = catComparisons.filter(c => c.winner === 'go').length
    const catTsWins = catComparisons.filter(c => c.winner === 'ts').length
    const catTies = catComparisons.filter(c => c.winner === 'tie').length
    
    const avgCatGoNs = catComparisons.reduce((s, c) => s + c.goNsPerOp, 0) / catComparisons.length
    const avgCatTsNs = catComparisons.reduce((s, c) => s + c.tsNsPerOp, 0) / catComparisons.length
    const catRatio = avgCatGoNs / avgCatTsNs
    const catWinner = Math.abs(catRatio - 1) < 0.05 ? 'tie' : (catRatio > 1 ? 'ts' : 'go')
    const catSpeedup = catRatio > 1 ? catRatio : 1 / catRatio
    const catIcon = catWinner === 'go' ? '🟢' : catWinner === 'ts' ? '🔵' : '⚪'
    
    md += `### ${category}\n\n`
    md += `${catIcon} **${catWinner === 'tie' ? 'Similar' : catWinner === 'go' ? `Go ${catSpeedup.toFixed(2)}x faster` : `TS ${catSpeedup.toFixed(2)}x faster`}**\n\n`
    md += `Benchmarks: ${catComparisons.length} | Go wins: ${catGoWins} | TS wins: ${catTsWins} | Ties: ${catTies}\n\n`
  }

  // Memory Analysis (Go only)
  md += '---\n\n'
  md += '## Memory Analysis (Go)\n\n'
  md += '| Benchmark | Bytes/op | Allocs/op |\n'
  md += '|-----------|----------|----------|\n'
  
  for (const c of comparisons) {
    md += `| ${c.benchmark} | ${formatNumber(c.goBytesPerOp, 0)} | ${formatNumber(c.goAllocsPerOp, 0)} |\n`
  }
  
  md += '\n'

  // Detailed Raw Data
  md += '---\n\n'
  md += '## Detailed Raw Data\n\n'
  md += '| Benchmark | Suite | Go ns/op | TS ns/op | Go ops/s | TS ops/s | Ratio | Winner |\n'
  md += '|-----------|-------|----------|----------|----------|----------|-------|--------|\n'
  
  for (const c of comparisons) {
    const icon = c.winner === 'go' ? '🟢' : c.winner === 'ts' ? '🔵' : '⚪'
    md += `| ${c.benchmark} | ${c.suite} | ${formatNumber(c.goNsPerOp, 0)} | ${formatNumber(c.tsNsPerOp, 0)} | ${formatNumber(c.goOpsPerSec, 0)} | ${formatNumber(c.tsOpsPerSec, 0)} | ${c.ratio.toFixed(3)} | ${icon} |\n`
  }

  // Methodology
  md += '\n---\n\n'
  md += '## Methodology\n\n'
  md += '### Test Environment\n\n'
  md += '- **Network:** Anvil (Mainnet fork)\n'
  md += '- **Go Benchmark:** `go test -bench=. -benchmem -benchtime=10s -count=5`\n'
  md += '- **TS Benchmark:** `vitest bench` with 10s per benchmark\n\n'
  
  md += '### Measurement Notes\n\n'
  md += '- **ns/op:** Nanoseconds per operation (lower is better)\n'
  md += '- **ops/s:** Operations per second (higher is better)\n'
  md += '- **Ratio:** Go time / TS time (>1 means TS is faster)\n'
  md += '- **Tie:** Within 5% of each other\n\n'
  
  md += '### Caveats\n\n'
  md += '- Network latency dominates most benchmarks (RPC calls)\n'
  md += '- Results may vary based on network conditions\n'
  md += '- CPU-bound operations may show different characteristics\n'

  return md
}

// ============================================================================
// Main
// ============================================================================

async function main() {
  const resultsDir = join(import.meta.dir, 'results')
  
  // Determine result file paths
  let goResultsPath: string
  let tsResultsPath: string
  let outputBaseName: string
  
  if (mode === 'single' && args.bench) {
    // Single benchmark mode
    goResultsPath = join(import.meta.dir, args['go-results'] || `results/${args.bench}-go-results.txt`)
    tsResultsPath = join(import.meta.dir, args['ts-results'] || `results/${args.bench}-ts-results.txt`)
    outputBaseName = args.bench
    console.log(`\nSingle benchmark mode: ${args.bench}`)
  } else {
    // Full/default mode
    goResultsPath = join(resultsDir, 'go-results.txt')
    tsResultsPath = join(resultsDir, 'ts-results.txt')
    outputBaseName = mode === 'full' ? 'full-report' : 'comparison'
  }

  // Check files exist
  if (!existsSync(goResultsPath)) {
    console.error(`Error: Go results not found at ${goResultsPath}`)
    console.error('Run benchmarks first.')
    process.exit(1)
  }

  if (!existsSync(tsResultsPath)) {
    console.error(`Error: TypeScript results not found at ${tsResultsPath}`)
    console.error('Run benchmarks first.')
    process.exit(1)
  }

  // Parse results
  const goContent = readFileSync(goResultsPath, 'utf-8')
  const tsContent = readFileSync(tsResultsPath, 'utf-8')

  const goResults = parseGoResults(goContent)
  const tsResults = parseTSResults(tsContent)

  console.log(`Parsed ${goResults.length} Go benchmarks`)
  console.log(`Parsed ${tsResults.length} TypeScript benchmarks`)

  if (goResults.length === 0) {
    console.error('Error: No Go benchmark results found')
    process.exit(1)
  }

  if (tsResults.length === 0) {
    console.error('Error: No TypeScript benchmark results found')
    console.error('TypeScript results content (first 500 chars):')
    console.error(tsContent.slice(0, 500))
    process.exit(1)
  }

  // Compare
  const comparisons = compareResults(goResults, tsResults)

  if (comparisons.length === 0) {
    console.error('Error: No matching benchmarks found')
    console.error('Go benchmarks:', goResults.map((r) => r.name))
    console.error('TS benchmarks:', tsResults.map((r) => r.name))
    process.exit(1)
  }

  // Calculate statistics
  const stats = calculateOverallStats(comparisons)

  // Generate console report
  generateConsoleReport(comparisons, stats)

  // Generate and save markdown report
  let mdReport: string
  if (mode === 'full') {
    mdReport = generateFullReport(comparisons, stats)
  } else {
    mdReport = generateMarkdownReport(comparisons, stats)
  }
  
  const reportPath = join(resultsDir, `${outputBaseName}.md`)
  writeFileSync(reportPath, mdReport)
  console.log(`Markdown report saved to: ${reportPath}`)
  
  // In full mode, also generate the standard comparison.md
  if (mode === 'full') {
    const standardReport = generateMarkdownReport(comparisons, stats)
    writeFileSync(join(resultsDir, 'comparison.md'), standardReport)
    console.log(`Standard comparison saved to: ${join(resultsDir, 'comparison.md')}`)
  }
}

main().catch((err) => {
  console.error('Error:', err)
  process.exit(1)
})
