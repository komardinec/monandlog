run/web:
	go run ./cmd/web

db/mysql:
	mysql -u root -p 

db/migration:
	@echo "Creating migration files for ${name}"
	migrate create -seq -ext=.sql -dir=./migrations ${name}

db/up:
	@echo "Running up migrations..."
	migrate -path=./migrations -database="${MYSQL_DB_DSN}" up
