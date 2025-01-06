.PHONY: migrate-up migrate-down docker-up docker-down

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate-up:
	docker-compose exec app migrate -path ./migrations -database "postgresql://postgres:postgres@db:5432/gomanager?sslmode=disable" up

migrate-down:
	docker-compose exec app migrate -path ./migrations -database "postgresql://postgres:postgres@db:5432/gomanager?sslmode=disable" down
