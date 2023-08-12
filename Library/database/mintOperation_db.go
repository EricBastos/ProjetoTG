package database

import (
	"github.com/EricBastos/ProjetoTG/Library/entities"
	"gorm.io/gorm"
)

type MintOperationsDB struct {
	DB *gorm.DB
}

func NewMintOperationsDB(db *gorm.DB) *MintOperationsDB {
	return &MintOperationsDB{DB: db}
}

func (u *MintOperationsDB) Create(op *entities.MintOp) error {
	return u.DB.Create(op).Error
}

func (u *MintOperationsDB) CreateEmit(op *entities.MintOp, emitFunction func() error) error {
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

func (u *MintOperationsDB) Get(id string) (*entities.MintOp, error) {
	var op entities.MintOp
	err := u.DB.Preload("SmartContractOps").Preload("SmartContractOps.Feedback").First(&op, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &op, nil
}
