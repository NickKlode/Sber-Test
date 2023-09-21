build:
	docker-compose build
run:
	docker-compose up
migrate:
	migrate -path ./schema -database 'postgres://postgres:postgres@0.0.0.0:5432/postgres?sslmode=disable' up
test:
	go test -v -cover ./...