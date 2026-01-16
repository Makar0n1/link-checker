# Backend Changelog

Лог изменений Backend команды.

---

## 🔴 СРОЧНОЕ ЗАДАНИЕ - 2026-01-16

### Проблема: Index Service и Health Service не интегрированы

**Статус контейнеров:** Auth (8081) и Backlink (8082) работают. Index и Health - НЕ ЗАПУЩЕНЫ (закомментированы в docker-compose.yml)

**Что нужно сделать:**

1. **Раскомментировать сервисы в docker-compose.yml:**
   - index-service (порт 8083)
   - health-service (порт 8084)

2. **Добавить роуты в nginx.conf:**
```nginx
location /api/v1/platforms/ {
    proxy_pass http://index_service;
    # ... остальные proxy headers
}

location /api/v1/sites/ {
    proxy_pass http://health_service;
    # ... остальные proxy headers
}
```

3. **Добавить upstreams в nginx.conf:**
```nginx
upstream index_service {
    server index-service:8083;
}

upstream health_service {
    server health-service:8084;
}
```

4. **Применить миграции:**
```bash
cat services/index-service/migrations/001_init.up.sql | docker exec -i linktracker-postgres psql -U linktracker -d linktracker
cat services/health-service/migrations/001_init.up.sql | docker exec -i linktracker-postgres psql -U linktracker -d linktracker
```

5. **Пересобрать и запустить:**
```bash
docker compose up -d --build index-service health-service nginx
```

**Тестирование:**
```bash
curl http://localhost:8083/health  # напрямую
curl http://localhost:8080/api/v1/platforms  # через nginx (с токеном)
curl http://localhost:8080/api/v1/sites  # через nginx (с токеном)
```

**Файлы:**
- docker-compose.yml - раскомментировать index-service, health-service
- infrastructure/nginx/nginx.conf - добавить роуты

---

### 2026-01-16 - Index Service + Health Service
**Branch:** feature/backend/index-health-services
**Status:** Done

#### Index Service (порт 8083)
Сервис проверки индексации URL в поисковых системах.

Endpoints:
- GET /api/v1/platforms - список платформ с пагинацией и фильтрами
- POST /api/v1/platforms - добавить URL
- GET /api/v1/platforms/{id} - получить платформу
- PUT /api/v1/platforms/{id} - обновить платформу
- DELETE /api/v1/platforms/{id} - удалить платформу
- POST /api/v1/platforms/bulk - массовое добавление (до 100)
- POST /api/v1/platforms/{id}/check - запустить проверку индексации
- GET /health, /ready - health checks

Database:
- Таблица `platforms`: id, user_id, url, domain, index_status, is_indexed, first_indexed_at, last_checked_at, check_count, potential_score, is_must_have, notes, created_at, updated_at
- ENUM тип: index_status (pending, indexed, not_indexed, error)

#### Health Service (порт 8084)
Сервис мониторинга здоровья сайтов.

Endpoints:
- GET /api/v1/sites - список сайтов
- POST /api/v1/sites - добавить сайт
- GET /api/v1/sites/{id} - получить сайт
- PUT /api/v1/sites/{id} - обновить сайт
- DELETE /api/v1/sites/{id} - удалить сайт
- POST /api/v1/sites/{id}/check - запустить проверку здоровья
- GET /api/v1/sites/{id}/history - история проверок
- GET /health, /ready - health checks

Database:
- Таблица `monitored_sites`: id, user_id, url, domain, http_status, is_alive, response_time_ms, allows_indexing, robots_txt_status, has_noindex, pages_indexed, last_checked_at, created_at
- Таблица `site_check_history`: id, site_id, http_status, is_alive, response_time_ms, checked_at

#### Files
Index Service:
- services/index-service/go.mod
- services/index-service/cmd/main.go
- services/index-service/internal/config/config.go
- services/index-service/internal/model/platform.go
- services/index-service/internal/model/dto.go
- services/index-service/internal/repository/platform_repository.go
- services/index-service/internal/service/platform_service.go
- services/index-service/internal/handler/platform_handler.go
- services/index-service/internal/handler/health.go
- services/index-service/migrations/001_init.up.sql
- services/index-service/migrations/001_init.down.sql
- services/index-service/Dockerfile

Health Service:
- services/health-service/go.mod
- services/health-service/cmd/main.go
- services/health-service/internal/config/config.go
- services/health-service/internal/model/site.go
- services/health-service/internal/model/dto.go
- services/health-service/internal/repository/site_repository.go
- services/health-service/internal/service/site_service.go
- services/health-service/internal/handler/site_handler.go
- services/health-service/internal/handler/health.go
- services/health-service/migrations/001_init.up.sql
- services/health-service/migrations/001_init.down.sql
- services/health-service/Dockerfile

Docs:
- docs/api/index-service.yaml
- docs/api/health-service.yaml
- docs/examples/index/*.json
- docs/examples/health/*.json

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
