
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
