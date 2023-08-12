package database

import (
	"github.com/EricBastos/ProjetoTG/Library/entities"
	"gorm.io/gorm"
)

type SmartcontractOperationDB struct {
	DB *gorm.DB
}

func NewSmartcontractOperationDB(db *gorm.DB) *SmartcontractOperationDB {
	return &SmartcontractOperationDB{DB: db}
}

func (u *SmartcontractOperationDB) Create(op *entities.SmartcontractOperation) error {
	return u.DB.Create(op).Error
}
