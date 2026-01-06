import { createServer } from 'http'
import { render } from './ssr-dist/ssr.js'

const server = createServer(async (req, res) => {
  if (req.method !== 'POST') {
    res.writeHead(405)
    res.end('Method Not Allowed')
    return
  }

  let body = ''
  req.on('data', chunk => {
    body += chunk.toString()
  })

  req.on('end', async () => {
    try {
      const page = JSON.parse(body)
      const rendered = await render(page)

      res.writeHead(200, { 'Content-Type': 'application/json' })
      res.end(JSON.stringify(rendered))
    } catch (error) {
      console.error('SSR Error:', error)
      res.writeHead(500, { 'Content-Type': 'application/json' })
      res.end(JSON.stringify({ error: error.message }))
    }
  })
})

const PORT = process.env.SSR_PORT || 13714
server.listen(PORT, () => {
  console.log(`SSR server listening on port ${PORT}`)
})
