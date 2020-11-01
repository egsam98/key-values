package main

import (
	"log"
	"os/exec"

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
	if err := runSwaggerUI(); err != nil {
		log.Println(err)
	}
	server.StartEchoServer(ADDR, sqlite, true)
}

func runSwaggerUI() error {
	return exec.Command("make", "swagger-ui").Start()
}
