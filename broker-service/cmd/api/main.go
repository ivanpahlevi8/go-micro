package main

import (
	"fmt"
	"log"
	"net/http"
)

// create const
const port = "80"

// create receiver app
type AppConfig struct{}

func main() {
	// create app config
	app := AppConfig{}

	// create server
	serv := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	log.Printf("starting server back-end at port : %s\n", port)

	// start server
	err := serv.ListenAndServe()

	// check for an error
	if err != nil {
		log.Printf("errro when starting server : %s\n", err.Error())
		return
	}
}
