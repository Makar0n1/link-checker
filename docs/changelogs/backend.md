# Backend Changelog

Лог изменений Backend команды.

---

### 2026-01-16 18:30 (GMT+3) - Index Service + Health Service Infrastructure
**Branch:** main
**Status:** Done

#### Изменения
- Раскомментированы index-service и health-service в docker-compose.yml
- Добавлены upstreams в nginx.conf для index_service и health_service
- Добавлены routes /api/v1/platforms/ и /api/v1/sites/ в nginx

---

### 2026-01-16 17:00 (GMT+3) - Index Service + Health Service
**Branch:** main
**Status:** Done

## INDEX SERVICE (порт 8083)
Сервис проверки индексации URL в поисковых системах.

### GET /api/v1/platforms
Получить список платформ с пагинацией и фильтрами.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Query Parameters:**
| Параметр | Тип | Описание |
|----------|-----|----------|
| page | int | Номер страницы (default: 1) |
| per_page | int | Записей на странице (default: 20, max: 100) |
| index_status | string | Фильтр по статусу: pending, indexed, not_indexed, error |
| is_indexed | bool | Фильтр по индексации: true/false |
| is_must_have | bool | Фильтр по приоритету: true/false |
| domain | string | Поиск по домену (ILIKE) |

**Response 200:**
```json
{
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "url": "https://example.com/article/seo-guide",
      "domain": "example.com",
      "index_status": "indexed",
      "is_indexed": true,
      "first_indexed_at": "2024-01-15T12:00:00Z",
      "last_checked_at": "2024-01-15T14:00:00Z",
      "check_count": 3,
      "potential_score": 85,
      "is_must_have": true,
      "notes": "High authority domain",
      "created_at": "2024-01-15T10:30:00Z",
      "updated_at": "2024-01-15T14:00:00Z"
    }
  ],
  "page": 1,
  "per_page": 20,
  "total": 1,
  "total_pages": 1
}
```

**Response 401:** `{"error": "unauthorized"}`

**Пример curl:**
```bash
curl -X GET "http://localhost:8080/api/v1/platforms/?page=1&per_page=20&index_status=indexed" \
  -H "Authorization: Bearer <token>"
```

---

### POST /api/v1/platforms
Добавить новую платформу для отслеживания.

**Headers:**
```
Authorization: Bearer <access_token>
Content-Type: application/json
```

**Request:**
```json
{
  "url": "https://example.com/article/seo-guide",  // string, required
  "potential_score": 85,                            // int, optional (0-100)
  "is_must_have": true,                             // bool, optional
  "notes": "High authority domain"                  // string, optional
}
```

**Response 201:**
```json
{
  "id": 1,
  "user_id": 1,
  "url": "https://example.com/article/seo-guide",
  "domain": "example.com",
  "index_status": "pending",
  "is_indexed": false,
  "first_indexed_at": null,
  "last_checked_at": null,
  "check_count": 0,
  "potential_score": 85,
  "is_must_have": true,
  "notes": "High authority domain",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z"
}
```

**Response 400:** `{"error": "url is required"}`
**Response 401:** `{"error": "unauthorized"}`

**Пример curl:**
```bash
curl -X POST http://localhost:8080/api/v1/platforms/ \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com/page","potential_score":75,"is_must_have":false}'
```

---

### GET /api/v1/platforms/{id}
Получить платформу по ID.

**Response 200:** Platform object
**Response 401:** `{"error": "unauthorized"}`
**Response 403:** `{"error": "not platform owner"}`
**Response 404:** `{"error": "platform not found"}`

**Пример curl:**
```bash
curl -X GET http://localhost:8080/api/v1/platforms/1 \
  -H "Authorization: Bearer <token>"
```

---

### PUT /api/v1/platforms/{id}
Обновить платформу.

**Request:**
```json
{
  "url": "https://new-url.com/page",      // string, optional
  "potential_score": 90,                   // int, optional
  "is_must_have": true,                    // bool, optional
  "notes": "Updated notes"                 // string, optional
}
```

**Response 200:** Updated Platform object
**Response 401:** `{"error": "unauthorized"}`
**Response 403:** `{"error": "not platform owner"}`
**Response 404:** `{"error": "platform not found"}`

**Пример curl:**
```bash
curl -X PUT http://localhost:8080/api/v1/platforms/1 \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"potential_score":95,"notes":"Updated"}'
```

---

### DELETE /api/v1/platforms/{id}
Удалить платформу.

**Response 204:** No Content
**Response 401:** `{"error": "unauthorized"}`
**Response 403:** `{"error": "not platform owner"}`
**Response 404:** `{"error": "platform not found"}`

**Пример curl:**
```bash
curl -X DELETE http://localhost:8080/api/v1/platforms/1 \
  -H "Authorization: Bearer <token>"
```

---

### POST /api/v1/platforms/bulk
Массовое создание платформ (до 100).

**Request:**
```json
{
  "platforms": [
    {
      "url": "https://site1.com/article",
      "potential_score": 70,
      "is_must_have": false
    },
    {
      "url": "https://site2.com/blog",
      "potential_score": 90,
      "is_must_have": true
    }
  ]
}
```

**Response 200:**
```json
{
  "success": 2,
  "failed": 0,
  "errors": [],
  "created": [
    { "id": 1, "url": "https://site1.com/article", ... },
    { "id": 2, "url": "https://site2.com/blog", ... }
  ]
}
```

**Пример curl:**
```bash
curl -X POST http://localhost:8080/api/v1/platforms/bulk \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"platforms":[{"url":"https://site1.com"},{"url":"https://site2.com"}]}'
```

---

### POST /api/v1/platforms/{id}/check
Запустить проверку индексации.

**Response 200:**
```json
{
  "platform_id": 1,
  "url": "https://example.com/article",
  "http_status": 200,
  "is_indexed": true,
  "index_status": "indexed",
  "checked_at": "2024-01-15T14:00:00Z"
}
```

**Response с ошибкой:**
```json
{
  "platform_id": 1,
  "url": "https://example.com/article",
  "http_status": 0,
  "is_indexed": false,
  "index_status": "error",
  "checked_at": "2024-01-15T14:00:00Z",
  "error": "connection timeout"
}
```

**Пример curl:**
```bash
curl -X POST http://localhost:8080/api/v1/platforms/1/check \
  -H "Authorization: Bearer <token>"
```

---

## HEALTH SERVICE (порт 8084)
Сервис мониторинга здоровья сайтов.

### GET /api/v1/sites
Получить список сайтов с пагинацией и фильтрами.

**Query Parameters:**
| Параметр | Тип | Описание |
|----------|-----|----------|
| page | int | Номер страницы (default: 1) |
| per_page | int | Записей на странице (default: 20, max: 100) |
| is_alive | bool | Фильтр по доступности: true/false |
| has_noindex | bool | Фильтр по noindex: true/false |
| domain | string | Поиск по домену (ILIKE) |

**Response 200:**
```json
{
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "url": "https://mysite.com",
      "domain": "mysite.com",
      "http_status": 200,
      "is_alive": true,
      "response_time_ms": 245,
      "allows_indexing": true,
      "robots_txt_status": "allow",
      "has_noindex": false,
      "pages_indexed": 150,
      "last_checked_at": "2024-01-15T14:00:00Z",
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "page": 1,
  "per_page": 20,
  "total": 1,
  "total_pages": 1
}
```

**Пример curl:**
```bash
curl -X GET "http://localhost:8080/api/v1/sites/?is_alive=true" \
  -H "Authorization: Bearer <token>"
```

---

### POST /api/v1/sites
Добавить сайт для мониторинга.

**Request:**
```json
{
  "url": "https://mysite.com"  // string, required
}
```

**Response 201:**
```json
{
  "id": 1,
  "user_id": 1,
  "url": "https://mysite.com",
  "domain": "mysite.com",
  "http_status": null,
  "is_alive": false,
  "response_time_ms": null,
  "allows_indexing": null,
  "robots_txt_status": "",
  "has_noindex": false,
  "pages_indexed": 0,
  "last_checked_at": null,
  "created_at": "2024-01-15T10:30:00Z"
}
```

**Пример curl:**
```bash
curl -X POST http://localhost:8080/api/v1/sites/ \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"url":"https://mysite.com"}'
```

---

### GET /api/v1/sites/{id}
Получить сайт по ID.

**Response 200:** MonitoredSite object
**Response 404:** `{"error": "site not found"}`

---

### PUT /api/v1/sites/{id}
Обновить сайт.

**Request:**
```json
{
  "url": "https://new-site.com",  // string, optional
  "pages_indexed": 200             // int, optional
}
```

**Response 200:** Updated MonitoredSite object

---

### DELETE /api/v1/sites/{id}
Удалить сайт.

**Response 204:** No Content

---

### POST /api/v1/sites/{id}/check
Запустить проверку здоровья сайта.

Проверяет:
- HTTP статус и время ответа
- robots.txt (allow/disallow)
- meta noindex в HTML

**Response 200:**
```json
{
  "site_id": 1,
  "url": "https://mysite.com",
  "http_status": 200,
  "is_alive": true,
  "response_time_ms": 245,
  "allows_indexing": true,
  "robots_txt_status": "allow",
  "has_noindex": false,
  "checked_at": "2024-01-15T14:00:00Z"
}
```

**Пример curl:**
```bash
curl -X POST http://localhost:8080/api/v1/sites/1/check \
  -H "Authorization: Bearer <token>"
```

---

### GET /api/v1/sites/{id}/history
Получить историю проверок сайта.

**Query Parameters:**
| Параметр | Тип | Описание |
|----------|-----|----------|
| page | int | Номер страницы (default: 1) |
| per_page | int | Записей на странице (default: 20, max: 100) |

**Response 200:**
```json
{
  "data": [
    {
      "id": 3,
      "site_id": 1,
      "http_status": 200,
      "is_alive": true,
      "response_time_ms": 245,
      "checked_at": "2024-01-15T14:00:00Z"
    },
    {
      "id": 2,
      "site_id": 1,
      "http_status": 200,
      "is_alive": true,
      "response_time_ms": 312,
      "checked_at": "2024-01-15T12:00:00Z"
    }
  ],
  "page": 1,
  "per_page": 20,
  "total": 2,
  "total_pages": 1
}
```

**Пример curl:**
```bash
curl -X GET "http://localhost:8080/api/v1/sites/1/history?page=1" \
  -H "Authorization: Bearer <token>"
```

---

#### Database Schema

**platforms (index-service):**
```sql
CREATE TYPE index_status AS ENUM ('pending', 'indexed', 'not_indexed', 'error');

CREATE TABLE platforms (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    url TEXT NOT NULL,
    domain VARCHAR(255),
    index_status index_status DEFAULT 'pending',
    is_indexed BOOLEAN DEFAULT FALSE,
    first_indexed_at TIMESTAMP,
    last_checked_at TIMESTAMP,
    check_count INTEGER DEFAULT 0,
    potential_score INTEGER DEFAULT 0,
    is_must_have BOOLEAN DEFAULT FALSE,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);
```

**monitored_sites & site_check_history (health-service):**
```sql
CREATE TABLE monitored_sites (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    url TEXT NOT NULL,
    domain VARCHAR(255),
    http_status INTEGER,
    is_alive BOOLEAN DEFAULT FALSE,
    response_time_ms INTEGER,
    allows_indexing BOOLEAN,
    robots_txt_status VARCHAR(50),
    has_noindex BOOLEAN DEFAULT FALSE,
    pages_indexed INTEGER DEFAULT 0,
    last_checked_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE site_check_history (
    id BIGSERIAL PRIMARY KEY,
    site_id BIGINT REFERENCES monitored_sites(id) ON DELETE CASCADE,
    http_status INTEGER,
    is_alive BOOLEAN,
    response_time_ms INTEGER,
    checked_at TIMESTAMP DEFAULT NOW()
);
```

#### Files
- services/index-service/* (все файлы сервиса)
- services/health-service/* (все файлы сервиса)
- docs/api/index-service.yaml
- docs/api/health-service.yaml
- docs/examples/index/*.json
- docs/examples/health/*.json

---

### 2026-01-16 15:00 (GMT+3) - Shared Go Module
**Branch:** main
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

### 2026-01-16 14:00 (GMT+3) - Backlink Service Implementation
**Branch:** main
**Status:** Done

## BACKLINK SERVICE (порт 8082)

### GET /api/v1/projects
Получить список проектов пользователя.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response 200:**
```json
{
  "data": [
    {
      "id": 1,
      "name": "My SEO Project",
      "user_id": 1,
      "google_sheet_id": "1abc...",
      "created_at": "2024-01-15T10:00:00Z"
    }
  ],
  "page": 1,
  "per_page": 20,
  "total": 1,
  "total_pages": 1
}
```

**Пример curl:**
```bash
curl -X GET http://localhost:8080/api/v1/projects/ \
  -H "Authorization: Bearer <token>"
```

---

### POST /api/v1/projects
Создать проект.

**Request:**
```json
{
  "name": "My SEO Project",        // string, required
  "google_sheet_id": "1abc..."     // string, optional
}
```

**Response 201:**
```json
{
  "id": 1,
  "name": "My SEO Project",
  "user_id": 1,
  "google_sheet_id": "1abc...",
  "created_at": "2024-01-15T10:00:00Z"
}
```

**Пример curl:**
```bash
curl -X POST http://localhost:8080/api/v1/projects/ \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"name":"My Project"}'
```

---

### GET /api/v1/backlinks
Получить список бэклинков с пагинацией и фильтрами.

**Query Parameters:**
| Параметр | Тип | Описание |
|----------|-----|----------|
| page | int | Номер страницы |
| per_page | int | Записей на странице (max: 100) |
| project_id | int | Фильтр по проекту |
| status | string | pending, active, broken, removed, nofollow |
| link_type | string | dofollow, nofollow, sponsored, ugc |
| source_url | string | Поиск по source URL |
| target_url | string | Поиск по target URL |

**Response 200:**
```json
{
  "data": [
    {
      "id": 1,
      "project_id": 1,
      "source_url": "https://example.com/blog/seo-tips",
      "target_url": "https://mysite.com/services",
      "anchor_text": "SEO services",
      "status": "active",
      "link_type": "dofollow",
      "http_status": 200,
      "last_checked_at": "2024-01-15T12:00:00Z",
      "created_at": "2024-01-15T10:35:00Z"
    }
  ],
  "page": 1,
  "per_page": 20,
  "total": 1,
  "total_pages": 1
}
```

**Пример curl:**
```bash
curl -X GET "http://localhost:8080/api/v1/backlinks/?project_id=1&status=active" \
  -H "Authorization: Bearer <token>"
```

---

### POST /api/v1/backlinks
Создать бэклинк.

**Request:**
```json
{
  "project_id": 1,                              // int, required
  "source_url": "https://example.com/blog",     // string, required
  "target_url": "https://mysite.com/page",      // string, required
  "anchor_text": "click here",                  // string, optional
  "link_type": "dofollow"                       // string, optional: dofollow, nofollow, sponsored, ugc
}
```

**Response 201:**
```json
{
  "id": 1,
  "project_id": 1,
  "source_url": "https://example.com/blog",
  "target_url": "https://mysite.com/page",
  "anchor_text": "click here",
  "status": "pending",
  "link_type": "dofollow",
  "http_status": null,
  "last_checked_at": null,
  "created_at": "2024-01-15T10:35:00Z"
}
```

**Пример curl:**
```bash
curl -X POST http://localhost:8080/api/v1/backlinks/ \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"project_id":1,"source_url":"https://example.com","target_url":"https://mysite.com"}'
```

---

### POST /api/v1/backlinks/bulk
Массовое создание бэклинков (до 100).

**Request:**
```json
{
  "backlinks": [
    {
      "project_id": 1,
      "source_url": "https://site1.com/article",
      "target_url": "https://mysite.com/page1",
      "anchor_text": "link text 1",
      "link_type": "dofollow"
    },
    {
      "project_id": 1,
      "source_url": "https://site2.com/blog",
      "target_url": "https://mysite.com/page2",
      "anchor_text": "link text 2",
      "link_type": "nofollow"
    }
  ]
}
```

**Response 200:**
```json
{
  "success": 2,
  "failed": 0,
  "errors": []
}
```

---

### DELETE /api/v1/backlinks/bulk
Массовое удаление бэклинков (до 100).

**Request:**
```json
{
  "ids": [1, 2, 3]
}
```

**Response 200:**
```json
{
  "success": 3,
  "failed": 0,
  "errors": []
}
```

---

#### Database Schema

```sql
CREATE TYPE link_status AS ENUM ('pending', 'active', 'broken', 'removed', 'nofollow');
CREATE TYPE link_type AS ENUM ('dofollow', 'nofollow', 'sponsored', 'ugc');

CREATE TABLE projects (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    user_id BIGINT NOT NULL,
    google_sheet_id VARCHAR(255),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE backlinks (
    id BIGSERIAL PRIMARY KEY,
    project_id BIGINT REFERENCES projects(id) ON DELETE CASCADE,
    source_url TEXT NOT NULL,
    target_url TEXT NOT NULL,
    anchor_text TEXT,
    status link_status DEFAULT 'pending',
    link_type link_type DEFAULT 'dofollow',
    http_status INTEGER,
    last_checked_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

### 2026-01-16 12:00 (GMT+3) - Auth Service Implementation
**Branch:** main
**Status:** Done

## AUTH SERVICE (порт 8081)

### POST /api/v1/auth/register
Регистрация нового пользователя.

**Request:**
```json
{
  "email": "user@example.com",    // string, required, valid email
  "password": "MyPassword123",    // string, required, min 8 chars
  "name": "John Doe"              // string, required
}
```

**Response 201:**
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  "role": "user",
  "created_at": "2024-01-15T10:00:00Z"
}
```

**Response 400:** `{"error": "invalid request body"}`
**Response 409:** `{"error": "user already exists"}`

**Пример curl:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"Test12345","name":"Test User"}'
```

---

### POST /api/v1/auth/login
Аутентификация и получение JWT токенов.

**Request:**
```json
{
  "email": "user@example.com",
  "password": "MyPassword123"
}
```

**Response 200:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 900,
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe",
    "role": "user"
  }
}
```

**Response 401:** `{"error": "invalid credentials"}`

**Пример curl:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@test.com","password":"Test12345"}'
```

---

### POST /api/v1/auth/refresh
Обновление access и refresh токенов.

**Request:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response 200:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 900
}
```

**Response 401:** `{"error": "invalid refresh token"}`

**Пример curl:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<refresh_token>"}'
```

---

### GET /api/v1/auth/me
Получение информации о текущем пользователе.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Response 200:**
```json
{
  "id": 1,
  "email": "user@example.com",
  "name": "John Doe",
  "role": "user",
  "created_at": "2024-01-15T10:00:00Z"
}
```

**Response 401:** `{"error": "unauthorized"}`

**Пример curl:**
```bash
curl -X GET http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer <access_token>"
```

---

### POST /api/v1/auth/logout
Выход и инвалидация refresh токена.

**Headers:**
```
Authorization: Bearer <access_token>
```

**Request:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response 200:**
```json
{
  "message": "logged out successfully"
}
```

**Пример curl:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<refresh_token>"}'
```

---

#### Database Schema

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE refresh_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

<!-- Template:
### YYYY-MM-DD HH:MM (GMT+3) - Название задачи
**Branch:** main
**Status:** Done/In Progress

#### Что сделано
- пункт 1
- пункт 2

#### Файлы
- путь/к/файлу1
- путь/к/файлу2

---
-->
