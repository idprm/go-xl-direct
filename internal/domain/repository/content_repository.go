package repository

import (
	"database/sql"
	"log"

	"github.com/idprm/go-xl-direct/internal/domain/entity"
)

const (
	queryCountContentByName      = "SELECT COUNT(*) as count FROM contents WHERE service_id = $1 AND name = $2"
	queryCountContentBySequence  = "SELECT COUNT(*) as count FROM contents WHERE service_id = $1 AND sequence = $2"
	querySelectContentByName     = "SELECT value, tid, sequence FROM contents WHERE service_id = $1 AND name = $2 LIMIT 1"
	querySelectContentBySequence = "SELECT value, tid, sequence FROM contents WHERE service_id = $1 AND sequence = $2 LIMIT 1"
)

type ContentRepository struct {
	db *sql.DB
}

type IContentRepository interface {
	Count(int, string) (int, error)
	CountBySequence(int, int) (int, error)
	Get(int, string) (*entity.Content, error)
	GetBySequence(int, int) (*entity.Content, error)
}

func NewContentRepository(db *sql.DB) *ContentRepository {
	return &ContentRepository{
		db: db,
	}
}

func (r *ContentRepository) Count(serviceId int, name string) (int, error) {
	var count int
	err := r.db.QueryRow(queryCountContentByName, serviceId, name).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (r *ContentRepository) CountBySequence(serviceId, sequence int) (int, error) {
	var count int
	err := r.db.QueryRow(queryCountContentBySequence, serviceId, sequence).Scan(&count)
	if err != nil {
		log.Println(err)
		return count, err
	}
	return count, nil
}

func (r *ContentRepository) Get(serviceId int, name string) (*entity.Content, error) {
	var content entity.Content
	err := r.db.QueryRow(querySelectContentByName, serviceId, name).Scan(&content.Value, &content.Tid, &content.Sequence)
	if err != nil {
		log.Println(err)
		return &content, err
	}
	return &content, nil
}

func (r *ContentRepository) GetBySequence(serviceId, sequence int) (*entity.Content, error) {
	var content entity.Content
	err := r.db.QueryRow(querySelectContentBySequence, serviceId, sequence).Scan(&content.Value, &content.Tid, &content.Sequence)
	if err != nil {
		log.Println(err)
		return &content, err
	}
	return &content, nil
}
