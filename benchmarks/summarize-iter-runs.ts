/**
 * Summarize per-iteration benchmark runs.
 *
 * Reads:  <runDir>/iter-<N>/_overall/comparison.md
 * Writes: <runDir>/summary.md
 */

import { existsSync, readdirSync, readFileSync, writeFileSync } from 'fs'
import { join } from 'path'
import { parseArgs } from 'util'

type IterSummary = {
  iter: number
  overallWinner: 'go' | 'ts' | 'tie' | 'unknown'
  overallSpeedup: number | null
  goWins: number | null
  tsWins: number | null
  total: number | null
}

function parseOverallLine(md: string): { winner: IterSummary['overallWinner']; speedup: number | null } {
  // Examples:
  // **🏆 Go is 15.98x faster overall** (geometric mean)
  // **🏆 TypeScript is 2.10x faster overall** (geometric mean)
  // **🤝 Performance is similar between Go and TypeScript**
  const go = md.match(/\*\*🏆\s+Go is\s+([\d.]+)x faster overall\*\*/i)
  if (go) return { winner: 'go', speedup: Number(go[1]) }
  const ts = md.match(/\*\*🏆\s+TypeScript is\s+([\d.]+)x faster overall\*\*/i)
  if (ts) return { winner: 'ts', speedup: Number(ts[1]) }
  const tie = md.match(/\*\*🤝\s+Performance is similar between Go and TypeScript\*\*/i)
  if (tie) return { winner: 'tie', speedup: null }
  return { winner: 'unknown', speedup: null }
}

function parseWinsRow(md: string): { goWins: number | null; tsWins: number | null; total: number | null } {
  // | Wins | 59/59 | 0/59 |
  const m = md.match(/\|\s*Wins\s*\|\s*(\d+)\s*\/\s*(\d+)\s*\|\s*(\d+)\s*\/\s*(\d+)\s*\|/i)
  if (!m) return { goWins: null, tsWins: null, total: null }
  const goWins = Number(m[1])
  const total = Number(m[2])
  const tsWins = Number(m[3])
  return { goWins, tsWins, total }
}

function geometricMean(values: number[]): number | null {
  const xs = values.filter((v) => Number.isFinite(v) && v > 0)
  if (xs.length === 0) return null
  const logSum = xs.reduce((s, v) => s + Math.log(v), 0)
  return Math.exp(logSum / xs.length)
}

function format(n: number | null, decimals = 2): string {
  if (n === null || !Number.isFinite(n)) return 'N/A'
  return n.toLocaleString('en-US', { minimumFractionDigits: decimals, maximumFractionDigits: decimals })
}

const { values: args } = parseArgs({
  options: {
    'run-dir': { type: 'string' },
  },
  allowPositionals: true,
})

const runDir = args['run-dir']
if (!runDir) {
  console.error('Missing --run-dir')
  process.exit(1)
}

if (!existsSync(runDir)) {
  console.error(`Run dir does not exist: ${runDir}`)
  process.exit(1)
}

const iterDirs = readdirSync(runDir, { withFileTypes: true })
  .filter((d) => d.isDirectory() && d.name.startsWith('iter-'))
  .map((d) => d.name)
  .map((name) => ({ name, iter: Number(name.replace(/^iter-/, '')) }))
  .filter((x) => Number.isFinite(x.iter))
  .sort((a, b) => a.iter - b.iter)

const summaries: IterSummary[] = []

for (const d of iterDirs) {
  const comparisonPath = join(runDir, d.name, '_overall', 'comparison.md')
  if (!existsSync(comparisonPath)) continue
  const md = readFileSync(comparisonPath, 'utf-8')
  const { winner, speedup } = parseOverallLine(md)
  const wins = parseWinsRow(md)
  summaries.push({
    iter: d.iter,
    overallWinner: winner,
    overallSpeedup: speedup,
    goWins: wins.goWins,
    tsWins: wins.tsWins,
    total: wins.total,
  })
}

if (summaries.length === 0) {
  console.warn('No per-iteration overall comparisons found; skipping summary.md')
  process.exit(0)
}

const goWinsCount = summaries.filter((s) => s.overallWinner === 'go').length
const tsWinsCount = summaries.filter((s) => s.overallWinner === 'ts').length
const tieCount = summaries.filter((s) => s.overallWinner === 'tie').length

const geoMean = geometricMean(summaries.map((s) => s.overallSpeedup ?? NaN))
const avg = summaries
  .map((s) => s.overallSpeedup)
  .filter((v): v is number => typeof v === 'number' && Number.isFinite(v))
  .reduce((a, b) => a + b, 0) / summaries.filter((s) => typeof s.overallSpeedup === 'number').length

let concluding: string
if (goWinsCount > tsWinsCount) concluding = `Go wins in ${goWinsCount}/${summaries.length} iteration levels.`
else if (tsWinsCount > goWinsCount) concluding = `TypeScript wins in ${tsWinsCount}/${summaries.length} iteration levels.`
else concluding = `No clear winner across iteration levels (${goWinsCount} Go wins, ${tsWinsCount} TS wins, ${tieCount} ties).`

let out = ''
out += '# Benchmark Run Summary (all iteration levels)\n\n'
out += `Run dir: \`${runDir}\`\n\n`
out += `- Iteration levels: **${summaries.length}**\n`
out += `- Overall: **${concluding}**\n`
out += `- Geometric mean of per-level speedup: **${format(geoMean, 2)}x**\n`
out += `- Arithmetic mean of per-level speedup: **${format(avg, 2)}x**\n\n`

out += '## Per-iteration results\n\n'
out += '| Iterations | Winner | Speedup (geo mean) | Go wins | TS wins | Total |\n'
out += '|-----------:|--------|-------------------:|--------:|--------:|------:|\n'
for (const s of summaries) {
  const winner =
    s.overallWinner === 'go' ? 'Go' : s.overallWinner === 'ts' ? 'TypeScript' : s.overallWinner === 'tie' ? 'Tie' : 'Unknown'
  out += `| ${s.iter} | ${winner} | ${s.overallSpeedup === null ? 'N/A' : format(s.overallSpeedup, 2)} | ${s.goWins ?? 'N/A'} | ${s.tsWins ?? 'N/A'} | ${s.total ?? 'N/A'} |\n`
}
out += '\n'

out += '## Notes\n\n'
out += '- Per-iteration overall comparisons are read from `iter-<N>/_overall/comparison.md`.\n'
out += '- Per-suite raw outputs are stored in `iter-<N>/<suite>/go.md` and `iter-<N>/<suite>/ts.md`.\n'

const summaryPath = join(runDir, 'summary.md')
writeFileSync(summaryPath, out)
console.log(`Wrote summary: ${summaryPath}`)

