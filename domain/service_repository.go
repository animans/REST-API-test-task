package domain

// ServiceRepository ...
type ServiceRepository interface {
	Save(s *Service) (int, error)
	GetByID(id string) (*Service, error)
	UpdateByID(sid string, s *Service) error
	DeleteByID(sid string) error
	ListByFilter(ListFilterService) (ListResult, error)
	SumByFilter(SumFilterService) (SumResult, error)
}
