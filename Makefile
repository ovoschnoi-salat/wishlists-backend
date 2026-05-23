# Команды для управления сервисами

.PHONY: restart-prod
restart-prod:
ifdef service
	docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d $(service)
else
	docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d
endif
	@echo "services restarted!"
	@echo ""

# Собирает и перезапускает сервисы
.PHONY: deploy-prod
deploy-prod:
ifdef service
	docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build $(service)
else
	docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build
endif
	@echo "services deployed!"
	@echo ""

.PHONY: stop-prod
stop-prod:
ifdef service
	docker compose -f docker-compose.yml -f docker-compose.prod.yml down $(service)
else
	docker compose -f docker-compose.yml -f docker-compose.prod.yml down
endif
	@echo "services stopped!"
	@echo ""


# Собирает и перезапускает сервисы
.PHONY: deploy-dev
deploy-dev:
ifdef service
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build $(service)
else
	docker compose -f docker-compose.yml -f docker-compose.dev.yml up -d --build
endif
	@echo "services deployed!"
	@echo ""

.PHONY: stop-dev
stop-dev:
ifdef service
	docker compose -f docker-compose.yml -f docker-compose.dev.yml down $(service)
else
	docker compose -f docker-compose.yml -f docker-compose.dev.yml down
endif
	@echo "services stopped!"
	@echo ""


.PHONY: install-utils
install-utils:
	go install tool

.PHONY: generate-swagger
generate-swagger: install-utils
	swag init  --generalInfo cmd/main.go


PG_DSN:=postgres://postgres:postgres@localhost:5432/wishlists?sslmode=disable&timezone=UTC

.PHONY: generate-sql
generate-sql: install-utils
	sqlc generate


# migrations
.PHONY: create-migration
create-migration: install-utils
ifndef NAME
	$(error NAME не задан. Использование: make create-migration NAME=<имя>)
endif
	(cd migrations; goose -s create $(NAME) sql)

.PHONY: migrate-up
migrate-up: install-utils
	(cd migrations; goose postgres "${PG_DSN}" up)

.PHONY: migrate-down
migrate-down: install-utils
	(cd migrations; goose postgres "${PG_DSN}" down)


ifneq (,$(wildcard ./.env))
	include .env
	export
endif

.PHONY: db-dump
db-dump:
	mkdir backup
	docker exec -t wishlists-db pg_dump -U ${POSTGRES_USER} wishlists > backup/dump.sql
