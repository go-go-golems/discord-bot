function scoreMessage(text) {
  const content = normalizeWhitespace(text);
  if (!content) {
    return 0;
  }
  let score = 0.2;
  if (/```/.test(content)) score += 0.35;
  if (
    /^\s*(use|run|set|install|configure|fix|remember|note|save|export|deploy)\b/im.test(
      content,
    )
  )
    score += 0.25;
  if (
    /\b(the fix is|the answer is|we should|you can|for example|try this)\b/i.test(
      content,
    )
  )
    score += 0.15;
  if (/https?:\/\//i.test(content)) score += 0.1;
  if (
    /\b(sqlite|database|bot|slack|javascript|goja|command|workflow|review)\b/i.test(
      content,
    )
  )
    score += 0.05;
  if (content.length > 160) score += 0.05;
  return Math.min(score, 1);
}

function deriveTitle(text) {
  const normalized = normalizeWhitespace(text);
  if (!normalized) {
    return "Untitled knowledge";
  }
  const firstLine = normalized.split(/[\n\.\?!]/)[0].trim();
  const title =
    firstLine.length > 72 ? `${firstLine.slice(0, 69).trim()}...` : firstLine;
  return title || "Untitled knowledge";
}

function deriveSummary(text) {
  const normalized = normalizeWhitespace(text);
  if (!normalized) {
    return "";
  }
  return normalized.length > 260
    ? `${normalized.slice(0, 257).trim()}...`
    : normalized;
}

function deriveTags(text) {
  const content = normalizeWhitespace(text).toLowerCase();
  const tags = new Set();
  if (/slack|slash command|workspace|channel/.test(content)) tags.add("slack");
  if (/sqlite|database|sql/.test(content)) tags.add("database");
  if (/javascript|js|goja/.test(content)) tags.add("javascript");
  if (/command|modal|button|component|autocomplete/.test(content))
    tags.add("interaction");
  if (/review|verify|stale|draft|queue/.test(content)) tags.add("curation");
  if (/runbook|workflow|how to|remember|teach/.test(content))
    tags.add("knowledge");
  if (/bot|automation/.test(content)) tags.add("bot");
  return Array.from(tags);
}

function deriveAliases(text, title) {
  const aliases = new Set();
  if (title) {
    aliases.add(title.toLowerCase());
  }
  const content = normalizeWhitespace(text).toLowerCase();
  const match = content.match(
    /(?:remember|teach|search|review|capture|record)\s+([a-z0-9][a-z0-9\s-]{2,48})/i,
  );
  if (match && match[1]) {
    aliases.add(normalizeWhitespace(match[1]).toLowerCase());
  }
  return Array.from(aliases);
}

function parseIdList(raw) {
  if (Array.isArray(raw)) {
    return raw.map(normalizeWhitespace).filter(Boolean);
  }
  const text = normalizeWhitespace(raw);
  if (!text) {
    return [];
  }
  return text
    .split(/[\s,;]+/g)
    .map((part) => normalizeWhitespace(part))
    .filter(Boolean);
}

function normalizeWhitespace(value) {
  if (value === undefined || value === null) {
    return "";
  }
  return String(value).replace(/\s+/g, " ").trim();
}

module.exports = {
  scoreMessage,
  deriveTitle,
  deriveSummary,
  deriveTags,
  deriveAliases,
};
