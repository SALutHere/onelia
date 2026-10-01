# Подключение файла переменных окружения
include .env
export


# Экспортирование переменной директории проекта
export PROJECT_ROOT=${shell pwd}


# Сценарии работы с окружением

## Поднять окружение
env-up:
	@docker compose up -d onelia-postgres

## Остановить окружение
env-down:
	@docker compose down onelia-postgres

## Очистить окружение
env-cleanup:
	@read -p "Вы действительно хотите очистить все volume-файлы окружения? Данные будут утеряны. [y/n]: " ans; \
	if [ $$ans = "y" ]; then \
		docker compose down onelia-postgres -v && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

## Пробросить порт Postgres на основную систему (5434)
env-port-forward:
	@docker compose up -d port-forwarder

## Перекрыть проброс проста Postgres (5434)
env-port-close:
	@docker compose down port-forwarder


# Сценарии работы с миграциями

## Создать up- и down-файлы новой миграции с заданием имени
migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует необходимый параметр seq. Пример: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm onelia-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

## Накатить миграции
migrate-up:
	@make migrate-action action=up

## Откатить миграции
migrate-down:
	@make migrate-action action=down

## (не для прямого использования)
## Обобщённое действие с миграциями
migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр action. Пример: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm onelia-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@onelia-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"
