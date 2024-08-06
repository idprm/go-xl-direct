package services

import (
	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/idprm/go-xl-direct/internal/domain/repository"
)

type ContentService struct {
	contentRepo repository.IContentRepository
}

type IContentService interface {
	IsContent(int, string) bool
	IsContentSequence(int, int) bool
	GetContentName(int, string, string) (*entity.Content, error)
	GetContentSequence(int, int, string) (*entity.Content, error)
}

func NewContentService(contentRepo repository.IContentRepository) *ContentService {
	return &ContentService{
		contentRepo: contentRepo,
	}
}

func (s *ContentService) IsContent(serviceId int, name string) bool {
	count, _ := s.contentRepo.Count(serviceId, name)
	return count > 0
}

func (s *ContentService) IsContentSequence(serviceId, sequence int) bool {
	count, _ := s.contentRepo.CountBySequence(serviceId, sequence)
	return count > 0
}

func (s *ContentService) GetContentName(serviceId int, name string, pin string) (*entity.Content, error) {
	result, err := s.contentRepo.Get(serviceId, name)
	if err != nil {
		return nil, err
	}

	var content entity.Content

	if result != nil {
		content = entity.Content{
			Value:    result.Value,
			Tid:      result.Tid,
			Sequence: result.Sequence,
		}
		content.SetPIN(pin)
	}
	return &content, nil
}

func (s *ContentService) GetContentSequence(serviceId, sequence int, pin string) (*entity.Content, error) {
	result, err := s.contentRepo.GetBySequence(serviceId, sequence)
	if err != nil {
		return nil, err
	}

	var content entity.Content

	if result != nil {
		content = entity.Content{
			Value:    result.Value,
			Tid:      result.Tid,
			Sequence: result.Sequence,
		}
		content.SetPIN(pin)
	}
	return &content, nil
}
