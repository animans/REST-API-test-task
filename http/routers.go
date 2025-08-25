package http

import "github.com/gorilla/mux"

// Register ...
func Register(r *mux.Router, h *Handlers) {
	api := r.NewRoute().Subrouter()
	api.HandleFunc("/service", h.Create).Methods("POST")
	api.HandleFunc("/service/{id}", h.Get).Methods("GET")
	api.HandleFunc("/service/{id}", h.Put).Methods("PUT")
	api.HandleFunc("/service/{id}", h.Delete).Methods("DELETE")
	api.HandleFunc("/service", h.List).Methods("GET")
}
