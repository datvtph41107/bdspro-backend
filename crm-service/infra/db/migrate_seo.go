package db

import (
	"embed"
	"log"
	"strings"
	"unicode"
)

//go:embed migrate_v2/*.up.sql
var seoMigrateFS embed.FS

// seoMigrationFiles — thứ tự apply schema SEO (migrate_v2).
// Bỏ qua support_ticket / admin_opportunity (không thuộc SEO).
var seoMigrationFiles = []string{
	"000001_seo_domain_registry.up.sql",
	"000002_seo_links.up.sql",
	"000003_seo_relative.up.sql",
	"000004_seo_trend.up.sql",
	"000005_seo_generation_log.up.sql",
	"000006_seo_source_snapshot.up.sql",
	"000007_seo_render_fields.up.sql",
	"000008_seo_measurement_feedback.up.sql",
	"000009_seo_objective_foundation.up.sql",
	"000010_seo_audience_user_need.up.sql",
	"000011_seo_query_intent_cluster.up.sql",
	"000012_seo_content_taxonomy_placement.up.sql",
	"000013_seo_m8_measurement_control.up.sql",
	"000014_seo_generation_log_running_status.up.sql",
	"000015_seo_runtime_event.up.sql",
	"000016_seo_typed_source_identity.up.sql",
	"000019_planning_news_visibility_compatibility.up.sql",
	"000023_user_saved_planning_news.up.sql",
	"000024_simplify_seo_to_pages.up.sql",
	"000025_seo_domain_category_id.up.sql",
	"000026_seo_count.up.sql",
	"000027_seo_count_aggregate.up.sql",
	"000028_seo_domain_image_url.up.sql",
	"000029_seo_domain_province_id.up.sql",
	"000030_restore_seo_category.up.sql",
}

func ensureSEOMigrations() {
	if err := execSQL(`
		CREATE TABLE IF NOT EXISTS crm_seo_schema_migrations (
			filename VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`); err != nil {
		log.Printf("⚠️  ensureSEOMigrations: không tạo được crm_seo_schema_migrations: %v", err)
	}

	// Bootstrap một lần: schema đã có sẵn mà chưa có ledger → đánh dấu các file hiện tại.
	var ledgerCount int
	_ = scanSQL(`SELECT COUNT(1) FROM crm_seo_schema_migrations`, &ledgerCount)
	if ledgerCount == 0 && tableExists("seo_category") && tableExists("seo_content_category_assignment") {
		for _, name := range seoMigrationFiles {
			_ = execSQL(
				`INSERT INTO crm_seo_schema_migrations(filename) VALUES (?) ON CONFLICT DO NOTHING`,
				name,
			)
		}
		log.Println("✅ SEO schema đã có — bootstrap crm_seo_schema_migrations")
	}

	for _, name := range seoMigrationFiles {
		var applied int
		if err := scanSQL(
			`SELECT COUNT(1) FROM crm_seo_schema_migrations WHERE filename = ?`,
			&applied,
			name,
		); err == nil && applied > 0 {
			continue
		}

		path := "migrate_v2/" + name
		raw, err := seoMigrateFS.ReadFile(path)
		if err != nil {
			log.Printf("⚠️  ensureSEOMigrations: đọc %s: %v", name, err)
			continue
		}
		if err := execSQLScript(string(raw), name); err != nil {
			log.Printf("⚠️  ensureSEOMigrations: %s: %v", name, err)
			continue
		}
		_ = execSQL(
			`INSERT INTO crm_seo_schema_migrations(filename) VALUES (?) ON CONFLICT DO NOTHING`,
			name,
		)
		log.Printf("✅ SEO migrate: %s", name)
	}
}

func execSQLScript(script, label string) error {
	stmts := splitSQLStatements(script)
	for i, stmt := range stmts {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		upper := strings.ToUpper(stmt)
		if upper == "BEGIN" || upper == "COMMIT" || upper == "ROLLBACK" {
			continue
		}
		if err := execSQL(stmt); err != nil {
			// Clear aborted-tx state on pooled connection before continue/return.
			_ = execSQL("ROLLBACK")
			if isIgnorableMigrateError(err) {
				log.Printf("⚠️  ensureSEOMigrations[%s#%d]: bỏ qua: %v", label, i+1, err)
				continue
			}
			return err
		}
	}
	return nil
}

func isIgnorableMigrateError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "already exists") ||
		strings.Contains(msg, "duplicate key") ||
		strings.Contains(msg, "is not a view")
}

// splitSQLStatements tách statement theo ';' ngoài string / dollar-quote / comment.
func splitSQLStatements(script string) []string {
	var (
		out       []string
		buf       strings.Builder
		inSingle  bool
		inLineCmt bool
		inBlock   bool
		dollarTag string // non-empty when inside $tag$...$tag$
		i         int
	)

	flush := func() {
		s := strings.TrimSpace(buf.String())
		buf.Reset()
		if s != "" {
			out = append(out, s)
		}
	}

	for i < len(script) {
		ch := script[i]

		if inLineCmt {
			buf.WriteByte(ch)
			if ch == '\n' {
				inLineCmt = false
			}
			i++
			continue
		}

		if inBlock {
			buf.WriteByte(ch)
			if ch == '*' && i+1 < len(script) && script[i+1] == '/' {
				buf.WriteByte('/')
				i += 2
				inBlock = false
				continue
			}
			i++
			continue
		}

		if dollarTag != "" {
			buf.WriteByte(ch)
			if ch == '$' && strings.HasPrefix(script[i:], dollarTag) {
				buf.WriteString(dollarTag[1:])
				i += len(dollarTag)
				dollarTag = ""
				continue
			}
			i++
			continue
		}

		if inSingle {
			buf.WriteByte(ch)
			if ch == '\'' {
				if i+1 < len(script) && script[i+1] == '\'' {
					buf.WriteByte('\'')
					i += 2
					continue
				}
				inSingle = false
			}
			i++
			continue
		}

		// line comment
		if ch == '-' && i+1 < len(script) && script[i+1] == '-' {
			buf.WriteByte(ch)
			buf.WriteByte('-')
			i += 2
			inLineCmt = true
			continue
		}

		// block comment
		if ch == '/' && i+1 < len(script) && script[i+1] == '*' {
			buf.WriteByte(ch)
			buf.WriteByte('*')
			i += 2
			inBlock = true
			continue
		}

		// dollar quote $tag$
		if ch == '$' {
			j := i + 1
			for j < len(script) {
				c := script[j]
				if c == '$' {
					dollarTag = script[i : j+1]
					buf.WriteString(dollarTag)
					i = j + 1
					break
				}
				if !(unicode.IsLetter(rune(c)) || unicode.IsDigit(rune(c)) || c == '_') {
					buf.WriteByte(ch)
					i++
					dollarTag = ""
					break
				}
				j++
			}
			if j >= len(script) {
				buf.WriteByte(ch)
				i++
				dollarTag = ""
			}
			continue
		}

		if ch == '\'' {
			inSingle = true
			buf.WriteByte(ch)
			i++
			continue
		}

		if ch == ';' {
			flush()
			i++
			continue
		}

		buf.WriteByte(ch)
		i++
	}
	flush()
	return out
}
