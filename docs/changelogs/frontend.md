# Frontend Changelog

Лог изменений Frontend команды.

---

### 2026-01-16 16:15 (GMT+3) - Исправление Auth Guard для Docker
**Branch:** feature/frontend/api-integration
**Status:** Done

#### Проблема
- При переходе на /login или /register происходил редирект 307 на /backlinks
- Причина: старая/невалидная cookie `auth-token` в браузере
- Middleware проверял только наличие cookie, не валидность токена

#### Что сделано
- Добавлена валидация JWT формата в middleware
- Автоматическое удаление невалидных/старых cookie
- Отклонение mock-token (использовался ранее для тестов)

#### Логика middleware
```
1. Проверяем cookie "auth-token"
2. Валидируем JWT формат (3 части, base64)
3. Если токен невалидный + пытаемся зайти в dashboard → редирект на /login + очистка cookie
4. Если токен валидный + на auth странице → редирект на /backlinks
5. Если токен невалидный + на auth странице → очистка cookie
```

#### Файлы
- frontend/web-app/src/middleware.ts (JWT validation, cookie cleanup)

#### Как тестировать
```bash
# 1. Очистить cookies в браузере для localhost
# 2. Открыть http://localhost:8080/login
# 3. Должна открыться страница логина без редиректа
# 4. После логина - редирект на /backlinks
```

---

### 2026-01-16 14:00 (GMT+3) - Интеграция с Backend API
**Branch:** feature/frontend/api-integration
**Status:** Done

#### Что сделано
- Создан API клиент с управлением токенами (access_token, refresh_token)
- Автоматический refresh токена при истечении (30 сек буфер)
- Хранение токенов в localStorage + cookie `auth-token` для middleware
- React Query hooks для всех API операций
- Интеграция всех страниц с реальным API

#### API Client (src/lib/api/client.ts)
```typescript
// Конфигурация - относительные URL для Docker (nginx proxy)
const AUTH_API_URL = process.env.NEXT_PUBLIC_API_URL || "";       // → /api/v1/auth/*
const BACKLINK_API_URL = process.env.NEXT_PUBLIC_BACKLINK_API_URL || ""; // → /api/v1/*

// Token storage
localStorage: access_token, refresh_token, token_expiry
cookie: auth-token (для middleware)

// Auto-refresh: если токен истекает через <30 сек, автоматический refresh
```

#### React Query Hooks

**useAuth (src/hooks/use-auth.ts)**
```typescript
useLogin()      // POST /api/v1/auth/login → setTokens → redirect /backlinks
useRegister()   // POST /api/v1/auth/register → redirect /login
useLogout()     // POST /api/v1/auth/logout → clearTokens → redirect /login
useCurrentUser() // GET /api/v1/auth/me (enabled only if token exists)
```

**useProjects (src/hooks/use-projects.ts)**
```typescript
useProjects()                    // GET /api/v1/projects
useProject(id)                   // GET /api/v1/projects/:id
useCreateProject()               // POST /api/v1/projects
useUpdateProject()               // PUT /api/v1/projects/:id
useDeleteProject()               // DELETE /api/v1/projects/:id
```

**useBacklinks (src/hooks/use-backlinks.ts)**
```typescript
useBacklinks(params)             // GET /api/v1/backlinks?project_id=X&page=1&per_page=20&status=active
useBacklink(id)                  // GET /api/v1/backlinks/:id
useCreateBacklink()              // POST /api/v1/backlinks
useUpdateBacklink()              // PUT /api/v1/backlinks/:id
useDeleteBacklink()              // DELETE /api/v1/backlinks/:id
useBulkCreateBacklinks()         // POST /api/v1/backlinks/bulk
useBulkDeleteBacklinks()         // DELETE /api/v1/backlinks/bulk
```

#### TypeScript Types (src/types/api.ts)
```typescript
interface User {
  id: number;
  email: string;
  name: string;
  role: "user" | "admin";
  created_at: string;
}

interface AuthResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;  // секунды (обычно 900 = 15 мин)
}

interface Project {
  id: number;
  name: string;
  user_id: number;
  google_sheet_id: string | null;
  created_at: string;
}

type LinkStatus = "pending" | "active" | "broken" | "removed" | "nofollow";
type LinkType = "dofollow" | "nofollow" | "sponsored" | "ugc";

interface Backlink {
  id: number;
  project_id: number;
  source_url: string;
  target_url: string;
  anchor_text: string;
  status: LinkStatus;
  link_type: LinkType;
  http_status: number | null;
  last_checked_at: string | null;
  created_at: string;
}

interface PaginatedResponse<T> {
  data: T[];
  page: number;
  per_page: number;
  total: number;
  total_pages: number;
}
```

#### Страницы

**Login (/login)**
- Форма: email, password
- При submit: useLogin().mutate({ email, password })
- Toast при ошибке
- Redirect на /backlinks при успехе

**Register (/register)**
- Форма: name, email, password, confirmPassword
- Валидация: пароли совпадают, min 8 символов
- При submit: useRegister().mutate({ email, password, name })
- Redirect на /login при успехе

**Backlinks (/backlinks)**
- Селектор проекта (useProjects)
- Фильтр по статусу (all/active/pending/broken/removed/nofollow)
- Таблица с данными (useBacklinks с пагинацией)
- Inline edit anchor_text → useUpdateBacklink
- Bulk delete → useBulkDeleteBacklinks
- Refresh button

**Header**
- useCurrentUser() для отображения имени/email
- useLogout() для выхода

#### Конфигурация
```bash
# .env.local (для локальной разработки)
NEXT_PUBLIC_API_URL=http://localhost:8081
NEXT_PUBLIC_BACKLINK_API_URL=http://localhost:8082

# В Docker через nginx - пустые (относительные URL)
NEXT_PUBLIC_API_URL=
NEXT_PUBLIC_BACKLINK_API_URL=
```

#### Файлы
- frontend/web-app/src/lib/api/client.ts (API клиент)
- frontend/web-app/src/lib/api/auth.ts (Auth API)
- frontend/web-app/src/lib/api/projects.ts (Projects API)
- frontend/web-app/src/lib/api/backlinks.ts (Backlinks API)
- frontend/web-app/src/lib/api/index.ts (exports)
- frontend/web-app/src/hooks/use-auth.ts (Auth hooks)
- frontend/web-app/src/hooks/use-projects.ts (Projects hooks)
- frontend/web-app/src/hooks/use-backlinks.ts (Backlinks hooks)
- frontend/web-app/src/types/api.ts (API типы)
- frontend/web-app/src/components/providers.tsx (QueryClientProvider)
- frontend/web-app/src/components/ui/select.tsx (shadcn Select)
- frontend/web-app/src/app/layout.tsx (providers integration)
- frontend/web-app/src/app/(auth)/login/page.tsx (API integration)
- frontend/web-app/src/app/(auth)/register/page.tsx (API integration)
- frontend/web-app/src/app/(dashboard)/backlinks/page.tsx (API integration)
- frontend/web-app/src/components/layout/header.tsx (logout API)
- frontend/web-app/src/components/tables/backlinks-columns.tsx (new API types)
- frontend/web-app/.env.local
- frontend/web-app/.env.example

---

### 2026-01-16 - Landing Page и Auth Guard
**Branch:** feature/frontend/ui-setup
**Status:** Done

#### Что сделано
- Создана публичная Landing Page с Hero, Features и CTA секциями
- Добавлен Auth Guard middleware для защиты dashboard роутов
- Login устанавливает auth-token cookie после успешной авторизации
- Logout очищает cookie и редиректит на /login
- Dashboard доступен только авторизованным пользователям

#### Файлы
- frontend/web-app/src/app/page.tsx (Landing Page)
- frontend/web-app/src/middleware.ts (Auth Guard)
- frontend/web-app/src/app/(auth)/login/page.tsx (cookie set)
- frontend/web-app/src/components/layout/header.tsx (logout)

---

### 2026-01-16 - UI Kit и базовый Layout
**Branch:** feature/frontend/ui-setup
**Status:** Done

#### Что сделано
- Инициализирован Next.js 14 проект с App Router, TypeScript, Tailwind CSS
- Установлен и настроен shadcn/ui (Button, Input, Card, Table, Dialog, DropdownMenu, Toast, Badge, Tabs, Checkbox)
- Установлены зависимости: TanStack Table, TanStack Query, Zustand, Recharts, Lucide React
- Созданы layout компоненты: Sidebar (навигация), Header (user menu, notifications), MainLayout (wrapper)
- Созданы страницы-заглушки: /login, /register, /backlinks, /index-checker, /site-health
- Создан DataTable компонент с возможностями:
  - Сортировка по колонкам
  - Фильтрация по тексту
  - Inline редактирование ячеек (EditableCell)
  - Выбор строк (чекбоксы)
  - Bulk actions (удаление)
  - Пагинация
  - Переключение видимости колонок
- Добавлены mock данные для демонстрации (12 backlinks)
- Типы для Backlink модели

#### Файлы
- frontend/web-app/src/app/layout.tsx
- frontend/web-app/src/app/page.tsx
- frontend/web-app/src/app/(auth)/layout.tsx
- frontend/web-app/src/app/(auth)/login/page.tsx
- frontend/web-app/src/app/(auth)/register/page.tsx
- frontend/web-app/src/app/(dashboard)/layout.tsx
- frontend/web-app/src/app/(dashboard)/backlinks/page.tsx
- frontend/web-app/src/app/(dashboard)/index-checker/page.tsx
- frontend/web-app/src/app/(dashboard)/site-health/page.tsx
- frontend/web-app/src/components/layout/sidebar.tsx
- frontend/web-app/src/components/layout/header.tsx
- frontend/web-app/src/components/layout/main-layout.tsx
- frontend/web-app/src/components/tables/data-table.tsx
- frontend/web-app/src/components/tables/editable-cell.tsx
- frontend/web-app/src/components/tables/backlinks-columns.tsx
- frontend/web-app/src/lib/mock-data.ts
- frontend/web-app/src/types/backlink.ts

---

<!-- Template:
### YYYY-MM-DD - Название задачи
**Branch:** feature/frontend/xxx
**Status:** Done/In Progress

#### Что сделано
- пункт 1
- пункт 2

#### Файлы
- путь/к/файлу1
- путь/к/файлу2

---
-->
