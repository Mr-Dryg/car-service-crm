# CRM для автосервиса

## Запуск

### 1. Создать `.env` в корне проекта

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=db_user
DB_PASSWORD=db_password
DB_NAME=car_service
DB_SSLMODE=disable
```

### 2. Запустить котейнер
```
$ docker compose up -d --build
``` 