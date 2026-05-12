// Пакет engine — ядро семантического событийного движка EventFlowEngine.
//
// Философский фундамент: субъектная онтология А.В. Болдачева.
// Архитектурная основа: Parallax (Frank Horrigan) / EventFlow (Alexander Boldachev).
//
// Принципы:
//   - Всё есть событие: любые данные хранятся как события в графе
//   - Immutable: события не изменяются после записи
//   - Append-only: запись только добавлением в конец
//   - Event Sourcing: полное проигрывание графа с любого момента
package engine

import "time"

// EventData представляет базовое событие в событийном графе.
//
// Вся информация в системе (акторы, сущности, атрибуты, модели, индивиды,
// значения свойств) хранится как события. Агрегаты собираются из событий
// по запросу, а не хранятся отдельно.
//
// Инварианты:
//   - ID — монотонно возрастает, равен EngineBase.Position в момент записи
//   - BaseEventID < ID — ссылка на существующее событие-родитель
//   - ValueID — один из предопределённых типов генезис-графа
//   - ActorEventID — ссылка на событие актора, создавшего запись
type EventData struct {
	// ID — уникальный идентификатор события (позиция в графе).
	// Назначается движком автоматически.
	ID int

	// BaseEventID — событие, на котором зафиксировано данное.
	// Например, для атрибута имени это будет событие индивида.
	BaseEventID int

	// ValueID — семантический тип события.
	// Определяет, что именно фиксируется: атрибут, отношение,
	// модель, индивид, ограничение и т.д.
	ValueID int

	// Conditions — условие, при котором событие может быть создано.
	// Аналог Causals из спецификации EventFlow.
	Conditions ConditionRule

	// ActorEventID — событие актора, создавшего запись.
	// Обеспечивает аудит и контроль доступа.
	ActorEventID int

	// Value — строковое значение события.
	// Для атрибутов — значение атрибута.
	// Для моделей — label модели.
	// Для ограничений — значение ограничения ("1", "0", ID и т.д.).
	Value string

	// Date — время создания события (UTC).
	Date time.Time
}

// Константы типов генезис-графа (Genesis Graph).
//
// Эти константы определяют мета-онтологию движка — базовые семантические
// типы, из которых строятся все остальные. Соответствуют StaticEvent
// из Parallax/AuroraCore.
const (
	// === Базовые типы (Core Types) ===

	// TypeEvent — корневой тип, основание всего графа.
	// Первое событие в любом экземпляре движка.
	TypeEvent = 1

	// TypeSubEvent — подтип события.
	// Используется для создания иерархии типов.
	TypeSubEvent = 2

	// === Примитивы онтологии (Ontology Primitives) ===

	// TypeActor — актор: авторизованный источник событий.
	// Аналог субъекта в субъектной онтологии.
	TypeActor = 3

	// TypeEntity — сущность: семантический тип индивида.
	// Например: "Товар", "Запрос", "Документ".
	TypeEntity = 4

	// TypeRelation — отношение: связь между индивидами.
	// Например: "является частью", "принадлежит".
	TypeRelation = 5

	// TypeAttribute — атрибут: свойство с фиксированным типом данных.
	// Например: "цвет", "цена", "имя".
	TypeAttribute = 6

	// TypeAttributeConstraint — ограничение атрибута.
	// Определяет допустимые значения.
	TypeAttributeConstraint = 7

	// TypeAttributeValue — значение атрибута (для enum-типов).
	TypeAttributeValue = 8

	// TypeModel — модель: описание сущности или действия.
	// Содержит список свойств с ограничениями.
	TypeModel = 9

	// TypeIndividual — индивид: конкретный экземпляр сущности.
	// Создаётся по модели.
	TypeIndividual = 10

	// TypeRole — роль: набор прав доступа.
	TypeRole = 11

	// TypeValueProperty — свойство значения: базовый тип для ограничений.
	TypeValueProperty = 12

	// === Ограничения (Restrictions / ValueProperty subtypes) ===

	// TypeDataType — тип данных атрибута (string, int, enum...).
	TypeDataType = 13

	// TypeCardinality — количество допустимых значений.
	// -1 = неограниченно, 0 = опционально, 1 = единичное.
	TypeCardinality = 15

	// TypeRequired — обязательность значения (1 = обязательно, 0 = нет).
	TypeRequired = 16

	// TypePermission — права доступа (ID актора или роли).
	TypePermission = 17

	// TypeSet — значение по умолчанию.
	TypeSet = 18

	// TypeMutable — разрешение изменять значение (1 = можно, 0 = нельзя).
	TypeMutable = 19

	// === Предустановленные модели (Predefined Models) ===

	// TypeEventModel — модель для событий.
	TypeEventModel = 22

	// TypeEntityModel — модель для сущностей.
	TypeEntityModel = 23

	// TypeRelationModel — модель для отношений.
	TypeRelationModel = 24

	// TypeDataTypeModel — модель для типов данных.
	TypeDataTypeModel = 25

	// TypeAttributeModel — модель для атрибутов.
	TypeAttributeModel = 26

	// TypeActorModel — модель для акторов.
	TypeActorModel = 27

	// TypeRoleModel — модель для ролей.
	TypeRoleModel = 28
)

// Предопределённые ID для генезис-графа.
// Используются при инициализации движка.
const (
	// ID констант событий генезис-графа
	GenesisIDBasicType       = 14 // basic_type — простой тип данных (строка, число)
	GenesisIDEnumType        = 25 // enum_type — перечисляемый тип данных
	GenesisIDAttrName        = 17 // Name — атрибут имени
	GenesisIDMainActor       = 21 // Actor_Main — главный актор (системный)
	GenesisDataTypeBasicType = 14 // ID типа данных "basic_type"
	GenesisDataTypeEnumType  = 25 // ID типа данных "enum_type"
)

// EventContext содержит контекст для проверки условий выполнения событий.
// Передаётся в ConditionRule.IsMet() для принятия решения.
type EventContext struct {
	// Event — текущее проверяемое событие.
	Event *EventData

	// Properties — значения свойств в текущем контексте.
	// Ключ — ID свойства, значение — текущее значение.
	Properties map[int]string

	// ActorID — ID актора, выполняющего действие.
	ActorID int
}
