package controllers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/krasish/payment-system/internal/models"
)

type Transaction struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	UUID          string
	BelongsToUUID *string

	Type   string
	Status string
	Amount float64

	MerchantEmail string
	CustomerEmail string
	CustomerPhone string
}

// Notice that since I decided to "reverse" the relation direction in my implementation
// the logic has some differences with what is described in the task.
func (t *Transaction) getModelStatus(_type models.TransactionType, belongsToModel *models.Transaction) (models.TransactionStatus, error) {
	if belongsToModel != nil {
		if belongsToModel.Status != models.StatusApproved {
			return models.StatusError, nil
		}

		switch _type {
		case models.TypeAuthorize:
			if belongsToModel != nil {
				return "", errors.New("authorize transactions cannot have belongs to relations")
			}
		case models.TypeReversal:
			fallthrough
		case models.TypeCharge:
			if belongsToModel != nil && belongsToModel.Type != models.TypeAuthorize {
				return "", errors.New("transaction can only be related to authorize transaction")
			}
		case models.TypeRefund:
			if belongsToModel != nil && belongsToModel.Type != models.TypeCharge {
				return "", errors.New("refund transactions can only be reacted to charge transaction ")
			}
		}

	}
	return models.NewTransactionStatus(t.Status)
}

func (t *Transaction) toModel(merchantID uint, belongsToModel *models.Transaction) (*models.Transaction, error) {
	var (
		_type         models.TransactionType
		status        models.TransactionStatus
		belongsToID   *uint
		amount        = models.ToCurrency(t.Amount)
		customerEmail = t.CustomerEmail
		customerPhone = t.CustomerPhone
		err           error
	)

	_type, err = models.NewTransactionType(t.Type)
	if err != nil {
		return nil, err
	}
	status, err = t.getModelStatus(_type, belongsToModel)
	if err != nil {
		return nil, err
	}
	if belongsToModel != nil {
		belongsToID = &belongsToModel.ID
		amount = belongsToModel.Amount
		customerEmail = belongsToModel.CustomerEmail
		customerPhone = belongsToModel.CustomerPhone
	}

	return models.NewTransaction(t.UUID, amount, _type, status, customerEmail, customerPhone, merchantID, belongsToID)
}

func (t *Transaction) fromModel(model *models.Transaction) {
	var (
		belongsToUUID *string
	)
	if model.BelongsToID != nil && model.BelongsTo != nil {
		belongsToUUID = &model.BelongsTo.ExternalID
	}

	t.CreatedAt = model.CreatedAt
	t.UpdatedAt = model.UpdatedAt
	t.UUID = model.ExternalID
	t.BelongsToUUID = belongsToUUID
	t.Type = string(model.Type)
	t.Status = string(model.Status)
	t.Amount = model.Amount.Float64()
	t.MerchantEmail = model.Merchant.Email
	t.CustomerEmail = model.CustomerEmail
	t.CustomerPhone = model.CustomerPhone
}

type TransactionController struct {
	transactionStore *models.TransactionStore
	merchantStore    *models.MerchantStore
}

func NewTransactionController(transactionStore *models.TransactionStore, merchantStore *models.MerchantStore) *TransactionController {
	return &TransactionController{transactionStore: transactionStore, merchantStore: merchantStore}
}

func (c *TransactionController) CreateTransaction(ctx context.Context, t *Transaction) error {
	var (
		belongsToModel *models.Transaction
		merchant, err  = c.merchantStore.GetMerchantByEmail(ctx, t.MerchantEmail)
	)
	if err != nil {
		return fmt.Errorf("while getting merchant during transaciton creation: %w", err)
	}
	if t.BelongsToUUID != nil {
		belongsToModel, err = c.transactionStore.GetTransactionByUUID(ctx, *t.BelongsToUUID)
		if err != nil {
			return fmt.Errorf("while getting referenced transaction during transaciton creation: %w", err)
		}
	}
	model, err := t.toModel(merchant.UserID, belongsToModel)
	if err != nil {
		return err
	}
	return c.transactionStore.CreateTransaction(ctx, model)
}

func (c *TransactionController) GetTransactions(ctx context.Context) ([]*Transaction, error) {
	transactions, err := c.transactionStore.GetAllTransactions(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]*Transaction, len(transactions))
	for i := range transactions {
		res[i] = &Transaction{}
		res[i].fromModel(transactions[i])
	}
	return res, nil
}
