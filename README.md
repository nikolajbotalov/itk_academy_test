# Wallet Service

**REST API для управления балансом кошельков с нагрузкой 1000 RPS.**

---

## Требования

- Golang 1.23+
- PostgreSQL 15
- Docker + docker-compose

---

## Особенности

| Функция | Реализация |
|--------|-----------|
| **1000 RPS на один кошелёк** | `SELECT ... FOR UPDATE` |
| **Ни одного 50X** | Транзакции + retry |
| **Миграции** | `golang-migrate` |

---

## API

### `POST /api/v1/wallet`

```json
{
  "wallet_id": "a3f1c9e7-2b4d-4f8a-9e6c-1d5f8b7a3e2c",
  "operation_type": "DEPOSIT",
  "amount": 100
}
```

### `GET /api/v1/wallet/{wallet_uuid}`

```json
{
  "wallet_id": "a3f1c9e7-2b4d-4f8a-9e6c-1d5f8b7a3e2c",
  "balance": 10000,
  "created_at": "2025-01-01T10:00:00Z",
  "updated_at": "2025-01-01T10:00:00Z"
}
```

### Клонировать репозиторий
git clone https://github.com/nikobotal/wallet-service.git

### Запуск
docker-compose up --build