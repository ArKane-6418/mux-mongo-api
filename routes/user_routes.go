package routes

import (
	"mux-mongo-api/controllers"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func UserRoute(router *mux.Router) {
	// All routes related to users come here
	router.HandleFunc("/user", controllers.CreateUser()).Methods("POST")
	router.HandleFunc("/user/{userId}", controllers.GetUser()).Methods("GET")
	router.HandleFunc("/user/{userId}", controllers.UpdateUser()).Methods("PUT")
	router.HandleFunc("/user/{userId}", controllers.DeleteUser()).Methods("DELETE")
	router.HandleFunc("/users", controllers.GetAllUsers()).Methods("GET")
	router.HandleFunc("/users", controllers.DeleteAllUsers()).Methods("DELETE")
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}
