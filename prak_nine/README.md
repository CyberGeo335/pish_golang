# Практическое задание №9: Реализация регистрации и входа пользователей. Хэширование паролей с bcrypt

`ФИО`: Козин Георгий Александрович

`Группа`: ПИМО-01-25

## Цель:
>🎯 Цель: реализовать регистрацию и входа пользователей. Хэширование паролей с bcrypt
---

> 💡 Important: В целях оптимизации рабочего пространства, было принято решение создать один репозиторий с дирикториями
> домашних работ.
---

## Задание:
* Научиться безопасно хранить пароли (bcrypt), валидировать вход и обрабатывать ошибки;
* Реализовать эндпоинты POST /auth/register и POST /auth/login;
* Закрепить работу с БД (PostgreSQL + GORM или database/sql) и валидацией ввода;

### Описание проекта и требования:
#### Структура проекта:
```
├── prak_nine
│   ├── README.md
│   ├── assets
│   ├── cmd
│   │   └── main.go
│   ├── go.mod
│   ├── go.sum
│   └── internal
│       ├── app
│       │   └── app.go
│       ├── core
│       │   └── user.go
│       ├── http
│       │   └── handlers
│       │       └── auth.go
│       ├── platform
│       │   └── config
│       │       └── config.go
│       └── repo
│           ├── postgres.go
│           └── user_repo.go
```
#### Запуск проекта:
1) Клоним репозиторий:
```bash
git clone https://github.com/CyberGeo335/pish_golang.git
```
2) Проверяем что Go и Git есть:
```bash
prak_nine % go version
go version go1.23.2 darwin/arm64
prak_nine % git --version
git version 2.39.5 (Apple Git-154)
prak_nine % 
```
3) Переходим в девятую работу:
```bash
cd prak_nine/cmd
```

4) Пример нашего `.env`:
```bash
DB_DSN=postgres://root:root@http://address:5432/pz9_bcrypt?sslmode=disable
```


3) Проверим, что наши ручки работают:
```bash
curl http://localhost:8081/hello
curl http://localhost:8081/user
```
![Скриншот запуска](./assets/Снимок%20экрана%202025-10-09%20в%2014.16.05.png)

Ручки отработали, ответ получен. Так же была сделана проверка на форматирование кода:

![Скриншот запуска](./assets/Снимок%20экрана%202025-10-09%20в%2013.07.05.png)