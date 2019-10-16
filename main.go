package main

import (
	"github.com/alexeyklyukin/rev4/pkg/controller"
	log "github.com/sirupsen/logrus"
	"net/http"

	"github.com/alexeyklyukin/rev4/pkg/db"
	"github.com/alexeyklyukin/rev4/pkg/routes"
)



func main() {
	pg, err := db.GetConnection("");
	if err != nil {
		log.Fatalf("could not establish database connection: %v", err)
	}

	defer pg.CloseConnections()
	ctl := controller.NewController(pg)

	router := routes.Routes(ctl)
	log.Fatal(http.ListenAndServe(":8080", router))
}

