package models

import (
	"context"
	"database/sql/driver"
	"fmt"
	"time"

	"gorm.io/gorm"
)

const UserTableName = "payment_system_user"

type UserRole string

const (
	RoleMerchant UserRole = "MERCHANT"
	RoleAdmin    UserRole = "ADMIN"
)

func NewUserRole(s string) (UserRole, error) {
	return enumFactory(s, RoleMerchant, RoleAdmin)
}

func (ur *UserRole) Scan(value interface{}) error {
	return scanEnumValue(ur, value)
}

func (ur UserRole) Value() (driver.Value, error) {
	return string(ur), nil
}

type UserStatus string

const (
	StatusActive   UserStatus = "ACTIVE"
	StatusInactive UserStatus = "INACTIVE"
)

func NewUserStatus(s string) (UserStatus, error) {
	return enumFactory(s, StatusActive, StatusInactive)
}

func (us *UserStatus) Scan(value interface{}) error {
	return scanEnumValue(us, value)
}

func (us UserStatus) Value() (driver.Value, error) {
	return string(us), nil
}

type User struct {
	ID        uint      `gorm:"primaryKey;->"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Role   UserRole   `gorm:"column:_role;type:user_role"`
	Status UserStatus `gorm:"type:user_status"`
}

func NewUser(role UserRole, status UserStatus) *User {
	return &User{Role: role, Status: status}
}

func (u User) TableName() string {
	return UserTableName
}

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{db: db}
}

func (s *UserStore) CreateUser(ctx context.Context, u *User) error {
	return createSingleGorm(ctx, u, s.db)
}

func (s *UserStore) CreateUsers(ctx context.Context, us []*User) error {
	return createMultipleGorm(ctx, us, s.db)
}

func (s *UserStore) GetAllUsers(ctx context.Context) ([]*User, error) {
	var us []*User
	res := s.db.WithContext(ctx).Find(&us)
	if err := res.Error; err != nil {
		return nil, fmt.Errorf("while getting users: %w", err)
	}
	return us, nil
}
