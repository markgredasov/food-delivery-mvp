# Авито.Кухня — MVP

MVP-агрегатор доставки еды: API для заведений (приём заказов) и API для
клиентской части (каталог, заказы), а также сервис-заглушка одного заведения
`restaurant-simulator`, демонстрирующий интеграцию по этому API.

Сделано в рамках тестового задания для стажёра Avito Backend.

## Содержание

- [Быстрый старт](#быстрый-старт)
- [Основные сценарии (CJM)](#основные-сценарии-cjm)
- [Архитектура (C4)](#архитектура-c4)
- [Схема БД](#схема-бд)
- [API](#api)
- [MVP-упрощения](#mvp-упрощения-и-их-обоснование)

## Быстрый старт

```bash
docker compose up --build
```

Поднимутся три контейнера:

- `postgres` — PostgreSQL 18;
- `app` — основной сервис (порт `8080`), сам применяет миграции при старте;
- `restaurant-simulator` — сервис-заглушка одного заведения (порт `8081`),
  синхронизирует меню и обрабатывает заказы.

Проверка, что всё поднялось:

```bash
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/restaurants
```

Локальная разработка без Docker — см. `Makefile` (`make generate`, `make
migrate`, `make lint`, `make test`).

## Основные сценарии (CJM)

Диаграммы в `docs/diagrams/` (PlantUML, code-generated):

- [`cjm-user.puml`](docs/diagrams/cjm-user.puml) — путь пользователя: список
  заведений → меню → добавление позиций в заказ → создание заказа → отслеживание
  статуса до `delivered`/`rejected_by_restaurant`. Учтены краевые случаи:
  пустой список заведений, пустое меню, недоступная позиция (400/404),
  неактивное заведение (409).
- [`cjm-restaurant.puml`](docs/diagrams/cjm-restaurant.puml) — путь заведения:
  синхронизация меню → получение заказа (webhook, с fallback через polling) →
  обработка дубликатов → accept/reject → прогресс статусов до
  `delivered`, либо авто-отказ платформой по таймауту, если заведение не
  ответило за `AUTO_REJECT_AFTER` (по умолчанию 5 минут).

## Архитектура (C4)

- [`c4-context.puml`](docs/diagrams/c4-context.puml) — уровень 1 (Context):
  пользователь, заведение, сервис Авито.Кухня, внешняя система заведения.
- [`c4-container.puml`](docs/diagrams/c4-container.puml) — уровень 2
  (Container): API service (Go), PostgreSQL,
Основной сервис написан в стиле DDD-слоёв без ORM:

```
handler (HTTP, generated ServerInterface)
   -> service (бизнес-логика, транзакции)
      -> repository (pgxpool, без ORM)
         -> PostgreSQL
domain (money, address, restaurant, menu, order — сущности и инварианты,
        не зависят от остальных слоёв)
```

`restaurant-simulator` — независимый Go-модуль (свой `go.mod`), запускается
отдельным процессом/контейнером и общается с основным сервисом только по
HTTP — как это делал бы реальный сторонний интегратор.

## Схема БД

ER-диаграмма: [`docs/diagrams/db-schema.puml`](docs/diagrams/db-schema.puml).

Таблицы:

- **restaurants** — `id, name, address, status (active|inactive), service_url`.
  Только `active`-заведения принимают заказы.
- **categories** — `id, name, icon`.
- **menu_items** — `id, restaurant_id, category_id, name, price, available`.
  `price` — `NUMERIC(12,2)`, без отдельного поля валюты (см. упрощения).
- **orders** — `id, restaurant_id, user_id, delivery_address, status,
  total_amount`. `status` — текстовый enum с CHECK-constraint'ом,
  повторяющий стейт-машину заказа (см. `internal/domain/order`).
- **order_items** — `id, order_id, menu_item_id, quantity, price`.

Индексы: `menu_items(restaurant_id)`, `orders(restaurant_id, status)`,
`orders(created_at)` — под основные паттерны выборки (меню заведения,
активные заказы заведения, поиск просроченных для auto-reject).

## API

OpenAPI-спецификация: [`docs/swagger.yaml`](docs/swagger.yaml).

Идентификация вызывающей стороны — заголовками (полноценная авторизация не
входит в задание):

- `X-Restaurant-ID` — на эндпоинтах `/api/v1/restaurant/*`;
- `X-User-ID` — опционально, на `POST /api/v1/orders` (анонимные заказы
  допускаются).

Ключевые эндпоинты:

| Метод | Путь | Кто | Назначение |
|---|---|---|---|
| GET | `/api/v1/restaurants` | пользователь | список активных заведений |
| GET | `/api/v1/restaurants/{id}` | пользователь | детали заведения |
| GET | `/api/v1/restaurants/{id}/menu` | пользователь | меню заведения |
| GET | `/api/v1/categories` | пользователь | список категорий |
| POST | `/api/v1/orders` | пользователь | создать заказ |
| GET | `/api/v1/orders/{id}` | пользователь | статус заказа |
| GET | `/api/v1/restaurant/orders` | заведение | свои pending/активные заказы |
| PUT | `/api/v1/restaurant/menu` | заведение | заменить меню целиком |
| POST | `/api/v1/restaurant/orders/{id}/accept` | заведение | принять заказ |
| POST | `/api/v1/restaurant/orders/{id}/reject` | заведение | отклонить заказ |
| PATCH | `/api/v1/restaurant/orders/{id}/status` | заведение | продвинуть статус |

Состояния заказа (`internal/domain/order`):

```
pending ──────────────► sent_to_restaurant ──► accepted ──► preparing ──► ready ──► in_delivery ──► delivered
   │                            │
   └──────────► rejected_by_restaurant ◄───────┘
        (руками заведением, либо автоматически платформой по таймауту)
```

## MVP-упрощения и их обоснование

Явно вынесенные упрощения (для продакшена требуют доработки):

- **Нет авторизации/аутентификации** — по условиям задания не требуется;
  идентификация вызывающей стороны — заголовками `X-Restaurant-Id` /
  `X-User-Id`. В продакшене — полноценная авторизация.
- **Регистрация заведения — через seed-миграцию**, а не через отдельный API
  саморегистрации: в MVP список заведений закрытый ("некоторый список
  заведений", по условиям задания).
- **Нет поля валюты** — везде подразумевается RUB, `price`/`total_amount` —
  `NUMERIC(12,2)` без ISO-кода валюты.
- **`PUT /restaurant/menu` — soft-replace**: позиции, отсутствующие в новом
  payload, помечаются `available=false`, а не удаляются — чтобы не порвать
  FK у уже существующих `order_items`.
- **`restaurant-simulator` — in-memory хранилище** (`internal/store`,
  мьютекс + map) для идемпотентности webhook/poll — это заглушка, а не
  система записи; при перезапуске контейнера состояние теряется (заказ
  может быть обработан повторно, что безопасно за счёт идемпотентности на
  стороне основного сервиса — accept/reject/status проверяют текущий
  статус заказа).
- **Нет интеграционных тестов** (testcontainers и т.п.) — только unit-тесты
  домена и сервисного слоя (на фейковых репозиториях). Ограничение по
  времени задания; отмечено как возможное улучшение.