package main

import (
	"key-value-store/db"
	"key-value-store/server"
)

const (
	ADDR         = "localhost:8080"
	DatabaseName = "db.sqlite"
)

func main() {
	sqlite, err := db.NewSQLite(DatabaseName)
	if err != nil {
		panic(err)
	}
	server.StartEchoServer(ADDR, sqlite, true)
}
