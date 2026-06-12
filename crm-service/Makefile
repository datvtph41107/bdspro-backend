migrateup:
	migrate -path infra/db/migrate_v2 -database "postgresql://hieubap:Matkhau.1@14.225.210.29:5432/db_crm?sslmode=disable" -verbose up

migratedown:
	migrate -path infra/db/migrate_v2 -database "postgresql://hieubap:Matkhau.1@14.225.210.29:5432/db_crm?sslmode=disable" -verbose down

.PHONY: migrateup migratedown

# migrate create -ext sql -dir infra/db/migrate_v2 -seq seo_domain_schema