package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/idprm/go-xl-direct/internal/domain/entity"
	"github.com/redis/go-redis/v9"
)

var (
	name string = "session"
)

type SessionRepository struct {
	rds *redis.Client
}

type ISessionRepository interface {
	Set(*entity.Session) error
	Get() (*entity.Session, error)
}

func NewSessionRepository(rds *redis.Client) *VerifyRepository {
	return &VerifyRepository{
		rds: rds,
	}
}

func (r *SessionRepository) Set(t *entity.Session) error {
	s, _ := json.Marshal(t)
	err := r.rds.Set(context.TODO(), name, string(s), 50*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *SessionRepository) Get() (*entity.Session, error) {
	val, err := r.rds.Get(context.TODO(), name).Result()
	if err != nil {
		return nil, err
	}
	var s *entity.Session
	json.Unmarshal([]byte(val), &s)
	return s, nil
}
