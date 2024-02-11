package repository

import (
	"fmt"
	"github.com/PRC-36/amikompedia-fiber/domain/entity"
	"gorm.io/gorm"
	"log"
)

type SessionRepository interface {
	SessionCreate(tx *gorm.DB, value *entity.Session) error
	FindByID(tx *gorm.DB, sessionID string) (*entity.Session, error)
}

type sessionRepositoryImpl struct {
}

func NewSessionRepository() SessionRepository {
	return &sessionRepositoryImpl{}
}

func (s *sessionRepositoryImpl) SessionCreate(tx *gorm.DB, value *entity.Session) error {
	result := tx.Create(value)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when create session : %v", result.Error))
		return result.Error
	}

	return nil
}

func (s *sessionRepositoryImpl) FindByID(tx *gorm.DB, sessionID string) (*entity.Session, error) {
	var session entity.Session

	result := tx.Where("id = ?", sessionID).First(&session)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find session by id : %v", result.Error))
		return nil, result.Error
	}

	return &session, nil

}
