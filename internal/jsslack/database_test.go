package jsslack

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-go-golems/discord-bot/pkg/slackbot"
	"github.com/stretchr/testify/require"
)

func TestDatabaseInspectionAndReopen(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "state.sqlite")
	source := `const db=require("database");let ready=false;module.exports=require("slack").defineBot(({configure,command})=>{configure({name:"persistent",run:{fields:{dbPath:{type:"string",required:true}}}});command("/test",{description:"Counter"},ctx=>{if(!ready){db.configure("sqlite3",ctx.config.dbPath);db.exec("CREATE TABLE IF NOT EXISTS counts (value INTEGER)");ready=true;}db.exec("INSERT INTO counts(value) VALUES (?)",1);return {text:String(db.query("SELECT count(*) AS n FROM counts")[0].n)};});});`
	path := script(t, source)
	_, err := Inspect(context.Background(), path, 0)
	require.NoError(t, err)
	_, err = os.Stat(dbPath)
	require.True(t, os.IsNotExist(err))
	for _, expected := range []string{"1", "2"} {
		h, err := Load(context.Background(), path, Options{Config: map[string]any{"dbPath": dbPath}})
		require.NoError(t, err)
		rec := &slackbot.Recorder{}
		require.NoError(t, h.Dispatch(context.Background(), command(), rec))
		require.Equal(t, expected, rec.Operations()[0].Reply.Text)
		require.NoError(t, h.Close(context.Background()))
		require.NoError(t, h.database.Close())
	}
}
func TestDatabaseCannotOpenDuringInspection(t *testing.T) {
	path := script(t, `require("database").configure("sqlite3",":memory:");`)
	_, err := Inspect(context.Background(), path, 0)
	require.ErrorContains(t, err, "only available inside handlers")
}
