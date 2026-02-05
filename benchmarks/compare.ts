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
  category: string
  goNsPerOp: number
  goOpsPerSec: number
  tsNsPerOp: number
  tsOpsPerSec: number
  winner: 'go' | 'ts' | 'tie'
  ratio: number // go/ts ratio (>1 means TS is faster)
  speedup: number // how much faster the winner is (always >= 1)
  speedupStr: string // "Go 1.5x faster" or "TS 1.2x faster"
}

interface OverallStats {
  totalBenchmarks: number
  goWins: number
  tsWins: number
  ties: number
  avgGoNsPerOp: number
  avgTsNsPerOp: number
  avgRatio: number
  overallWinner: 'go' | 'ts' | 'tie'
  overallSpeedup: number
  overallSummary: string
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

// Categorize benchmark by name
function categorizeBenchmark(name: string): string {
  if (name.includes('Basic') || name.includes('basic')) return 'Basic Calls'
  if (name.includes('Data') || name.includes('data')) return 'With Parameters'
  if (name.includes('Account') || name.includes('account')) return 'With Account'
  if (name.includes('Multiple') || name.includes('multiple')) return 'Batch Operations'
  return 'Other'
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

      // Calculate speedup (how much faster the winner is)
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
        category: categorizeBenchmark(benchName),
        goNsPerOp: goBench.nsPerOp,
        goOpsPerSec,
        tsNsPerOp,
        tsOpsPerSec,
        winner,
        ratio,
        speedup,
        speedupStr,
      })
    }
  }

  return comparisons
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

  return {
    totalBenchmarks: comparisons.length,
    goWins,
    tsWins,
    ties,
    avgGoNsPerOp,
    avgTsNsPerOp,
    avgRatio,
    overallWinner,
    overallSpeedup,
    overallSummary,
  }
}

// Format number with commas
function formatNumber(n: number, decimals = 0): string {
  return n.toLocaleString('en-US', {
    minimumFractionDigits: decimals,
    maximumFractionDigits: decimals,
  })
}

// Generate console report
function generateConsoleReport(comparisons: ComparisonResult[], stats: OverallStats): void {
  console.log('\n' + '='.repeat(90))
  console.log('  BENCHMARK COMPARISON: viem-go vs viem TypeScript')
  console.log('='.repeat(90) + '\n')

  // Overall Summary Box
  console.log('┌' + '─'.repeat(88) + '┐')
  console.log('│' + stats.overallSummary.padStart(45 + stats.overallSummary.length / 2).padEnd(88) + '│')
  console.log('└' + '─'.repeat(88) + '┘')
  console.log('')

  // Detailed Results Table
  console.log('DETAILED RESULTS')
  console.log('─'.repeat(90))
  console.log(
    '| Benchmark'.padEnd(28) +
      '| Go (ns/op)'.padEnd(14) +
      '| TS (ns/op)'.padEnd(14) +
      '| Go (ops/s)'.padEnd(12) +
      '| TS (ops/s)'.padEnd(12) +
      '| Result'.padEnd(20) + '|'
  )
  console.log('|' + '-'.repeat(27) + '|' + '-'.repeat(13) + '|' + '-'.repeat(13) + '|' + '-'.repeat(11) + '|' + '-'.repeat(11) + '|' + '-'.repeat(19) + '|')

  for (const c of comparisons) {
    console.log(
      '| ' +
        c.benchmark.padEnd(26) +
        '| ' +
        formatNumber(c.goNsPerOp, 0).padEnd(12) +
        '| ' +
        formatNumber(c.tsNsPerOp, 0).padEnd(12) +
        '| ' +
        formatNumber(c.goOpsPerSec, 0).padEnd(10) +
        '| ' +
        formatNumber(c.tsOpsPerSec, 0).padEnd(10) +
        '| ' +
        c.speedupStr.padEnd(18) +
        '|'
    )
  }

  console.log('')

  // Category Breakdown
  const categories = [...new Set(comparisons.map((c) => c.category))]
  if (categories.length > 1) {
    console.log('BY CATEGORY')
    console.log('─'.repeat(50))
    for (const category of categories) {
      const catComparisons = comparisons.filter((c) => c.category === category)
      const catStats = calculateOverallStats(catComparisons)
      const winnerIcon = catStats.overallWinner === 'go' ? '🟢' : catStats.overallWinner === 'ts' ? '🔵' : '⚪'
      console.log(`  ${winnerIcon} ${category}: ${catStats.overallWinner === 'tie' ? 'Similar' : catStats.overallWinner === 'go' ? `Go ${catStats.overallSpeedup.toFixed(2)}x faster` : `TS ${catStats.overallSpeedup.toFixed(2)}x faster`}`)
    }
    console.log('')
  }

  // Summary Statistics
  console.log('SUMMARY')
  console.log('─'.repeat(50))
  console.log(`  Total benchmarks: ${stats.totalBenchmarks}`)
  console.log(`  Go wins:          ${stats.goWins} (${((stats.goWins / stats.totalBenchmarks) * 100).toFixed(0)}%)`)
  console.log(`  TS wins:          ${stats.tsWins} (${((stats.tsWins / stats.totalBenchmarks) * 100).toFixed(0)}%)`)
  console.log(`  Ties:             ${stats.ties} (${((stats.ties / stats.totalBenchmarks) * 100).toFixed(0)}%)`)
  console.log('')
  console.log(`  Avg Go:           ${formatNumber(stats.avgGoNsPerOp, 0)} ns/op (${formatNumber(1_000_000_000 / stats.avgGoNsPerOp, 0)} ops/s)`)
  console.log(`  Avg TS:           ${formatNumber(stats.avgTsNsPerOp, 0)} ns/op (${formatNumber(1_000_000_000 / stats.avgTsNsPerOp, 0)} ops/s)`)
  console.log('')

  // Legend
  console.log('LEGEND')
  console.log('─'.repeat(50))
  console.log('  🟢 Go faster  |  🔵 TS faster  |  ⚪ Similar')
  console.log('  ns/op = nanoseconds per operation (lower is better)')
  console.log('  ops/s = operations per second (higher is better)')
  console.log('')
}

// Generate markdown report
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

  // Detailed Results
  md += '## Detailed Results\n\n'
  md += '| Benchmark | Go (ns/op) | TS (ns/op) | Go (ops/s) | TS (ops/s) | Result |\n'
  md += '|-----------|------------|------------|------------|------------|--------|\n'

  for (const c of comparisons) {
    const resultIcon = c.winner === 'go' ? '🟢' : c.winner === 'ts' ? '🔵' : '⚪'
    md += `| ${c.benchmark} | ${formatNumber(c.goNsPerOp, 0)} | ${formatNumber(c.tsNsPerOp, 0)} | ${formatNumber(c.goOpsPerSec, 0)} | ${formatNumber(c.tsOpsPerSec, 0)} | ${resultIcon} ${c.speedupStr} |\n`
  }

  // Category Breakdown
  const categories = [...new Set(comparisons.map((c) => c.category))]
  if (categories.length > 1) {
    md += '\n## By Category\n\n'
    for (const category of categories) {
      const catComparisons = comparisons.filter((c) => c.category === category)
      const catStats = calculateOverallStats(catComparisons)
      const icon = catStats.overallWinner === 'go' ? '🟢' : catStats.overallWinner === 'ts' ? '🔵' : '⚪'
      const result = catStats.overallWinner === 'tie' ? 'Similar' : catStats.overallWinner === 'go' ? `Go ${catStats.overallSpeedup.toFixed(2)}x faster` : `TS ${catStats.overallSpeedup.toFixed(2)}x faster`
      md += `- ${icon} **${category}**: ${result}\n`
    }
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

  // Calculate overall statistics
  const stats = calculateOverallStats(comparisons)

  // Generate reports
  generateConsoleReport(comparisons, stats)

  // Save markdown report
  const mdReport = generateMarkdownReport(comparisons, stats)
  const reportPath = join(resultsDir, 'comparison.md')
  writeFileSync(reportPath, mdReport)
  console.log(`Markdown report saved to: ${reportPath}`)
}

main().catch((err) => {
  console.error('Error:', err)
  process.exit(1)
})
