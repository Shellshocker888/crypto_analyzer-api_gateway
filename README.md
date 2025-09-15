# Crypto Analyzer — API Gateway

## Table of Contents
- [Overview](#overview)
- [Project Structure](#project-structure)
- [Requirements](#requirements)
- [API Reference](#api-reference)
- [Architecture](#architecture)
- [Security](#security)
- [Monitoring](#monitoring)


## Overview

Crypto Analyzer API Gateway — высокопроизводительный шлюз для микросервисов Auth и Portfolio, реализованный на Go с использованием Fiber и gRPC.

Основные возможности:
- Роутинг запросов к микросервисам Auth и Portfolio
- Интеграция с gRPC-сервисами
- Обработка аутентификации и авторизации через middleware
- Логирование и трассировка запросов
- Поддержка публичных и приватных портфелей


## Project Structure

````
├── cmd/                   # Точка входа (main.go)
├── gen/                   # Сгенерированные файлы gRPC
│   ├── go/auth/
│   └── go/portfolio/
├── internal/
│   ├── app/               # Сборка и запуск сервиса
│   ├── config/            # Конфигурация и модели
│   ├── controller/        # HTTP/GRPC контроллеры и middleware
│   │   ├── auth/
│   │   ├── middleware/
│   │   └── portfolio/
│   ├── domain/            # Сущности и бизнес-ошибки
│   │   ├── auth/
│   │   └── portfolio/
│   ├── infrastructure/    # Логгер, gRPC клиенты, Postgres, Redis
│   │   ├── auth/grpc/
│   │   ├── logger/
│   │   └── portfolio/grpc/
│   └── usecase/           # Бизнес-логика
│       ├── auth/
│       └── portfolio/
├── logs/                  # Логи сервиса
├── migrations/            # SQL миграции (если есть)
├── proto/                 # Описание gRPC API
├── tests/                 # Интеграционные тесты
├── .env
├── .gitignore
├── go.mod
├── go.sum
└── README.md
````


## Requirements

- Go 1.20+
- PostgreSQL 13+
- Redis 6+
- Protoc + gRPC plugin
- Prometheus 2+ (для сбора метрик)
- Alertmanager 0.27+ (для обработки алертов)
- Grafana 9+ (для визуализации метрик)



## API Reference
````
Auth endpoints (через gateway, проксируются на Auth Service):

GET /auth/ping — проверка работоспособности Auth Service

Portfolio endpoints:

POST /portfolios — создать новый портфель
GET /portfolio/:id — получить содержимое портфеля
POST /portfolio/:id/asset — добавить или обновить актив
DELETE /portfolio/:id/asset — удалить актив
GET /portfolios — получить все портфели пользователя
GET /portfolio/:id/history — история стоимости портфеля
GET /portfolio/public/:username — публичные портфели другого пользователя

Все защищённые методы используют middleware AuthVerify для проверки токена.
````

## Architecture
````
Сервис построен по принципам DDD и чистой архитектуры:
Controller — обработка HTTP-запросов и интеграция с middleware
Usecase / Service — бизнес-логика для портфелей и пользователей
Domain — сущности и ошибки бизнес-логики
Infrastructure — клиенты gRPC, логирование, Postgres, Redis
Middleware — аутентификация, логирование, трассировка
````

## Security
````
JWT Access и Refesh токены через Auth Service
Middleware проверяет Authorization header
Логирование trace-id для каждого запроса
````

## Monitoring

Для мониторинга состояния API Gateway используется **Prometheus** и **Alertmanager**.  
Все метрики и алерты описаны подробно в отдельной документации: [docs/monitoring.md](docs/monitoring.md)
