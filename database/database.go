package database

import (
	"database/sql"
	"fmt"
	"log"
	"github.com/gobuffalo/packr/v2"
	migrate "github.com/rubenv/sql-migrate"
)

var DbConnection *sql.DB

func DbMigrate(db *sql.DB) {
	migrations := &migrate.PackrMigrationSource{
		Box: packr.New("migrations", "./sql_migration"),
	}

	n, err := migrate.Exec(db, "postgres", migrations, migrate.Up)
	if err != nil {
		log.Fatal(err)
	}

	DbConnection = db
	fmt.Println("Applied", n, "migrations!")
}

