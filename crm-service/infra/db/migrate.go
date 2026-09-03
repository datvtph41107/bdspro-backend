package db

import (
	_db "common/db"
	"log"
)

func execSQL(stmt string, args ...any) error {
	return _db.DB.Exec(stmt, args...).Error
}

func scanSQL(stmt string, dest any, args ...any) error {
	return _db.DB.Raw(stmt, args...).Scan(dest).Error
}

func tableExists(name string) bool {
	var n int
	_ = scanSQL(`
		SELECT COUNT(1)
		FROM information_schema.tables
		WHERE table_schema = 'public' AND table_name = ?
	`, &n, name)
	return n > 0
}

func AutoMigrate() {
	// ensureSEOMigrations()

	// if err := _db.DB.AutoMigrate(&seo_domain.SeoRef{}); err != nil {
	// 	log.Printf("⚠️  AutoMigrate SeoRef: %v", err)
	// }

	log.Println("Migrate success")
}
