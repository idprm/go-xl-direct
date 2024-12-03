package services

import (
	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/idprm/go-xl-direct/internal/domain/repository"
)

type VerifyService struct {
	verifyRepo repository.IVerifyRepository
}

type IVerifyService interface {
	Set(*entity.Verify) error
	Get(string) (*entity.Verify, error)
}

func NewVerifyService(verifyRepo repository.IVerifyRepository) *VerifyService {
	return &VerifyService{
		verifyRepo: verifyRepo,
	}
}

func (s *VerifyService) Set(t *entity.Verify) error {
	return s.verifyRepo.Set(t)
}

func (s *VerifyService) Get(trxId string) (*entity.Verify, error) {
	return s.verifyRepo.Get(trxId)
}
