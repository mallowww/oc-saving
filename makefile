.PHONY: docker-up tidy run
# backend cmd
tidy:
	cd backend && go mod tidy

run:
	cd backend && go run main.go

# frontend cmd

# docker
docker-up:
	cd mock && docker-compose up -d
	@echo "Postgresql: localhost:5432"
	@echo "Redis: localhost:6379"
	@echo "Kafka UI: http://localhost:8080"
	@echo "Adminer (DB UI): http://localhost:8081"
	@echo "Redis Commander: http://localhost:8082"
	@echo "finish setup"

docker-down:
	cd mock && docker-compose down



