package repositories

import (
    "eventmaster-go/internal/models"
    "time"

    "gorm.io/gorm"
)

// SessionRepository defines data access methods for sessions
type SessionRepository interface {
    Create(session *models.Session) error
    FindByToken(token string) (*models.Session, error)
    DeleteByToken(token string) error
    DeleteExpired(before time.Time) error
}

type sessionRepository struct {
    db *gorm.DB
}

// NewSessionRepository creates a new session repository
func NewSessionRepository(db *gorm.DB) SessionRepository {
    return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(session *models.Session) error {
    return r.db.Create(session).Error
}

func (r *sessionRepository) FindByToken(token string) (*models.Session, error) {
    var session models.Session
    if err := r.db.Preload("User").First(&session, "token = ?", token).Error; err != nil {
        return nil, err
    }
    return &session, nil
}

func (r *sessionRepository) DeleteByToken(token string) error {
    return r.db.Where("token = ?", token).Delete(&models.Session{}).Error
}

func (r *sessionRepository) DeleteExpired(before time.Time) error {
    return r.db.Where("expires_at < ?", before).Delete(&models.Session{}).Error
}
