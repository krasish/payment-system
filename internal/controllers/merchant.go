package controllers

import (
	"errors"
	"fmt"
	"time"

	"github.com/krasish/payment-system/internal/models"
)

type Merchant struct {
	CreatedAt time.Time
	UpdatedAt time.Time

	Name        string
	Description string
	Email       string
	Status      string
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

type MerchantController struct {
	store *models.MerchantStore
}

func NewMerchantController(store *models.MerchantStore) *MerchantController {
	return &MerchantController{store: store}
}

func (c *MerchantController) CreateMerchants(ms []*Merchant) error {
	modelMerchants := make([]*models.Merchant, 0, len(ms))
	for _, m := range ms {
		model, err := m.toModel()
		if err != nil {
			return fmt.Errorf("while converting merchants to model DTO: %w", err)
		}
		modelMerchants = append(modelMerchants, model)
	}
	return c.store.CreateMerchants(modelMerchants)
}
