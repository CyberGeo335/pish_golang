# Практическое задание №10: JWT-аутентификация: создание и проверка токенов. Middleware для авторизации

`ФИО`: Козин Георгий Александрович

`Группа`: ПИМО-01-25

## Цель:
>🎯 Цель: создание и проверка токенов. Middleware для авторизации
---

> 💡 Important: В целях оптимизации рабочего пространства, было принято решение создать один репозиторий с дирикториями
> домашних работ.
---

## Задание:
* Понять устройство JWT и где его уместно применять в REST API.
* Сгенерировать и проверить JWT в Go (HS256), передавать его в Authorization: Bearer ….
* Реализовать middleware-аутентификацию (достаёт токен, валидирует, кладёт клеймы в context).
* Добавить middleware-авторизацию (RBAC/права на эндпоинты).
* Встроить это в уже знакомую архитектуру HTTP-сервиса/роутера.


### Описание проекта и требования:
#### Структура проекта:
```bash
├── auth.log
├── bin
│   └── auth
├── cmd
│   └── server
│       └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── core
│   │   ├── service.go
│   │   └── user.go
│   ├── http
│   │   ├── httputil
│   │   │   └── httputil.go
│   │   ├── middleware
│   │   │   ├── authn.go
│   │   │   ├── authz.go
│   │   │   └── logging.go
│   │   └── router.go
│   ├── platform
│   │   ├── config
│   │   │   └── config.go
│   │   └── jwt
│   │       └── jwt.go
│   └── repo
│       └── user_mem.go
└── keys
    ├── key1_private.pem
    ├── key1_public.pem
    ├── key2_private.pem
    └── key2_public.pem
```
#### Запуск проекта:
1) Клоним репозиторий:
```bash
git clone https://github.com/CyberGeo335/pish_golang.git
```
2) Проверяем что Go и Git есть:
```bash
prak_ten % go version
go version go1.23.2 darwin/arm64
prak_ten % git --version
git version 2.39.5 (Apple Git-154)
prak_ten % 
```
#### Команды:
##### Login как admin

![Скриншот запуска](./assets/Снимок%20экрана%202025-12-07%20в%2018.50.55.png)

```
curl -s -X POST "http://178.72.139.210:8084/api/v1/login" \
  -H "Content-Type: application/json" \
  -d '{"Email":"admin@example.com","Password":"secret123"}'
```
Ожидаемый JSON:

{"access_token":"...","refresh_token":"...","token_type":"Bearer"}
```bash
{
  "access_token":"eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleTEiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJwcmFrX3Rlbl9jbGllbnRzIiwiZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsImV4cCI6MTc2NTEyNDAzMSwiaWF0IjoxNzY1MTIzMTMxLCJpc3MiOiJwcmFrX3RlbiIsInJvbGUiOiJhZG1pbiIsInN1YiI6MSwidHlwIjoiYWNjZXNzIn0.VM6j3MfbicpQdufm2c4mxpJuX0oPvfb3__8o07xFUoMe8sJtyIkCKQJbHNDzIhKZ5_QiE9WMYanyhz_mxf1fsebgvkh2BaDBXHq_INGmaYFONbEzXjlIArxRyWuTttqNWUCH5xJff6G7clwAOTzVMlX2IKSAcSd0iZCRrbVGoOmuhJGv-YLml56A8Eu6v3k1rAaZdJEoq4waeVcY64dGxM1TZFh4_uTsx9Wwm7Z_FnCZkrHefak7FsOBH5t0FvVIu-GxYMpvQY88kGjgmlTeyMmANx_CQiFlOJf2QtROdM7N4Vda8fnLw0rHpdTiISPuI47CtDcsmU8u2c7MirQK9A",
  "refresh_token":"eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleTEiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJwcmFrX3Rlbl9jbGllbnRzIiwiZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsImV4cCI6MTc2NTcyNzkzMSwiaWF0IjoxNzY1MTIzMTMxLCJpc3MiOiJwcmFrX3RlbiIsInJvbGUiOiJhZG1pbiIsInN1YiI6MSwidHlwIjoicmVmcmVzaCJ9.VptafuWkZfIJregJ7PK3FtGeh8AV9DkxEqCDD4oOynEgwYzd7EEtyE3_SRsAnNAx3p9Bok4OIraGL6iTssPiIPuNxqyc8cw2Zpmam_rBr7RadUOiJYIzka43ask2RXTCYjDx5DJ5C7dGkjSpe2IORpZTgjf6oKkQe3MjMLuO43lcP8xZD0GOku1XVzyZH6q8D_o1wr_9NIlkfihub7lWO_KCNHWU4q5bzRnXrfo5ZtBeFEeKJMTp5ZKJvVMlauCQzeV4GKUVLUK3UFG4ApV_XBrZVmhu-2OIWAAdI6IPGPiUyRgdhWUfq1crQuJJM01U6507NSfrNXVJIedwl99iGg",
  "token_type":"Bearer"
  }
```

##### me от имени admin
```bash
curl -s "http://178.72.139.210:8084/api/v1/me" \
-H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleTEiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJwcmFrX3Rlbl9jbGllbnRzIiwiZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsImV4cCI6MTc2NTEyNDAzMSwiaWF0IjoxNzY1MTIzMTMxLCJpc3MiOiJwcmFrX3RlbiIsInJvbGUiOiJhZG1pbiIsInN1YiI6MSwidHlwIjoiYWNjZXNzIn0.VM6j3MfbicpQdufm2c4mxpJuX0oPvfb3__8o07xFUoMe8sJtyIkCKQJbHNDzIhKZ5_QiE9WMYanyhz_mxf1fsebgvkh2BaDBXHq_INGmaYFONbEzXjlIArxRyWuTttqNWUCH5xJff6G7clwAOTzVMlX2IKSAcSd0iZCRrbVGoOmuhJGv-YLml56A8Eu6v3k1rAaZdJEoq4waeVcY64dGxM1TZFh4_uTsx9Wwm7Z_FnCZkrHefak7FsOBH5t0FvVIu-GxYMpvQY88kGjgmlTeyMmANx_CQiFlOJf2QtROdM7N4Vda8fnLw0rHpdTiISPuI47CtDcsmU8u2c7MirQK9A"
```

![Скриншот запуска](./assets/Снимок%20экрана%202025-12-07%20в%2019.02.51.png)


##### /api/v1/admin/stats (admin-only)
```bash
curl -s "http://178.72.139.210:8084/api/v1/admin/stats" \
-H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleTEiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJwcmFrX3Rlbl9jbGllbnRzIiwiZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsImV4cCI6MTc2NTEyNDAzMSwiaWF0IjoxNzY1MTIzMTMxLCJpc3MiOiJwcmFrX3RlbiIsInJvbGUiOiJhZG1pbiIsInN1YiI6MSwidHlwIjoiYWNjZXNzIn0.VM6j3MfbicpQdufm2c4mxpJuX0oPvfb3__8o07xFUoMe8sJtyIkCKQJbHNDzIhKZ5_QiE9WMYanyhz_mxf1fsebgvkh2BaDBXHq_INGmaYFONbEzXjlIArxRyWuTttqNWUCH5xJff6G7clwAOTzVMlX2IKSAcSd0iZCRrbVGoOmuhJGv-YLml56A8Eu6v3k1rAaZdJEoq4waeVcY64dGxM1TZFh4_uTsx9Wwm7Z_FnCZkrHefak7FsOBH5t0FvVIu-GxYMpvQY88kGjgmlTeyMmANx_CQiFlOJf2QtROdM7N4Vda8fnLw0rHpdTiISPuI47CtDcsmU8u2c7MirQK9A"
```

![Скриншот запуска](./assets/Снимок%20экрана%202025-12-07%20в%2019.04.00.png)

##### Логин user: получить access/refresh
```bash
curl -s -X POST "http://178.72.139.210:8084/api/v1/login" \
  -H "Content-Type: application/json" \
  -d '{"Email":"user@example.com","Password":"secret123"}'
```

![Скриншот запуска](./assets/Снимок%20экрана%202025-12-07%20в%2019.04.57.png)

```bash
{
  "access_token":"eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleTEiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJwcmFrX3Rlbl9jbGllbnRzIiwiZW1haWwiOiJ1c2VyQGV4YW1wbGUuY29tIiwiZXhwIjoxNzY1MTI0MzkwLCJpYXQiOjE3NjUxMjM0OTAsImlzcyI6InByYWtfdGVuIiwicm9sZSI6InVzZXIiLCJzdWIiOjIsInR5cCI6ImFjY2VzcyJ9.tyD2hEPyC3Sqk859qMYatitkKuUpwtaYVEhVgmc1DvZpAOc8C2krvtgNPYi1tt6ZUc8TG5de-w3FkyKopiSy096gmHITh00IUb3MYE4FWgTURyikbZkROduxeInRenwXiyMhrcJcaU1BytkRm82DH451nnWTYrOAlScImcRZTNfc2d3qhBL4u8buhbwykrtzeyLapLJI5Myct_KlsHRG9Cf852UmEuSwy4_Yz8Lb6JpZ13Aa620tE9-vedeHn432AZ-wbrdcmjJHSL4NIIMzccxaxm884mNibtT5fOdSg8DcdevJx6GXcb-fPD3Odcd2O7TwjdpJkRvSPRG37ziTzg",
  "refresh_token":"eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleTEiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJwcmFrX3Rlbl9jbGllbnRzIiwiZW1haWwiOiJ1c2VyQGV4YW1wbGUuY29tIiwiZXhwIjoxNzY1NzI4MjkwLCJpYXQiOjE3NjUxMjM0OTAsImlzcyI6InByYWtfdGVuIiwicm9sZSI6InVzZXIiLCJzdWIiOjIsInR5cCI6InJlZnJlc2gifQ.Z377lMQXXF4JdnX1IQW3xgnubwo6k8BWtPWoSY9c2PB2yvrFN8aijfB9sbKM5SB32l5qPQKC94l_okFG7FeySVZ-OgIhbZmUlIhPSPqJxD_IC3enWhKe8P-cr2ImFZdrv0-Uvt6WdlfnaB3AQUW7hNcxBt9PXxVttP9H_8WkpnaYPuyoRbnWG_KG0bxmATdnA6rw2iyzSCM75ZFtMJeMTndlh4louKS1uHOSBGkU7EW6kl2a8pJ9Dh0qH5rAwKyCVqMO2C6YYoIk-gA2iKJWx2QUlI_DqVEjhFq9RtgfUATiNvxSFky9bEMAmS3jk8MGPEOBt3y4sj0rD2Zm5w1X7g",
  "token_type":"Bearer"
}
```

##### Попытка залезть к /api/v1/users/1 (admin) — должно быть 403 (5 рис)
```bash
curl -i "http://178.72.139.210:8084/api/v1/users/1" \
  -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleTEiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJwcmFrX3Rlbl9jbGllbnRzIiwiZW1haWwiOiJ1c2VyQGV4YW1wbGUuY29tIiwiZXhwIjoxNzY1MTI0MzkwLCJpYXQiOjE3NjUxMjM0OTAsImlzcyI6InByYWtfdGVuIiwicm9sZSI6InVzZXIiLCJzdWIiOjIsInR5cCI6ImFjY2VzcyJ9.tyD2hEPyC3Sqk859qMYatitkKuUpwtaYVEhVgmc1DvZpAOc8C2krvtgNPYi1tt6ZUc8TG5de-w3FkyKopiSy096gmHITh00IUb3MYE4FWgTURyikbZkROduxeInRenwXiyMhrcJcaU1BytkRm82DH451nnWTYrOAlScImcRZTNfc2d3qhBL4u8buhbwykrtzeyLapLJI5Myct_KlsHRG9Cf852UmEuSwy4_Yz8Lb6JpZ13Aa620tE9-vedeHn432AZ-wbrdcmjJHSL4NIIMzccxaxm884mNibtT5fOdSg8DcdevJx6GXcb-fPD3Odcd2O7TwjdpJkRvSPRG37ziTzg"
```

![Скриншот запуска](./assets/Снимок%20экрана%202025-12-07%20в%2019.12.43.png)

##### Admin может читать любого пользователя (6 рис)
```bash
curl -s "http://178.72.139.210:8084/api/v1/users/1" \
  -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleTEiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJwcmFrX3Rlbl9jbGllbnRzIiwiZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsImV4cCI6MTc2NTEyNTI4NywiaWF0IjoxNzY1MTI0Mzg3LCJpc3MiOiJwcmFrX3RlbiIsInJvbGUiOiJhZG1pbiIsInN1YiI6MSwidHlwIjoiYWNjZXNzIn0.I4nq344la1WatTBGlumUvoxAcM-ss0-VHylfyXJIljyV_MgBwciY1Eop9__-yuILNwCbklV9emzeBxuJCq6TYMyKzY3Cu1H4sRxLZ8pYuz0rLToaB9R_EuGZQpGAsLFDNJkEB7CU0h-GTbPeBdedTH2LatX9s8rLfUjwQV8-qPhBcDZgOHPDTqsRLuL_19qaj0MlkZ_mpD-uvBoCgp6SuN4MJ0oJxb8YTh8IbZ_mSCpiMlJ09lkNjbiJOpkLarBG7pEDURdY1s_E2b68HdRk_cMI5yZtFYWNyG_Su0RVeJQ3avn4pZlghS_GOtkZcpUUxX3eSC_sLEWIeCLgHP9MlQ"```
```

![Скриншот запуска](./assets/Снимок%20экрана%202025-12-07%20в%2019.21.22.png)

##### Вызов /api/v1/refresh (рис 7)

```bash
curl -s -X POST "http://178.72.139.210:8084/api/v1/refresh" \
  -H "Content-Type: application/json" \
  -d "{\"refresh_token\":\"eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleTEiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJwcmFrX3Rlbl9jbGllbnRzIiwiZW1haWwiOiJ1c2VyQGV4YW1wbGUuY29tIiwiZXhwIjoxNzY1NzI4MjkwLCJpYXQiOjE3NjUxMjM0OTAsImlzcyI6InByYWtfdGVuIiwicm9sZSI6InVzZXIiLCJzdWIiOjIsInR5cCI6InJlZnJlc2gifQ.Z377lMQXXF4JdnX1IQW3xgnubwo6k8BWtPWoSY9c2PB2yvrFN8aijfB9sbKM5SB32l5qPQKC94l_okFG7FeySVZ-OgIhbZmUlIhPSPqJxD_IC3enWhKe8P-cr2ImFZdrv0-Uvt6WdlfnaB3AQUW7hNcxBt9PXxVttP9H_8WkpnaYPuyoRbnWG_KG0bxmATdnA6rw2iyzSCM75ZFtMJeMTndlh4louKS1uHOSBGkU7EW6kl2a8pJ9Dh0qH5rAwKyCVqMO2C6YYoIk-gA2iKJWx2QUlI_DqVEjhFq9RtgfUATiNvxSFky9bEMAmS3jk8MGPEOBt3y4sj0rD2Zm5w1X7g\"}"
```

![Скриншот запуска](./assets/Снимок%20экрана%202025-12-07%20в%2019.29.47.png)

```bash
{
  "access_token":"eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleTEiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJwcmFrX3Rlbl9jbGllbnRzIiwiZW1haWwiOiJ1c2VyQGV4YW1wbGUuY29tIiwiZXhwIjoxNzY1MTI1NTQxLCJpYXQiOjE3NjUxMjQ2NDEsImlzcyI6InByYWtfdGVuIiwicm9sZSI6InVzZXIiLCJzdWIiOjIsInR5cCI6ImFjY2VzcyJ9.QprKQ8j90gCKIT8ykBqI9A740TczkeJIhJff7mZBBHp_azcdVpOLBpqjhe7cdr068F7NsDaNqLoqloRzzC19O8PVrwg3Q0w-dpEMaMW9k1pqSu6h4fzdQgV-h6_LoDGJaqNmzaQJ46vRjxUINQPduut3-i5fjseE3MVKopDDV5D-zkkZJ6ENiJwm_kv3S6k56DGoTT0ajhrWO27-UUtcI0BvX9UrdyuZ85xniTftDzqfutalVNP-lQMqyAeWhn86acA8G-WR68qJcYU3c9MrSvOI0oqgrtKseuDI6I1cqlLAdxKiVIXY4-A3aPmJgwQMnbjPfE45arA4tk8RKpPuWA",
  "refresh_token":"eyJhbGciOiJSUzI1NiIsImtpZCI6ImtleTEiLCJ0eXAiOiJKV1QifQ.eyJhdWQiOiJwcmFrX3Rlbl9jbGllbnRzIiwiZW1haWwiOiJ1c2VyQGV4YW1wbGUuY29tIiwiZXhwIjoxNzY1NzI5NDQxLCJpYXQiOjE3NjUxMjQ2NDEsImlzcyI6InByYWtfdGVuIiwicm9sZSI6InVzZXIiLCJzdWIiOjIsInR5cCI6InJlZnJlc2gifQ.p7yeZx96wLg5KLfAKhh2_h-3p9Zk-owndeFnEc0YqUyuogVXJT46ypfgTaExW6KSO7rO9StDXr2iCvQ3D5EgswRUHhJnlItTX5A29-u9Z-2BnzNZHWnghCxPLSMfJ8nnDBfdgBs-S6Z8gxY_7XUq9a4eftARwunu7AzEC4UqtuxI-gwXnAdd21utmj6fX580PqCSJDlhr8qkx0W3ACNNnco2KNlltqkuK1K3yIJysq4TZdpU1JXk7Mn8aYoBwiVjkSBQhvDcvHJl23zdHJhCE84KFTN0Vx2vY3_4J0W6XGfoUknK6F-pB5RwQk7XwAfticVsoaVxcYi362GLJlgptQ",
  "token_type":"Bearer"
}
```