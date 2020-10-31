package server

import (
	"log"
	"net"
	"net/http"

	"github.com/labstack/echo/v4"

	"key-value-store/controllers"
	"key-value-store/db"
)

const CacheSize = 1024

type EchoServer struct {
	*http.Server
	addr string
}

func StartEchoServer(addr string, dbtx db.DBTX, block bool) *EchoServer {
	e := echo.New()
	e.HTTPErrorHandler = func(err error, ctx echo.Context) {
		log.Printf("%+v\n", err)
		e.DefaultHTTPErrorHandler(err, ctx)
	}
	initRoutes(e, dbtx)

	ready := make(chan struct{})
	join := make(chan struct{})

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	go func() {
		ready <- struct{}{}
		if err := http.Serve(lis, e); err != nil {
			e.Logger.Fatal(err)
		}
		<-join
	}()

	<-ready
	if block {
		join <- struct{}{}
	}
	return &EchoServer{Server: e.Server, addr: addr}
}

func initRoutes(e *echo.Echo, dbtx db.DBTX) {
	c := controllers.NewKeyValueController(CacheSize, dbtx)
	e.GET("/kv/:key", c.Get)
	e.PUT("/kv/:key", c.Put)
}
