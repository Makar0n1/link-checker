# Team Guidelines - Link Checker Project

## Обязательные правила для всех команд

### 1. Перед началом работы

**ВСЕГДА смотри документацию перед тем как приступать к задаче!**

```bash
# Проверь последние изменения
git pull origin develop

# Изучи changelogs других команд
cat docs/changelogs/backend.md
cat docs/changelogs/frontend.md
cat docs/changelogs/devops.md
cat docs/changelogs/tests.md
```

### 2. Структура документации

```
docs/
├── api/                    # OpenAPI спецификации
│   ├── auth-service.yaml
│   ├── backlink-service.yaml
│   ├── index-service.yaml
│   └── health-service.yaml
├── changelogs/             # Логи изменений по командам
│   ├── backend.md
│   ├── frontend.md
│   ├── devops.md
│   └── tests.md
└── examples/               # Примеры запросов/ответов API
```

### 3. Ведение changelog

После КАЖДОГО коммита обновляй changelog своей команды:

```markdown
## [Sprint X] - YYYY-MM-DD

### Added
- Что добавлено

### Changed
- Что изменено

### Fixed
- Что исправлено
```

### 4. Git workflow

1. Работай только в своей feature ветке
2. Всегда от develop: `git checkout -b feature/{team}/{task-name} develop`
3. Коммиты с префиксом: `feat(service):`, `fix(service):`, `test(qa):`
4. Перед PR - rebase на develop

### 5. Порты сервисов

| Service          | Port |
|------------------|------|
| Auth Service     | 8081 |
| Backlink Service | 8082 |
| Index Service    | 8083 |
| Health Service   | 8084 |
| Crawler Service  | 8085 |
| Frontend         | 3000 |
| PostgreSQL       | 5432 |
| Redis            | 6379 |

### 6. Коммуникация между командами

**Бекенд → Фронтенд:**
- Документируй API в `docs/api/{service}.yaml`
- Добавляй примеры в `docs/examples/`
- Указывай breaking changes в changelog

**Фронтенд → Бекенд:**
- Сообщай о нужных endpoint'ах
- Указывай формат данных который ожидаешь

**DevOps → Все:**
- Документируй изменения в docker-compose.yml
- Указывай новые env переменные
- Обновляй CI/CD pipeline описание

**Тестеры → Все:**
- Сообщай о найденных багах
- Документируй тест-кейсы
- Указывай покрытие тестами

### 7. Checklist перед коммитом

- [ ] Код работает локально
- [ ] Обновлен changelog
- [ ] Нет конфликтов с develop
- [ ] Добавлена документация (если нужно)
- [ ] Тесты проходят (если есть)

### 8. Структура сервисов

```
services/{service-name}/
├── cmd/main.go              # Entry point
├── internal/
│   ├── config/              # Конфигурация
│   ├── handler/             # HTTP handlers
│   ├── model/               # Модели и DTO
│   ├── repository/          # Работа с БД
│   └── service/             # Бизнес-логика
├── migrations/              # SQL миграции
├── Dockerfile
├── go.mod
└── go.sum
```

### 9. Shared модуль

Общий код в `shared/go/pkg/`:
- `middleware/jwt.go` - JWT валидация
- `response/json.go` - HTTP ответы
- `models/claims.go` - JWT claims

Используй в сервисах:
```go
import "github.com/link-tracker/shared/pkg/middleware"
```

### 10. База данных

- Каждый сервис имеет свои миграции
- Миграции применяются вручную (пока)
- Все сервисы используют одну БД `linktracker`

---

**ВАЖНО: Эти правила обязательны для всех команд. Нарушение приводит к конфликтам и багам!**
