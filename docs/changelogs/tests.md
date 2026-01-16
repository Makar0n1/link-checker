# QA/Tests Changelog

Лог изменений QA команды.

---

### 2026-01-16 - E2E и API тесты
**Branch:** feature/tests/e2e-api
**Status:** Done

#### Что сделано

**E2E тесты (Playwright):**
- Настроен Playwright для E2E тестирования
- Созданы тесты для Authentication flow:
  - Отображение landing page
  - Навигация на login/register
  - Регистрация нового пользователя
  - Валидация паролей (несовпадение, короткий пароль)
  - Логин и редирект на dashboard
  - Ошибка при неверном пароле
  - Защита dashboard routes (auth guard)
  - Logout
  - Редирект авторизованного пользователя со страниц auth
  - Навигация между страницами auth

**API тесты (Vitest):**

Auth Service:
- POST /api/v1/auth/register → 201 (регистрация)
- POST /api/v1/auth/register → 409 (дублирующий email)
- POST /api/v1/auth/register → 400 (невалидные данные)
- POST /api/v1/auth/login → 200 + tokens (успешный логин)
- POST /api/v1/auth/login → 401 (неверный пароль)
- GET /api/v1/auth/me → 200 (с токеном)
- GET /api/v1/auth/me → 401 (без токена)
- POST /api/v1/auth/refresh → 200 (обновление токена)
- POST /api/v1/auth/refresh → 401 (невалидный токен)
- POST /api/v1/auth/logout → 200

Backlink Service:
- POST /api/v1/projects → 201 (создание проекта)
- GET /api/v1/projects → 200 (список проектов)
- GET /api/v1/projects/:id → 200 (получение проекта)
- PUT /api/v1/projects/:id → 200 (обновление проекта)
- DELETE /api/v1/projects/:id → 204 (удаление проекта)
- POST /api/v1/backlinks → 201 (создание backlink)
- GET /api/v1/backlinks → 200 (список backlinks с пагинацией)
- GET /api/v1/backlinks/:id → 200 (получение backlink)
- PUT /api/v1/backlinks/:id → 200 (обновление backlink)
- DELETE /api/v1/backlinks/:id → 204 (удаление backlink)
- POST /api/v1/backlinks/bulk → 201 (bulk создание)
- DELETE /api/v1/backlinks/bulk → 200 (bulk удаление)
- Все endpoints → 401 (unauthorized requests)

#### Файлы
- tests/e2e/package.json
- tests/e2e/playwright.config.ts
- tests/e2e/specs/auth.spec.ts
- tests/api/package.json
- tests/api/vitest.config.ts
- tests/api/helpers/api-client.ts
- tests/api/auth.test.ts
- tests/api/backlinks.test.ts

#### Запуск тестов

E2E тесты:
```bash
cd tests/e2e
npm install
npm test
```

API тесты:
```bash
cd tests/api
npm install
npm test
```

---

<!-- Template:
### YYYY-MM-DD - Название задачи
**Branch:** feature/tests/xxx
**Status:** Done/In Progress

#### Что сделано
- пункт 1
- пункт 2

#### Файлы
- путь/к/файлу1
- путь/к/файлу2

---
-->
