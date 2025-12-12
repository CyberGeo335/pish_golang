# Практическое задание №14: Оптимизация запросов к БД. Использование connection pool

`ФИО`: Козин Георгий Александрович

`Группа`: ПИМО-01-25

## Цель:
>🎯 Цель: Оптимизация запросов к БД. Использование connection pool
---

> 💡 Important: В целях оптимизации рабочего пространства, было принято решение создать один репозиторий с дирикториями
> домашних работ.
---

## Задание:
1.	Научиться находить «узкие места» в SQL-запросах и устранять их (индексы, переписывание запросов, пагинация, батчинг).
2.	Освоить настройку пула подключений (connection pool) в Go и параметры его тюнинга.
3.	Научиться использовать EXPLAIN/ANALYZE, базовые метрики (pg_stat_statements), подготовленные запросы и транзакции.
4.	Применить техники уменьшения N+1 запросов и сокращения аллокаций на горячем пути.


### Описание проекта и требования:
#### Структура проекта:
```
├── README.md
├── cmd
│   └── api
│       └── main.go
├── docker-compose.yaml
├── go.mod
├── go.sum
└── internal
    ├── config
    │   └── config.go
    ├── model
    │   └── note.go
    ├── pagination
    │   └── cursor.go
    ├── storage
    │   ├── postgres
    │   │   ├── queries.go
    │   │   └── repo.go
    │   └── redis
    │       └── cache.go
    └── transport
        └── http
            ├── handlers.go
            ├── respond.go
            └── server.go
```

#### Запуск проекта:
1) Клоним репозиторий:
```bash
git clone https://github.com/CyberGeo335/pish_golang.git
```
2) Проверяем что Go и Git есть:

```bash
pish_golang % cd prak_fourteen
prak_fourteen % go version
go version go1.23.2 darwin/arm64
prak_fourteen % git --version
git version 2.39.5 (Apple Git-154)
prak_fourteen % 
```

3) Переходим в девятую работу:

```bash
cd prak_fourteen/cmd
```

4) Пример нашего `.env`:
```bash
# Remote Postgres
DB_DSN=postgres://root:root@http://address:5432/pz9_bcrypt?sslmode=disable

# HTTP
HTTP_ADDR=:8087

# Redis (локально через docker compose)
REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=kek
REDIS_DB=0
CACHE_TTL_SECONDS=45

```

P.S: in process ))))