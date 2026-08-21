migrateup:
	migrate -path infra/db/migrate_v2 -database "${CRM_DATABASE_URL}" -verbose up

migratedown:
	migrate -path infra/db/migrate_v2 -database "${CRM_DATABASE_URL}" -verbose down

.PHONY: migrateup migratedown

# migrate create -ext sql -dir infra/db/migrate_v2 -seq seo_domain_schema