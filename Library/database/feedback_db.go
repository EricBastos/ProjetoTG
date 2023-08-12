package database

import (
	"errors"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	"gorm.io/gorm"
)

type FeedbackDB struct {
	DB *gorm.DB
}

func NewFeedbackDB(db *gorm.DB) *FeedbackDB {
	return &FeedbackDB{DB: db}
}

func (u *FeedbackDB) Create(account *entities.Feedback) error {
	return u.DB.Create(account).Error
}

func (u *FeedbackDB) FindById(id string) (*entities.Feedback, error) {
	var feedback entities.Feedback
	if err := u.DB.Where("id = ?", id).First(&feedback).Error; err != nil {
		return nil, err
	} else if feedback.ID == nil {
		return nil, errors.New("mint feedback not found")
	}
	return &feedback, nil
}
