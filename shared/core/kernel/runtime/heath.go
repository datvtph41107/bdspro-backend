package runtime

import (
	"net/http"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func HealthHandler(db *gorm.DB, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db != nil {
			sqlDB, err := db.DB()
			if err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}

			if err := sqlDB.Ping(); err != nil {
				http.Error(w, "Database error", http.StatusInternalServerError)
				return
			}
		}

		if rdb != nil {
			if err := rdb.Ping(r.Context()).Err(); err != nil {
				http.Error(w, "Redis error", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
