package controllers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/krasish/payment-system/internal/models"
)

type User struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Role      string
	Status    string
}

func (m *User) CSVUnmarshal(record []string) error {
	if len(record) != 1 {
		return fmt.Errorf("wrong number of records in user CSV record: %v", record)
	} else if m == nil {
		return errors.New("cannot unmarshal CSV to a nil user")
	}
	m.Status = record[0]
	m.Role = string(models.RoleAdmin)
	return nil
}

func (u *User) toModel() (*models.User, error) {
	role, err := models.NewUserRole(u.Role)
	if err != nil {
		return nil, err
	}
	status, err := models.NewUserStatus(u.Status)
	if err != nil {
		return nil, err
	}
	return models.NewUser(role, status), nil
}

type UserController struct {
	store *models.UserStore
}

func NewUserController(store *models.UserStore) *UserController {
	return &UserController{store: store}
}

func (c *UserController) CreateUsers(ctx context.Context, us []*User) error {
	modelUsers := make([]*models.User, 0, len(us))
	for _, u := range us {
		model, err := u.toModel()
		if err != nil {
			return fmt.Errorf("while converting users to model DTO: %w", err)
		}
		modelUsers = append(modelUsers, model)
	}

	return c.store.CreateUsers(ctx, modelUsers)
}
