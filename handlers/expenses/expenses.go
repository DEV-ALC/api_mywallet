package expenses

import (
	"encoding/json"
	"meuprojeto/db"
	"net/http"
)

// Struct com campos extras pra sync
type Expense struct {
	Id          int     `json:"id"`
	Description string  `json:"description"`
	Amount      float32 `json:"amount"`
	Tag_id      int     `json:"tag_id,omitempty"` // usado em outras rotas
	UserID      int     `json:"user_id,omitempty"`
	CreatedAt   string  `json:"created_at,omitempty"`
	UpdatedAt   string  `json:"updated_at,omitempty"`
	DeletedAt   string  `json:"deleted_at,omitempty"`
}

// Buscar todas as despesas não deletadas
func BuscarDespesas(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
		SELECT id, description, amount, tag_id, user_id
		FROM expenses
		WHERE deleted_at IS NULL
	`)
	if err != nil {
		http.Error(w, "Erro na query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var expenses []Expense
	for rows.Next() {
		var expense Expense
		if err := rows.Scan(&expense.Id, &expense.Description, &expense.Amount, &expense.Tag_id, &expense.UserID); err != nil {
			http.Error(w, "Erro no scan", http.StatusInternalServerError)
			return
		}
		expenses = append(expenses, expense)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(expenses)
}

// Buscar despesas de um usuário específico
func BuscarDespesaUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "ID do usuário é obrigatório", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query(`
		SELECT id, description, amount, tag_id, user_id
		FROM expenses
		WHERE user_id = ? AND deleted_at IS NULL
	`, userID)
	if err != nil {
		http.Error(w, "Erro na query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var expenses []Expense
	for rows.Next() {
		var expense Expense
		if err := rows.Scan(&expense.Id, &expense.Description, &expense.Amount, &expense.Tag_id, &expense.UserID); err != nil {
			http.Error(w, "Erro no scan", http.StatusInternalServerError)
			return
		}
		expenses = append(expenses, expense)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(expenses)
}

// Buscar despesas modificadas após a data `last`
func DespesasSync(w http.ResponseWriter, r *http.Request) {
	last := r.URL.Query().Get("last")
	user_id := r.URL.Query().Get("id")

	if user_id == "" || last == "" || last == "0000-00-00 00:00:00" {
		last = "1970-01-01 00:00:00"
	}

	rows, err := db.DB.Query(`
		SELECT id, description, amount, updated_at, deleted_at
		FROM expenses
		WHERE user_id = ?
		AND (
			updated_at > ? OR (deleted_at IS NOT NULL AND deleted_at > ?)
		)
	`, user_id, last, last)
	if err != nil {
		http.Error(w, "Erro ao buscar despesas", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var expenses []Expense
	for rows.Next() {
		var expense Expense
		if err := rows.Scan(&expense.Id, &expense.Description, &expense.Amount, &expense.UpdatedAt, &expense.DeletedAt); err != nil {

			http.Error(w, "Erro ao processar dados", http.StatusInternalServerError)
			return
		}
		expenses = append(expenses, expense)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(expenses)
}
