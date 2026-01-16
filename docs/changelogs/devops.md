# DevOps Changelog

Лог изменений DevOps команды.

---

### 2026-01-16 - Dockerfile для сервисов
**Branch:** feature/devops/dockerize-services
**Status:** Done

#### Что сделано
- Создан Dockerfile для auth-service (multi-stage build, Go 1.22)
- Создан Dockerfile для backlink-service (с поддержкой shared модуля)
- Обновлён docker-compose.yml — раскомментированы auth-service и backlink-service
- Добавлены health checks для Go сервисов (endpoint /health)
- Добавлены environment переменные для сервисов
- Настроены depends_on с условием service_healthy
- Удалён deprecated атрибут version из docker-compose.yml

#### Файлы
- services/auth-service/Dockerfile
- services/backlink-service/Dockerfile
- docker-compose.yml

---

### 2026-01-16 - Инициализация репозитория
**Branch:** main
**Status:** Done

#### Что сделано
- Инициализирован git репозиторий
- Создана структура папок для микросервисной архитектуры
- Настроен docker-compose.yml с PostgreSQL 16 и Redis 7
- Создан docker-compose.prod.yml для production
- Настроен .gitignore для Go, Node.js, IDE файлов
- Создан README.md с описанием проекта
- Создан Makefile с базовыми командами
- Созданы шаблоны changelog файлов

#### Файлы
- docker-compose.yml
- docker-compose.prod.yml
- .gitignore
- README.md
- Makefile
- docs/changelogs/backend.md
- docs/changelogs/frontend.md
- docs/changelogs/devops.md
- docs/changelogs/tests.md

---
