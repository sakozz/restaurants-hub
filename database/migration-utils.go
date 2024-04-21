package database

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigrations() {
	migrationUrl := os.Getenv("MIGRATION_URL")

	migration, err := migrate.New("file://"+migrationUrl, DbConnectionString())

	if err != nil {
		log.Fatal("Error connecting migration URL: ", err)
	}

	if err = migration.Up(); err != nil {
		if err != migrate.ErrNoChange {
			log.Fatal("Failied to run migration UP: ", err)
		}
	}

	fmt.Println("Database migrated successfully.")
}
