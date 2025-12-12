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

5) Запуск на удалённом сервере:
```bash
go build -o bin/prak_fourteen ./cmd/api/

nohup ./bin/prak_fourteen >prak_fourteen.log 2>&1 &
```

6) проверка что всё равботает (пикча 3):
```bash
curl -s http://178.72.139.210:8087/health
```

![Скриншот](./assets/Снимок%20экрана%202025-12-12%20в%2012.13.13.png)
![Скриншот](./assets/Снимок%20экрана%202025-12-12%20в%2013.18.17.png)
![Скриншот](./assets/Снимок%20экрана%202025-12-12%20в%2013.26.25.png)



7) Создадим заметку (пикча 4):
```bash
curl -s -X POST http://178.72.139.210:8087/notes \
  -H 'Content-Type: application/json' \
  -d '{"title":"redis tips","content":"..."}'
```

![Скриншот](./assets/Снимок%20экрана%202025-12-12%20в%2013.27.44.png)

8) Получить по Id, где второй должен кэшироваться  (пикча 5):
```bash
curl -s http://178.72.139.210:8087/notes/1
curl -s http://178.72.139.210:8087/notes/1
```

![Скриншот](./assets/Снимок%20экрана%202025-12-12%20в%2013.29.21.png)

9) Список (keyset) + поиск по title (FTS) (пикча 6):
```bash
curl -s "http://178.72.139.210:8087/notes?limit=20&q=redis"
curl -s "http://178.72.139.210:8087/notes?limit=20&q=redis&cursor="eyJjcmVhdGVkX2F0IjoiMjAyNS0xMi0xMlQxMDoyNzozNy43NjU3WiIsImlkIjoxfQ"
```

![Скриншот](./assets/Снимок%20экрана%202025-12-12%20в%2013.44.37.png)

10) Список (keyset) + поиск по title (FTS) (пикча 7):
```bash
curl -s -X PATCH http://178.72.139.210:8087/notes/1 \
  -H 'Content-Type: application/json' \
  -d '{"title":"redis tips v2","content":"updated"}'
```

![Скриншот](./assets/Снимок%20экрана%202025-12-12%20в%2013.45.02.png)

11) Батч вместо N+1 (пикча 8):
```bash
curl -s "http://178.72.139.210:8087/notes/batch?ids=1,2,3,4"
```

![Скриншот](./assets/Снимок%20экрана%202025-12-12%20в%2013.46.32.png)


### Пояснялки:
#### 2
До оптимизации было три неприятных пункта:
* Первая — пагинация через OFFSET: чем дальше страница, тем больше строк базе приходится просканировать и, поэтому время
растёт почти линейно с номером страницы. 
* Вторая — N+1: сами запросы по первичному ключу выглядят быстрыми, но когда их N штук, стоимость сети и количества 
round-trip’ов доминирует и p95/p99 резко растут. 
* Третья — поиск по title: если запрос написан не в точности под выражение индекса, Postgres не может использовать GIN 
и уходит в Seq Scan, что в EXPLAIN видно сразу: Seq Scan + Filter и много BUFFERS read.
#### 3
Применили два вида изменений: 
* Переписывание запросов и индексы. Запросы: OFFSET заменили на keyset-пагинацию по паре (created_at, id), чтобы база 
продолжала чтение “с места”, а не пересчитывала всю глубину OFFSET. N+1 заменили на один батч-запрос по массиву 
id через ANY, чтобы сделать один запрос вместо N. Поиск по title привели к форме, которая точно совпадает с индексом,
иначе индекс просто не применится. Из индексов: оставили/использовали композитный btree на (created_at, id) под keyset и
GIN на tsvector под FTS. Если поиск становится очень горячим, следующий шаг — хранить tsvector как STORED/generated 
column и индексировать уже колонку: это убирает постоянный пересчёт to_tsvector на каждом запросе и часто снижает CPU. 
#### 4
Пул настроили конкретно так: 
* MaxConns=20, 
* MinConns=5, 
* MaxConnLifetime=1h
Оптимально
