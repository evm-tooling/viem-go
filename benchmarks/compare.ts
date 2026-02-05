/**
 * Benchmark Comparison Script
 *
 * Parses Go and TypeScript benchmark results and generates a comparison report.
 *
 * Usage: bun run compare.ts
 */

import { readFileSync, existsSync, writeFileSync } from 'fs'
import { join } from 'path'

interface GoBenchmark {
  name: string
  iterations: number
  nsPerOp: number
  bytesPerOp: number
  allocsPerOp: number
}

interface TSBenchmark {
  name: string
  hz: number // operations per second
  mean: number // mean time in ms
  samples: number
}

interface ComparisonResult {
  benchmark: string
  goNsPerOp: number
  goOpsPerSec: number
  tsNsPerOp: number
  tsOpsPerSec: number
  winner: 'go' | 'ts' | 'tie'
  ratio: number // go/ts ratio (>1 means TS is faster)
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

// Parse TypeScript vitest benchmark output
function parseTSResults(content: string): TSBenchmark[] {
  const results: TSBenchmark[] = []
  const lines = content.split('\n')

  for (const line of lines) {
    // Match vitest bench table output:
    // "  · viem-ts: call (basic)               3,592.85  0.1888  6.0242  ..."
    // Format: · name    hz    min    max    mean    p75    p99    p995    p999    rme    samples
    // Note: hz can have commas (e.g., 3,592.85)
    const match = line.match(
      /[·✓]\s*(viem-ts:\s*call\s*\([^)]+\))\s+([\d,]+\.?\d*)\s+([\d.]+)\s+([\d.]+)\s+([\d.]+)/
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
      })
      
      // Warn about failed benchmarks
      if (hz === 0) {
        console.warn(`Warning: Benchmark "${name}" has 0 hz (failed or no samples)`)
      }
    }
  }

  return results
}

// Map Go benchmark names to TypeScript names
function mapBenchmarkName(goName: string): string {
  const mapping: Record<string, string> = {
    BenchmarkCall_Basic: 'viem-ts: call (basic)',
    BenchmarkCall_WithData: 'viem-ts: call (with data)',
    BenchmarkCall_WithAccount: 'viem-ts: call (with account)',
    BenchmarkCall_Decimals: 'viem-ts: call (decimals)',
    BenchmarkCall_Symbol: 'viem-ts: call (symbol)',
    BenchmarkCall_BalanceOfMultiple: 'viem-ts: call (balanceOf multiple)',
  }
  return mapping[goName] || goName
}

// Compare results
function compareResults(
  goResults: GoBenchmark[],
  tsResults: TSBenchmark[]
): ComparisonResult[] {
  const comparisons: ComparisonResult[] = []

  for (const goBench of goResults) {
    const tsName = mapBenchmarkName(goBench.name)
    const tsBench = tsResults.find((ts) => ts.name === tsName)

    if (tsBench) {
      const goOpsPerSec = 1_000_000_000 / goBench.nsPerOp
      const tsOpsPerSec = tsBench.hz
      const tsNsPerOp = 1_000_000_000 / tsOpsPerSec

      const ratio = goBench.nsPerOp / tsNsPerOp
      let winner: 'go' | 'ts' | 'tie'
      if (Math.abs(ratio - 1) < 0.05) {
        winner = 'tie'
      } else {
        winner = ratio > 1 ? 'ts' : 'go'
      }

      comparisons.push({
        benchmark: goBench.name.replace('Benchmark', ''),
        goNsPerOp: goBench.nsPerOp,
        goOpsPerSec,
        tsNsPerOp,
        tsOpsPerSec,
        winner,
        ratio,
      })
    }
  }

  return comparisons
}

// Format number with commas
function formatNumber(n: number, decimals = 0): string {
  return n.toLocaleString('en-US', {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  })
}

// Generate console report
function generateConsoleReport(comparisons: ComparisonResult[]): void {
  console.log('\n' + '='.repeat(80))
  console.log('  BENCHMARK COMPARISON: viem-go vs viem TypeScript')
  console.log('='.repeat(80) + '\n')

  // Header
  console.log(
    '| Benchmark'.padEnd(30) +
      '| Go (ns/op)'.padEnd(15) +
      '| TS (ns/op)'.padEnd(15) +
      '| Ratio'.padEnd(10) +
      '| Winner |'
  )
  console.log('|' + '-'.repeat(29) + '|' + '-'.repeat(14) + '|' + '-'.repeat(14) + '|' + '-'.repeat(9) + '|' + '-'.repeat(8) + '|')

  for (const c of comparisons) {
    const ratioStr = c.ratio.toFixed(2) + 'x'
    const winnerStr = c.winner === 'go' ? 'Go' : c.winner === 'ts' ? 'TS' : 'Tie'

    console.log(
      '| ' +
        c.benchmark.padEnd(28) +
        '| ' +
        formatNumber(c.goNsPerOp, 0).padEnd(13) +
        '| ' +
        formatNumber(c.tsNsPerOp, 0).padEnd(13) +
        '| ' +
        ratioStr.padEnd(8) +
        '| ' +
        winnerStr.padEnd(7) +
        '|'
    )
  }

  console.log('')

  // Summary
  const goWins = comparisons.filter((c) => c.winner === 'go').length
  const tsWins = comparisons.filter((c) => c.winner === 'ts').length
  const ties = comparisons.filter((c) => c.winner === 'tie').length

  console.log('Summary:')
  console.log(`  Go wins:  ${goWins}`)
  console.log(`  TS wins:  ${tsWins}`)
  console.log(`  Ties:     ${ties}`)
  console.log('')

  // Interpretation
  console.log('Ratio interpretation:')
  console.log('  > 1.0x = TypeScript is faster')
  console.log('  < 1.0x = Go is faster')
  console.log('  ≈ 1.0x = Similar performance')
  console.log('')
}

// Generate markdown report
function generateMarkdownReport(comparisons: ComparisonResult[]): string {
  let md = '# Benchmark Comparison: viem-go vs viem TypeScript\n\n'
  md += `Generated: ${new Date().toISOString()}\n\n`

  md += '## Results\n\n'
  md += '| Benchmark | Go (ns/op) | TS (ns/op) | Ratio | Winner |\n'
  md += '|-----------|------------|------------|-------|--------|\n'

  for (const c of comparisons) {
    const ratioStr = c.ratio.toFixed(2) + 'x'
    const winnerStr = c.winner === 'go' ? '**Go**' : c.winner === 'ts' ? '**TS**' : 'Tie'

    md += `| ${c.benchmark} | ${formatNumber(c.goNsPerOp, 0)} | ${formatNumber(c.tsNsPerOp, 0)} | ${ratioStr} | ${winnerStr} |\n`
  }

  md += '\n## Summary\n\n'

  const goWins = comparisons.filter((c) => c.winner === 'go').length
  const tsWins = comparisons.filter((c) => c.winner === 'ts').length
  const ties = comparisons.filter((c) => c.winner === 'tie').length

  md += `- Go wins: ${goWins}\n`
  md += `- TS wins: ${tsWins}\n`
  md += `- Ties: ${ties}\n`

  md += '\n## Notes\n\n'
  md += '- Ratio > 1.0x means TypeScript is faster\n'
  md += '- Ratio < 1.0x means Go is faster\n'
  md += '- Benchmarks run against the same Anvil instance for fair comparison\n'

  return md
}

// Main
async function main() {
  const resultsDir = join(import.meta.dir, 'results')
  const goResultsPath = join(resultsDir, 'go-results.txt')
  const tsResultsPath = join(resultsDir, 'ts-results.txt')

  // Check files exist
  if (!existsSync(goResultsPath)) {
    console.error(`Error: Go results not found at ${goResultsPath}`)
    console.error('Run "make bench" first.')
    process.exit(1)
  }

  if (!existsSync(tsResultsPath)) {
    console.error(`Error: TypeScript results not found at ${tsResultsPath}`)
    console.error('Run "make bench" first.')
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
    console.error('TypeScript results content:')
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

  // Generate reports
  generateConsoleReport(comparisons)

  // Save markdown report
  const mdReport = generateMarkdownReport(comparisons)
  const reportPath = join(resultsDir, 'comparison.md')
  writeFileSync(reportPath, mdReport)
  console.log(`Markdown report saved to: ${reportPath}`)
}

main().catch((err) => {
  console.error('Error:', err)
  process.exit(1)
})
