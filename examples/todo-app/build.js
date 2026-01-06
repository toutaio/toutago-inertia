import * as esbuild from 'esbuild'
import { promises as fs } from 'fs'
import path from 'path'

const isWatch = process.argv.includes('--watch')
const isSSR = process.argv.includes('--ssr')

const buildClient = async () => {
  const context = await esbuild.context({
    entryPoints: ['views/app.ts'],
    bundle: true,
    outdir: 'public/build',
    format: 'esm',
    splitting: true,
    sourcemap: true,
    metafile: true,
    loader: {
      '.vue': 'ts',
    },
    plugins: [
      {
        name: 'vue',
        setup(build) {
          const vuePlugin = require('@vitejs/plugin-vue')
          // Simplified Vue plugin integration
        },
      },
    ],
  })

  if (isWatch) {
    await context.watch()
    console.log('Watching for client changes...')
  } else {
    await context.rebuild()
    await context.dispose()
    console.log('Client build complete')
  }

  // Generate manifest
  const result = await esbuild.build({
    entryPoints: ['views/app.ts'],
    bundle: true,
    write: false,
    metafile: true,
  })

  const manifest = {}
  for (const [file, info] of Object.entries(result.metafile.outputs)) {
    const key = file.replace('public/build/', '')
    manifest[key] = {
      file: key,
      imports: info.imports?.map(i => i.path.replace('public/build/', '')) || [],
    }
  }

  await fs.writeFile(
    'public/build/manifest.json',
    JSON.stringify(manifest, null, 2)
  )
}

const buildSSR = async () => {
  const context = await esbuild.context({
    entryPoints: ['views/ssr.ts'],
    bundle: true,
    outfile: 'ssr-dist/ssr.js',
    format: 'esm',
    platform: 'node',
    sourcemap: true,
    external: ['vue', 'vue/server-renderer'],
  })

  if (isWatch) {
    await context.watch()
    console.log('Watching for SSR changes...')
  } else {
    await context.rebuild()
    await context.dispose()
    console.log('SSR build complete')
  }
}

// Main build
;(async () => {
  try {
    await fs.mkdir('public/build', { recursive: true })
    await buildClient()
    
    if (isSSR) {
      await fs.mkdir('ssr-dist', { recursive: true })
      await buildSSR()
    }
    
    if (!isWatch) {
      console.log('Build complete!')
    }
  } catch (error) {
    console.error('Build failed:', error)
    process.exit(1)
  }
})()
