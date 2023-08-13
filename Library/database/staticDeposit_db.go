package database

import (
	"errors"
	"github.com/EricBastos/ProjetoTG/Library/entities"
	"github.com/EricBastos/ProjetoTG/Library/utils"
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

func (u *StaticDepositDB) GetLogs(taxId, responsibleUser string, page, pageSize int) ([]entities.StaticDepositAPI, error) {
	var allDepositInfo []entities.StaticDeposit

	err := u.DB.
		Preload("MintOps", func(db *gorm.DB) *gorm.DB {
			return db.Order("mint_ops.created_at DESC")
		}).
		Preload("MintOps.SmartContractOps", func(db *gorm.DB) *gorm.DB {
			return db.Order("smartcontract_operations.created_at DESC")
		}).
		Preload("MintOps.SmartContractOps.Feedback").
		Scopes(utils.Paginate(page, pageSize)).
		Order("created_at DESC").
		Find(&allDepositInfo, "responsible_user = ? and tax_id = ?", responsibleUser, taxId).Error

	if err != nil {
		return nil, err
	}

	var res []entities.StaticDepositAPI

	for _, b := range allDepositInfo {

		var mOps []entities.MintOpAPI

		for _, m := range b.MintOps {

			var sOps []entities.SmartcontractOperationAPI

			for _, s := range m.SmartContractOps {

				var f *entities.FeedbackAPI

				if s.Feedback != nil {
					f = &entities.FeedbackAPI{
						ID:                       s.Feedback.ID,
						SmartcontractOperationId: s.Feedback.SmartcontractOperationId,
						Success:                  s.Feedback.Success,
						ErrorMsg:                 s.Feedback.ErrorMsg,
						CreatedAt:                s.Feedback.CreatedAt,
					}
				}

				newS := entities.SmartcontractOperationAPI{
					ID:            s.ID,
					OperationName: s.OperationName,
					Executed:      s.Executed,
					Tx:            s.Tx,
					Reason:        s.Reason,
					CreatedAt:     s.CreatedAt,
					Feedback:      f,
				}
				sOps = append(sOps, newS)
			}

			newM := entities.MintOpAPI{
				Id:               m.Id,
				Amount:           m.Amount,
				Reason:           m.Reason,
				CreatedAt:        m.CreatedAt,
				SmartContractOps: sOps,
			}
			mOps = append(mOps, newM)
		}

		newB := entities.StaticDepositAPI{
			Chain:         b.Chain,
			WalletAddress: b.WalletAddress,
			Amount:        b.Amount,
			TaxId:         b.TaxId,
			Due:           b.Due,
			Id:            b.Id,
			CreatedAt:     b.CreatedAt,
			Status:        b.Status,
			UpdatedAt:     b.UpdatedAt,
			MintOps:       mOps,
		}

		res = append(res, newB)
	}

	if res == nil {
		return []entities.StaticDepositAPI{}, nil
	}

	return res, err
}
