package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/RodrigoGonzalez78/go_analyzer/analyzer"
)

type AnalyzeCommandResponse struct {
	Success bool        `json:"success"`
	AST     interface{} `json:"ast,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Analysis map[string]interface{} `json:"analysis,omitempty"`
}

func AnalyzeCommand(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Command string `json:"command"`
	}

	var request Request

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Error al decodificar el contenido", http.StatusBadRequest)
		return
	}

	if request.Command == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AnalyzeCommandResponse{
			Success: false,
			Error: map[string]interface{}{
				"type":    "EMPTY_COMMAND",
				"message": "No se envió ningún comando",
				"position": 0,
			},
		})
		return
	}

	parsedAction, analyzeErr := analyzer.CreateAction(request.Command)
	if analyzeErr != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(AnalyzeCommandResponse{
			Success: false,
			Error: map[string]interface{}{
				"type":    "SYNTAX_ERROR",
				"message": analyzeErr.Error(),
				"position": 0,
			},
		})
		return
	}

	ast := buildAST(request.Command, parsedAction)
	
	// Crear información del análisis
	analysis := map[string]interface{}{
		"command": request.Command,
		"verb": parsedAction.Verbo,
		"words": parsedAction.Palabras,
		"date": parsedAction.Fecha,
		"time": parsedAction.Hora,
		"description": strings.Join(parsedAction.Palabras, " "),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AnalyzeCommandResponse{
		Success: true,
		AST:     ast,
		Analysis: analysis,
	})
} 