package database

import (
	"errors"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	"gorm.io/gorm"
	"time"
)

type StaticDepositDB struct {
	DB *gorm.DB
}

func NewStaticDepositDB(db *gorm.DB) *StaticDepositDB {
	return &StaticDepositDB{DB: db}
}

func (u *StaticDepositDB) Create(deposit *entities.StaticDeposit) error {
	return u.DB.Create(deposit).Error
}

func (u *StaticDepositDB) Update(updatedDeposit *entities.StaticDeposit) error {
	return u.DB.Save(updatedDeposit).Error
}

func (u *StaticDepositDB) FindById(id string) (*entities.StaticDeposit, error) {
	var deposit entities.StaticDeposit
	if err := u.DB.Where("id = ?", id).First(&deposit).Error; err != nil {
		return nil, err
	} else if deposit.Id == nil {
		return nil, errors.New("static deposit not found")
	}
	return &deposit, nil
}

func (u *StaticDepositDB) FindUnpaidByTaxIdAndAmount(taxId string, amount int) (*entities.StaticDeposit, error) {
	var deposit entities.StaticDeposit
	err := u.DB.
		Where("tax_id = ? and amount = ? and status = ? and due > ?", taxId, amount, "UNPAID", time.Now()).
		Order("created_at ASC").
		Limit(1).
		Find(&deposit).
		Error
	if err != nil {
		return nil, err
	} else if deposit.Id == nil {
		return nil, errors.New("entry not found")
	}
	return &deposit, nil
}

func (u *StaticDepositDB) FindUnpaidByTaxId(taxId string) (*entities.StaticDeposit, error) {
	var deposit entities.StaticDeposit
	err := u.DB.
		Where("tax_id = ?  and status = ? and due > ?", taxId, "UNPAID", time.Now()).
		Order("created_at ASC").
		Limit(1).
		Find(&deposit).
		Error
	if err != nil {
		return nil, err
	} else if deposit.Id == nil {
		return nil, errors.New("entry not found")
	}
	return &deposit, nil
}
