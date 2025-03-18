package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
)

func getItems(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items := store.GetItems()
		err := json.NewEncoder(w).Encode(items)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func getItem(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		item, ok := store.GetItem(id)
		if !ok {
			http.Error(w, "item not found", http.StatusNotFound)
			return
		}

		err = json.NewEncoder(w).Encode(item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func addItem(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var item Item
		err := json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		store.AddItem(item)
		w.WriteHeader(http.StatusCreated)
	}
}

func updateItem(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var item Item
		err = json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ok := store.UpdateItem(id, item)
		if !ok {
			http.Error(w, "item not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func deleteItem(store *Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ok := store.DeleteItem(id)
		if !ok {
			http.Error(w, "item not found", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s %s\n", r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func main() {
	store := NewStore()

	r := chi.NewRouter()
	r.Use(loggingMiddleware)
	r.Route("/items", func(r chi.Router) {
		r.Get("/", getItems(store))
		r.Post("/", addItem(store))
		r.Get("/{id}", getItem(store))
		r.Put("/{id}", updateItem(store))
		r.Delete("/{id}", deleteItem(store))
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
