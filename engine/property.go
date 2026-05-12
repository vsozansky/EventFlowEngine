package engine

// PropertyProvider (провайдер свойств) — определяет набор свойств
// для модели или контейнера.
//
// Провайдер содержит три категории свойств:
//   - Attributes: атрибуты с ограничениями (Required, Cardinality, и т.д.)
//   - Relations: отношения с ограничениями
//   - Events: системные события, привязанные к провайдеру
//
// Провайдеры образуют иерархию: каждый AttachedAttribute или
// AttachedRelation может содержать вложенный PropertyProvider
// для поддержки составных свойств (например, у атрибута "Предложение"
// есть вложенные атрибуты "Комментарий" и "Решение").
type PropertyProvider struct {
	// ID — идентификатор провайдера.
	ID int

	// Attributes — карта ID атрибута → AttachedAttribute (с ограничениями).
	Attributes map[int]*AttachedAttribute

	// Relations — карта ID отношения → AttachedRelation (с ограничениями).
	Relations map[int]*AttachedRelation

	// Events — карта ID свойства → AttachedEvent (системные события).
	Events map[int]*AttachedEvent
}

// AttachedAttribute — атрибут, привязанный к модели с ограничениями.
//
// Содержит ссылку на определение атрибута (Attribute) и набор
// ограничений, определяющих правила работы с этим атрибутом
// в контексте конкретной модели.
//
// Поддерживает вложенные свойства через PropertyProvider:
// например, у атрибута "Предложение" (Offer) могут быть вложенные
// атрибуты "Комментарий", "Решение", "Подтверждение".
type AttachedAttribute struct {
	// Attribute — определение атрибута (ссылка).
	Attribute *Attribute

	// Required — true, если значение обязательно для заполнения.
	Required bool

	// Mutable — true, если значение можно изменять после установки.
	Mutable bool

	// Cardinality — количество допустимых значений:
	//   -1 = неограниченно, 0 = опционально, 1 = единичное
	Cardinality int

	// Permission — ID актора или роли, имеющей право на запись.
	// nil означает отсутствие ограничений.
	Permission *int

	// DefaultValue — значение по умолчанию.
	DefaultValue *BoxedValue

	// Conditions — условие, при котором атрибут доступен для заполнения.
	Conditions ConditionRule

	// PropertyProvider — вложенные свойства (для составных атрибутов).
	PropertyProvider *PropertyProvider
}

// AttachedRelation — отношение, привязанное к модели с ограничениями.
//
// Аналогично AttachedAttribute, но для отношений между индивидами.
// Значением отношения является ссылка на другой индивид.
type AttachedRelation struct {
	// Relation — определение отношения (ссылка).
	Relation *Relation

	// Required — true, если отношение обязательно.
	Required bool

	// Mutable — true, если значение можно изменить.
	Mutable bool

	// Cardinality — количество допустимых значений.
	Cardinality int

	// Permission — ID актора или роли с правом на запись.
	Permission *int

	// Conditions — условие доступности отношения.
	Conditions ConditionRule

	// PropertyProvider — вложенные свойства отношения.
	PropertyProvider *PropertyProvider
}

// AttachedEvent — системное событие, привязанное к провайдеру.
//
// Используется для системных действий: создание индивида,
// отправка уведомления, автоматическое вычисление значения.
type AttachedEvent struct {
	// Property — общее свойство (ID + Label).
	Property *Property

	// Value — значение системного события.
	Value string

	// Conditions — условие выполнения системного события.
	Conditions ConditionRule
}

// Property — общее свойство (атрибут или отношение).
// Используется для системных событий, где не важен тип свойства.
type Property struct {
	// PropertyID — ID свойства.
	PropertyID int

	// Label — название свойства.
	Label string
}

// PropertyContainer — контейнер значений свойств индивида.
//
// В отличие от PropertyProvider (который определяет, какие свойства
// должны быть у модели), PropertyContainer хранит фактические значения
// свойств для конкретного индивида.
//
// Например, для индивида "Иванов" (создан по модели "Сотрудник"):
//
//	Container:
//	  Attribute "Name" → "Иванов"
//	  Attribute "Age" → "30"
//	  Relation "Department" → ID индивида "Бухгалтерия"
type PropertyContainer struct {
	// Provider — провайдер, определяющий структуру контейнера.
	Provider *PropertyProvider

	// Attributes — значения атрибутов: ID атрибута → список контейнеров.
	Attributes map[int][]*PropertyContainer

	// Relations — значения отношений: ID отношения → список контейнеров.
	Relations map[int][]*PropertyContainer

	// Value — текущее значение (для вложенных свойств).
	Value *BoxedValue

	// Fixed — true, если значение фиксировано (не может быть изменено).
	Fixed bool
}

// BoxedValue — значение с метаданными.
//
// Содержит не только само значение, но и информацию о его источнике:
//   - PlainValue: строковое представление значения
//   - ShownValue: отображаемое значение (для enum — текст, а не ID)
//   - Date: когда было установлено
//   - EventID: ID события, установившего значение
type BoxedValue struct {
	// PlainValue — строковое представление значения.
	PlainValue string

	// ShownValue — отображаемое значение (может отличаться от PlainValue).
	ShownValue string

	// Date — время установки значения.
	Date int64

	// EventID — ID события, установившего это значение.
	EventID int
}
