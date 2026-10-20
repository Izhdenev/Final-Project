# Final-Project — Планировщик задач

## Описание проекта

Веб-сервер на Go, реализующий бэкенд для планировщика задач (todo-приложения).
Фронтенд уже был предоставлен в директории `web/`, задача состояла в том,
чтобы написать REST API, работающее с SQLite-базой данных.

Сервер умеет:
- добавлять, редактировать, удалять и отмечать выполненными задачи;
- возвращать список ближайших задач, отсортированных по дате;
- искать задачи по подстроке в заголовке/комментарии или по конкретной дате;
- вычислять следующую дату для повторяющихся задач по правилам вида `d <n>` (каждые n дней) и `y` (ежегодно);
- поддерживать аутентификацию по паролю с использованием JWT.

## Выполненные задания со звёздочкой

- **Поиск задач.** Реализован поиск в `/api/tasks?search=...`:
    - если `search` парсится как дата в формате `02.01.2006` — фильтр по колонке `date`;
    - иначе — поиск по подстроке в `title` и `comment`.
- **Аутентификация.** Реализован обработчик `/api/signin`, выдающий JWT-токен,
  и middleware `auth`, защищающий `/api/task`, `/api/tasks`, `/api/task/done`.
  Проверка активируется только при непустой переменной окружения `TODO_PASSWORD`.
- **Docker.** Собран multi-stage `Dockerfile` (сборка в `golang:1.26`, запуск в `alpine:latest`).
  SQLite-база монтируется с хоста, фронтенд и бинарник кладутся внутрь образа.

## Локальный запуск

### Требования

- Go 1.26 или новее;
- SQLite-база создаётся автоматически при первом запуске.

### Запуск

```bash
# 1. Скачать зависимости
go mod tidy

# 2. (опционально) задать переменные окружения
export TODO_PORT=7540
export TODO_DBFILE=tests/scheduler.db
export TODO_PASSWORD=12345

# 3. Запустить сервер
go run main.go
```

- **Без пароля:** открой **http://localhost:7540/**
- **С паролем:** открой **http://localhost:7540/login.html**, введи пароль,
  после чего установится кука `token` и ты попадёшь в планировщик.

## Инструкция по запуску тестов

Тесты требуют запущенного сервера. В одном терминале:

```bash
rm -f tests/scheduler.db
go run main.go
```

### Параметры в `tests/settings.go`

```bash
var Port = 7540
var DBFile = "scheduler.db"
var FullNextDate = false    // доп. кейсы для правил w и m(не реализовывал)
var Search = true          // проверка /api/tasks?search=
var Token = ""             // заполнить, только если задан TODO_PASSWORD
var WebDir = "web"
```

## Сборка и запуск через Docker

```bash
# Сборка образа
docker build -t scheduler:latest .

# Папка под БД на хосте
mkdir -p ~/scheduler-data

# Запуск с паролем
docker run --rm -it \
  --name scheduler \
  -p 7540:7540 \
  -v ~/scheduler-data:/data \
  -e TODO_PASSWORD=12345 \
  scheduler:latest

# Запуск без пароля (аутентификация отключена)
docker run --rm -it \
  --name scheduler \
  -p 7540:7540 \
  -v ~/scheduler-data:/data \
  scheduler:latest
```

- Внутри контейнера: `TODO_PORT=7540`, `TODO_DBFILE=/data/scheduler.db`.
- SQLite-файл живёт на хосте в `~/scheduler-data/scheduler.db`.
- В браузере: **http://localhost:7540/** (или `/login.html` при заданном пароле).

В фоне: `docker run -d ...`, логи — `docker logs -f scheduler`.