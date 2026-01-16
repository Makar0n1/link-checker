# Frontend Changelog

Лог изменений Frontend команды.

---

### 2026-01-16 - Интеграция с Backend API
**Branch:** feature/frontend/api-integration
**Status:** Done

#### Что сделано
- Создан API клиент с управлением токенами (access_token, refresh_token)
- Автоматический refresh токена при истечении
- Хранение токенов в localStorage + cookie для middleware
- React Query hooks для аутентификации (useLogin, useRegister, useLogout, useCurrentUser)
- React Query hooks для проектов (useProjects, useCreateProject, useUpdateProject, useDeleteProject)
- React Query hooks для бэклинков (useBacklinks, useCreateBacklink, useUpdateBacklink, useDeleteBacklink, useBulkCreateBacklinks, useBulkDeleteBacklinks)
- Интеграция login/register страниц с реальным API
- Интеграция backlinks страницы с реальным API (загрузка, фильтрация, пагинация, CRUD)
- Интеграция logout с реальным API
- Добавлен QueryClientProvider для React Query
- Добавлен компонент Select от shadcn/ui
- Типы для API (User, Project, Backlink, и др.)
- Конфигурация .env.local для API endpoints

#### Файлы
- frontend/web-app/src/lib/api/client.ts (API клиент)
- frontend/web-app/src/lib/api/auth.ts (Auth API)
- frontend/web-app/src/lib/api/projects.ts (Projects API)
- frontend/web-app/src/lib/api/backlinks.ts (Backlinks API)
- frontend/web-app/src/hooks/use-auth.ts (Auth hooks)
- frontend/web-app/src/hooks/use-projects.ts (Projects hooks)
- frontend/web-app/src/hooks/use-backlinks.ts (Backlinks hooks)
- frontend/web-app/src/types/api.ts (API типы)
- frontend/web-app/src/components/providers.tsx (QueryClientProvider)
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
