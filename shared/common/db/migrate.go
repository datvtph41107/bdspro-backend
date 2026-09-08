package _db

import (
	"log"
	"reflect"
)

// // AutoMigrate tạo bảng nếu chưa có
// Hàm AutoMigrate với Reflection
func AutoMigrate(e interface{}) {
	// Kiểm tra kiểu của e và truyền vào cho AutoMigrate
	val := reflect.ValueOf(e)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		log.Fatal("❌ Đối tượng truyền vào phải là một struct")
	}

	err := DB.AutoMigrate(e)
	if err != nil {
		log.Fatal("❌ Lỗi khi migrate database:", err)
	}
	log.Println("✅ Đã migrate database thành công!")
}
