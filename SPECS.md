# EventFlowEngine — Спецификация проекта

> **Версия:** 0.1.0 (архитектурная)  
> **Язык:** Go 1.22+  
> **Модуль:** `github.com/eventflow/engine`  
> **Статус:** Проектирование

---

## 1. Общая спецификация

### 1.1. Назначение

EventFlowEngine — библиотека (Go package) для создания и обработки семантических событийных графов. Предоставляет:

- Хранение событий в виде направленного ациклического графа (DAG)
- Семантическую типизацию через генезис-граф (мета-онтологию)
- Создание моделей сущностей и действий
- Создание индивидов по моделям
- Валидацию событий по моделям и ограничениям
- Исполнение бизнес-логики через систему условий (Condition)
- Полный event sourcing (immutable append-only)

### 1.2. Область применения

- Семантические платформы для работы с потоками событий
- Системы документооборота с онтологическим моделированием
- P2P-платформы с событийной архитектурой
- Системы с полным аудитом и event sourcing

### 1.3. Ограничения (scope)

**Входит в первую фазу:**
- ✅ In-memory хранилище
- ✅ Генезис-граф (мета-онтология)
- ✅ Базовая система условий (EventCondition, Equality, Conjunction, Disjunction)
- ✅ Модели и индивиды
- ✅ Базовые транзакции

**НЕ входит в первую фазу:**
- ❌ P2P-сеть
- ❌ Консенсус
- ❌ Распределённое хранение
- ❌ Язык запросов (QueryEngine)
- ❌ UI / API
- ❌ Персистентное хранилище (PostgreSQL, IPFS)

---

## 2. Объектная модель

### 2.1. EventData (базовое событие)

```go
// EventData — единственная примитивная сущность в системе.
// Всё остальное (акторы, сущности, атрибуты, модели, индивиды) —
// это события со ссылками на другие события.
type EventData struct {
    ID           int           // уникальный идентификатор
    BaseEventID  int           // событие-родитель (на чём фиксируется)
    ValueID      int           // тип значения (семантический тип)
    Conditions   ConditionRule // условие/причина (Causals)
    ActorEventID int           // актор, создавший событие
    Value        string        // значение события
    Date         time.Time     // время создания (UTC)
}
```

**Инварианты:**
- `ID` — монотонно возрастает (Position движка)
- `BaseEventID` — всегда ссылается на существующее событие
- `ValueID` — всегда один из предопределённых типов генезис-графа
- `ActorEventID` — ссылается на событие актора
- `Date` — устанавливается в момент создания, UTC

### 2.2. Типы событий (StaticEvent)

```go
// Базовые типы генезис-графа
const (
    Event    = 1  // корень графа
    SubEvent = 2  // подтип события

    Actor              = 3   // актор
    Entity             = 4   // сущность
    Relation           = 5   // отношение
    Attribute          = 6   // атрибут
    AttributeConstraint = 7  // ограничение атрибута
    AttributeValue     = 8   // значение атрибута
    Model              = 9   // модель
    Individual         = 10  // индивид
    Role               = 11  // роль
    ValueProperty      = 12  // свойство значения

    // Ограничения (подтипы ValueProperty)
    DataType    = 13
    Cardinality = 15
    Required    = 16
    Permission  = 17
    Set         = 18
    Mutable     = 19

    // Модели
    EventModel      = 22
    EntityModel     = 23
    RelationModel   = 24
    DataTypeModel   = 25
    AttributeModel  = 26
    ActorModel      = 27
    RoleModel       = 28
)
```

### 2.3. ConditionRule (система условий)

```go
type ConditionRule interface {
    IsMet(ctx *EventContext) bool
}

// EventConditionRule — событие X должно существовать
type EventConditionRule struct {
    EventID int
}

// PropertyEqualityRule — свойство = значение
type PropertyEqualityRule struct {
    PropertyID int
    Value      string
}

// PropertyInequalityRule — свойство ≠ значение
type PropertyInequalityRule struct {
    PropertyID int
    Value      string
}

// ConjunctionRule — И (все)
type ConjunctionRule struct {
    Values []ConditionRule
}

// DisjunctionRule — ИЛИ (любое)
type DisjunctionRule struct {
    Values []ConditionRule
}
```

### 2.4. Агрегаты

```go
// Entity — семантический тип индивида
type Entity struct {
    ID          int
    Label       string
    Models      []*Model
    Individuals []*Individual
}

// Model — описание сущности/действия (список свойств с ограничениями)
type Model struct {
    ID              int
    Label           string
    ParentModel     int
    EventBase       int
    PropertyProvider *PropertyProvider
}

// Individual — конкретный экземпляр сущности/действия
type Individual struct {
    ID              int
    Label           string
    Actor           int
    Valid           bool
    Model           *Model
    PropertyProvider *PropertyProvider
}
```

### 2.5. Свойства

```go
// PropertyProvider — определяет набор свойств (для модели или контейнера)
type PropertyProvider struct {
    ID         int
    Attributes map[int]*AttachedAttribute
    Relations  map[int]*AttachedRelation
    Events     map[int]*AttachedEvent
}

// AttachedAttribute — атрибут, привязанный к модели с ограничениями
type AttachedAttribute struct {
    Attribute       *Attribute
    Required        bool
    Mutable         bool
    Cardinality     int     // -1 = ∞, 0 = опционально, 1 = единичное
    Permission      *int    // ID актора/роли
    DefaultValue    *BoxedValue
    Conditions      ConditionRule
    PropertyProvider *PropertyProvider // вложенные свойства
}

// AttachedRelation — отношение, привязанное к модели
type AttachedRelation struct {
    Relation        *Relation
    Required        bool
    Mutable         bool
    Cardinality     int
    Permission      *int
    Conditions      ConditionRule
    PropertyProvider *PropertyProvider
}

// Attribute — определение атрибута (тип данных, возможные значения)
type Attribute struct {
    ID       int
    Label    string
    DataType *DataType
    IsBoxed  bool
    Values   []*AttributeValue
}

// Relation — определение отношения (связь с другими индивидами)
type Relation struct {
    ID     int
    Label  string
}

// PropertyContainer — контейнер значений свойств индивида
type PropertyContainer struct {
    Provider   *PropertyProvider
    Attributes map[int][]*PropertyContainer
    Relations  map[int][]*PropertyContainer
    Value      *BoxedValue
    Fixed      bool
}

// BoxedValue — значение с метаданными
type BoxedValue struct {
    PlainValue string
    ShownValue string
    Date       time.Time
    EventID    int
}
```

---

## 3. Storage (интерфейс хранилища)

```go
type Storage interface {
    // Управление транзакциями
    Begin() error
    Commit() error
    Rollback() error

    // Запись (append-only)
    AppendEvent(event EventData) error
    NextID() int

    // Чтение событий
    GetEvent(id int) (*EventData, error)
    GetEventsByBase(baseEventID int) ([]*EventData, error)
    GetEventsByValue(valueID int) ([]*EventData, error)
    GetEventsByActor(actorEventID int) ([]*EventData, error)

    // Агрегаты (строятся из событий)
    GetEntity(entityID int) (*Entity, error)
    GetEntities() ([]*Entity, error)
    
    GetModel(modelID int) (*Model, error)
    
    GetIndividual(individualID int) (*Individual, error)
    GetIndividuals() ([]*Individual, error)

    // Провайдеры свойств
    GetPropertyProvider(providerID int) (*PropertyProvider, error)
    GetPropertyProviderAttribute(providerID, attrID int) (*AttachedAttribute, error)

    // Контейнеры
    GetPropertyContainer(containerID int) (*PropertyContainer, error)
}
```

---

## 4. EngineBase (ядро)

```go
type EngineBase struct {
    Position int          // счётчик событий (следующий ID)
    Storage  Storage      // хранилище
    Types    *TypeManager // менеджер типов
}

// Создание движка
func NewEngine(storage Storage) *EngineBase

// Инициализация генезис-графа
func (e *EngineBase) InitGenesis() error

// Обработка события (валидация + запись)
func (e *EngineBase) ProcessEvent(event EventData) error
```

**Валидация ProcessEvent:**
1. `event.ID == Position` (иначе ошибка)
2. `event.BaseEventID < event.ID` (ссылка только на существующее)
3. `event.ValueID` — валидный тип из TypeManager
4. `event.ActorEventID` — валидный актор
5. Запись в Storage.AppendEvent
6. `Position++`

---

## 5. Транзакции (CRUD-операции)

```go
// Сущности
func (e *EngineBase) CreateEntity(actorID int, label string) (int, error)

// Атрибуты
func (e *EngineBase) CreateAttribute(actorID int, label string, dataType int) (int, error)
func (e *EngineBase) AssignAttributeDataType(actorID, attrID, dataType int) error
func (e *EngineBase) AssignAttributeValue(actorID, attrID int, value string) error

// Модели
func (e *EngineBase) CreateModel(actorID, eventBase, parentModel int, label string) (int, error)
func (e *EngineBase) AssignProviderAttribute(actorID, providerID, attrID int, conditions ConditionRule) (int, error)
func (e *EngineBase) AssignProviderRelation(actorID, providerID, relationID int, conditions ConditionRule) (int, error)

// Ограничения
func (e *EngineBase) SetRequired(actorID, assignationID, propertyID int, required bool) error
func (e *EngineBase) SetCardinality(actorID, assignationID, propertyID int, cardinality int) error
func (e *EngineBase) SetMutable(actorID, assignationID, propertyID int, mutable bool) error
func (e *EngineBase) SetPermission(actorID, assignationID, propertyID, permission int) error
func (e *EngineBase) SetDefaultValue(actorID, assignationID, propertyID int, value string) error

// Индивиды
func (e *EngineBase) CreateIndividual(actorID, eventBase, modelID int, label string) (int, error)
func (e *EngineBase) AssignContainerProperty(actorID, containerID, propertyID int, value string) (int, error)
```

---

## 6. Генезис-граф (GenesisData)

При инициализации движка создаются следующие события (порядок важен!):

```
# Базовые типы
1: 1: Event → "Event"              (Event: Event: "Event")
2: 1: SubEvent → "SubEvent"        (Event → SubEvent)

# Основные сущности онтологии
3: 1: SubEvent → "Actor"
4: 1: SubEvent → "Entity"
5: 1: SubEvent → "Relation"
6: 1: SubEvent → "Attribute"
7: 1: SubEvent → "AttributeConstraint"
9: 1: SubEvent → "Model"
10: 1: SubEvent → "Individual"
11: 1: SubEvent → "Role"
12: 1: SubEvent → "ValueProperty"

# Ограничения (как подтипы ValueProperty)
13: 12: SubEvent → "DataType"
15: 12: SubEvent → "Cardinality"
16: 12: SubEvent → "Required"
17: 12: SubEvent → "Permission"
18: 12: SubEvent → "Set"
19: 12: SubEvent → "Mutable"

# Модели
22: 1: Model → "Model_Event"
23: 4: Model → "Model_Entity"
24: 5: Model → "Model_Relation"
25: 13: Model → "Model_DataType"
26: 6: Model → "Model_Attribute"
27: 3: Model → "Model_Actor"
28: 11: Model → "Model_Role"

# Предустановленные индивиды
14: 13: Individual "basic_type" [Model=DataTypeModel]
    → AttributeConstraint (16: 14, property=14)
25: 13: Individual "enum_type" [Model=DataTypeModel]
17: 6: Individual "Name" [Model=AttributeModel]
    → DataType (18: 17, dataType=14)
    → ModelAttribute (19: 22, attr=17)
21: 3: Individual "Actor_Main" [Model=ActorModel]
    → IndividualAttribute (22: 21, attr=17, value="Main Actor")
```

---

## 7. ConditionEngine (исполнитель условий)

```go
type ConditionEngine struct {
    storage Storage
}

func NewConditionEngine(storage Storage) *ConditionEngine

// Evaluate — проверяет выполнение условия в контексте
func (ce *ConditionEngine) Evaluate(rule ConditionRule, ctx *EventContext) bool

// EventContext — контекст для проверки условий
type EventContext struct {
    Event        *EventData
    Properties   map[int]string // propertyID → value
    ActorID      int
}
```

**Правила проверки:**
- `EventConditionRule` → событие с ID существует в Storage
- `PropertyEqualityRule` → значение свойства в контексте = value
- `PropertyInequalityRule` → значение свойства ≠ value
- `ConjunctionRule` → все дочерние правила истинны
- `DisjunctionRule` → хотя бы одно дочернее правило истинно

---

## 8. Типы данных (DataType)

```go
type DataType struct {
    ID      int
    Label   string
    IsBoxed bool // boxed = enum (выбор из списка)
}

// Предустановленные типы данных
var BuiltinDataTypes = []*DataType{
    {ID: 14, Label: "basic_type", IsBoxed: false},  // string, int и т.д.
    {ID: 25, Label: "enum_type",  IsBoxed: true},   // выбор из списка
}
```

---

## 9. Схема отношений

```
Entity
  └── Model (valueID=Model, value=label)
       └── PropertyProvider (modelID = providerID)
            ├── AttachedAttribute[]
            │    ├── Attribute (ссылка на определение)
            │    ├── Required / Cardinality / Mutable / Permission
            │    ├── Condition (условие)
            │    └── PropertyProvider (вложенные)
            └── AttachedRelation[]
                 ├── Relation (ссылка на определение)
                 ├── Required / Cardinality / Mutable / Permission
                 └── PropertyProvider (вложенные)

Individual (valueID=Individual)
  └── PropertyContainer (individualID = containerID)
       ├── Attributes[attrID] → PropertyContainer[]
       └── Relations[relID] → PropertyContainer[]
```

---

## 10. Критерии готовности (Definition of Done)

1. Все типы из StaticEvent определены в константах
2. EventData и ConditionRule реализованы
3. MemoryStorage проходит базовые тесты на запись/чтение
4. EngineBase.ProcessEvent работает с валидацией
5. InitGenesis создаёт генезис-граф (22+ события)
6. CreateEntity / CreateModel / CreateIndividual работают
7. AssignProviderAttribute / AssignProviderRelation работают
8. Условия (ConditionRule) проверяются корректно
9. Пакет собирается без ошибок
10. Тесты покрывают ключевые сценарии (≥70%)
