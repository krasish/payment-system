package models

import (
	"context"
	"database/sql/driver"
	"errors"
	"fmt"
	"net/mail"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

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

func NewTransactionStatus(s string) (TransactionStatus, error) {
	return enumFactory(s, StatusApproved, StatusReversed, StatusRefunded, StatusError)
}

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

func NewTransactionType(s string) (TransactionType, error) {
	return enumFactory(s, TypeAuthorize, TypeCharge, TypeRefund, TypeReversal)
}

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
	Merchant   Merchant

	BelongsToID *uint `gorm:"column:belongs_to"`
	BelongsTo   *Transaction
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	user := &User{}
	res := tx.Model(&User{}).Where("id = ?", t.MerchantID).First(user)
	if res.Error != nil {
		return fmt.Errorf("while getting user in transaction before create hook: %w", err)
	}
	if user.Status != StatusActive {
		return errors.New("while creating transaction: user not in active status")
	}
	return err
}

func (t *Transaction) AfterCreate(tx *gorm.DB) (err error) {
	if t.BelongsToID != nil {
		switch t.Type { //nolint:exhaustive
		case TypeRefund:
			res := tx.Model(&Transaction{}).Where("id = ?", *t.BelongsToID).Update("status", StatusRefunded)
			if err := res.Error; err != nil {
				return fmt.Errorf("while updating refund referenced transaction in after create hook: %w", err)
			}
		case TypeReversal:
			res := tx.Model(&Transaction{}).Where("id = ?", *t.BelongsToID).Update("status", StatusReversed)
			if err := res.Error; err != nil {
				return fmt.Errorf("while updating reversal referenced transaction in adter create hook: %w", err)
			}
		}
	}
	return nil
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
		CustomerEmail: strings.ToLower(customerEmail),
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

func (s *TransactionStore) CreateTransaction(ctx context.Context, t *Transaction) error {
	return createSingleGorm[Transaction](ctx, t, s.db)
}

func (s *TransactionStore) CreateTransactions(ctx context.Context, ts []*Transaction) error {
	return createMultipleGorm[Transaction](ctx, ts, s.db)
}

func (s *TransactionStore) GetTransactionByUUID(ctx context.Context, extID string) (*Transaction, error) {
	var t Transaction
	err := s.db.WithContext(ctx).Model(&Transaction{}).Where("ext_uuid = ?", extID).Preload("Merchant").Preload("BelongsTo").First(&t).Error
	if err != nil {
		return nil, fmt.Errorf("while getting all transctions: %w", err)
	}
	return &t, nil
}

func (s *TransactionStore) GetAllTransactions(ctx context.Context) ([]*Transaction, error) {
	var ts []*Transaction
	err := s.db.WithContext(ctx).Preload("Merchant").Find(&ts).Error
	if err != nil {
		return nil, fmt.Errorf("while getting all transctions: %w", err)
	}
	s.buildTransactionRelations(ts)
	return ts, nil
}

func (s *TransactionStore) DeleteTransaction(ctx context.Context, t *Transaction) error {
	return deleteSingleGorm(ctx, t, s.db)
}

type TransactionPeriodicJob func(context.Context)

func (s *TransactionStore) GetPeriodicJobDeleter(deleteOlderThan, jobExecutionInterval time.Duration) TransactionPeriodicJob {
	return func(ctx context.Context) {
		ticker := time.NewTicker(jobExecutionInterval)
		for {
			select {
			case <-ticker.C:
				res := s.db.Where("created_at < ?", time.Now().Add(-deleteOlderThan)).Delete(&Transaction{})
				if err := res.Error; err != nil {
					logrus.Warnf("periodic transaction deletion job failed: %v", err)
				}
			case <-ctx.Done():
				return
			}
		}
	}
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
