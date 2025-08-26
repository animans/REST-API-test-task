package http

import (
	"encoding/json"
	"log/slog"
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
	env, ok := os.LookupEnv("BIND_ADDR")
	if !ok {
		env = "8080"
	}
	router := mux.NewRouter()
	Register(router, h)
	slog.Info("Starting api")
	return http.ListenAndServe(env, router)
}

// Create
// @Summary      Create service
// @Description  Создать запись подписки
// @Tags         service
// @Accept       json
// @Produce      json
// @Param        input body     domain.CreatedRequest true "service payload"
// @Success      201   {object} domain.CreatedRequest
// @Failure      400   {string} string "bad request"
// @Failure      500   {string} string "internal error"
// @Router       /service [post]
func (h *Handlers) Create(w http.ResponseWriter, r *http.Request) {
	slog.Info("Create start")
	var in domain.CreatedRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		slog.Error("invalid json", "err", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(in.Name) == "" {
		slog.Error("invalid name", "name", in.Name)
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	if in.Price < 0 {
		slog.Error("invalid price", "price", in.Price)
		http.Error(w, "price must be >= 0", http.StatusBadRequest)
		return
	}

	sdate, err := time.Parse("01-2006", in.StartDate)
	if err != nil {
		slog.Error("invalid sdate", "err", err)
		http.Error(w, "invalid time", http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(in.Uuid)
	if err != nil {
		slog.Error("invalid uuid", "uuid", uuid)
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	ser := domain.NewService(in.Name, in.Price, uuid, sdate)
	id, err := h.Repo.Save(ser)
	if err != nil {
		slog.Error("invalid id", "err", err)
		http.Error(w, "save error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	out := CreatedResponseID{ID: id}
	_ = json.NewEncoder(w).Encode(out)
	slog.Info("Create done", "out", out)
}

// Get
// @Summary      Get service by ID
// @Tags         service
// @Produce      json
// @Param        id   path integer true "Service ID" format(integer)
// @Success      200  {object} domain.CreatedRequest
// @Failure      404  {string} string "not found"
// @Router       /service/{id} [get]
func (h *Handlers) Get(w http.ResponseWriter, r *http.Request) {
	slog.Info("Get start", "mux.Vars(r)", mux.Vars(r))
	id := mux.Vars(r)["id"]
	ser, err := h.Repo.GetByID(id)
	if err != nil {
		slog.Error("invalid id", "err", err)
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
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
	slog.Info("Get done", "out", out)
}

// Update
// @Summary      Update service
// @Tags         service
// @Accept       json
// @Param        id    path  integer                 true "Service ID" format(integer)
// @Param        input body  domain.CreatedRequest  true "update payload"
// @Success      204
// @Failure      400   {string} string "bad request"
// @Failure      404   {string} string "not found"
// @Router       /service/{id} [put]
func (h *Handlers) Put(w http.ResponseWriter, r *http.Request) {
	slog.Info("Put start", "mux.Vars(r)", mux.Vars(r))
	id := mux.Vars(r)["id"]
	var in domain.CreatedRequest
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		slog.Error("invalid json", "err", err)
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(in.Name) == "" {
		slog.Error("invalid name", "name", in.Name)
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	if in.Price < 0 {
		slog.Error("invalid price", "price", in.Price)
		http.Error(w, "price must be >= 0", http.StatusBadRequest)
		return
	}

	sdate, err := time.Parse("01-2006", in.StartDate)
	if err != nil {
		slog.Error("invalid sdate", "err", err)
		http.Error(w, "invalid date (want MM-YYYY)", http.StatusBadRequest)
		return
	}

	uuid, err := uuid.Parse(in.Uuid)
	if err != nil {
		slog.Error("invalid uuid", "err", err)
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	ser := domain.NewService(in.Name, in.Price, uuid, sdate)
	if err := h.Repo.UpdateByID(id, ser); err != nil {
		slog.Error("update error", "err", err)
		http.Error(w, "update error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	slog.Info("Put done")
}

// Delete
// @Summary      Delete service
// @Tags         service
// @Param        id path integer true "Service ID" format(integer)
// @Success      204 {string} string "deleted"
// @Failure      404 {string} string "not found"
// @Router       /service/{id} [delete]
func (h *Handlers) Delete(w http.ResponseWriter, r *http.Request) {
	slog.Info("Delete start", "mux.Vars(r)", mux.Vars(r))
	id := mux.Vars(r)["id"]
	if err := h.Repo.DeleteByID(id); err != nil {
		slog.Error("delete error", "err", err)
		http.Error(w, "delete error", http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusNoContent)
	slog.Info("Delete done")
}

// List
// @Summary      List services
// @Tags         service
// @Produce      json
// @Param        name    query string false "filter by service name (contains)"
// @Param        user_id query string false "User UUID" format(uuid)
// @Param        price   query string false "Price"
// @Param        from    query string false "From month (MM-YYYY)" example(01-2024)
// @Param        to      query string false "To month   (MM-YYYY)" example(03-2024)
// @Param        sort    query string false "sort by (service_created_at, service_price, service_name)"
// @Param        dir     query string false "order by (asc, desc)"
// @Param        limit   query string false "limit  (1 <= limit <= 100)" example(50)
// @Success      200 {object} domain.ListResult
// @Router       /service [get]
func (h *Handlers) List(w http.ResponseWriter, r *http.Request) {
	slog.Info("List start", "r.URL.Query()", r.URL.Query())
	q := r.URL.Query()
	var f domain.ListFilterService

	if name := q.Get("name"); name != "" {
		f.Name = name
	}

	if s := q.Get("user_id"); s != "" {
		user_id, err := uuid.Parse(s)
		if err != nil {
			slog.Error("invalid user_id", "err", err)
			http.Error(w, "bad user_id", http.StatusBadRequest)
			return
		}
		f.Uuid = &user_id
	}

	if s := q.Get("price"); s != "" {
		price, err := strconv.Atoi(s)
		if err != nil {
			slog.Error("invalid price", "err", err)
			http.Error(w, "bad price", http.StatusBadRequest)
			return
		}
		f.Price = price
	}

	if s := q.Get("from"); s != "" {
		fromStartDate, err := time.Parse("01-2006", s)
		if err != nil {
			slog.Error("invalid fromStartDate", "err", err)
			http.Error(w, "bad fromDate (MM-YYY)", http.StatusBadRequest)
			return
		}
		f.FromStartDate = &fromStartDate
	}
	if s := q.Get("to"); s != "" {
		toStartDate, err := time.Parse("01-2006", s)
		if err != nil {
			slog.Error("invalid ToStartDate", "err", err)
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
		slog.Error("invalid res", "err", err)
		http.Error(w, "internal err", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
	slog.Info("List done", "res", res)
}

// Summary
// @Summary      Sum price by period
// @Description  Суммарная стоимость подписок за период с фильтрами
// @Tags         service
// @Produce      json
// @Param        name    query string false "service name (contains)"
// @Param        user_id query string false "User UUID" format(uuid)
// @Param        from    query string false "From month (MM-YYYY)" example(01-2024)
// @Param        to      query string false "To month   (MM-YYYY)" example(03-2024)
// @Success      200 {object} domain.SumResult
// @Failure      400 {string} string "bad request"
// @Router       /service/summary [get]
func (h *Handlers) ListSum(w http.ResponseWriter, r *http.Request) {
	slog.Info("ListSum start", "r.URL.Query()", r.URL.Query())
	var f domain.SumFilterService

	q := r.URL.Query()

	if s := q.Get("name"); s != "" {
		f.Name = s
	}
	if s := q.Get("user_id"); s != "" {
		uuid, err := uuid.Parse(s)
		if err != nil {
			slog.Error("invalid uuid", "err", err)
			http.Error(w, "bad user_id", http.StatusBadRequest)
			return
		}
		f.Uuid = &uuid
	}
	if s := q.Get("from"); s != "" {
		fromStartDate, err := time.Parse("01-2006", s)
		if err != nil {
			slog.Error("invalid fromStartDate", "err", err)
			http.Error(w, "bad fromDate (MM-YYYY)", http.StatusBadRequest)
			return
		}
		f.FromStartDate = &fromStartDate
	}
	if s := q.Get("to"); s != "" {
		toStartDate, err := time.Parse("01-2006", s)
		if err != nil {
			slog.Error("invalid toStartDate", "err", err)
			http.Error(w, "bad toDate (MM-YYYY)", http.StatusBadRequest)
			return
		}
		f.ToStartDate = &toStartDate
	}

	out, err := h.Repo.SumByFilter(f)
	if err != nil {
		slog.Error("invalid out", "err", err)
		http.Error(w, "internal err", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(out)
	slog.Info("ListSum", "out", out)
}
