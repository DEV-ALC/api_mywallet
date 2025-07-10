package expenses

import (
	"encoding/json"
	"fmt"
	"meuprojeto/db"
	"net/http"
)

type Expense struct {
	Id          int     `json:"id"`
	Description string  `json:"description"`
	Amount      float32 `json:"amount"`
	Tag_id      int     `json:"tag_id"`
	UserID      int     `json:"user_id"`
}

func BuscarDespesas(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
		SELECT id, description, amount, tag_id, user_id
		FROM expenses
		WHERE deleted_at = '0000-00-00 00:00:00'
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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(expenses)

}

func BuscarDespesaUser(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("id")
	if userID == "" {
		http.Error(w, "ID do usuário é obrigatório", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query(`
		SELECT id, description, amount
		FROM expenses
		WHERE user_id = ? AND deleted_at = '0000-00-00 00:00:00'
	`, userID)
	if err != nil {
		http.Error(w, "Erro na query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var desc string
		var amount float64
		if err := rows.Scan(&id, &desc, &amount); err != nil {
			http.Error(w, "Erro no scan", http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Despesa %d: %s - R$%.2f\n", id, desc, amount)
	}
}

func DespesasSync(w http.ResponseWriter, r *http.Request) {
	last := r.URL.Query().Get("last")
	if last == "" {
		http.Error(w, "Parâmetro 'last' é obrigatório", http.StatusBadRequest)
		return
	}

	rows, err := db.DB.Query(`
		SELECT id, description, amount, updated_at, deleted_at
		FROM expenses
		WHERE (updated_at > ? OR deleted_at > ?)`, last, last)
	if err != nil {
		http.Error(w, "Erro na query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var desc string
		var amount float64
		var updatedAt, deletedAt *string

		if err := rows.Scan(&id, &desc, &amount, &updatedAt, &deletedAt); err != nil {
			http.Error(w, "Erro no scan", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Despesa %d: %s - R$%.2f (updated: %v, deleted: %v)\n",
			id, desc, amount, updatedAt, deletedAt)
	}
}
