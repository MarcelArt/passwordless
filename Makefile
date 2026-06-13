swag:
	@swag init --parseDependency --parseInternal

dev: swag
	@air serve

go: swag
	@go run main.go serve

migrate:
	@go run main.go migrate

migrate-force:
	@go run main.go migrate --drop

docker-build:
	@docker build -t marcelaritonang/passwordless:latest .

docker-run: docker-build
	@docker run -p 7002:7002 --env-file .env marcelaritonang/passwordless:latest

docker-push: docker-build
	@docker push marcelaritonang/passwordless:latest

compose-down:
	@docker compose down

compose: compose-down
	@docker compose up -d