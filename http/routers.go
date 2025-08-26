package http

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Register ...
func Register(r *mux.Router, h *Handlers) {
	api := r.NewRoute().Subrouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	api.HandleFunc("/service", h.Create).Methods("POST")
	api.HandleFunc("/service", h.List).Methods("GET")
	api.HandleFunc("/service/summary", h.ListSum).Methods("GET")
	api.HandleFunc("/service/{id}", h.Get).Methods("GET")
	api.HandleFunc("/service/{id}", h.Put).Methods("PUT")
	api.HandleFunc("/service/{id}", h.Delete).Methods("DELETE")

}
