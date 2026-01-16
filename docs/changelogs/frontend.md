# Frontend Changelog

Лог изменений Frontend команды.

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
