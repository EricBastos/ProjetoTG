package database

import (
	"errors"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	"gorm.io/gorm"
)

type TransferDB struct {
	DB *gorm.DB
}

func NewTransferDB(db *gorm.DB) *TransferDB {
	return &TransferDB{DB: db}
}

func (u *TransferDB) Create(transfer *entities.Transfer) error {
	return u.DB.Create(transfer).Error
}

func (u *TransferDB) FindById(id string) (*entities.Transfer, error) {
	var transfer entities.Transfer
	if err := u.DB.Where("id = ?", id).First(&transfer).Error; err != nil {
		return nil, err
	} else if transfer.Id == "" {
		return nil, errors.New("transfer not found")
	}
	return &transfer, nil
}
