package database

import (
	"errors"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	"gorm.io/gorm"
)

type StaticDepositFeedbackDB struct {
	DB *gorm.DB
}

func NewStaticDepositFeedbackDB(db *gorm.DB) *StaticDepositFeedbackDB {
	return &StaticDepositFeedbackDB{DB: db}
}

func (u *StaticDepositFeedbackDB) Create(feedback *entities.StaticDepositFeedback) error {
	return u.DB.Create(feedback).Error
}

func (u *StaticDepositFeedbackDB) FindById(id string) (*entities.StaticDepositFeedback, error) {
	var feedback entities.StaticDepositFeedback
	if err := u.DB.Where("id = ?", id).First(&feedback).Error; err != nil {
		return nil, err
	} else if feedback.ID == "" {
		return nil, errors.New("static deposit feedback not found")
	}
	return &feedback, nil
}
