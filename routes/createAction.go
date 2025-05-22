package routes

import (
	"encoding/json"
	"net/http"

	"github.com/RodrigoGonzalez78/go_analyzer/analyzer"
	"github.com/RodrigoGonzalez78/go_analyzer/db"
	"github.com/RodrigoGonzalez78/go_analyzer/models"
)

func CreateAction(w http.ResponseWriter, r *http.Request) {

	claim, _ := r.Context().Value("userData").(*models.Claim)

	type Request struct {
		Comand string `json:"comand" `
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

	action, err := analyzer.CreateAction(comand.Comand)

	if err != nil {
		http.Error(w, "Error analisando el comando", http.StatusInternalServerError)
		return
	}

	action.UserID = claim.UserID

	err = db.CreateAction(action)
	if err != nil {
		http.Error(w, "Error creando la accion", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
