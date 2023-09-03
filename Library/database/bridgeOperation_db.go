package database

import (
	"github.com/EricBastos/ProjetoTG/Library/entities"
	"github.com/EricBastos/ProjetoTG/Library/utils"
	"gorm.io/gorm"
)

type BridgeOperationsDB struct {
	DB *gorm.DB
}

func NewBridgeOperationsDB(db *gorm.DB) *BridgeOperationsDB {
	return &BridgeOperationsDB{DB: db}
}

func (u *BridgeOperationsDB) Create(op *entities.BridgeOp) error {
	return u.DB.Create(op).Error
}

func (u *BridgeOperationsDB) CreateEmit(op *entities.BridgeOp, emitFunction func() error) error {
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

func (u *BridgeOperationsDB) Get(id string) (*entities.BridgeOp, error) {
	var op entities.BridgeOp
	err := u.DB.Preload("SmartContractOps").Preload("SmartContractOps.Feedback").First(&op, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &op, nil
}

func (u *BridgeOperationsDB) GetLogs(responsibleUser string, page, pageSize int) ([]entities.BridgeOpAPI, error) {
	var allBurnInfo []entities.BridgeOp

	err := u.DB.
		Preload("SmartContractOps", func(db *gorm.DB) *gorm.DB {
			return db.Order("smartcontract_operations.created_at DESC")
		}).
		Preload("SmartContractOps.Feedback").
		Scopes(utils.Paginate(page, pageSize)).
		Order("created_at DESC").
		Find(&allBurnInfo, "responsible_user = ?", responsibleUser).Error

	if err != nil {
		return nil, err
	}

	var res []entities.BridgeOpAPI

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

		newB := entities.BridgeOpAPI{
			Id:               b.Id,
			InputChain:       b.InputChain,
			OutputChain:      b.OutputChain,
			WalletAddress:    b.WalletAddress,
			Amount:           b.Amount,
			Permit:           b.Permit,
			SmartContractOps: sOps,
			CreatedAt:        b.CreatedAt,
		}

		res = append(res, newB)
	}

	if res == nil {
		return []entities.BridgeOpAPI{}, nil
	}

	return res, err
}
