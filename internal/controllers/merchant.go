package controllers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/krasish/payment-system/internal/models"
)

type Merchant struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	Name                string
	Description         string
	Email               string
	Status              string
	TotalTransactionSum float64
}

func (m *Merchant) CSVUnmarshal(record []string) error {
	if len(record) != 4 {
		return fmt.Errorf("wrong number of records in merchant CSV record: %v", record)
	} else if m == nil {
		return errors.New("cannot unmarshal CSV to a nil merchant")
	}
	m.Name = record[0]
	m.Description = record[1]
	m.Email = record[2]
	m.Status = record[3]
	return nil
}

func (m *Merchant) toModel() (*models.Merchant, error) {
	status, err := models.NewUserStatus(m.Status)
	if err != nil {
		return nil, err
	}
	return models.NewMerchant(m.Name, m.Description, m.Email, status)
}

func (m *Merchant) fromModel(model *models.Merchant) {
	m.CreatedAt = model.User.CreatedAt
	m.UpdatedAt = model.User.UpdatedAt
	m.Name = model.Name
	m.Description = model.Description
	m.Email = model.Email
	m.Status = string(model.User.Status)
	m.TotalTransactionSum = model.TotalTransactionSum.Float64()
}

type MerchantController struct {
	store *models.MerchantStore
}

func NewMerchantController(store *models.MerchantStore) *MerchantController {
	return &MerchantController{store: store}
}

func (c *MerchantController) CreateMerchants(ctx context.Context, ms []*Merchant) error {
	modelMerchants := make([]*models.Merchant, 0, len(ms))
	for _, m := range ms {
		model, err := m.toModel()
		if err != nil {
			return fmt.Errorf("while converting merchants to model in create merchant: %w", err)
		}
		modelMerchants = append(modelMerchants, model)
	}
	return c.store.CreateMerchants(ctx, modelMerchants)
}

func (c *MerchantController) GetMerchantByMail(ctx context.Context, email string) (*Merchant, error) {
	merchantModel, err := c.store.GetMerchantByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	m := new(Merchant)
	m.fromModel(merchantModel)
	return m, nil
}

func (c *MerchantController) GetMerchants(ctx context.Context) ([]*Merchant, error) {
	merchants, err := c.store.GetAllMerchants(ctx)
	if err != nil {
		return nil, err
	}
	res := make([]*Merchant, len(merchants))
	for i := range merchants {
		res[i] = &Merchant{}
		res[i].fromModel(merchants[i])
	}
	return res, nil
}

func (c *MerchantController) UpdateMerchant(ctx context.Context, merchant *Merchant) error {
	model, err := merchant.toModel()
	if err != nil {
		return fmt.Errorf("while converting merchant to model in update merchant")
	}
	if err = c.store.UpdateMerchant(ctx, model); err != nil {
		return err
	}
	return nil
}

func (c *MerchantController) DeleteMerchant(ctx context.Context, merchantEmail string) error {
	if err := c.store.DeleteMerchant(ctx, merchantEmail); err != nil {
		return err
	}
	return nil
}
