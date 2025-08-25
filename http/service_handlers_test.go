package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/animans/REST-API-test-task/domain"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func mustJSON(t *testing.T, v any) *bytes.Buffer {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return bytes.NewBuffer(b)
}

func wantStatus(t *testing.T, rec *httptest.ResponseRecorder, code int) {
	t.Helper()
	if rec.Code != code {
		t.Fatalf("status: got=%d want=%d body=%q", rec.Code, code, rec.Body.String())
	}
}

func wantBodyContains(t *testing.T, rec *httptest.ResponseRecorder, substr string) {
	t.Helper()
	if !strings.Contains(rec.Body.String(), substr) {
		t.Fatalf("body doesn't contain %q; body=%q", substr, rec.Body.String())
	}
}

func wantBodyEmpty(t *testing.T, rec *httptest.ResponseRecorder) {
	t.Helper()
	if rec.Body.Len() != 0 {
		t.Fatalf("expected empty body, got: %q", rec.Body.String())
	}
}

type fakeRepo struct {
	saved   *domain.Service
	saveErr error
}

func (f *fakeRepo) GetByID(id string) (*domain.Service, error) {
	sdate, err := time.Parse("01-2006", "08-2025")
	if err != nil {
		return nil, err
	}

	fService := map[string]*domain.Service{
		"1": domain.NewService(
			"Yandex Plus",
			400,
			uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			sdate,
		),
	}
	fser, ok := fService[id]
	if !ok {
		f.saveErr = errors.New("db invalid id")
		return &domain.Service{}, f.saveErr
	}
	return fser, nil
}

func (f *fakeRepo) Save(s *domain.Service) (int, error) {
	f.saved = s
	return 1, f.saveErr
}

func (f *fakeRepo) UpdateByID(id string, s *domain.Service) error {
	fService := map[string]*domain.Service{
		"1": f.saved,
	}
	_, ok := fService[id]
	if !ok {
		f.saveErr = errors.New("db invalid id")
		return f.saveErr
	}
	f.saved = s
	return nil
}

func (f *fakeRepo) DeleteByID(id string) error {
	fService := map[string]*domain.Service{
		"1": f.saved,
	}
	_, ok := fService[id]
	if !ok {
		f.saveErr = errors.New("db invalid id")
		return f.saveErr
	}
	delete(fService, id)
	f.saved = nil
	return nil
}

func TestCreate(t *testing.T) {
	frepo := &fakeRepo{}
	h := NewHandlers(frepo)
	id := uuid.MustParse("00000000-0000-0000-0000-000000000001")

	load := struct {
		Name      string `json:"service_name"`
		Price     int    `json:"price"`
		Uuid      string `json:"user_id"`
		StartDate string `json:"start_date"`
	}{
		Name:      "Yandex Plus",
		Price:     500,
		Uuid:      id.String(),
		StartDate: "03-2012",
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/service", mustJSON(t, load))
	req.Header.Set("Content-Type", "application/json")

	h.Create(rec, req)
	wantStatus(t, rec, http.StatusCreated)
	wantBodyContains(t, rec, "1")

	if frepo.saved == nil {
		t.Fatalf("repo.Save was not called")
	}
	if frepo.saved.GetName() != load.Name {
		t.Fatalf("saved.Name: got=%q want=%q", frepo.saved.GetName(), load.Name)
	}
	if frepo.saved.GetPrice() != load.Price {
		t.Fatalf("saved.Price: got=%d want=%d", frepo.saved.GetPrice(), load.Price)
	}
	if frepo.saved.GetUUID() != id { // предполагаю поле UserID в domain.Service
		t.Fatalf("saved.UserID: got=%s want=%s", frepo.saved.GetUUID(), id)
	}
	if frepo.saved.GetStartDate().Format("01-2006") != load.StartDate {
		t.Fatalf("saved.StartDate: got=%s want %s", frepo.saved.GetStartDate().Format("01-2006"), load.StartDate)
	}
}

func TestGet(t *testing.T) {
	sdate, err := time.Parse("01-2006", "08-2025")
	if err != nil {
		t.Fatal(err)
	}

	frepo := &fakeRepo{
		saved: domain.NewService(
			"Yandex Plus",
			400,
			uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			sdate,
		),
	}

	want := struct {
		Name      string `json:"service_name"`
		Price     int    `json:"price"`
		Uuid      string `json:"user_id"`
		StartDate string `json:"start_date"`
	}{
		Name:      frepo.saved.GetName(),
		Price:     frepo.saved.GetPrice(),
		Uuid:      frepo.saved.GetUUID().String(),
		StartDate: frepo.saved.GetStartDate().Format("01-2006"),
	}

	h := &Handlers{Repo: frepo}
	rec := httptest.NewRecorder()

	req := httptest.NewRequest(http.MethodGet, "/service/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	h.Get(rec, req)

	wantStatus(t, rec, 200)

	var got struct {
		Name      string `json:"service_name"`
		Price     int    `json:"price"`
		Uuid      string `json:"user_id"`
		StartDate string `json:"start_date"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v; body=%s", err, rec.Body.String())
	}

	if got != want {
		t.Fatalf("mismatch:\n got: %#v\nwant: %#v", got, want)
	}
}

func TestPut(t *testing.T) {
	sdate, err := time.Parse("01-2006", "08-2025")
	if err != nil {
		t.Fatal(err)
	}

	frepo := &fakeRepo{
		saved: domain.NewService(
			"Yandex Plus",
			400,
			uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			sdate,
		),
	}
	load := struct {
		Name      string `json:"service_name"`
		Price     int    `json:"price"`
		Uuid      string `json:"user_id"`
		StartDate string `json:"start_date"`
	}{
		Name:      "GPT Plus",
		Price:     500,
		Uuid:      "00000000-0000-0000-0000-000000000002",
		StartDate: "05-2025",
	}
	h := NewHandlers(frepo)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/service/1", mustJSON(t, load))
	req = mux.SetURLVars(req, map[string]string{"id": "1"})

	h.Put(rec, req)

	wantStatus(t, rec, 204)
	wantBodyEmpty(t, rec)
	switch {
	case frepo.saved.GetName() != load.Name:
		t.Error("update name error")
	case frepo.saved.GetPrice() != load.Price:
		t.Error("update price error")
	case frepo.saved.GetUUID().String() != load.Uuid:
		t.Error("update uuid error")
	case frepo.saved.GetStartDate().Format("01-2006") != load.StartDate:
		t.Error("update start_date error")
	}
}

func TestDelete(t *testing.T) {
	sdate, err := time.Parse("01-2006", "08-2025")
	if err != nil {
		t.Fatal(err)
	}
	frepoo := &fakeRepo{
		saved: domain.NewService(
			"Yandex Plus",
			400,
			uuid.MustParse("00000000-0000-0000-0000-000000000001"),
			sdate,
		),
	}
	h := NewHandlers(frepoo)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/service/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	h.Delete(rec, req)

	wantStatus(t, rec, 204)
	wantBodyEmpty(t, rec)
	if frepoo.saved != nil {
		t.Fatal("delete error")
	}
}

// func TestCreate(t *testing.T) {
// 	casetest := []struct {
// 		name       string
// 		body       string
// 		wantStatus int
// 	}{
// 		{"ok", `"ID":1`, 201},
// 		{"bad_json", `{bad}`, 400},
// 	}
// 	for _, c := range casetest {
// 		t.Run(c.name, func(string))
// 	}
// 	repo := infastructure.NewServiceRepoPG()
// 	err := repo.Open()
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	s := NewHandlers(repo)

// 	uuid, _ := uuid.NewRandom()
// 	ser := domain.NewService("Yandex Plus", 400, uuid, time.Now())
// 	json, err := json.Marshal(ser)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	rec := httptest.NewRecorder()
// 	req, _ := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(json))
// 	s.Create(rec, req)

// }
