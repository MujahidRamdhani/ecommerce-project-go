package cronjob

import (
	"database/sql"
	"log"
	"github.com/robfig/cron/v3"
)

func StartTransactionExpiryChecker(db *sql.DB) {
	c := cron.New()

	c.AddFunc("@every 1m", func() {
		log.Println("Checking for expired transactions...")

		query := `UPDATE transactions SET status = 'Failed' WHERE expired_at <= NOW() AND status = 'Unpaid'`
		result, err := db.Exec(query)
		if err != nil {
			log.Println("Failed to update expired transactions:", err)
			return
		}
		rowsAffected, _ := result.RowsAffected()
		log.Println("Expired transactions updated:", rowsAffected)
	})

	c.Start()
}

