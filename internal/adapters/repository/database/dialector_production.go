//go:build production

package database

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func openDialector(driver, dsn string) (gorm.Dialector, error) {
	switch driver {
	case "postgres":
		return postgres.Open(dsn), nil
	default:
		return nil, fmt.Errorf("unsupported database driver in production build: %s", driver)
	}
}
