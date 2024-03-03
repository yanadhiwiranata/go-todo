package todo

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type TodosResource struct {
	DB *gorm.DB
}

type Todo struct {
	ID        int    `gorm:"primaryKey" json:"id"`
	Title     string `gorm:"index" json:"title"`
	Done      bool   `gorm:"index" json:"done"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
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
	})

	return r
}

func (rs TodosResource) List(w http.ResponseWriter, r *http.Request) {
	var todos []Todo
	result := rs.DB.Find(&todos)

	if result.Error != nil {
		HttpResponse(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	HttpResponse(w, todos, http.StatusOK)
}

func (rs TodosResource) Create(w http.ResponseWriter, r *http.Request) {
	var requestBody TodoRequest
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		HttpResponse(w, err.Error(), http.StatusBadRequest)
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
		HttpResponse(w, "todos create error", http.StatusInternalServerError)
		return
	}

	HttpResponse(w, []byte("todos create"), http.StatusOK)
}

func (rs TodosResource) Get(w http.ResponseWriter, r *http.Request) {
	ids := chi.URLParam(r, "id")
	if ids == "" {
		HttpResponse(w, GenericResponse{ResponseMessage: "Invalid id"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(ids)
	if err != nil {
		HttpResponse(w, GenericResponse{ResponseMessage: "Invalid id"}, http.StatusBadRequest)
		return
	}

	todo := &Todo{ID: id}
	result := rs.DB.First(&todo)

	if result.Error != nil {
		HttpResponse(w, GenericResponse{ResponseMessage: "Todo not found"}, http.StatusInternalServerError)
		return
	}

	HttpResponse(w, todo, http.StatusOK)
}

func (rs TodosResource) Update(w http.ResponseWriter, r *http.Request) {
	ids := chi.URLParam(r, "id")
	if ids == "" {
		HttpResponse(w, GenericResponse{ResponseMessage: "Invalid id"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(ids)
	if err != nil {
		HttpResponse(w, GenericResponse{ResponseMessage: "Invalid id"}, http.StatusBadRequest)
		return
	}

	var requestBody TodoRequest
	err = json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		HttpResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(requestBody.Title) == 0 {
		HttpResponse(w, GenericResponse{ResponseMessage: "Title is required"}, http.StatusBadRequest)
		return
	}

	todo := &Todo{ID: id}
	result := rs.DB.First(&todo)

	if result.Error != nil {
		HttpResponse(w, GenericResponse{ResponseMessage: "Todo not found"}, http.StatusNotFound)
		return
	}

	todo.Title = requestBody.Title
	todo.Done = requestBody.Done

	rs.DB.Save(&todo)

	HttpResponse(w, todo, http.StatusOK)
}

func (rs TodosResource) Delete(w http.ResponseWriter, r *http.Request) {
	ids := chi.URLParam(r, "id")
	if ids == "" {
		HttpResponse(w, GenericResponse{ResponseMessage: "Invalid id"}, http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(ids)
	if err != nil {
		HttpResponse(w, GenericResponse{ResponseMessage: "Invalid id"}, http.StatusBadRequest)
		return
	}

	todo := &Todo{ID: id}
	result := rs.DB.First(&todo)

	if result.Error != nil {
		HttpResponse(w, GenericResponse{ResponseMessage: "Todo not found"}, http.StatusNotFound)
		return
	}

	rs.DB.Delete(&todo)

	HttpResponse(w, todo, http.StatusOK)
}

func HttpResponse(w http.ResponseWriter, body interface{}, code int) {
	if code == 0 {
		code = http.StatusOK
	}

	if body != nil {
		w.Header().Set("Content-Type", "application/json")
		bytes, err := json.Marshal(body)
		if err != nil {
			response := GenericResponse{ResponseMessage: "Internal Server Error"}
			bytes, _ = json.Marshal(response)
		}
		w.Write(bytes)
	}
	w.WriteHeader(code)
}

type GenericResponse struct {
	ResponseMessage string `json:"response_message"`
}
