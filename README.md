# Link Tracker

Микросервисный инструмент для SEO-команды: мониторинг бэклинков, проверка индексации URL, мониторинг здоровья сайтов.

## Возможности

- **Мониторинг бэклинков** — проверка наличия ссылок и анкоров на внешних сайтах
- **Проверка индексации** — мониторинг индексации URL в поисковых системах
- **Мониторинг здоровья** — отслеживание доступности и состояния сайтов

## Архитектура

Проект построен на микросервисной архитектуре:

| Сервис | Порт | Описание |
|--------|------|----------|
| api-gateway | 8080 | API Gateway, точка входа |
| auth-service | 8081 | Аутентификация и авторизация |
| backlink-service | 8082 | Мониторинг бэклинков |
| index-service | 8083 | Проверка индексации |
| health-service | 8084 | Мониторинг здоровья сайтов |
| crawler-service | 8085 | Краулер для сбора данных |
| scheduler-service | 8086 | Планировщик задач |
| web-app | 3000 | Frontend (Next.js) |

### Инфраструктура

- **PostgreSQL 16** — основная база данных (порт 5432)
- **Redis 7** — кэширование и очереди (порт 6379)

## Структура проекта

```
link-tracker/
├── services/                 # Микросервисы (Go)
│   ├── api-gateway/
│   ├── auth-service/
│   ├── backlink-service/
│   ├── index-service/
│   ├── health-service/
│   ├── crawler-service/
│   └── scheduler-service/
├── frontend/                 # Frontend приложения
│   └── web-app/              # Next.js приложение
├── shared/                   # Общие компоненты
│   ├── proto/                # Protocol Buffers
│   └── contracts/            # API контракты
├── infrastructure/           # Инфраструктура
│   ├── docker/               # Docker конфигурации
│   └── scripts/              # Скрипты деплоя
├── docs/                     # Документация
│   ├── api/                  # API документация
│   ├── architecture/         # Архитектурные решения
│   ├── changelogs/           # Логи изменений по командам
│   └── examples/             # Примеры использования
├── .github/
│   └── workflows/            # GitHub Actions
├── docker-compose.yml        # Локальная разработка
├── docker-compose.prod.yml   # Production конфигурация
└── Makefile                  # Команды для разработки
```

## Быстрый старт

### Требования

- Docker & Docker Compose
- Go 1.21+
- Node.js 20+
- Make

### Запуск инфраструктуры

```bash
# Запуск PostgreSQL и Redis
make infra-up

# Или напрямую через docker-compose
docker-compose up -d postgres redis
```

### Остановка

```bash
make infra-down
# или
docker-compose down
```

## Разработка

### Команды Make

```bash
make help          # Показать все доступные команды
make infra-up      # Запустить инфраструктуру
make infra-down    # Остановить инфраструктуру
make logs          # Показать логи
```

### Git Workflow

- `main` — production-ready код
- `develop` — ветка для интеграции
- `feature/*` — ветки для фич
  - `feature/backend/*` — backend задачи
  - `feature/frontend/*` — frontend задачи
  - `feature/devops/*` — DevOps задачи

### Changelogs

Каждая команда ведёт свой changelog в `docs/changelogs/`:
- `backend.md` — Backend изменения
- `frontend.md` — Frontend изменения
- `devops.md` — DevOps изменения
- `tests.md` — QA изменения

## Лицензия

Proprietary
