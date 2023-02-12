package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/indramahaarta/golang-authentication-jwt/controller/authcontroller"
	"github.com/indramahaarta/golang-authentication-jwt/controller/maincontroller"
	"github.com/indramahaarta/golang-authentication-jwt/database"
	"github.com/indramahaarta/golang-authentication-jwt/middlewares"
)

var DB *sql.DB

func main() {
	mux := chi.NewRouter()

	// routes
	mux.Get("/", maincontroller.Home)

	mux.Post("/register", authcontroller.Register)
	mux.Post("/login", authcontroller.Login)
	mux.Get("/logout", authcontroller.Logout)

	mux.Group(func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			r.Use(middlewares.JWT)
			r.Get("/product", maincontroller.Product)
		})
	})

	// database connection
	db, err := database.ConnectDatabase()
	if err != nil {
		log.Panic(err.Error())
	}
	DB = db

	// listen and server routes
	log.Println("Starting Application on port :8080")
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Panic(err)
	}

}
