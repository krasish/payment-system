package models

import (
	"fmt"
	"net/mail"

	"gorm.io/gorm"
)

type Merchant struct {
	UserID uint `gorm:"primaryKey"`
	User   User

	Name                string
	TotalTransactionSum Currency `gorm:"-"`
	Description         string
	Email               string

	Transactions []Transaction `gorm:"->;foreignKey:MerchantID;references:UserID"`
}

func NewMerchant(name string, description string, email string, status UserStatus) (*Merchant, error) {
	_, err := mail.ParseAddress(email)
	if err != nil {
		return nil, fmt.Errorf("while creating merchant: %q is not a valid email address: %w", email, err)
	}
	return &Merchant{Name: name, Description: description, Email: email, User: User{
		Role:   RoleMerchant,
		Status: status,
	}}, nil
}

func (m *Merchant) calculateTTS() {
	m.TotalTransactionSum = Currency(0)
	for _, t := range m.Transactions {
		if t.Status == StatusApproved && t.Type == TypeCharge {
			m.TotalTransactionSum += t.Amount
		}
	}
}

func (m *Merchant) buildTransactionRelations() {
	tm := make(map[uint]*Transaction, 0)
	for i := range m.Transactions {
		tm[m.Transactions[i].ID] = &m.Transactions[i]
	}
	for i := range m.Transactions {
		if btID := m.Transactions[i].BelongsToID; btID != nil {
			m.Transactions[i].BelongsTo = tm[*btID]
		}
	}
}

type MerchantStore struct {
	db *gorm.DB
}

func NewMerchantStore(db *gorm.DB) *MerchantStore {
	return &MerchantStore{db: db}
}

func (s *MerchantStore) CreateMerchant(m *Merchant) error {
	return createSingleGorm(m, s.db)
}

func (s *MerchantStore) CreateMerchants(ms []*Merchant) error {
	return createMultipleGorm(ms, s.db)
}

func (s *MerchantStore) GetAllMerchants() ([]*Merchant, error) {
	var ms []*Merchant
	err := s.db.Model(&Merchant{}).Preload("User").Preload("Transactions").Find(&ms).Error
	if err != nil {
		return nil, fmt.Errorf("while getting all merchants: %w", err)
	}
	for i := range ms {
		if ms[i] != nil {
			ms[i].calculateTTS()
			ms[i].buildTransactionRelations()
		}
	}
	return ms, nil
}

func (s *MerchantStore) GetMerchantById(id uint) (*Merchant, error) {
	var m *Merchant
	err := s.db.Model(&Merchant{}).Where("user_id = ?", id).Preload("User").Preload("Transactions").First(&m).Error
	if err != nil {
		return nil, fmt.Errorf("while getting all merchants: %w", err)
	}
	m.calculateTTS()
	m.buildTransactionRelations()
	return m, nil
}

func (s *MerchantStore) DeleteMerchant(m *Merchant) error {
	return deleteSingleGorm(m, s.db)
}
