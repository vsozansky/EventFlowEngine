package engine

import "fmt"

// TypeManager — менеджер семантических типов событий.
//
// Хранит реестр всех известных типов генезис-графа и их строковые
// обозначения (labels). Используется для валидации ValueID при
// записи новых событий.
type TypeManager struct {
	// types — карта: ID типа → его строковое название
	types map[int]string

	// labels — обратная карта: название → ID типа
	labels map[string]int
}

// NewTypeManager создаёт менеджер типов со всеми предустановленными
// типами генезис-графа. Эти типы являются мета-онтологией движка
// и не могут быть изменены.
func NewTypeManager() *TypeManager {
	tm := &TypeManager{
		types:  make(map[int]string),
		labels: make(map[string]int),
	}

	// Регистрируем все базовые типы генезис-графа
	tm.mustRegister(TypeEvent, "Event")
	tm.mustRegister(TypeSubEvent, "SubEvent")
	tm.mustRegister(TypeActor, "Actor")
	tm.mustRegister(TypeEntity, "Entity")
	tm.mustRegister(TypeRelation, "Relation")
	tm.mustRegister(TypeAttribute, "Attribute")
	tm.mustRegister(TypeAttributeConstraint, "AttributeProperty")
	tm.mustRegister(TypeAttributeValue, "AttributeConstraint")
	tm.mustRegister(TypeModel, "Model")
	tm.mustRegister(TypeIndividual, "Individual")
	tm.mustRegister(TypeRole, "Role")
	tm.mustRegister(TypeValueProperty, "ValueProperty")
	tm.mustRegister(TypeDataType, "DataType")
	tm.mustRegister(TypeCardinality, "Cardinality")
	tm.mustRegister(TypeRequired, "Required")
	tm.mustRegister(TypePermission, "Permission")
	tm.mustRegister(TypeSet, "Set")
	tm.mustRegister(TypeMutable, "Mutable")
	tm.mustRegister(TypeEventModel, "Model_Event")
	tm.mustRegister(TypeEntityModel, "Model_Entity")
	tm.mustRegister(TypeRelationModel, "Model_Relation")
	tm.mustRegister(TypeDataTypeModel, "Model_DataType")
	tm.mustRegister(TypeAttributeModel, "Model_Attribute")
	tm.mustRegister(TypeActorModel, "Model_Actor")
	tm.mustRegister(TypeRoleModel, "Model_Role")

	return tm
}

// mustRegister добавляет тип в реестр. Паникует при дубликате ID или label.
// Используется только при инициализации, поэтому паника допустима.
func (tm *TypeManager) mustRegister(id int, name string) {
	if _, ok := tm.types[id]; ok {
		panic(fmt.Sprintf("тип с ID %d уже зарегистрирован: %s", id, tm.types[id]))
	}
	if _, ok := tm.labels[name]; ok {
		panic(fmt.Sprintf("тип с именем %s уже зарегистрирован: %d", name, tm.labels[name]))
	}
	tm.types[id] = name
	tm.labels[name] = id
}

// IsValid проверяет, существует ли тип с указанным ID в реестре.
func (tm *TypeManager) IsValid(typeID int) bool {
	_, ok := tm.types[typeID]
	return ok
}

// Label возвращает строковое название типа по его ID.
// Возвращает пустую строку, если тип не найден.
func (tm *TypeManager) Label(typeID int) string {
	return tm.types[typeID]
}

// TypeID возвращает ID типа по его строковому названию.
// Возвращает 0 и false, если название не найдено.
func (tm *TypeManager) TypeID(name string) (int, bool) {
	id, ok := tm.labels[name]
	return id, ok
}

// All возвращает копию карты всех зарегистрированных типов.
func (tm *TypeManager) All() map[int]string {
	result := make(map[int]string, len(tm.types))
	for k, v := range tm.types {
		result[k] = v
	}
	return result
}
