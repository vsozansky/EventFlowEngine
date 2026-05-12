package engine

// Entity (сущность) — семантический тип индивида.
//
// Сущность определяет, к какому классу относится индивид.
// Например: "Товар", "Запрос", "Документ", "Организация".
//
// Сущность собирается из событий графа: базовое событие с
// ValueID = TypeEntity и все связанные с ним модели.
type Entity struct {
	// ID — идентификатор сущности (ID события с ValueID=TypeEntity).
	ID int

	// Label — название сущности (из Value события).
	Label string

	// Models — модели, доступные для этой сущности.
	Models []*Model

	// Individuals — индивиды этой сущности.
	Individuals []*Individual
}

// Model (модель) — описание сущности или действия.
//
// Модель содержит список свойств (атрибутов и отношений) с ограничениями:
// обязательность, множественность, изменяемость, права доступа, условия.
//
// Каждая модель привязана к базовому событию (EventBase) и имеет
// родительскую модель (ParentModel), от которой наследует структуру.
//
// Пример модели бизнес-процесса "Обработка запроса о товаре":
//
//	Processing request: Model
//	  : Relation: Subject       (Permission: customer)
//	  : Attribute: Offer        (Permission: employee)
//	    :: Attribute: Comment   (Permission: manager)
//	    :: Attribute: Solution  (Permission: manager)
//	    :: Attribute: Confirmation (Permission: customer)
//	  : Attribute: Status
type Model struct {
	// ID — идентификатор модели (ID события с ValueID=TypeModel).
	ID int

	// Label — название модели.
	Label string

	// ParentModel — ID родительской модели (0, если нет родителя).
	ParentModel int

	// EventBase — ID базового события, к которому привязана модель.
	EventBase int

	// PropertyProvider — провайдер свойств модели.
	// Содержит атрибуты и отношения с их ограничениями.
	PropertyProvider *PropertyProvider
}

// Individual (индивид) — конкретный экземпляр сущности или действия.
//
// Индивид создаётся по модели. Например, если есть сущность "Товар"
// и модель "Товар", то конкретный товар "Насос" — это индивид.
//
// Индивид содержит значения свойств (через PropertyProvider),
// которые могут быть введены актором или вычислены автоматически.
type Individual struct {
	// ID — идентификатор индивида (ID события с ValueID=TypeIndividual).
	ID int

	// Label — название индивида.
	Label string

	// Actor — ID актора, создавшего индивида.
	Actor int

	// Valid — true, если все обязательные свойства заполнены.
	Valid bool

	// Model — модель, по которой создан индивид.
	Model *Model

	// PropertyProvider — провайдер свойств индивида.
	// Содержит значения свойств, привязанные к этому индивиду.
	PropertyProvider *PropertyProvider
}

// DataType (тип данных) — семантический тип данных атрибута.
//
// Определяет, какие значения может принимать атрибут:
//   - basic_type: строка, число, дата (произвольное значение)
//   - enum_type: выбор из предопределённого списка
type DataType struct {
	// ID — идентификатор типа данных.
	ID int

	// Label — название типа данных.
	Label string

	// IsBoxed — true, если это enum (значение выбирается из списка).
	IsBoxed bool
}

// AttributeValue — возможное значение атрибута (для enum-типов).
type AttributeValue struct {
	// EventID — ID события, фиксирующего это значение.
	EventID int

	// Value — строковое представление значения.
	Value string
}

// Attribute (атрибут) — определение свойства с типом данных.
//
// Атрибут — это многократно используемое определение свойства.
// Например, атрибут "Name" со строковым типом может использоваться
// в разных моделях.
//
// Поддерживает:
//   - Базовые типы (basic_type): string, int, float
//   - Перечисляемые типы (enum_type): выбор из списка
type Attribute struct {
	// ID — идентификатор атрибута.
	ID int

	// Label — название атрибута (например, "Name", "Price", "Status").
	Label string

	// DataType — тип данных атрибута.
	DataType *DataType

	// IsBoxed — true, если атрибут enum (выбор из значений).
	IsBoxed bool

	// Values — возможные значения (для enum-типов).
	Values []*AttributeValue
}

// Relation (отношение) — связь между индивидами.
//
// Отношение фиксирует, что один индивид связан с другим.
// Например: "является частью", "принадлежит организации", "создан актором".
type Relation struct {
	// ID — идентификатор отношения.
	ID int

	// Label — название отношения.
	Label string
}
