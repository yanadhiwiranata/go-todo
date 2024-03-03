package todo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type TodosResource struct {
	DB *gorm.DB
}

// gorm.Model definition
type Todo struct {
	gorm.Model
	Title string `gorm:"index"`
	Done  bool   `gorm:"index"`
}

type TodoRequest struct {
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

// Routes creates a REST router for the todos resource
func (rs TodosResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Get("/", rs.List)    // GET /todos - read a list of todos
	r.Post("/", rs.Create) // POST /todos - create a new todo and persist it
	r.Delete("/", rs.Delete)

	r.Route("/{id}", func(r chi.Router) {
		// r.Use(rs.TodoCtx) // lets have a todos map, and lets actually load/manipulate
		r.Get("/", rs.Get)       // GET /todos/{id} - read a single todo by :id
		r.Put("/", rs.Update)    // PUT /todos/{id} - update a single todo by :id
		r.Delete("/", rs.Delete) // DELETE /todos/{id} - delete a single todo by :id
		r.Get("/sync", rs.Sync)
	})

	return r
}

func (rs TodosResource) List(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todos list of stuff.."))
}

func (rs TodosResource) Create(w http.ResponseWriter, r *http.Request) {
	var requestBody TodoRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(requestBody.Title) == 0 {
		w.Write([]byte("Title  is required"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	todo := &Todo{Title: requestBody.Title, Done: requestBody.Done}

	result := rs.DB.Create(todo)
	if result.Error != nil {
		fmt.Println(result.Error.Error())
		w.Write([]byte("todos create error"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write([]byte("todos create"))
}

func (rs TodosResource) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo get"))
}

func (rs TodosResource) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo update"))
}

func (rs TodosResource) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo delete"))
}

func (rs TodosResource) Sync(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("todo sync"))
}

type GenericResponse struct {
	ResponseMessage string `json:"response_message"`
}
