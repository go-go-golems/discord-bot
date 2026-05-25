const discord = require("discord")
const { defineBot } = discord
const fs = require("fs")
const express = require("express")

const app = express.app()

app.get("/", (_req, res) => {
  res.json({
    ok: true,
    bot: "fs-express-smoke",
    routes: ["GET /", "POST /say"],
    message: "xgoja Express HTTP provider is running"
  })
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
    return { content: "express HTTP is available from the xgoja go-go-goja-http provider at GET / and POST /say" }
  })
})
