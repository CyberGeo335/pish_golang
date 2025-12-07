# Практическое задание №12: Подключение Swagger/OpenAPI. Автоматическая генерация документации
`ФИО`: Козин Георгий Александрович

`Группа`: ПИМО-01-25

## Цель:
>🎯 Подключение Swagger/OpenAPI. Автоматическая генерация документации
---

> 💡 Important: В целях оптимизации рабочего пространства, было принято решение создать один репозиторий с дирикториями
> домашних работ.
---

## Задание:
1.	Освоить основы спецификации OpenAPI (Swagger) для REST API.
2.	Подключить автогенерацию документации к проекту из ПЗ 11 (notes-api).
3.	Научиться публиковать интерактивную документацию (Swagger UI / ReDoc) на эндпоинте GET /docs.
4.	Синхронизировать код и спецификацию (комментарии-аннотации → генерация) и/или «schema-first» (генерация кода из openapi.yaml).
5.	Подготовить процесс обновления документации (Makefile/скрипт).


### Описание проекта и требования:
#### Структура проекта:
```
├── README.md
├── assets
├── cmd
│   └── api
│       └── main.go
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── go.mod
├── go.sum
└── internal
    ├── api
    │   └── openapi.yaml
    ├── core
    │   ├── note.go
    │   └── service
    │       └── note_service.go
    ├── http
    │   ├── handlers
    │   │   └── notes.go
    │   └── router.go
    └── repo
        └── note_mem.go


```
#### Запуск проекта:
1) Клоним репозиторий:
```bash
git clone https://github.com/CyberGeo335/pish_golang.git
```
2) Проверяем что Go и Git есть:
```bash
pish_golang % cd prak_twelve

prak_twelve % go version
go version go1.23.2 darwin/arm64

prak_twelve % git --version
git version 2.39.5 (Apple Git-154)

prak_twelve % 
```

2) создаём SWAGGER:
```bash
swag init -g cmd/api/main.go -o docs
```
Стоит отметить, что есть ошибка, связанная с `LeftDelim` и `RightDelim` (cносим их к чёртовому С В А Г Г Е Р У)

#### Демонстрация:
1) URL сваггера:
```bash
http://178.72.139.210:8085/docs/index.html#/
```
2) Скриншот работающей страницы Swagger UI 

![Скриншот](./assets/Снимок%20экрана%202025-12-07%20в%2022.11.27.png)
![Скриншот](./assets/Снимок%20экрана%202025-12-07%20в%2022.12.41.png)
![Скриншот](./assets/Снимок%20экрана%202025-12-07%20в%2022.13.29.png)