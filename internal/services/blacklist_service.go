package services

import "github.com/idprm/go-xl-direct/internal/domain/repository"

type BlacklistService struct {
	blacklistRepo repository.IBlacklistRepository
}

type IBlacklistService interface {
	IsBlacklist(msisdn string) bool
}

func NewBlacklistService(blacklistRepo repository.IBlacklistRepository) *BlacklistService {
	return &BlacklistService{
		blacklistRepo: blacklistRepo,
	}
}

func (s *BlacklistService) IsBlacklist(msisdn string) bool {
	count, _ := s.blacklistRepo.Count(msisdn)
	return count > 0
}
