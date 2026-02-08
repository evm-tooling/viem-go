import { defineConfig } from 'vitest/config'

const suite = process.env.BENCH_SUITE
const include = suite ? [`**/${suite}.bench.ts`] : ['**/*.bench.ts']

export default defineConfig({
  test: {
    benchmark: {
      include,
      reporters: ['default'],
      outputFile: {
        json: '../results/ts-results.json',
      },
    },
    testTimeout: 120_000,
    hookTimeout: 60_000,
  },
})
