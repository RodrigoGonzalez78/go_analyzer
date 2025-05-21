package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RodrigoGonzalez78/go_analyzer/db"
	"github.com/RodrigoGonzalez78/go_analyzer/routes"
)

func main() {

	db.StartDB()
	r := http.NewServeMux()

	r.HandleFunc("POST /login", routes.Login)
	r.HandleFunc("POST /register", routes.Register)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Servidor iniciado en el puerto: " + port)
	http.ListenAndServe(":"+port, r)
}
