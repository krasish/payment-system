package models

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"

	"gorm.io/gorm/clause"

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
	return &Merchant{Name: name, Description: description, Email: strings.ToLower(email), User: User{
		Role:   RoleMerchant,
		Status: status,
	}}, nil
}

func (m *Merchant) BeforeUpdate(tx *gorm.DB) error {
	merchantWithID := Merchant{}
	res := tx.Model(&merchantWithID).Where("email = ?", m.Email).First(&merchantWithID)
	if err := res.Error; err != nil {
		return fmt.Errorf("while preloading user_id for merchant in before update hook: %w", err)
	}
	res = tx.Model(User{}).Where("id = ?", merchantWithID.UserID).Update("status", m.User.Status)
	if err := res.Error; err != nil {
		return fmt.Errorf("while updating user in merchant before update hook: %w", err)
	}
	return nil
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

func (s *MerchantStore) CreateMerchant(ctx context.Context, m *Merchant) error {
	return createSingleGorm(ctx, m, s.db)
}

func (s *MerchantStore) CreateMerchants(ctx context.Context, ms []*Merchant) error {
	return createMultipleGorm(ctx, ms, s.db)
}

func (s *MerchantStore) GetAllMerchants(ctx context.Context) ([]*Merchant, error) {
	var ms []*Merchant
	err := s.db.WithContext(ctx).Model(&Merchant{}).Preload("User").Preload("Transactions").Find(&ms).Error
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

func (s *MerchantStore) GetMerchantById(ctx context.Context, id uint) (*Merchant, error) {
	return s.getMerchantByCondition(ctx, "user_id = ?", id)

}

func (s *MerchantStore) GetMerchantByEmail(ctx context.Context, email string) (*Merchant, error) {
	return s.getMerchantByCondition(ctx, "email = ?", strings.ToLower(email))
}

func (s *MerchantStore) UpdateMerchant(ctx context.Context, m *Merchant) error {
	res := s.db.WithContext(ctx).Model(m).Omit(clause.Associations).Where("email = ?", m.Email).Updates(map[string]interface{}{"name": m.Name, "description": m.Description, "email": m.Email})
	if err := res.Error; err != nil {
		return fmt.Errorf("while updating merchant: %w", err)
	}
	return nil
}

func (s *MerchantStore) getMerchantByCondition(ctx context.Context, condition string, arg any) (*Merchant, error) {
	var m *Merchant
	err := s.db.WithContext(ctx).Model(&Merchant{}).Where(condition, arg).Preload("User").Preload("Transactions").First(&m).Error
	if err != nil {
		return nil, fmt.Errorf("while getting merchantwith condition %q: %w", condition, err)
	}
	m.calculateTTS()
	m.buildTransactionRelations()
	return m, nil
}

func (s *MerchantStore) DeleteMerchant(ctx context.Context, email string) error {
	res := s.db.WithContext(ctx).Where("email = ?", email).Delete(&Merchant{})
	if err := res.Error; err != nil {
		if pqError, ok := err.(*pgconn.PgError); ok && pqError.Code == "23503" {
			return errors.New("merchant is still referenced by transactions")
		}
		return fmt.Errorf("while deleting merchant: %w", err)
	}
	return nil
}
