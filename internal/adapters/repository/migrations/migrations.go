package migrations

import (
	"fmt"

	"gorm.io/gorm"

	"soporte/internal/adapters/repository/models"
)

func AutoMigrateModels(db *gorm.DB) error {
	if err := db.AutoMigrate(models.All()...); err != nil {
		return fmt.Errorf("auto migrate models: %w", err)
	}

	return nil
}
