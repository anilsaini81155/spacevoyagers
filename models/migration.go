package models

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

// Migration represents a migration script with a name and SQL query.
type Migration struct {
	Name  string
	Query string
}

// Global DB variable to hold the database connection
var DB *sql.DB

// SetDB allows the main function to set the global DB instance
func SetDB(db *sql.DB) {
	DB = db
}

var migrations = []Migration{
	{
		Name: "create_exoplanets_table",
		Query: `
            CREATE TABLE IF NOT EXISTS exoplanets (
                id INT AUTO_INCREMENT,
                name VARCHAR(255) NOT NULL,
                description TEXT,
                distance FLOAT,
                radius FLOAT,
                mass FLOAT DEFAULT NULL,
                type VARCHAR(50),
                PRIMARY KEY (id)
            );
        `,
	},
	{
		Name: "add_gravity_column_to_exoplanets",
		Query: `
            ALTER TABLE exoplanets ADD COLUMN gravity FLOAT DEFAULT NULL;
        `,
	},
}

// RunMigrations applies all pending migrations
func RunMigrations() {

	if DB == nil {
		panic("database connection is not initialized")
	}

	for _, migration := range migrations {
		applyMigration(migration)
	}
}

// applyMigration runs a specific migration if it hasn't been applied yet
func applyMigration(migration Migration) {
	if !hasMigrationBeenApplied(migration.Name) {
		_, err := DB.Exec(migration.Query)
		if err != nil {
			log.Fatalf("Error applying migration %s: %v", migration.Name, err)
		}
		log.Printf("Migration %s applied successfully", migration.Name)
		recordMigration(migration.Name)
	}
}

// hasMigrationBeenApplied checks if a migration has already been applied
func hasMigrationBeenApplied(name string) bool {
	query := `SELECT COUNT(*) FROM migrations WHERE name = ?`
	var count int
	err := DB.QueryRow(query, name).Scan(&count)
	if err != nil {
		// If the `migrations` table doesn't exist, create it.
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1146 {
			log.Println("Migrations table doesn't exist. Creating it now...")
			createMigrationsTable()
			return false
		}
		log.Fatalf("Error checking migration status: %v", err)
	}
	return count > 0
}

// recordMigration marks a migration as applied by adding it to the `migrations` table
func recordMigration(name string) {
	query := `INSERT INTO migrations (name) VALUES (?)`
	_, err := DB.Exec(query, name)
	if err != nil {
		log.Fatalf("Error recording migration %s: %v", name, err)
	}
}

// createMigrationsTable creates the migrations table if it doesn't exist
func createMigrationsTable() {
	query := `
        CREATE TABLE IF NOT EXISTS migrations (
            id INT AUTO_INCREMENT,
            name VARCHAR(255) NOT NULL,
            applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            PRIMARY KEY (id)
        );
    `
	_, err := DB.Exec(query)
	if err != nil {
		log.Fatalf("Error creating migrations table: %v", err)
	}
}
