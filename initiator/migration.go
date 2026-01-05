package initiator

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

// InitMigration performs database migrations for all schemas in the specified folder.
func InitMigration(dbUrl, schemasFolder string) error {
	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer db.Close()

	// Try to run migrations from the folder directly first
	fullPath := fmt.Sprintf("file://%s", schemasFolder)
	driver, err := postgres.WithInstance(db, &postgres.Config{
		SchemaName: "public",
	})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fullPath,
		viper.GetString("db.name"),
		driver,
	)
	if err == nil {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Printf("Migration failed for root folder: %v", err)
		} else {
			log.Printf("Successfully applied migrations from %s", schemasFolder)
			return nil
		}
	}

	// Fallback to subdirectories if root folder didn't work (original behavior)
	schemaDirs, err := os.ReadDir(schemasFolder)
	if err != nil {
		return fmt.Errorf("failed to read schemas folder: %w", err)
	}

	for _, schemaDir := range schemaDirs {
		if !schemaDir.IsDir() {
			continue
		}

		schemaPath := filepath.Join(schemasFolder, schemaDir.Name())
		fullPath := fmt.Sprintf("file://%s", schemaPath)

		driver, err := postgres.WithInstance(db, &postgres.Config{
			SchemaName: "public",
		})
		if err != nil {
			return fmt.Errorf("failed to create driver for schema %s: %w", schemaDir.Name(), err)
		}

		m, err := migrate.NewWithDatabaseInstance(
			fullPath,
			viper.GetString("db.name"),
			driver,
		)
		if err != nil {
			return fmt.Errorf("failed to create migration instance for schema %s: %w", schemaDir.Name(), err)
		}

		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			return fmt.Errorf("migration failed for schema %s: %w", schemaDir.Name(), err)
		}

		log.Printf("Successfully applied migrations for schema %s", schemaDir.Name())
	}

	return nil
}
