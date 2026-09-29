package main

import (
	"testing"

	"github.com/jinzhu/copier"
)

// Struct gốc
type User struct {
	Name     string
	Age      int
	Email    string
	Location string
	Phone    string
}

// 👇 Copy thủ công
func copyManual(u User) User {
	return User{
		Name:     u.Name,
		Age:      u.Age,
		Email:    u.Email,
		Location: u.Location,
		Phone:    u.Phone,
	}
}

// 👇 Copy dùng copier
func copyWithCopier(u User) User {
	var dst User
	copier.Copy(&dst, &u)
	return dst
}

func BenchmarkCopyManual(b *testing.B) {
	src := User{Name: "Alice", Age: 30, Email: "a@example.com", Location: "HN", Phone: "123456"}
	for i := 0; i < b.N; i++ {
		_ = copyManual(src)
	}
}

func BenchmarkCopyWithCopier(b *testing.B) {
	src := User{Name: "Alice", Age: 30, Email: "a@example.com", Location: "HN", Phone: "123456"}
	for i := 0; i < b.N; i++ {
		_ = copyWithCopier(src)
	}
}
