package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jeypc/go-jwt-mux/controller/authController"
	"github.com/jeypc/go-jwt-mux/controller/importController"
	insertcontroller "github.com/jeypc/go-jwt-mux/controller/insertController"
	"github.com/jeypc/go-jwt-mux/middleware"
	"github.com/jeypc/go-jwt-mux/models"
)

func main() {

	models.ConnectDatabase()
	r := mux.NewRouter()

	r.HandleFunc("/login", authController.Login).Methods("POST")
	r.HandleFunc("/register", authController.Register).Methods("POST")
	r.HandleFunc("/logout", authController.Logout).Methods("GET")

	// membuat group yang dibatasi oleh middlware
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/change-password", authController.ChangePassword).Methods("POST")
	api.HandleFunc("/import-data", importController.ImportFile).Methods("POST")
	api.HandleFunc("/insert-report", insertcontroller.InsertMultiple).Methods("POST")
	api.Use(middleware.JWTMiddleware)

	log.Fatal(http.ListenAndServe(":8080", r))
}
