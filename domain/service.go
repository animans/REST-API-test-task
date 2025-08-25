package domain

import (
	"time"

	"github.com/google/uuid"
)

// Service ...
type Service struct {
	name      string
	price     int
	uuid      uuid.UUID
	startDate time.Time
}

// ListFilterService ...
type ListFilterService struct {
	Name          string
	Price         int
	Uuid          uuid.UUID
	FromStartDate time.Time
	ToStartDate   time.Time
	SortBy        string
	SortDir       string
	Limit         int
}

// CreatedRequest ...
type CreatedRequest struct {
	Name      string `json:"service_name"`
	Price     int    `json:"price"`
	Uuid      string `json:"user_id"`
	StartDate string `json:"start_date"`
}

// ListResult ...
type ListResult struct {
	Items []CreatedRequest
}

// NewService ...
func NewService(sn string, sp int, uuid uuid.UUID, sd time.Time) *Service {
	return &Service{
		name:      sn,
		price:     sp,
		uuid:      uuid,
		startDate: sd,
	}
}

// GetName ...
func (s *Service) GetName() string {
	return s.name
}

// GetPrice
func (s *Service) GetPrice() int {
	return s.price
}

// GetUUID
func (s *Service) GetUUID() uuid.UUID {
	return s.uuid
}

// GetStartDate
func (s *Service) GetStartDate() time.Time {
	return s.startDate
}
