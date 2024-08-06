package services

import (
	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/idprm/go-xl-direct/internal/domain/repository"
)

type HistoryService struct {
	historyRepo repository.IHistoryRepository
}

type IHistoryService interface {
	SaveHistory(*entity.History) error
}

func NewHistoryService(historyRepo repository.IHistoryRepository) *HistoryService {
	return &HistoryService{
		historyRepo: historyRepo,
	}
}

func (s *HistoryService) SaveHistory(t *entity.History) error {
	err := s.historyRepo.Save(t)
	if err != nil {
		return err
	}
	return nil
}
