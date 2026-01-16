# DevOps Changelog

Лог изменений DevOps команды.

---

### 2026-01-16 - CI/CD Pipeline + Frontend Docker
**Branch:** feature/devops/ci-frontend
**Status:** Done

#### Что сделано
- Создан GitHub Actions CI Pipeline (.github/workflows/ci.yml)
  - Backend lint и тесты (Go 1.22, golangci-lint)
  - Frontend lint, type check, build (Node 20)
  - Docker build для всех сервисов
- Создан Dockerfile для frontend (Next.js standalone, multi-stage)
- Обновлён next.config.mjs — добавлен output: 'standalone'
- Добавлен frontend в docker-compose.yml
- Создан nginx как API Gateway (infrastructure/nginx/nginx.conf)
  - /api/v1/auth/* → auth-service:8081
  - /api/v1/projects/, /api/v1/backlinks/ → backlink-service:8082
  - /* → frontend:3000

#### Файлы
- .github/workflows/ci.yml
- frontend/web-app/Dockerfile
- frontend/web-app/next.config.mjs
- docker-compose.yml
- infrastructure/nginx/nginx.conf

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
