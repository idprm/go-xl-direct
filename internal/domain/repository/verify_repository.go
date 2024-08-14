package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/redis/go-redis/v9"
)

type VerifyRepository struct {
	rds *redis.Client
}

type IVerifyRepository interface {
	Set(*entity.Verify) error
	Get(string) (*entity.Verify, error)
}

func NewVerifyRepository(rds *redis.Client) *VerifyRepository {
	return &VerifyRepository{
		rds: rds,
	}
}

func (r *VerifyRepository) Set(t *entity.Verify) error {
	verify, _ := json.Marshal(t)
	err := r.rds.Set(context.TODO(), t.GetMsisdn(), string(verify), 1*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *VerifyRepository) Get(msisdn string) (*entity.Verify, error) {
	val, err := r.rds.Get(context.TODO(), msisdn).Result()
	if err != nil {
		return nil, err
	}
	var v *entity.Verify
	json.Unmarshal([]byte(val), &v)
	return v, nil
}
