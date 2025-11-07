# Практическое задание №4: Маршрутизация с chi. Создание небольшого CRUD-сервиса "Список задач"
`ФИО`: Козин Георгий Александрович

`Группа`: ПИМО-01-25

## Цель:
>🎯 Освоить базовую маршрутизацию HTTP-запросов в Go на примере роутера chi
>   Научиться строить REST-маршрутиы и обрабатывать методы GET/POST/PUT/DELETE
>   Реализовать небольшой CRUD-сервис "ToDo" без БД
>   Добавить простое middleware (логирование CORS)
>   Научиться тестировать API запросами через curl/Postman
---

> 💡 Important: В целях оптимизации рабочего пространства, было принято решение создать один репозиторий с дирикториями
> домашних работ.
---

## Задание:
* Освоить базовую маршрутизацию HTTP-запросов в Go на примере роутера chi 
* Научиться строить REST-маршрутиы и обрабатывать методы GET/POST/PUT/DELETE 
* Реализовать небольшой CRUD-сервис "ToDo" без БД 
* Добавить простое middleware (логирование CORS)
* Научиться тестировать API запросами через curl/Postman

### Описание проекта и требования:
#### Структура проекта:
```
└── prak_two
    ├── README.md
    ├── assets
    └── myapp
        ├── bin
        │   └── myapp
        ├── binmyapp.exe
        ├── cmd
        │   └── myapp
        │       └── main.go
        ├── go.mod
        ├── internal
        │   └── app
        │       ├── app.go
        │       └── handlers
        │           ├── fail.go
        │           ├── ping.go
        │           └── root.go
        └── utils
            ├── httpjson.go
            └── logger.go
```

#### Запуск проекта:
1) Клоним репозиторий:
```bash
git clone https://github.com/CyberGeo335/pish_golang.git
```
2) Проверяем что Go и Git есть:
```bash
prak_four % go version
go version go1.23.2 darwin/arm64
prak_four % git --version
git version 2.39.5 (Apple Git-154)
prak_four % 
```
3) Переходим в четвертую домашнюю работу:
```bash
cd prak_four/myapp
```

#### Проверка работоспособности:
1) Запустим наше приложение:
```bash
go run ./cmd/myapp
```
Покажем, что всё работает.

Наши ручки:
```bash
# first curl
curl http://localhost:8080/

# second curl
curl http://localhost:8080/ping

# third curl
curl http://localhost:8080/fail

# fourth curl
curl -i -H "X-Request-Id: demo-123" http://localhost:8080/ping
```

![Скриншот запуска](./assets/Снимок%20экрана%202025-10-09%20в%2015.44.11.png)