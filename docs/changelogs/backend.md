# Backend Changelog

Лог изменений Backend команды.

---

### 2026-01-16 - Shared Go Module
**Branch:** develop
**Status:** Done

#### Added
- shared/go/pkg/middleware/jwt.go - общий JWT middleware для всех сервисов
- shared/go/pkg/models/claims.go - JWT Claims структура
- shared/go/pkg/response/json.go - HTTP response helpers (JSON, Error, Paginated)

#### Files
- shared/go/go.mod
- shared/go/pkg/middleware/jwt.go
- shared/go/pkg/models/claims.go
- shared/go/pkg/response/json.go

---

### 2026-01-16 - Backlink Service Implementation
**Branch:** feature/backend/backlink-service
**Status:** Done

#### Added
Projects API:
- GET /api/v1/projects - список проектов пользователя
- POST /api/v1/projects - создать проект
- GET /api/v1/projects/:id - получить проект
- PUT /api/v1/projects/:id - обновить проект
- DELETE /api/v1/projects/:id - удалить проект

Backlinks API:
- GET /api/v1/backlinks - список с пагинацией и фильтрами (project_id, status, link_type, source_url, target_url)
- POST /api/v1/backlinks - создать бэклинк
- GET /api/v1/backlinks/:id - получить бэклинк
- PUT /api/v1/backlinks/:id - обновить бэклинк
- DELETE /api/v1/backlinks/:id - удалить бэклинк
- POST /api/v1/backlinks/bulk - массовое создание (до 100)
- DELETE /api/v1/backlinks/bulk - массовое удаление (до 100)
- POST /api/v1/backlinks/import - импорт из Google Sheets (заглушка)

#### Database Changes
- Таблица `projects`: id, name, user_id, google_sheet_id, created_at
- Таблица `backlinks`: id, project_id, source_url, target_url, anchor_text, status, link_type, http_status, last_checked_at, created_at
- ENUM типы: link_status (pending, active, broken, removed, nofollow), link_type (dofollow, nofollow, sponsored, ugc)
- Индексы: idx_projects_user_id, idx_backlinks_project_id, idx_backlinks_status, idx_backlinks_source_url, idx_backlinks_target_url, idx_backlinks_last_checked_at

#### Files
- services/backlink-service/go.mod
- services/backlink-service/cmd/main.go
- services/backlink-service/internal/config/config.go
- services/backlink-service/internal/model/backlink.go
- services/backlink-service/internal/model/dto.go
- services/backlink-service/internal/repository/project_repository.go
- services/backlink-service/internal/repository/backlink_repository.go
- services/backlink-service/internal/service/project_service.go
- services/backlink-service/internal/service/backlink_service.go
- services/backlink-service/internal/handler/project_handler.go
- services/backlink-service/internal/handler/backlink_handler.go
- services/backlink-service/internal/handler/health.go
- services/backlink-service/migrations/001_init.up.sql
- services/backlink-service/migrations/001_init.down.sql
- docs/api/backlink-service.yaml
- docs/examples/backlink/*.json

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
