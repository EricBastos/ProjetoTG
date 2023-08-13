package database

import (
	"github.com/EricBastos/ProjetoTG/Library/entities"
	"github.com/EricBastos/ProjetoTG/Library/utils"
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

func (u *BurnOperationsDB) GetLogs(docData, responsibleUser string, page, pageSize int) ([]entities.BurnOpAPI, error) {
	var allBurnInfo []entities.BurnOp

	err := u.DB.
		Preload("SmartContractOps", func(db *gorm.DB) *gorm.DB {
			return db.Order("smartcontract_operations.created_at DESC")
		}).
		Preload("SmartContractOps.Feedback").
		Preload("Transfers", func(db *gorm.DB) *gorm.DB {
			return db.Order("transfers.created_at DESC")
		}).
		Preload("Transfers.Feedbacks", func(db *gorm.DB) *gorm.DB {
			return db.Order("transfer_feedbacks.created_at DESC")
		}).
		Scopes(utils.Paginate(page, pageSize)).
		Order("created_at DESC").
		Find(&allBurnInfo, "responsible_user = ? and user_tax_id = ?", responsibleUser, docData).Error

	if err != nil {
		return nil, err
	}

	var res []entities.BurnOpAPI

	for _, b := range allBurnInfo {

		var sOps []entities.SmartcontractOperationAPI

		for _, s := range b.SmartContractOps {

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

		var ts []entities.TransferAPI

		for _, t := range b.Transfers {

			var fs []entities.TransferFeedbackAPI

			for _, f := range t.Feedbacks {

				newF := entities.TransferFeedbackAPI{
					ID:        f.ID,
					CreatedAt: f.CreatedAt,
				}

				fs = append(fs, newF)
			}

			newT := entities.TransferAPI{
				Amount:        t.Amount,
				Name:          t.Name,
				TaxId:         t.TaxId,
				BankCode:      t.BankCode,
				BranchCode:    t.BranchCode,
				AccountNumber: t.AccountNumber,
				Id:            t.Id,
				CreatedAt:     t.CreatedAt,
				Feedbacks:     fs,
			}
			ts = append(ts, newT)
		}

		newB := entities.BurnOpAPI{
			Id:               b.Id,
			Chain:            b.Chain,
			WalletAddress:    b.WalletAddress,
			Amount:           b.Amount,
			Permit:           b.Permit,
			SmartContractOps: sOps,
			Transfers:        ts,
			CreatedAt:        b.CreatedAt,
		}

		res = append(res, newB)
	}

	if res == nil {
		return []entities.BurnOpAPI{}, nil
	}

	return res, err
}
