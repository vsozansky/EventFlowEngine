# EventFlowEngine — Правила разработки

> **Для кого:** Разработчики Go-пакета `github.com/eventflow/engine`  
> **Назначение:** Единые принципы, требования и ограничения при создании кода

---

## 1. Фундаментальные принципы

### 1.1. Принцип «Всё есть событие»

В системе нет ничего, кроме событий. Акторы, сущности, атрибуты, модели, индивиды — это **не отдельные структуры**, а агрегаты, **собираемые из событий** по запросу.

```go
// НЕПРАВИЛЬНО: хранить отдельную таблицу сущностей
type EntityRow struct {
    ID    int
    Label string
}

// ПРАВИЛЬНО: сущность собирается из событий с ValueID=Entity
func (s *MemoryStorage) GetEntity(id int) (*Entity, error) {
    // Найти событие: id, valueID=Entity
    // Найти все дочерние события: baseEventID=id
    // Собрать Entity
}
```

### 1.2. Immutable (неизменяемость)

События **никогда не удаляются и не изменяются** после записи. Никаких UPDATE и DELETE. Если нужно изменить значение — создаётся новое событие, которое «перекрывает» предыдущее.

```go
// НЕПРАВИЛЬНО
event.Value = "new value"
storage.UpdateEvent(event)

// ПРАВИЛЬНО
newEvent := EventData{
    BaseEventID: event.ID,  // ссылка на предыдущее
    Value: "new value",
    // ... остальные поля
}
storage.AppendEvent(newEvent)
```

### 1.3. Append-only

Запись — только добавление в конец. `Storage.AppendEvent()` — единственный метод записи. ID события = `Position++` (монотонный счётчик).

### 1.4. Функциональный стиль

- Чистые функции без побочных эффектов (где возможно)
- Отсутствие глобального состояния
- Параметры передаются явно
- nil-значения обрабатываются как «отсутствие данных»

### 1.5. Семантическая типизация

Каждое событие имеет семантический тип (`ValueID`), определённый в генезис-графе. Типы не могут быть произвольными строками — только константы из `StaticEvent`.

---

## 2. Требования к коду

### 2.1. Ошибки

- Все функции, которые могут вернуть ошибку, **возвращают `error`**
- Ошибки **оборачиваются** контекстом через `fmt.Errorf("описание: %w", err)`
- Пакет `errors` используется для `errors.Is()` / `errors.As()`
- Определены sentinel errors для ключевых ситуаций:

```go
var (
    ErrEventNotFound      = errors.New("event not found")
    ErrInvalidPosition    = errors.New("invalid event position")
    ErrInvalidBaseRef     = errors.New("invalid base event reference")
    ErrInvalidValueType   = errors.New("invalid value type")
    ErrInvalidActor       = errors.New("invalid actor")
    ErrDuplicateEvent     = errors.New("duplicate event position")
    ErrStorageClosed      = errors.New("storage is closed")
)
```

### 2.2. Тестирование

- Table-driven тесты (`t.Run`)
- Тесты на `MemoryStorage` — все операции
- Тесты на `EngineBase.ProcessEvent` — валидация
- Тесты на `InitGenesis` — проверка количества и типов событий
- Тесты на транзакции — CreateEntity, CreateModel, CreateIndividual
- Тесты на ConditionEngine — все типы правил
- Тесты на граничные случаи: nil, пустые значения, отрицательные ID

```go
func TestCreateEntity(t *testing.T) {
    tests := []struct {
        name     string
        actorID  int
        label    string
        wantErr  bool
    }{
        {"valid entity", 21, "Product", false},
        {"empty label", 21, "", true},
        {"invalid actor", 999, "Test", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // ...
        })
    }
}
```

### 2.3. Документация

- **Все экспортируемые типы, функции, константы** документируются в GoDoc-стиле
- **Все неочевидные решения** комментируются в коде (почему, а не что)
- **Примеры использования** в `_test.go` или в `example_test.go`

```go
// EventData представляет собой базовое событие в графе.
// Любая информация в системе (акторы, сущности, атрибуты, модели,
// индивиды, значения) хранится как событие. Событие immutable
// после записи — никаких UPDATE или DELETE.
type EventData struct { ... }
```

### 2.4. Именование

| Элемент | Правило | Пример |
|---------|---------|--------|
| Пакеты | lowercase, одно слово | `engine`, `storage`, `condition` |
| Типы | PascalCase | `EventData`, `ConditionRule` |
| Функции | PascalCase (экспорт) / camelCase (приватные) | `NewEngine()`, `processEvent()` |
| Константы | PascalCase | `StaticEvent`, `ErrEventNotFound` |
| Переменные | camelCase | `eventData`, `storageAPI` |
| Тесты | `TestXxx` / `TestXxx_yyy` | `TestProcessEvent` |
| Бенчмарки | `BenchmarkXxx` | `BenchmarkProcessEvent` |

### 2.5. Форматирование

- `gofmt` / `go fmt` перед каждым коммитом
- `go vet` без предупреждений
- Максимальная длина строки — 120 символов (исключения для тестов)
- Импорты группируются: std → внешние → внутренние

```go
import (
    "errors"
    "fmt"
    "time"

    "github.com/some/external/lib"
)
```

---

## 3. Архитектурные правила

### 3.1. Слои

```
Storage (интерфейс) ← MemoryStorage (реализация)
     ↑
EngineBase (ядро, зависит от Storage interface)
     ↑
Transactions (CRUD-операции, методы EngineBase)
     ↑
ConditionEngine (исполнитель условий, зависит от Storage)
```

- `EngineBase` НЕ зависит от конкретной реализации Storage
- `MemoryStorage` НЕ зависит от EngineBase
- `ConditionEngine` получает Storage через интерфейс

### 3.2. Паника

Паника — только для случаев, когда программа не может продолжать:
- nil-указатель на Storage при создании Engine
- Ошибка инициализации генезис-графа
- Во всех остальных случаях — возврат `error`

### 3.3. Горутины

- Пакет потокобезопасен на уровне Storage (sync.RWMutex)
- EngineBase не горутинобезопасен (один поток)
- ConditionEngine не имеет состояния — может вызываться из多个 горутин

### 3.4. Nil-значения

- nil как ConditionRule означает «нет условия» (всегда истинно)
- nil как DefaultValue означает «нет значения по умолчанию»
- nil как Permission означает «нет ограничений по доступу»
- nil как PropertyProvider означает «нет вложенных свойств»

---

## 4. Ограничения (constraints)

### 4.1. Зависимости

**Внешние зависимости (разрешены):**
- Только стандартная библиотека Go (net, time, sync, fmt, errors...)

**Внешние зависимости (запрещены):**
- Любые сторонние библиотеки
- ORM, SQL, базы данных
- HTTP-серверы / клиенты
- Фреймворки

> Исключение: библиотеки для тестирования (testify/require — опционально)

### 4.2. Производительность

- Операции записи O(1) в среднем
- Операции чтения по ID O(1)
- Операции агрегации (GetEntity) — O(n) по числу дочерних событий
- `Position` — int64 (достаточно для 9×10^18 событий)

### 4.3. Безопасность

- Все входящие ID проверяются на положительность
- Все строки обрезаются до 64К
- Дата всегда UTC (локальные таймзоны запрещены)
- Актор всегда проверяется на существование

---

## 5. Процесс разработки

### 5.1. Порядок реализации

```
Фаза 1: База
  event.go → condition.go → storage.go → memory_storage.go
  → type_manager.go → engine.go → genesis.go
  → тесты

Фаза 2: Агрегаты
  entity.go → attribute.go → property.go → boxed_value.go
  → тесты

Фаза 3: Транзакции
  Расширение EngineBase методами CRUD
  → тесты

Фаза 4: ConditionEngine
  condition/evaluator.go
  → тесты
```

### 5.2. Правило одного файла

Если файл превышает 500 строк — декомпозировать.

### 5.3. Правило одного пакета

Пакет `engine` содержит ядро. Подпакеты:
- `engine` — ядро (основной)
- `condition` — ConditionEngine (исполнитель)
- `query` — QueryEngine (запросы, позже)

Никаких других подпакетов без обсуждения.

---

## 6. Пример допустимого кода

```go
package engine

import (
    "errors"
    "fmt"
    "time"
)

// EventData представляет базовое событие в событийном графе.
type EventData struct {
    ID           int
    BaseEventID  int
    ValueID      int
    Conditions   ConditionRule
    ActorEventID int
    Value        string
    Date         time.Time
}

// Validate проверяет корректность события перед записью.
func (e *EventData) Validate(currentPos int) error {
    if e.ID != currentPos {
        return fmt.Errorf("%w: got %d, expected %d", ErrInvalidPosition, e.ID, currentPos)
    }
    if e.BaseEventID >= e.ID {
        return fmt.Errorf("%w: base %d >= event %d", ErrInvalidBaseRef, e.BaseEventID, e.ID)
    }
    return nil
}
```

---

## 7. Чек-лист перед коммитом

- [ ] `go fmt ./...` — без изменений
- [ ] `go vet ./...` — без ошибок
- [ ] `go build ./...` — собирается
- [ ] `go test ./... -count=1` — все тесты проходят
- [ ] Нет лишних комментариев (TODO без даты — удалить)
- [ ] Нет закомментированного кода
- [ ] Экспортируемые типы документированы
- [ ] Ошибки обёрнуты контекстом
- [ ] Нет сторонних зависимостей
- [ ] Код следует принципам из раздела 1
