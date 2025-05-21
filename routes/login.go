package routes

import (
	"encoding/json"
	"net/http"

	"github.com/RodrigoGonzalez78/go_analyzer/db"
	"github.com/RodrigoGonzalez78/go_analyzer/models"
	"github.com/RodrigoGonzalez78/go_analyzer/utils"
)

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")

	var t models.User

	err := json.NewDecoder(r.Body).Decode(&t)

	if err != nil {
		http.Error(w, "Usuario y/o contraseña invalidas"+err.Error(), 400)
		return
	}

	if len(t.UserName) == 0 {
		http.Error(w, "El nombre de usuario es requerido", 400)
		return
	}

	if len(t.Password) == 0 {
		http.Error(w, "La contraseña es requerida", 400)
		return
	}

	user, err := db.GetUserByUserName(t.UserName)

	if err != nil {
		http.Error(w, "No se encontro el usuario", 400)
		return
	}

	passwordValid := utils.CheckPassword(user.Password, t.Password)

	if !passwordValid {
		http.Error(w, "Contraseña invalida", 400)
		return
	}

	jwtKey, err := utils.GenerateJWT(t.UserName)

	if err != nil {
		http.Error(w, "Ocurrio un error"+err.Error(), 400)
		return
	}

	resp := models.ResponseLogin{
		Token: jwtKey,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
