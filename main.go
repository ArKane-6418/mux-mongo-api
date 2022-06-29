package main

import (
	"log"
	"mux-mongo-api/configs"
	"mux-mongo-api/routes"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	configs.ConnectDB()

	routes.UserRoute(router)

	log.Fatal(http.ListenAndServe(":6000", router))
}
