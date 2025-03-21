package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/AdamBrutsaert/basic-go-http-server/internal/mux"
	intstore "github.com/AdamBrutsaert/basic-go-http-server/internal/store"
)

func getItems(store *intstore.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items := store.GetItems()
		err := json.NewEncoder(w).Encode(items)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func getItem(store *intstore.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
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

func addItem(store *intstore.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var item intstore.Item
		err := json.NewDecoder(r.Body).Decode(&item)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		store.AddItem(item)
		w.WriteHeader(http.StatusCreated)
	}
}

func updateItem(store *intstore.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var item intstore.Item
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

func deleteItem(store *intstore.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
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

func newItemsRouter(store *intstore.Store) *mux.MiddlewareMux {
	mux := mux.NewMiddlewareMux(loggingMiddleware)
	mux.HandleFunc("GET /", getItems(store))
	mux.HandleFunc("GET /{id}", getItem(store))
	mux.HandleFunc("POST /", addItem(store))
	mux.HandleFunc("PUT /{id}", updateItem(store))
	mux.HandleFunc("DELETE /{id}", deleteItem(store))
	return mux
}

func newRouter(store *intstore.Store) *mux.PrefixMux {
	mux := mux.NewPrefixMux()
	mux.Handle("/items", newItemsRouter(store))
	return mux
}
