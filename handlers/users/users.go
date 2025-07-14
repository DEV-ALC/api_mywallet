package users

import (
	"encoding/json"
	"fmt"
	"meuprojeto/db"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Usuario struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

func AutenticarUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	var login LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		http.Error(w, "Erro ao ler JSON", http.StatusBadRequest)
		return
	}

	var usuario Usuario
	var hash string

	err := db.DB.QueryRow(`
		SELECT id, username, email, password_hash
		FROM users
		WHERE username = ? AND access_enabled = 1 AND deleted_at IS NULL
	`, login.Username).Scan(&usuario.ID, &usuario.Username, &usuario.Email, &hash)

	if err != nil {
		http.Error(w, "Usuário não encontrado ou desativado", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(login.Password)); err != nil {
		http.Error(w, "Senha incorreta", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuario)
}

func UsersAtivo(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
		SELECT id, username, email
		FROM users
		WHERE access_enabled = true AND deleted_at = '0000-00-00 00:00:00'
	`)
	if err != nil {
		http.Error(w, "Erro na query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var username, email string
		if err := rows.Scan(&id, &username, &email); err != nil {
			http.Error(w, "Erro no scan", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Usuário %d: %s (%s)\n", id, username, email)
	}
}
