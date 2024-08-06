package services

import (
	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/idprm/go-xl-direct/internal/domain/repository"
)

type ServiceService struct {
	serviceRepo repository.IServiceRepository
}

type IServiceService interface {
	IsServiceByCode(string) bool
	GetServiceId(int) (*entity.Service, error)
	GetServiceByCode(string) (*entity.Service, error)
}

func NewServiceService(serviceRepo repository.IServiceRepository) *ServiceService {
	return &ServiceService{
		serviceRepo: serviceRepo,
	}
}

func (s *ServiceService) IsServiceByCode(code string) bool {
	count, _ := s.serviceRepo.CountByCode(code)
	return count > 0
}

func (s *ServiceService) GetServiceId(id int) (*entity.Service, error) {
	return s.serviceRepo.GetById(id)
}

func (s *ServiceService) GetServiceByCode(code string) (*entity.Service, error) {
	return s.serviceRepo.GetByCode(code)
}
