package main

import (
	"database/sql"
	"ecommerce-project-go/cronjob"
	"ecommerce-project-go/database"
	"ecommerce-project-go/routers"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	db  *sql.DB
	err error
)

func main() {
	err = godotenv.Load("config/.env")
	if err != nil {
		log.Fatal(err)
	}
	psqlInfo := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		os.Getenv("PGHOST"),
		os.Getenv("PGPORT"),
		os.Getenv("PGUSER"),
		os.Getenv("PGPASSWORD"),
		os.Getenv("PGDATABASE"))

	// connect ke database PostgreSQL
	db, err = sql.Open("postgres", psqlInfo)
	err = db.Ping()
	if err != nil {
		fmt.Println("Connection to database failed")
		panic(err)
	}

	fmt.Println("Successfully connected to database")

	// database migration
	database.DbMigrate(db)
	defer db.Close()

	//start cron job
	cronjob.StartTransactionExpiryChecker(db)

	// Start server menggunakan routers
	routers.StartServer(":"+os.Getenv("PORT"), db)
}
