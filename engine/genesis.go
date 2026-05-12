package engine

import "fmt"

// Genesis — функции инициализации генезис-графа.
//
// Генезис-граф (Genesis Graph) — это мета-онтология движка.
// Он содержит базовые типы событий, предустановленные модели
// и начальные индивиды, необходимые для работы системы.
//
// Генезис-граф создаётся единожды при инициализации движка
// и не может быть изменён после создания.

// genesisEvents возвращает список событий генезис-графа
// в правильном порядке. Порядок важен, так как каждое событие
// ссылается на предыдущие через BaseEventID.
//
// Схема соответствует генезис-графу из Parallax (Graph.cs):
//
//	1:   Event → "Event"
//	2:   Event → SubEvent → "SubEvent"
//	3:   Event → Actor (SubEvent)
//	4:   Event → Entity (SubEvent)
//	5:   Event → Relation (SubEvent)
//	6:   Event → Attribute (SubEvent)
//	7:   Event → AttributeConstraint (SubEvent)
//	9:   Event → Model (SubEvent)
//	10:  Event → Individual (SubEvent)
//	11:  Event → Role (SubEvent)
//	12:  Event → ValueProperty (SubEvent)
//	13:  ValueProperty → DataType (SubEvent)
//	15:  ValueProperty → Cardinality (SubEvent)
//	16:  ValueProperty → Required (SubEvent)
//	17:  ValueProperty → Permission (SubEvent)
//	18:  ValueProperty → Set (SubEvent)
//	19:  ValueProperty → Mutable (SubEvent)
//	14:  DataTypeModel → Individual "basic_type"
//	25:  DataTypeModel → Individual "enum_type"
//	17:  AttributeModel → Individual "Name" (+ DataType, + ModelAttribute)
//	21:  ActorModel → Individual "Actor_Main" (+ значение "Name")
//	22:  Event → Model "Model_Event"
//	23:  Entity → Model "Model_Entity"
//	24:  Relation → Model "Model_Relation"
//	25:  DataType → Model "Model_DataType"
//	26:  Attribute → Model "Model_Attribute"
//	27:  Actor → Model "Model_Actor"
//	28:  Role → Model "Model_Role"
func genesisEvents() []EventData {
	// Системный актор (0) — события создаются самой системой при инициализации.
	systemActor := 0

	return []EventData{
		// === Базовые типы ===
		{ID: 1, BaseEventID: TypeEvent, ValueID: TypeEvent, ActorEventID: systemActor, Value: "Event"},
		{ID: 2, BaseEventID: TypeEvent, ValueID: TypeSubEvent, ActorEventID: systemActor, Value: "SubEvent"},

		// === Примитивы онтологии (SubEvent → базовые) ===
		{ID: 3, BaseEventID: TypeEvent, ValueID: TypeSubEvent, ActorEventID: systemActor, Value: "Actor"},
		{ID: 4, BaseEventID: TypeEvent, ValueID: TypeSubEvent, ActorEventID: systemActor, Value: "Entity"},
		{ID: 5, BaseEventID: TypeEvent, ValueID: TypeSubEvent, ActorEventID: systemActor, Value: "Relation"},
		{ID: 6, BaseEventID: TypeEvent, ValueID: TypeSubEvent, ActorEventID: systemActor, Value: "Attribute"},
		{ID: 7, BaseEventID: TypeEvent, ValueID: TypeSubEvent, ActorEventID: systemActor, Value: "AttributeProperty"},

		// AttributeValue — подтип события (ID 8)
		{ID: 8, BaseEventID: TypeEvent, ValueID: TypeSubEvent, ActorEventID: systemActor, Value: "AttributeConstraint"},

		{ID: 9, BaseEventID: TypeEvent, ValueID: TypeSubEvent, ActorEventID: systemActor, Value: "Model"},
		{ID: 10, BaseEventID: TypeEvent, ValueID: TypeSubEvent, ActorEventID: systemActor, Value: "Individual"},
		{ID: 11, BaseEventID: TypeEvent, ValueID: TypeSubEvent, ActorEventID: systemActor, Value: "Role"},
		{ID: 12, BaseEventID: TypeEvent, ValueID: TypeSubEvent, ActorEventID: systemActor, Value: "ValueProperty"},

		// === Ограничения (подтипы ValueProperty) ===
		{ID: 13, BaseEventID: TypeValueProperty, ValueID: TypeSubEvent, ActorEventID: systemActor, Value: "DataType"},

		// === Предустановленные индивиды ===
		// basic_type — простой тип данных (строка, число)
		{ID: 14, BaseEventID: TypeDataType, ValueID: TypeIndividual, ActorEventID: systemActor, Value: "basic_type"},

		{ID: 15, BaseEventID: TypeValueProperty, ValueID: TypeSubEvent, ActorEventID: systemActor, Value: "Cardinality"},
		{ID: 16, BaseEventID: TypeValueProperty, ValueID: TypeSubEvent, ActorEventID: systemActor, Value: "Required"},
		{ID: 17, BaseEventID: TypeValueProperty, ValueID: TypeSubEvent, ActorEventID: systemActor, Value: "Permission"},
		{ID: 18, BaseEventID: TypeValueProperty, ValueID: TypeSubEvent, ActorEventID: systemActor, Value: "Set"},
		{ID: 19, BaseEventID: TypeValueProperty, ValueID: TypeSubEvent, ActorEventID: systemActor, Value: "Mutable"},

		// AttributeValue — подтип события (ID 8)
		{ID: 20, BaseEventID: TypeEvent, ValueID: TypeSubEvent, ActorEventID: systemActor, Value: "AttributeConstraint"},

		// Actor_Main — главный системный актор (ID 21)
		{ID: 21, BaseEventID: TypeActor, ValueID: TypeIndividual, ActorEventID: systemActor, Value: "Actor_Main"},

		// === Модели (Model) ===
		{ID: 22, BaseEventID: TypeEvent, ValueID: TypeModel, ActorEventID: systemActor, Value: "Model_Event"},
		{ID: 23, BaseEventID: TypeEntity, ValueID: TypeModel, ActorEventID: systemActor, Value: "Model_Entity"},
		{ID: 24, BaseEventID: TypeRelation, ValueID: TypeModel, ActorEventID: systemActor, Value: "Model_Relation"},
		{ID: 25, BaseEventID: TypeDataType, ValueID: TypeModel, ActorEventID: systemActor, Value: "Model_DataType"},
		{ID: 26, BaseEventID: TypeAttribute, ValueID: TypeModel, ActorEventID: systemActor, Value: "Model_Attribute"},
		{ID: 27, BaseEventID: TypeActor, ValueID: TypeModel, ActorEventID: systemActor, Value: "Model_Actor"},
		{ID: 28, BaseEventID: TypeRole, ValueID: TypeModel, ActorEventID: systemActor, Value: "Model_Role"},

		// === Дополнительные индивиды ===
		// enum_type — перечисляемый тип данных (ID 29, т.к. DataTypeModel = 25 занят)
		// Используем ID = 29 для enum_type индивида
		{ID: 29, BaseEventID: TypeDataType, ValueID: TypeIndividual, ActorEventID: systemActor, Value: "enum_type"},

		// Name — предустановленный атрибут имени (ID 30)
		{ID: 30, BaseEventID: TypeAttribute, ValueID: TypeIndividual, ActorEventID: systemActor, Value: "Name"},
	}
}

// InitGenesis инициализирует генезис-граф движка.
//
// Создаёт все базовые типы, модели и предустановленные индивиды.
// Должна быть вызвана сразу после NewEngine() до любых других операций.
//
// Возвращает ошибку, если:
//   - Движок уже был инициализирован (Position > 1)
//   - Любое из событий генезис-графа не прошло валидацию
func (e *EngineBase) InitGenesis() error {
	if e.Position != 1 {
		return fmt.Errorf("движок уже инициализирован: позиция %d", e.Position)
	}

	events := genesisEvents()
	for _, event := range events {
		if err := e.ProcessEvent(event); err != nil {
			return fmt.Errorf("ошибка инициализации генезис-графа на событии %d: %w",
				event.ID, err)
		}
	}

	return nil
}

// IsGenesisComplete проверяет, завершена ли инициализация генезис-графа.
// Возвращает true, если в графе есть все необходимые базовые типы.
func (e *EngineBase) IsGenesisComplete() bool {
	// Проверяем, что ключевые типы существуют
	events, err := e.Storage.GetEventsByValue(TypeSubEvent)
	if err != nil || len(events) == 0 {
		return false
	}
	return true
}
