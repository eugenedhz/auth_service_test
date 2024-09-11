package postgres

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"github.com/eugenedhz/auth_service_test/internal/domain/models"
	"github.com/eugenedhz/auth_service_test/internal/repository"
)

type SessionRepository struct {
	db *gorm.DB
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Save(ctx context.Context, session *models.Session) error {
	err := r.db.Save(session).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *SessionRepository) Get(ctx context.Context, userID string) (*models.Session, error) {
	session := &models.Session{}

	err := r.db.Where(&models.Session{UserID: userID}).First(session).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return session, repository.ErrSessionNotFound
		}

		return session, err
	}

	return session, nil
}
