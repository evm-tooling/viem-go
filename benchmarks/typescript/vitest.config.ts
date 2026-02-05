import { defineConfig } from 'vitest/config'

export default defineConfig({
  test: {
    benchmark: {
      include: ['**/*.bench.ts'],
      reporters: ['default'],
      outputFile: {
        json: '../results/ts-results.json',
      },
    },
    testTimeout: 60_000,
    hookTimeout: 60_000,
  },
})
