import { defineConfig } from 'tsup'

export default defineConfig({
  entry: ['src/index.ts'],
  format: ['cjs', 'esm'],
  dts: false,  // Skip for now due to Vue component type complexity
  clean: true,
  external: ['vue'],
  esbuildOptions(options) {
    options.loader = {
      '.vue': 'text'
    }
  }
})
