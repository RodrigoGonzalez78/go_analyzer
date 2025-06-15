package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RodrigoGonzalez78/go_analyzer/analyzer"
	"github.com/RodrigoGonzalez78/go_analyzer/db"
	"github.com/RodrigoGonzalez78/go_analyzer/middleware"
	"github.com/RodrigoGonzalez78/go_analyzer/routes"
)

func main() {

	db.StartDB()
	db.MigrateModels()

	analyzer.Ejemplo()

	r := http.NewServeMux()

	r.HandleFunc("POST /auth/login", routes.Login)
	r.HandleFunc("POST /auth/register", routes.Register)

	r.HandleFunc("POST /actions", middleware.Auth(routes.CreateAction))
	r.HandleFunc("GET /actions", middleware.Auth(routes.GetAllUserActions))
	r.HandleFunc("DELETE /actions/{id}", middleware.Auth(routes.DeleteAction))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Servidor iniciado en el puerto: " + port)
	http.ListenAndServe(":"+port, r)
}
