run:
	go run cmd/sso/main.go --config=./config/local.yaml

migrate:
	go run ./cmd/migrator \
		--storage-path=./storage/sso.db \
		--migrations-path=./migrations

test_migrate:
	go run ./cmd/migrator/main.go \
		--storage-path=./storage/sso.db \
		--migrations-path=./tests/migrations \
		--migration-table=migrations_test
