package database

import (
	"github.com/EricBastos/ProjetoTG/Library/entities"
	"gorm.io/gorm"
)

type BurnOperationsDB struct {
	DB *gorm.DB
}

func NewBurnOperationsDB(db *gorm.DB) *BurnOperationsDB {
	return &BurnOperationsDB{DB: db}
}

func (u *BurnOperationsDB) Create(op *entities.BurnOp) error {
	return u.DB.Create(op).Error
}

func (u *BurnOperationsDB) CreateEmit(op *entities.BurnOp, emitFunction func() error) error {
	err := u.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(op).Error; err != nil {
			return err
		}
		if err := emitFunction(); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (u *BurnOperationsDB) Get(id string) (*entities.BurnOp, error) {
	var op entities.BurnOp
	err := u.DB.Preload("SmartContractOps").Preload("SmartContractOps.Feedback").First(&op, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &op, nil
}
