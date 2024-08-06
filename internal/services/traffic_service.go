package services

import (
	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/idprm/go-xl-direct/internal/domain/repository"
)

type TrafficService struct {
	trafficRepo repository.ITrafficRepository
}

type ITrafficService interface {
	SaveCampaign(*entity.TrafficCampaign) error
	SaveMO(*entity.TrafficMO) error
	UpdateMOCharge(*entity.TrafficMO) error
}

func NewTrafficService(trafficRepo repository.ITrafficRepository) *TrafficService {
	return &TrafficService{
		trafficRepo: trafficRepo,
	}
}

func (s *TrafficService) SaveCampaign(t *entity.TrafficCampaign) error {
	err := s.trafficRepo.SaveCampaign(t)
	if err != nil {
		return err
	}
	return nil
}

func (s *TrafficService) SaveMO(t *entity.TrafficMO) error {
	err := s.trafficRepo.SaveMO(t)
	if err != nil {
		return err
	}
	return nil
}

func (s *TrafficService) UpdateMOCharge(t *entity.TrafficMO) error {
	err := s.trafficRepo.UpdateMOCharge(t)
	if err != nil {
		return err
	}
	return nil
}
