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

// corsMiddleware maneja los headers CORS para permitir conexiones desde el frontend
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permitir origen desde Next.js (puerto 3000 por defecto)
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		
		// Permitir métodos HTTP
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		
		// Permitir headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		
		// Permitir credenciales
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		
		// Manejar preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}

func main() {

	db.StartDB()
	db.MigrateModels()

	analyzer.Ejemplo()

	r := http.NewServeMux()

	r.HandleFunc("POST /auth/login", routes.Login)
	r.HandleFunc("POST /auth/register", routes.Register)

	r.HandleFunc("POST /analyze", routes.AnalyzeCommand)
	r.HandleFunc("POST /actions", middleware.Auth(routes.CreateAction))
	r.HandleFunc("GET /actions", middleware.Auth(routes.GetAllUserActions))
	r.HandleFunc("DELETE /actions/{id}", middleware.Auth(routes.DeleteAction))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println("Servidor iniciado en el puerto: " + port)
	// Aplicar middleware CORS al router
	http.ListenAndServe(":"+port, corsMiddleware(r))
}
