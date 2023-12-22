include .env 

.PHONY: clean
clean:
	docker-compose down

.PHONY: up
up: 
	docker compose up -d

.PHONY: down
down:
	docker compose down

