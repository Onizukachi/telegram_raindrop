include .env
export $(shell sed 's/=.*//' .env)

MIGRATE = migrate -path pkg/storage/migrations -database "$(DB_DSN)" -verbose

migrate-new:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir pkg/storage/migrations -seq $$name

## Применить все миграции
migrate-up:
	$(MIGRATE) up

## Применить N миграций (пример: make migrate-up-n n=1)
migrate-up-n:
	$(MIGRATE) up $(n)

## Откатить все миграции
migrate-down:
	$(MIGRATE) down

## Откатить N миграций (пример: make migrate-down-n n=1)
migrate-down-n:
	$(MIGRATE) down $(n)

## Показать текущую версию
migrate-version:
	$(MIGRATE) version

## Полностью сбросить БД (drop)
migrate-drop:
	$(MIGRATE) drop -f