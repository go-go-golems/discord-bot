const discord = require("discord")
const { defineBot } = discord
const fs = require("fs")
const express = require("express")

const app = express.app()

const pageHtml = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>xgoja Discord Say</title>
  <link rel="stylesheet" href="/retro.css">
</head>
<body>
  <main class="page-shell" aria-labelledby="title">
    <header class="masthead">
      <p class="eyebrow">xgoja express → discord</p>
      <h1 id="title">Send a Discord message</h1>
      <p class="lede">A tiny black-and-white System 1 inspired form. No window chrome, no menu bar.</p>
    </header>

    <form id="say-form" class="say-form" method="post" action="/say">
      <label>
        <span>Channel ID</span>
        <input name="channelId" type="text" inputmode="numeric" autocomplete="off" placeholder="123456789012345678" required>
      </label>

      <label>
        <span>Message</span>
        <textarea name="content" rows="5" placeholder="hello from xgoja express" required></textarea>
      </label>

      <div class="button-row">
        <button type="submit">Say it</button>
        <button type="reset" class="secondary">Clear</button>
      </div>
    </form>

    <pre id="result" class="result" aria-live="polite">Ready.</pre>
  </main>

  <script>
    const form = document.getElementById('say-form')
    const result = document.getElementById('result')
    form.addEventListener('submit', async (event) => {
      event.preventDefault()
      result.textContent = 'Sending…'
      const body = new URLSearchParams(new FormData(form))
      try {
        const response = await fetch('/say', {
          method: 'POST',
          headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
          body
        })
        const text = await response.text()
        result.textContent = text
      } catch (error) {
        result.textContent = String(error)
      }
    })
  </script>
</body>
</html>`

const retroCss = `:root {
  color-scheme: light;
  --ink: #000;
  --paper: #fff;
  --shade: #d9d9d9;
  --stripe: #f1f1f1;
}

* { box-sizing: border-box; }

body {
  margin: 0;
  min-height: 100vh;
  color: var(--ink);
  background:
    linear-gradient(45deg, var(--stripe) 25%, transparent 25%),
    linear-gradient(-45deg, var(--stripe) 25%, transparent 25%),
    linear-gradient(45deg, transparent 75%, var(--stripe) 75%),
    linear-gradient(-45deg, transparent 75%, var(--stripe) 75%),
    var(--paper);
  background-position: 0 0, 0 6px, 6px -6px, -6px 0;
  background-size: 12px 12px;
  font-family: Chicago, Geneva, Monaco, 'Courier New', monospace;
  font-size: 16px;
  line-height: 1.35;
}

.page-shell {
  width: min(680px, calc(100vw - 32px));
  margin: 48px auto;
  padding: 28px;
  background: var(--paper);
  border: 3px solid var(--ink);
  box-shadow: 8px 8px 0 var(--ink);
}

.masthead {
  padding-bottom: 18px;
  margin-bottom: 18px;
  border-bottom: 3px double var(--ink);
}

.eyebrow {
  margin: 0 0 8px;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  font-size: 12px;
}

h1 {
  margin: 0 0 10px;
  font-size: clamp(28px, 6vw, 44px);
  line-height: 0.95;
  letter-spacing: -0.04em;
}

.lede { margin: 0; max-width: 46rem; }

.say-form {
  display: grid;
  gap: 18px;
}

label span {
  display: inline-block;
  margin-bottom: 6px;
  padding: 0 4px;
  background: var(--paper);
  font-weight: bold;
}

input,
textarea {
  width: 100%;
  color: var(--ink);
  background: var(--paper);
  border: 3px solid var(--ink);
  border-radius: 0;
  padding: 10px 12px;
  font: inherit;
  outline: none;
  box-shadow: inset 2px 2px 0 var(--shade);
}

input:focus,
textarea:focus {
  outline: 3px solid var(--ink);
  outline-offset: 3px;
}

.button-row {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
}

button {
  color: var(--ink);
  background: var(--paper);
  border: 3px solid var(--ink);
  border-radius: 0;
  padding: 9px 18px;
  font: inherit;
  font-weight: bold;
  cursor: pointer;
  box-shadow: 4px 4px 0 var(--ink);
}

button:active {
  transform: translate(4px, 4px);
  box-shadow: none;
}

button.secondary {
  background: repeating-linear-gradient(45deg, var(--paper), var(--paper) 3px, var(--shade) 3px, var(--shade) 6px);
}

.result {
  min-height: 72px;
  margin: 22px 0 0;
  padding: 12px;
  white-space: pre-wrap;
  background: var(--paper);
  border: 3px double var(--ink);
  font: 14px Monaco, 'Courier New', monospace;
}`

app.get("/retro.css", (_req, res) => {
  res.type("text/css; charset=utf-8").send(retroCss)
})

app.get("/", (_req, res) => {
  res.type("text/html; charset=utf-8").send(pageHtml)
})

app.post("/say", async (req, res) => {
  const body = req.body || {}
  const channelId = body.channelId || req.query.channelId
  const content = body.content || body.message || "hello from xgoja express"
  if (!channelId) {
    res.status(400).json({ ok: false, error: "channelId is required" })
    return
  }
  await discord.channels.send(channelId, { content })
  res.json({ ok: true, channelId, content })
})

module.exports = defineBot(({ command, event, configure }) => {
  configure({
    name: "fs-express-smoke",
    description: "xgoja Discord bot smoke test using fs plus xgoja-owned Express HTTP routes"
  })

  event("ready", async (ctx) => {
    ctx.log.info("fs-express-smoke ready from generated xgoja runtime")
  })

  command("ping", { description: "Return a simple xgoja pong" }, async () => {
    return { content: "pong from xgoja discord-bot provider" }
  })

  command("read-config", { description: "Read a local file through require('fs')" }, async () => {
    const text = fs.readFileSync("./bot-data/message.txt", "utf8").trim()
    return { content: `config says: ${text}` }
  })

  command("express-status", { description: "Explain current express status" }, async () => {
    return { content: "express HTTP is available from the xgoja go-go-goja-http provider at GET /, GET /retro.css, and POST /say" }
  })
})
