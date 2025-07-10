package users

import (
	"fmt"
	"meuprojeto/db"
	"net/http"
)

func UsersAtivo(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
		SELECT id, username, email
		FROM users
		WHERE access_enabled = true AND deleted_at IS NULL
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
