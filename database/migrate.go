package database

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

// Migration represents a database migration
type Migration struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"unique;not null"`
	Applied  bool   `gorm:"default:false"`
	AppliedAt *gorm.DeletedAt
}

// RunMigrations executes all pending migrations
func RunMigrations() error {
	// Ensure migrations table exists
	err := DB.AutoMigrate(&Migration{})
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %v", err)
	}

	// Get all migration files
	migrationFiles, err := getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %v", err)
	}

	// Execute each migration
	for _, file := range migrationFiles {
		err := executeMigration(file)
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %v", file, err)
		}
	}

	log.Println("All migrations completed successfully")
	return nil
}

// getMigrationFiles returns a sorted list of migration files
func getMigrationFiles() ([]string, error) {
	migrationsDir := "database/migrations"
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return nil, err
	}

	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}

	// Sort files to ensure they run in order
	sort.Strings(migrationFiles)
	return migrationFiles, nil
}

// executeMigration executes a single migration file
func executeMigration(filename string) error {
	// Check if migration has already been applied
	var migration Migration
	result := DB.Where("name = ?", filename).First(&migration)
	
	if result.Error == nil && migration.Applied {
		log.Printf("Migration %s already applied, skipping", filename)
		return nil
	}

	// Read migration file
	filePath := filepath.Join("database/migrations", filename)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read migration file %s: %v", filename, err)
	}

	log.Printf("Applying migration: %s", filename)

	// Execute migration SQL
	err = DB.Exec(string(content)).Error
	if err != nil {
		return fmt.Errorf("failed to execute migration SQL: %v", err)
	}

	// Record migration as applied
	if result.Error != nil {
		// Create new migration record
		migration = Migration{
			Name:    filename,
			Applied: true,
		}
		err = DB.Create(&migration).Error
	} else {
		// Update existing migration record
		migration.Applied = true
		err = DB.Save(&migration).Error
	}

	if err != nil {
		return fmt.Errorf("failed to record migration: %v", err)
	}

	log.Printf("Migration %s applied successfully", filename)
	return nil
}

// CreateMigration creates a new migration file template
func CreateMigration(name string) error {
	// Get the next migration number
	migrationFiles, err := getMigrationFiles()
	if err != nil {
		return err
	}

	nextNum := len(migrationFiles) + 1
	filename := fmt.Sprintf("%03d_%s.sql", nextNum, strings.ReplaceAll(name, " ", "_"))
	filePath := filepath.Join("database/migrations", filename)

	// Create migration template
	template := fmt.Sprintf(`-- Migration: %s
-- Date: %s

-- Write your migration SQL here
-- Example:
-- ALTER TABLE table_name ADD COLUMN new_column VARCHAR(255);

-- Remember to use IF EXISTS/IF NOT EXISTS for idempotent migrations
`, name, "$(date +%Y-%m-%d)")

	err = ioutil.WriteFile(filePath, []byte(template), 0644)
	if err != nil {
		return fmt.Errorf("failed to create migration file: %v", err)
	}

	log.Printf("Migration file created: %s", filePath)
	return nil
}