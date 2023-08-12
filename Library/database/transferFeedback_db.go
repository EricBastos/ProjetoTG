package database

import (
	"errors"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	"gorm.io/gorm"
)

type TransferFeedbackDB struct {
	DB *gorm.DB
}

func NewTransferFeedbackDB(db *gorm.DB) *TransferFeedbackDB {
	return &TransferFeedbackDB{DB: db}
}

func (u *TransferFeedbackDB) Create(feedback *entities.TransferFeedback) error {
	return u.DB.Create(feedback).Error
}

func (u *TransferFeedbackDB) FindById(id string) (*entities.TransferFeedback, error) {
	var feedback entities.TransferFeedback
	if err := u.DB.Where("id = ?", id).First(&feedback).Error; err != nil {
		return nil, err
	} else if feedback.ID == nil {
		return nil, errors.New("transfer not found")
	}
	return &feedback, nil
}
