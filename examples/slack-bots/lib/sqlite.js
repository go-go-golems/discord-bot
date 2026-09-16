const db = require("database");
function createStore(schema) {
  let configured = false;
  function exec(sql, ...args) {
    const result = db.exec(sql, ...args);
    if (!result.success)
      throw new Error(result.error || "Database write failed");
    return result;
  }
  return {
    ensure(ctx) {
      if (configured) return;
      db.configure("sqlite3", ctx.config.dbPath);
      for (const statement of schema) exec(statement);
      configured = true;
    },
    exec,
    query: (sql, ...args) => db.query(sql, ...args),
  };
}
module.exports = { createStore };
