package models

import (
	"fmt"
	"gorm.io/gorm"
)

func scanEnumValue[T UserStatus | UserRole | TransactionStatus | TransactionType](to *T, value any) error {
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

func createSingleGorm[T User | Transaction | Merchant](entity *T, db *gorm.DB) error {
	res := db.Create(entity)
	if err := res.Error; err != nil {
		return fmt.Errorf("while creating %T: %w", entity, err)
	}
	return nil
}

func createMultipleGorm[T User | Transaction | Merchant](entities []*T, db *gorm.DB) error {
	res := db.Create(entities)
	if err := res.Error; err != nil {
		return fmt.Errorf("while creating multiple of type %T: %w", entities, err)
	}
	return nil
}

func deleteSingleGorm[T User | Transaction | Merchant](entity *T, db *gorm.DB) error {
	res := db.Delete(entity)
	if err := res.Error; err != nil {
		return fmt.Errorf("while deleting %T: %w", entity, err)
	}
	return nil
}
