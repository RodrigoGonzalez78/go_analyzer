package routes

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/RodrigoGonzalez78/go_analyzer/analyzer"
	"github.com/RodrigoGonzalez78/go_analyzer/db"
	"github.com/RodrigoGonzalez78/go_analyzer/models"
)

type CreateActionResponse struct {
	Success bool        `json:"success"`
	AST     interface{} `json:"ast,omitempty"`
	Error   interface{} `json:"error,omitempty"`
	Analysis map[string]interface{} `json:"analysis,omitempty"`
}

func CreateAction(w http.ResponseWriter, r *http.Request) {
	claim, _ := r.Context().Value("userData").(*models.Claim)

	type Request struct {
		Comand string `json:"comand"`
	}

	var comand Request

	err := json.NewDecoder(r.Body).Decode(&comand)
	if err != nil {
		http.Error(w, "Error al decodificar el contenido", http.StatusBadRequest)
		return
	}

	if comand.Comand == "" {
		http.Error(w, "No se envio ningun comando", http.StatusBadRequest)
		return
	}

	parsedAction, analyzeErr := analyzer.CreateAction(comand.Comand)
	if analyzeErr != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CreateActionResponse{
			Success: false,
			Error:   analyzeErr.Error(), // ya es *AnalyzerError
		})
		return
	}

	action, err := analyzer.TransformToAction(parsedAction, claim.UserName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(CreateActionResponse{
			Success: false,
			Error: map[string]interface{}{
				"type":    "TRANSFORM_ERROR",
				"message": err.Error(),
				"position": 0,
			},
		})
		return
	}

	action.UserName = claim.UserName

	err = db.CreateAction(action)
	if err != nil {
		http.Error(w, "Error creando la accion", http.StatusInternalServerError)
		return
	}

	ast := buildAST(comand.Comand, parsedAction)
	
	// Crear información del análisis
	analysis := map[string]interface{}{
		"command": comand.Comand,
		"verb": parsedAction.Verbo,
		"words": parsedAction.Palabras,
		"date": parsedAction.Fecha,
		"time": parsedAction.Hora,
		"description": strings.Join(parsedAction.Palabras, " "),
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CreateActionResponse{
		Success: true,
		AST:     ast,
		Analysis: analysis,
	})
}

// buildAST igual que en analyzeCommand.go
func buildAST(command string, parsed analyzer.ParsedAction) map[string]interface{} {
	root := map[string]interface{}{
		"name": "COMANDO",
		"attributes": map[string]interface{}{
			"value": command,
		},
		"children": []map[string]interface{}{
			{
				"name": "VERBO",
				"attributes": map[string]interface{}{
					"value": parsed.Verbo,
					"validVerbs": []string{"agendá", "anotá", "recordame"},
				},
			},
			{
				"name": "PALABRAS",
				"attributes": map[string]interface{}{
					"value": parsed.Palabras,
					"count": len(parsed.Palabras),
				},
			},
		},
	}

	// Agregar nodo de tiempo si existe
	if parsed.Fecha != "" || parsed.Hora != "" {
		tiempoNode := map[string]interface{}{
			"name": "TIEMPO",
			"children": []map[string]interface{}{},
		}
		
		if parsed.Fecha != "" {
			tiempoNode["children"] = append(tiempoNode["children"].([]map[string]interface{}), map[string]interface{}{
				"name": "FECHA",
				"attributes": map[string]interface{}{
					"value": parsed.Fecha,
					"type": getDateType(parsed.Fecha),
				},
			})
		}
		if parsed.Hora != "" {
			tiempoNode["children"] = append(tiempoNode["children"].([]map[string]interface{}), map[string]interface{}{
				"name": "HORA",
				"attributes": map[string]interface{}{
					"value": parsed.Hora,
					"format": "24h",
				},
			})
		}
		root["children"] = append(root["children"].([]map[string]interface{}), tiempoNode)
	}

	return root
}

func getDateType(fecha string) string {
	switch fecha {
	case "hoy", "mañana":
		return "relativa"
	case "lunes", "martes", "miércoles", "jueves", "viernes", "sábado", "domingo":
		return "día_semana"
	default:
		return "específica"
	}
}
