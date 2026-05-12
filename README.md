# EventFlowEngine

**Семантический event-driven движок на основе событийной онтологии**

> **Философский фундамент:** Субъектная онтология А.В. Болдачева  
> **Язык реализации:** Go  
> **Статус:** Архитектурное проектирование / Прототип  
> **Лицензия:** MIT

---

## О проекте

**EventFlowEngine** — это независимая реализация ядра семантического движка для работы с потоками событий (event flow) на основе онтологического моделирования. Движок позволяет:

- Создавать семантические модели сущностей и действий
- Формировать направленный ациклический граф событий (DAG)
- Валидировать события по моделям
- Исполнять бизнес-логику через условия (Condition)
- Обеспечивать полный event sourcing (immutable, append-only)
- Работать в event-driven архитектуре

Проект является Go-реализацией концепций, заложенных в платформах **Boldsea** / **EventFlow** (автор: Александр Болдачев) и **Parallax** (автор: Frank Horrigan).

---

## Философские основания

Движок построен на принципах **субъектной онтологии**:

| Принцип | Реализация в движке |
|---------|-------------------|
| Данность | Событийный граф — всё существует как событие |
| Субъект → Актор | Авторизованный источник событий |
| Объект → Индивид | Различаемый/изменяемый через события |
| Сознание → Граф | Единое пространство событий |
| Приватность | Permission, права доступа |
| Совпадение по указанию | Консенсус (в перспективе) |

---

## Архитектура

```
                    ┌─────────────────────┐
                    │   Приложения / API   │
                    └──────────┬──────────┘
                               │
                    ┌──────────▼──────────┐
                    │   EventFlowEngine    │
                    │     (Core Engine)    │
                    │                      │
                    │  ┌────────────────┐  │
                    │  │  EngineBase    │  │
                    │  │  • ProcessEvent│  │
                    │  │  • InitGenesis │  │
                    │  └───────┬────────┘  │
                    │          │           │
                    │  ┌───────▼────────┐  │
                    │  │    Storage     │  │
                    │  │  • Append-only │  │
                    │  │  • Event Graph │  │
                    │  │  • Queries     │  │
                    │  └───────┬────────┘  │
                    │          │           │
                    │  ┌───────▼────────┐  │
                    │  │   Condition    │  │
                    │  │   Engine       │  │
                    │  └────────────────┘  │
                    └──────────────────────┘
```

### Ключевые компоненты

| Компонент | Описание |
|-----------|----------|
| **EngineBase** | Ядро: обработка событий, валидация, запись в граф |
| **Storage** | Хранилище графа (in-memory, с возможностью смены) |
| **ConditionEngine** | Исполнитель условий бизнес-логики |
| **QueryEngine** | Язык запросов к графу |
| **TypeManager** | Менеджер семантических типов |

---

## Основные принципы

1. **Всё есть событие** — единый унифицированный формат
2. **Только по моделям** — каждое событие создаётся по модели
3. **Immutable / Append-only** — события не удаляются и не изменяются
4. **Event Sourcing** — полное проигрывание графа с любого момента
5. **Event-Driven** — реактивное взаимодействие через подписки
6. **Функциональный стиль** — чистые функции, рекурсия
7. **Семантическая типизация** — строгая через онтологию

---

## Быстрый старт

```go
package main

import (
    "fmt"
    "github.com/eventflow/engine"
)

func main() {
    // Создание хранилища и движка
    store := engine.NewMemoryStorage()
    eng := engine.NewEngine(store)
    
    // Инициализация генезис-графа
    err := eng.InitGenesis()
    if err != nil {
        panic(err)
    }
    
    fmt.Println("EventFlowEngine initialized at position:", eng.Position)
}
```

---

## Структура пакета

```
EventFlowEngine/
├── README.md              — описание проекта
├── SPECS.md               — полная спецификация
├── RULES.md               — правила разработки
├── go.mod / go.sum
├── engine/
│   ├── engine.go          — EngineBase, ProcessEvent
│   ├── event.go           — EventData, StaticEvent
│   ├── condition.go       — ConditionRule и реализации
│   ├── storage.go         — Storage interface
│   ├── memory_storage.go  — MemoryStorage
│   ├── type_manager.go    — TypeManager
│   ├── genesis.go         — GenesisData
│   ├── entity.go          — Entity, Model, Individual
│   ├── attribute.go       — Attribute, Relation
│   ├── property.go        — PropertyProvider, PropertyContainer
│   └── boxed_value.go     — BoxedValue
├── condition/
│   └── evaluator.go       — ConditionEngine
└── query/
    └── query.go           — QueryEngine
```

---

## Документация

Полная документация по архитектуре и объектной модели находится в `B:\ai-hermes\docs\`:

| Файл | Описание |
|------|----------|
| `Boldsea_Engine_Architecture.md` | Архитектура ядра |
| `Boldsea_BPMN_Example.md` | Пример бизнес-процесса |
| `Boldsea_Parallax_ObjectModel.md` | Объектная модель и Go-спецификация |

Исходные материалы (архив): `B:\ai-hermes\docs\Boldsea-ARHIVE\`

---

## Статус разработки

- [x] Архитектурное проектирование
- [ ] Базовая структура Go-пакета
- [ ] EventData и StaticEvent
- [ ] MemoryStorage
- [ ] EngineBase с ProcessEvent
- [ ] GenesisData и InitGenesis
- [ ] ConditionRule и Evaluator
- [ ] Entity / Model / Individual агрегаты
- [ ] PropertyProvider / PropertyContainer
- [ ] Транзакции (CRUD)
- [ ] QueryEngine
- [ ] Тесты
- [ ] Документация GoDoc
