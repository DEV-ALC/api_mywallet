package tags

import (
	"encoding/json"
	"meuprojeto/db"
	"net/http"
)

type Tag struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	UserID int    `json:"user_id"`
}

func BuscarTags(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
		SELECT id, name, user_id
		FROM tags
		WHERE deleted_at = '0000-00-00 00:00:00'
	`)
	if err != nil {
		http.Error(w, "Erro na query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tags []Tag

	for rows.Next() {
		var tag Tag
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.UserID); err != nil {
			http.Error(w, "Erro no scan", http.StatusInternalServerError)
			return
		}
		tags = append(tags, tag)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}
