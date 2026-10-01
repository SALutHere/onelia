# Подключение файла переменных окружения
include .env
export


# Экспортирование переменной директории проекта
export PROJECT_ROOT=${shell pwd}


# Сценарии работы с окружением

## Поднять контейнер с PostgreSQL
env-up:
	@docker compose up -d onelia-postgres

## Остановить контейнер с PostgreSQL
env-down:
	@docker compose down onelia-postgres

## Остановить и очистить контейнер с PostgreSQL
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


# Очистка логов
logs-cleanup:
	@read -p "Вы действительно хотите очистить все log-файлы? Логи будут утеряны. [y/n]: " ans; \
	if [ $$ans = "y" ]; then \
		rm -rf ${PROJECT_ROOT}/out/logs && \
		echo "Файлы логов очищены"; \
	else \
		echo "Очистка логов отменена"; \
	fi


# Сценарии работы с приложением

## Запустить приложение без контейнера
onelia-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	export POSTGRES_PORT=5434 && \
	go mod tidy && \
	go run ${PROJECT_ROOT}/cmd/onelia/main.go

## Поднять контейнер с приложением
onelia-deploy:
	@docker compose up -d --build onelia

## Остановить контейнер с приложением
onelia-undeploy:
	@docker compose down onelia
