package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yanadhiwiranata/go-todo/todo"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	dsn := "host=localhost user=hf dbname=gotodo port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.AutoMigrate(&todo.Todo{})
	if err != nil {
		fmt.Println("error open database: ", err)
	}
	r.Mount("/todo", todo.TodosResource{DB: db}.Routes())
	http.ListenAndServe(":3000", r)
}
