.PHONY: docker-up
# backend cmd
dev:
	cd backend && go mod tidy && go run main.go
	go mod tidy
	go run main.go

# frontend cmd

# docker
docker-up:
	docker-up: ## Start all docker services
	cd mock && docker-compose up -d
	@echo "Postgresql: localhost:5432"
	@echo "Redis: localhost:6379"
	@echo "Kafka UI: http://localhost:8080"
	@echo "Adminer (DB UI): http://localhost:8081"
	@echo "Redis Commander: http://localhost:8082"
	@echo "fnish setup"

docker-down:
	cd mock && docker-compose down



