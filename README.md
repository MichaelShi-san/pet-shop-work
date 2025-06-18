# PetShop API (Go + PostgreSQL)

RESTful API-проект для интернет-магазина товаров для питомцев. Реализован на языке Go с использованием PostgreSQL. Структура проекта разделена на модули с пакетами `internal` и `models`. Поддерживаются базовые CRUD-операции, заказы, транзакции и аналитика.

## 🚀 Как запустить проект

### 1. Клонировать репозиторий

```bash
git clone https://github.com/your-username/petshop-api.git
cd petshop-api

Как настроить базу данных:

Создайте базу данных PostgreSQL и пропишите параметры подключения в .env:

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=petshop

Запустите миграции:

go run cmd/migrate/main.go

Запустить сервер:

go run cmd/server/main.go

📦 Описание версий

✅ Версия v1 — CRUD для пользователей и товаров

- Реализованы таблицы: users, products.
- CRUD-операции через REST API (GET, POST, PUT, DELETE).

✅ Версия v2 — Заказы и позиции заказов

- Добавлены таблицы: orders, order_items.

- Реализован POST /orders с вложенными позициями.

- Добавлен GET /orders/{id} с JOIN для отображения позиций.

✅ Версия v3 — Транзакции и история заказов

- Добавлена таблица: transactions.

- Реализована ручка GET /users/{email}/history — история заказов с JOIN по транзакциям.

✅ Версия v4 — Аналитика популярных товаров

- Добавлена ручка GET /products/popular.

- Использован SQL-запрос с GROUP BY и ORDER BY для подсчета количества продаж.

📌 TODO

- Аутентификация (JWT).

- Разделение прав доступа (админ/пользователь).

- Покрытие тестами.

- Docker-окружение.

🛠️ Стек технологий

- Go
- PostgreSQL
- pgxpool
- net/http
- chi router

