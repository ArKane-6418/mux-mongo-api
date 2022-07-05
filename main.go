package main

import (
	"log"
	"net/http"

	"github.com/ArKane-6418/mux-mongo-api/configs"
	"github.com/ArKane-6418/mux-mongo-api/routes"

	"github.com/gorilla/mux"
)

// @title Golang Mux Mongo Users API
// @version 1.0
// @description This is a MongoDB-based Users API server

// @contact.name Joshua Ong
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:6060
// @BasePath /

// @schemes http

func main() {
	router := mux.NewRouter().StrictSlash(true)

	configs.ConnectDB()

	routes.UserRoute(router)

	log.Fatal(http.ListenAndServe(":6060", router))
}
