SHELL := /bin/bash

# Defaults (можно переопределить: make dev HTTP_ADDR=:9090)
HTTP_ADDR     ?= :8081
DB_HOST       ?= localhost
DB_PORT       ?= 5432
DB_USER       ?= postgres
DB_PASSWORD   ?= postgres
DB_NAME       ?= postgres
DB_SSLMODE    ?= disable
KAFKA_BROKERS ?= localhost:29092
KAFKA_TOPIC   ?= orders
KAFKA_GROUP_ID?= orders-consumer

.PHONY: dev up migrate run down stop build tidy

dev: up migrate run

up:
	docker compose -f ./docker-compose.yaml up -d postgres zookeeper kafka kafka-ui

migrate:
	docker compose -f ./docker-compose.yaml --profile migrate up migrate --exit-code-from migrate

run:
	@echo "Starting app on $(HTTP_ADDR)"
	HTTP_ADDR=$(HTTP_ADDR) \
	DB_HOST=$(DB_HOST) DB_PORT=$(DB_PORT) DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) DB_NAME=$(DB_NAME) DB_SSLMODE=$(DB_SSLMODE) \
	KAFKA_BROKERS=$(KAFKA_BROKERS) KAFKA_TOPIC=$(KAFKA_TOPIC) KAFKA_GROUP_ID=$(KAFKA_GROUP_ID) \
	go run ./cmd

down:
	docker compose -f ./docker-compose.yaml down -v

stop:
	docker compose -f ./docker-compose.yaml stop

build:
	go build ./...

tidy:
	go mod tidy


