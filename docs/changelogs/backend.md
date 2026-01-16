# Backend Changelog

Лог изменений Backend команды.

---

### 2026-01-16 - Auth Service Implementation
**Branch:** feature/backend/auth-service
**Status:** Done

#### Added
- POST /api/v1/auth/register - регистрация нового пользователя
- POST /api/v1/auth/login - аутентификация и получение JWT токенов
- POST /api/v1/auth/refresh - обновление access и refresh токенов
- GET /api/v1/auth/me - получение информации о текущем пользователе (protected)
- POST /api/v1/auth/logout - выход и инвалидация refresh токена
- GET /health - health check endpoint
- GET /ready - readiness check endpoint

#### Database Changes
- Таблица `users`: id, email, password_hash, name, role, created_at, updated_at
- Таблица `refresh_tokens`: id, user_id, token_hash, expires_at, created_at
- Индексы: idx_users_email, idx_refresh_tokens_user_id, idx_refresh_tokens_expires_at
- Триггер автообновления updated_at

#### Files
- services/auth-service/go.mod
- services/auth-service/cmd/main.go
- services/auth-service/internal/config/config.go
- services/auth-service/internal/model/user.go
- services/auth-service/internal/model/dto.go
- services/auth-service/internal/repository/user_repository.go
- services/auth-service/internal/repository/token_repository.go
- services/auth-service/internal/service/auth_service.go
- services/auth-service/internal/handler/auth_handler.go
- services/auth-service/internal/handler/health.go
- services/auth-service/internal/middleware/auth.go
- services/auth-service/migrations/001_init.up.sql
- services/auth-service/migrations/001_init.down.sql
- docs/api/auth-service.yaml
- docs/examples/auth/*.json

---

<!-- Template:
### YYYY-MM-DD - Название задачи
**Branch:** feature/backend/xxx
**Status:** Done/In Progress

#### Что сделано
- пункт 1
- пункт 2

#### Файлы
- путь/к/файлу1
- путь/к/файлу2

---
-->
