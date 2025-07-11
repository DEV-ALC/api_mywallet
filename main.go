package main

import (
	"fmt"
	"log"
	"meuprojeto/db"
	"meuprojeto/handlers/expenses"
	"meuprojeto/handlers/tags"
	"meuprojeto/handlers/users"
	"net/http"
)

func main() {
	db.Connect()

	//despesas
	http.HandleFunc("/despesas", expenses.BuscarDespesas)
	http.HandleFunc("/despesas/usuario", expenses.BuscarDespesaUser)
	http.HandleFunc("/despesas/usuario/sync", expenses.DespesasSync)

	//etiquetas
	http.HandleFunc("/tags", tags.BuscarTags)

	//usuarios
	http.HandleFunc("/usuarios/ativos", users.UsersAtivo)

	// Iniciar o servidor na porta 8080
	fmt.Println("🚀 Servidor rodando em http://localhost:81")
	log.Fatal(http.ListenAndServe(":81", nil))
}
