package repositories

import (
	"time"

	"crosshatch/internal/database/models"

	"gorm.io/gorm"
)

type SessionRepository struct {
	db *gorm.DB
}

func (r *SessionRepository) CreateSession(session *models.Session) error {
	return r.db.Create(session).Error
}

func (r *SessionRepository) GetSession(tokenHash string) (*models.Session, error) {
	session := models.Session{}
	err := r.db.First(&session, "token_hash = ?", tokenHash).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &session, nil
}

func (r *SessionRepository) DeleteSession(tokenHash string) error {
	return r.db.Delete(&models.Session{}, "token_hash = ?", tokenHash).Error
}

func (r *SessionRepository) DeleteExpired() error {
	return r.db.Delete(&models.Session{}, "expires_at < ?", time.Now()).Error
}

func (r *SessionRepository) DeleteSessionsForUser(userID uint) error {
	return r.db.Delete(&models.Session{}, "user_id = ?", userID).Error
}

func NewSessionRepository(db *gorm.DB) *SessionRepository {
	return &SessionRepository{db: db}
}
