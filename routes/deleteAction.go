package routes

import (
	"net/http"
	"strconv"

	"github.com/RodrigoGonzalez78/go_analyzer/db"
	"github.com/RodrigoGonzalez78/go_analyzer/models"
)

func DeleteAction(w http.ResponseWriter, r *http.Request) {
	claim, _ := r.Context().Value("userData").(*models.Claim)

	idStr := r.PathValue("id")
	actionID64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || actionID64 == 0 {
		http.Error(w, "ID de accion inválido", http.StatusBadRequest)
		return
	}
	actionID := uint(actionID64)

	action, err := db.GetActionByID(actionID)
	if err != nil {
		http.Error(w, "Error al obtener la accion: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if action.UserName != claim.UserName {
		http.Error(w, "No tienes permiso para eliminar esta accion", http.StatusForbidden)
		return
	}

	err = db.DeleteActionByID(actionID)
	if err != nil {
		http.Error(w, "Error al eliminar la accion: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
