package main

import (
	"log"
	"net/http"
	"os"

	"github.com/RodrigoGonzalez78/go_analyzer/db"
	// "github.com/RodrigoGonzalez78/go_analyzer/analyzer"
	"github.com/RodrigoGonzalez78/go_analyzer/middleware"
	"github.com/RodrigoGonzalez78/go_analyzer/routes"
	"github.com/rs/cors"
)

func main() {

	db.StartDB()
	db.MigrateModels()

	// Comentamos analyzer.Ejemplo() para testing
	// analyzer.Ejemplo()

	// Crear el router
	r := http.NewServeMux()

	r.HandleFunc("POST /auth/login", routes.Login)
	r.HandleFunc("POST /auth/register", routes.Register)

	r.HandleFunc("POST /actions", middleware.Auth(routes.CreateAction))
	r.HandleFunc("GET /actions", middleware.Auth(routes.GetAllUserActions))
	r.HandleFunc("DELETE /actions/{id}", middleware.Auth(routes.DeleteAction))

	// Configurar CORS para permitir solicitudes desde el frontend
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Origen del frontend
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           300, // Tiempo máximo (en segundos) que el navegador puede cachear los resultados de una solicitud preflight
	})

	// Aplicar middleware CORS a todas las rutas
	handler := corsHandler.Handler(r)

	// Probar con un puerto diferente para descartar problemas de permisos o conflictos
	port := "8081"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	log.Printf("Servidor con CORS habilitado iniciando en el puerto: %s\n", port)
	
	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Printf("Error tipo: %T\n", err)
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
