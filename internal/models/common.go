package models

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

type EnumsConstraint interface {
	UserRole | UserStatus | TransactionType | TransactionStatus
}

func enumFactory[T EnumsConstraint](s string, possibleValues ...T) (T, error) {
	if len(possibleValues) < 1 {
		//panic is intentional since meeting this condition shows
		//incorrect programming (an enum without possible values)
		panic("incorrectly constructed enum")
	}
	for _, value := range possibleValues {
		if s == string(value) {
			return value, nil
		}
	}
	return T(""), fmt.Errorf("%q is not a possible value for type %T", s, possibleValues[0])
}

func scanEnumValue[T EnumsConstraint](to *T, value any) error {
	switch typed := value.(type) {
	case []byte:
		*to = T(typed)
	case string:
		*to = T(typed)
	default:
		return fmt.Errorf("attempted to scan unsupported type: '%T'", value)
	}
	return nil
}

type TypesConstraint interface {
	User | Transaction | Merchant
}

func createSingleGorm[T TypesConstraint](ctx context.Context, entity *T, db *gorm.DB) error {
	res := db.WithContext(ctx).Create(entity)
	if err := res.Error; err != nil {
		return fmt.Errorf("while creating %T: %w", entity, err)
	}
	return nil
}

func createMultipleGorm[T TypesConstraint](ctx context.Context, entities []*T, db *gorm.DB) error {
	res := db.WithContext(ctx).Create(entities)
	if err := res.Error; err != nil {
		return fmt.Errorf("while creating multiple of type %T: %w", entities, err)
	}
	return nil
}

func deleteSingleGorm[T TypesConstraint](ctx context.Context, entity *T, db *gorm.DB) error {
	res := db.WithContext(ctx).Delete(entity)
	if err := res.Error; err != nil {
		return fmt.Errorf("while deleting %T: %w", entity, err)
	}
	return nil
}
