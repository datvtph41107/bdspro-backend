package query

import "gorm.io/gorm"

// Scope là một đơn vị filter — compose được, test được độc lập.
// Đây là "Specification" đúng nghĩa với GORM: không serialize ra string,
// không tự quản lý placeholder — để GORM làm việc đó.
type Scope func(*gorm.DB) *gorm.DB

// Apply gộp nhiều Scope lại, bỏ qua nil.
func Apply(db *gorm.DB, scopes ...Scope) *gorm.DB {
	active := make([]func(*gorm.DB) *gorm.DB, 0, len(scopes))
	for _, s := range scopes {
		if s != nil {
			active = append(active, s)
		}
	}
	return db.Scopes(active...)
}

// ─── Leaf scopes — reusable primitives ───────────────────────────────────────

func WhereEq(col string, val any) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(col+" = ?", val)
	}
}

func WhereContains(col, val string) Scope {
	if val == "" {
		return nil // nil scope = no-op
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(col+" ILIKE ?", "%"+val+"%")
	}
}

func WhereIn(col string, vals any) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(col+" IN ?", vals)
	}
}

func WhereGte(col string, val any) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(col+" >= ?", val)
	}
}

func WhereLte(col string, val any) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(col+" <= ?", val)
	}
}

func WhereBetween(col string, from, to any) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(col+" BETWEEN ? AND ?", from, to)
	}
}

// OrderBy với whitelist — field PHẢI đến từ whitelist, không bao giờ
// dùng trực tiếp user input trong ORDER BY (SQL injection vector khác với WHERE).
func OrderBy(field, dir string, allowed []string) Scope {
	for _, a := range allowed {
		if a == field {
			return func(db *gorm.DB) *gorm.DB {
				return db.Order(field + " " + dir)
			}
		}
	}
	return nil // unknown field → ignore
}

func Paginate(p Pagination) Scope {
	return func(db *gorm.DB) *gorm.DB {
		return db.Limit(p.Limit()).Offset(p.Offset())
	}
}
