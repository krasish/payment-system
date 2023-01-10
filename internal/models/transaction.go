package models

import (
	"database/sql/driver"
	"fmt"
	"net/mail"
	"time"

	"github.com/docker/distribution/uuid"

	"gorm.io/gorm"
)

type TransactionStatus string

const (
	StatusApproved TransactionStatus = "APPROVED"
	StatusReversed TransactionStatus = "REVERSED"
	StatusRefunded TransactionStatus = "REFUNDED"
	StatusError    TransactionStatus = "ERROR"
)

func (ts *TransactionStatus) Scan(value interface{}) error {
	return scanEnumValue(ts, value)
}

func (ts TransactionStatus) Value() (driver.Value, error) {
	return string(ts), nil
}

type TransactionType string

const (
	TypeAuthorize TransactionType = "AUTHORIZE"
	TypeCharge    TransactionType = "CHARGE"
	TypeRefund    TransactionType = "REFUND"
	TypeReversal  TransactionType = "REVERSAL"
)

func (tt *TransactionType) Scan(value interface{}) error {
	return scanEnumValue(tt, value)
}

func (tt TransactionType) Value() (driver.Value, error) {
	return string(tt), nil
}

type Transaction struct {
	ID        uint      `gorm:"primaryKey;->"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	ExternalID    string            `gorm:"column:ext_uuid;type:uuid"`
	Type          TransactionType   `gorm:"column:_type;type:transaction_type"`
	Amount        Currency          `gorm:"type:bigint"`
	Status        TransactionStatus `gorm:"type:transaction_status"`
	CustomerEmail string
	CustomerPhone string

	MerchantID uint
	Merchant   Merchant `gorm:"foreignKey:UserID"`

	BelongsToID *uint `gorm:"column:belongs_to"`
	BelongsTo   *Transaction
}

func NewTransaction(externalID string, amount Currency, Type TransactionType, status TransactionStatus, customerEmail string, customerPhone string, merchantID uint, belongsToID *uint) (*Transaction, error) {
	_, err := mail.ParseAddress(customerEmail)
	if err != nil {
		return nil, fmt.Errorf("while creating transaction: %q is not a valid email address: %w", customerEmail, err)
	}
	_, err = uuid.Parse(externalID)
	if err != nil {
		return nil, fmt.Errorf("while creating transaction: %q is not a valid uuid: %w", externalID, err)
	}
	transaction := &Transaction{
		ExternalID:    externalID,
		Type:          Type,
		Amount:        amount,
		Status:        status,
		CustomerEmail: customerEmail,
		CustomerPhone: customerPhone,
		MerchantID:    merchantID,
		BelongsToID:   belongsToID,
	}
	return transaction, nil
}

type TransactionStore struct {
	db *gorm.DB
}

func NewTransactionStore(db *gorm.DB) *TransactionStore {
	return &TransactionStore{db: db}
}

func (s *TransactionStore) CreateTransaction(t *Transaction) error {
	return createSingleGorm[Transaction](t, s.db)
}

func (s *TransactionStore) CreateTransactions(ts []*Transaction) error {
	return createMultipleGorm[Transaction](ts, s.db)
}

func (s *TransactionStore) UpdateStatus(t *Transaction) error {
	res := s.db.Model(t).Where("id = ?", t.ID).Update("status", t.Status)
	if err := res.Error; err != nil {
		return fmt.Errorf("while updating transaction: %w", err)
	}
	return nil
}

func (s *TransactionStore) GetAllTransactions() ([]*Transaction, error) {
	var ts []*Transaction
	err := s.db.Model(&Transaction{}).Preload("Merchant").Find(&ts).Error
	if err != nil {
		return nil, fmt.Errorf("while getting all transctions: %w", err)
	}
	s.buildTransactionRelations(ts)
	return ts, nil
}

func (s *TransactionStore) DeleteTransaction(t *Transaction) error {
	return deleteSingleGorm(t, s.db)
}

func (s *TransactionStore) buildTransactionRelations(ts []*Transaction) {
	tm := make(map[uint]*Transaction, 0)
	for i := range ts {
		tm[ts[i].ID] = ts[i]
	}
	for i := range ts {
		if btID := ts[i].BelongsToID; btID != nil {
			ts[i].BelongsTo = tm[*btID]
		}
	}
}
