package routes

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/RodrigoGonzalez78/go_analyzer/db"
	"github.com/RodrigoGonzalez78/go_analyzer/models"
)

func GetAllUserActions(w http.ResponseWriter, r *http.Request) {
	claim, _ := r.Context().Value("userData").(*models.Claim)

	query := r.URL.Query()

	pageStr := query.Get("page")
	pageSizeStr := query.Get("pageSize")

	page := 1
	pageSize := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	actions, err := db.GetUserActionsPaginated(claim.UserID, page, pageSize)
	if err != nil {
		http.Error(w, "Error al obtener las acciones", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actions)
}
