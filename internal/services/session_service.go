package services

import (
	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/idprm/go-xl-direct/internal/domain/repository"
)

type SessionService struct {
	sessionRepo repository.ISessionRepository
}

type ISessionService interface {
	Set(*entity.Session) error
	Get(string) (*entity.Session, error)
}

func NewSessionService(sessionRepo repository.ISessionRepository) *SessionService {
	return &SessionService{
		sessionRepo: sessionRepo,
	}
}

func (s *SessionService) Set(t *entity.Session) error {
	return s.sessionRepo.Set(t)
}

func (s *SessionService) Get() (*entity.Session, error) {
	return s.sessionRepo.Get()
}
