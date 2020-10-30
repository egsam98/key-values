package main

import (
	"log"

	"github.com/labstack/echo/v4"

	"key-value-store/controllers"
	"key-value-store/db"
)

const ADDR = ":8080"

func main() {
	e := echo.New()
	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		log.Printf("%+v\n", err)
		e.DefaultHTTPErrorHandler(err, ctx)
	}
	queries := initQueries(e)
	initRoutes(e, queries)
	e.Logger.Fatal(e.Start(ADDR))
}

func initQueries(e *echo.Echo) *db.Queries {
	sqlite, err := db.NewSQLite()
	if err != nil {
		e.Logger.Fatal(err)
	}
	return db.New(sqlite)
}

func initRoutes(e *echo.Echo, queries *db.Queries) {
	c := controllers.NewKeyValueController(queries)
	e.GET("/kv/:key", c.Get)
	e.PUT("/kv/:key", c.Put)
}
