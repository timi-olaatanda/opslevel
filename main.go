package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"opslevel/todo"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	r := mux.NewRouter()
	r.HandleFunc("/todo/add", func(w http.ResponseWriter, r *http.Request) {
		queries := mux.Vars(r)
		priority, err := strconv.Atoi(queries["priority"])
		if err != nil {
			log.Printf("Error converting priority. Priority: %v, Form: %v, Error: %s\r\n", priority, queries, err)
			http.Error(w, "Priority must be specified and a non-zero positive integer", http.StatusBadRequest)
			return
		}
		description := queries["description"]
		var todoos []*todo.TodoItem
		if todoos, err = todo.AddTask(priority, description); err != nil {
			http.Error(w, fmt.Sprintf("%s", err), http.StatusBadRequest)
			return
		}
		w.Header().Add("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.Encode(&struct {
			Todoos []*todo.TodoItem `json:"samePriorityTodo"`
		}{
			todoos,
		})
	}).Methods("POST").Queries("priority", "{priority:[1-9][0-9]{0,7}}", "description", `{description:[A-Za-z][A-Za-z .]*}`)

	r.HandleFunc("/todo/remove", func(w http.ResponseWriter, r *http.Request) {
		priority, err := strconv.Atoi(mux.Vars(r)["priority"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if tasks, err := todo.RemoveTasks(priority); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else {
			w.Header().Add("Content-Type", "application/json")
			enc := json.NewEncoder(w)
			enc.Encode(&struct {
				Tasks interface{} `json:"removedtasks"`
			}{
				tasks,
			})
		}
	}).Methods("DELETE").Queries("priority", "{priority:[1-9][0-9]{0,7}}")

	r.HandleFunc("/todo/missing", func(w http.ResponseWriter, r *http.Request) {
		priorities := todo.GetMissingPriorities()
		w.Header().Add("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.Encode(&struct {
			Priorities []int `json:"priorities"`
		}{
			priorities,
		})
	}).Methods("GET")

	r.HandleFunc("/todo/all", func(w http.ResponseWriter, r *http.Request) {
		todoItems := todo.GetAllTodoItems()
		w.Header().Add("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.Encode(&struct {
			Todoos []*todo.TodoItem `json:"todoItems"`
		}{
			todoItems,
		})
	}).Methods("GET")

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "The specified url does not exist", http.StatusNotFound)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Handler:           r,
		Addr:              fmt.Sprintf(":%s", port),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
	}

	// todo(later): handle graceful shutdown
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
