include .env
export

export PROJECT_ROOT:=$(shell pwd)

print-project-root:
	@echo $(PROJECT_ROOT)

env-up:
	@mkdir -p .out
	@chmod 777 .out
	@docker compose up -d todoapp-postgres

env-down:
	@docker compose down todoapp-postgres

env-cleanup:
	@read -p "Are you sure you want to remove the database volume? (y/n) " answer; \
	if [ "$$answer" = "y" ]; then \
		docker compose down todoapp-postgres port-forwarder && \
		docker volume rm golang-todoapp_todoapp-postgres-data 2>/dev/null || true && \
		echo "Database volume removed."; \
	else \
		echo "Operation canceled."; \
	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Ошибка: отсутствует необходимый параметр seq. Пример: make migrate-create seq=001"; \
		exit 1; \
	fi
	docker compose run --rm todoapp-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
		@if [ -z "$(action)" ]; then \
			echo "Ошибка: отсутствует необходимый параметр action. Пример: make migrate-action action=up"; \
			exit 1; \
		fi
		docker compose run --rm todoapp-postgres-migrate \
		-path /migrations \
		-database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@todoapp-postgres:5432/$(POSTGRES_DB)?sslmode=disable" \
		"$(action)"


env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder


todoapp-run:
	@export LOGGER_FOLDER=$(PROJECT_ROOT)/out/logs && \
	export POSTGRES_HOST=localhost && \
	go mod tidy && \
	go run $(PROJECT_ROOT)/cmd/todoapp/main.go
