# Практическое задание №14: Unit-тестирование функций (testing, testify)

`ФИО`: Козин Георгий Александрович

`Группа`: ПИМО-01-25

## Цель:
>🎯 Цель: Unit-тестирование функций (testing, testify)
---

> 💡 Important: В целях оптимизации рабочего пространства, было принято решение создать один репозиторий с дирикториями
> домашних работ.
---

## Задание:
* Освоить базовые приёмы unit-тестирования в Go с помощью стандартного пакета testing.
* Научиться писать табличные тесты, подзадачи t.Run, тестировать ошибки и паники.
* Освоить библиотеку утверждений testify (assert, require) для лаконичных проверок.
* Научиться измерять покрытие кода (go test -cover) и формировать html-отчёт покрытия.
* Подготовить минимальную структуру проектных тестов и общий чек-лист качества тестов.



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
pish_golang % cd prak_fifteen
prak_fifteen % go version
go version go1.23.2 darwin/arm64
prak_fifteen % git --version
git version 2.39.5 (Apple Git-154)
prak_fifteen % 
```

3) Переходим в пятнадцатую работу:

```bash
cd prak_fifteen/cmd
```

4) Запуск на удалённом сервере:
```bash
go build -o bin/prak_fifteen ./cmd/api/

nohup ./bin/prak_fifteen >prak_fourteen.log 2>&1 &
```
