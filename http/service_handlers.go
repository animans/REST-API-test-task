package http

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/animans/REST-API-test-task/domain"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Handlers ...
type Handlers struct {
	Repo domain.ServiceRepository
}

// NewHandlers ...
func NewHandlers(repo domain.ServiceRepository) *Handlers {
	return &Handlers{
		Repo: repo,
	}
}

// CreatedResponseID
type CreatedResponseID struct {
	ID int `json:"id"`
}

// CreatedResponse
type CreatedResponse struct {
	Name      string `json:"service_name"`
	Price     int    `json:"price"`
	Uuid      string `json:"user_id"`
	StartDate string `json:"start_date"`
}

// Start ...
func (h *Handlers) Start() error {
	env, ok := os.LookupEnv("bind_addr")
	if !ok {
		env = "8080"
	}
	router := mux.NewRouter()
	Register(router, h)
	log.Println("Starting api")
	return http.ListenAndServe(env, router)
}

// Create ...
func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	var in domain.CreatedRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(in.Name) == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	if in.Price < 0 {
		http.Error(w, "price must be >= 0", http.StatusBadRequest)
		return
	}

	sdate, err := time.Parse("01-2006", in.StartDate)
	if err != nil {
		http.Error(w, "invalid time", http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(in.Uuid)
	if err != nil {
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	ser := domain.NewService(in.Name, in.Price, uuid, sdate)
	id, err := h.Repo.Save(ser)
	if err != nil {
		http.Error(w, "save error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	out := CreatedResponseID{ID: id}
	_ = json.NewEncoder(w).Encode(out)
}

// Get ...
func (h *Handlers) Get(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	ser, err := h.Repo.GetByID(id)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	out := CreatedResponse{
		Name:      ser.GetName(),
		Price:     ser.GetPrice(),
		Uuid:      ser.GetUUID().String(),
		StartDate: ser.GetStartDate().Format("01-2006"),
	}
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(out)
}

// Put ...
func (h *Handlers) Put(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var in domain.CreatedRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(in.Name) == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	if in.Price < 0 {
		http.Error(w, "price must be >= 0", http.StatusBadRequest)
		return
	}

	sdate, err := time.Parse("01-2006", in.StartDate)
	if err != nil {
		http.Error(w, "invalid date (want MM-YYYY)", http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(in.Uuid)
	if err != nil {
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	ser := domain.NewService(in.Name, in.Price, uuid, sdate)
	if err := h.Repo.UpdateByID(id, ser); err != nil {
		http.Error(w, "update error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Delete ...
func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if err := h.Repo.DeleteByID(id); err != nil {
		http.Error(w, "delete error", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusNoContent)
}

// List ...
func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var f domain.ListFilterService

	if name := q.Get("name"); name != "" {
		f.Name = name
	}

	if s := q.Get("user_id"); s != "" {
		user_id, err := uuid.Parse(s)
		if err != nil {
			http.Error(w, "bad user_id", http.StatusBadRequest)
			return
		}
		f.Uuid = &user_id
	}

	if s := q.Get("price"); s != "" {
		price, err := strconv.Atoi(s)
		if err != nil {
			http.Error(w, "bad price", http.StatusBadRequest)
			return
		}
		f.Price = price
	}

	if s := q.Get("from"); s != "" {
		fromStartDate, err := time.Parse("01-2006", s)
		if err != nil {
			http.Error(w, "bad fromDate (MM-YYY)", http.StatusBadRequest)
			return
		}
		f.FromStartDate = &fromStartDate
	}
	if s := q.Get("to"); s != "" {
		toStartDate, err := time.Parse("01-2006", s)
		if err != nil {
			http.Error(w, "bad toDate (MM-YYY)", http.StatusBadRequest)
			return
		}
		f.ToStartDate = &toStartDate
	}
	s := strings.ToLower(q.Get("sort"))
	switch s {
	case "service_created_at", "service_price", "service_name":
		f.SortBy = s
	default:
		f.SortBy = "service_created_at"
	}
	s = strings.ToLower(q.Get("dir"))
	switch s {
	case "asc", "desc":
		f.SortDir = s
	default:
		f.SortDir = "desc"
	}
	if l := q.Get("limit"); l != "" {
		n, _ := strconv.Atoi(l)
		if n < 1 {
			n = 1
		}
		if n > 100 {
			n = 100
		}
		f.Limit = n
	} else {
		f.Limit = 50
	}

	res, err := h.Repo.ListByFilter(f)
	if err != nil {
		http.Error(w, "internal err", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}
