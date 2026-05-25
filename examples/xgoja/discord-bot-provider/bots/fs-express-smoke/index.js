const discord = require("discord")
const { defineBot } = discord
const fs = require("fs")
const express = require("express")

const app = express.app()

function readAsset(name) {
  return fs.readFileSync(`./web/${name}`, "utf8")
}

app.get("/retro.css", (_req, res) => {
  res.type("text/css; charset=utf-8").send(readAsset("retro.css"))
})

app.get("/", (_req, res) => {
  res.type("text/html; charset=utf-8").send(readAsset("index.html"))
})

app.get("/channels", async (req, res) => {
  const guildId = req.query.guildId
  const channels = guildId ? await discord.channels.list(guildId) : await discord.channels.list()
  const sendableTypes = new Set(["0", "5", "10", "11", "12"])
  const choices = channels
    .filter((channel) => channel && channel.id && channel.name && (sendableTypes.has(String(channel.type)) || channel.thread))
    .map((channel) => ({
      id: channel.id,
      name: channel.name,
      type: channel.type,
      parentId: channel.parentID || "",
      position: channel.position || 0
    }))
    .sort((a, b) => a.position - b.position || a.name.localeCompare(b.name))
  res.json({ ok: true, channels: choices })
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
    return { content: "express HTTP is available from the xgoja go-go-goja-http provider at GET /, GET /retro.css, GET /channels, and POST /say" }
  })
})
