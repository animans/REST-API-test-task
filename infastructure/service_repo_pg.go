package infastructure

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/animans/REST-API-test-task/domain"
	"github.com/google/uuid"
	_ "github.com/lib/pq" // ...
)

// ServiceRepoPG ...
type ServiceRepoPG struct {
	db *sql.DB
}

// NewServiceRepoPG ...
func NewServiceRepoPG() *ServiceRepoPG {
	return &ServiceRepoPG{}
}

// Open ...
func (r *ServiceRepoPG) Open() error {
	env, ok := os.LookupEnv("database_url")
	if !ok {
		env = "host=localhost user=baish password=postgres port=5432 dbname=REST-API-task-test_test sslmode=disable"
	}
	db, err := sql.Open("postgres", env)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	r.db = db

	return nil
}

// Close ...
func (r *ServiceRepoPG) Close() error {
	return r.db.Close()
}

// Save ...
func (r *ServiceRepoPG) Save(s *domain.Service) (int, error) {
	var id int

	if err := r.db.QueryRow(
		"INSERT INTO service_list (service_price, service_name, service_uuid, service_created_at) VALUES ($1, $2, $3, $4) RETURNING service_id",
		s.GetPrice(), s.GetName(), s.GetUUID(), s.GetStartDate(),
	).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

// repoService ...
type repoService struct {
	Name  string
	Price int
	Uuid  uuid.UUID
	Date  time.Time
}

// GetByID ...
func (r *ServiceRepoPG) GetByID(sid string) (*domain.Service, error) {
	var in repoService

	id, err := strconv.Atoi(sid)
	if err != nil {
		return &domain.Service{}, err
	}
	if err := r.db.QueryRow(
		"SELECT service_name, service_price, service_uuid, service_created_at FROM service_list WHERE service_id=$1",
		id,
	).Scan(&in.Name, &in.Price, &in.Uuid, &in.Date); err != nil {
		return &domain.Service{}, err
	}

	return domain.NewService(in.Name, in.Price, in.Uuid, in.Date), nil
}

// UpdateByID ...
func (r *ServiceRepoPG) UpdateByID(sid string, in *domain.Service) error {
	id, err := strconv.Atoi(sid)
	if err != nil {
		return err
	}
	res, err := r.db.Exec(
		"UPDATE service_list SET service_name=$1, service_price=$2, service_uuid=$3, service_created_at=$4 WHERE service_id=$5",
		in.GetName(), in.GetPrice(), in.GetUUID().String(), in.GetStartDate(),
		id,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("no rows updated, id=%d", id)
	}

	return nil
}

// DeleteByID ...
func (r *ServiceRepoPG) DeleteByID(sid string) error {
	id, err := strconv.Atoi(sid)
	if err != nil {
		return err
	}
	res, err := r.db.Exec(
		"DELETE FROM service_list WHERE service_id=$1",
		id,
	)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("no rows deleted, id=%d", id)
	}

	return nil
}

func (r *ServiceRepoPG) ListByFilter(domain.ListFilterService) (domain.ListResult, error) {

}
