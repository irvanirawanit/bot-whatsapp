package config

import (
	"database/sql"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func InitSqlite() {
	var err error
	DB, err = sql.Open("sqlite", "./users.db")
	if err != nil {
		panic(err)
	}

	sqlStmt := `CREATE TABLE IF NOT EXISTS users (id INTEGER NOT NULL PRIMARY KEY, name TEXT NOT NULL, token TEXT NOT NULL, webhook TEXT NOT NULL default "", jid TEXT NOT NULL default "", qrcode TEXT NOT NULL default "", connected INTEGER, expiration INTEGER, events TEXT NOT NULL default "All", timeouts TEXT);`
	_, err = DB.Exec(sqlStmt)
	if err != nil {
		panic(err)
	}

	sqlStmt = `INSERT OR IGNORE INTO users(ID, name, token, webhook, jid, qrcode, connected, expiration, events) VALUES(1, 'user1', 'user1', 'user1', 'user1', 'user1', 0, 0, 'All')`
	_, _ = DB.Exec(sqlStmt)

}
