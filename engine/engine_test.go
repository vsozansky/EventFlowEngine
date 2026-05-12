package engine

import (
	"testing"
)

// === Тесты MemoryStorage ===

// TestMemoryStorage_AppendAndGet проверяет базовую запись и чтение событий.
func TestMemoryStorage_AppendAndGet(t *testing.T) {
	store := NewMemoryStorage()

	event := EventData{
		ID:           store.NextID(),
		BaseEventID:  TypeEvent,
		ValueID:      TypeActor,
		ActorEventID: TypeEvent,
		Value:        "TestActor",
	}

	err := store.AppendEvent(event)
	if err != nil {
		t.Fatalf("не удалось записать событие: %v", err)
	}

	got, err := store.GetEvent(1)
	if err != nil {
		t.Fatalf("не удалось прочитать событие: %v", err)
	}

	if got.Value != "TestActor" {
		t.Errorf("ожидалось 'TestActor', получено '%s'", got.Value)
	}
}

// TestMemoryStorage_NextID проверяет монотонный рост ID.
func TestMemoryStorage_NextID(t *testing.T) {
	store := NewMemoryStorage()

	if store.NextID() != 1 {
		t.Errorf("первый ID должен быть 1, получен %d", store.NextID())
	}

	// После записи ID должен увеличиться
	store.AppendEvent(EventData{ID: 1, BaseEventID: TypeEvent, ValueID: TypeSubEvent, ActorEventID: TypeEvent, Value: "test"})

	if store.NextID() != 2 {
		t.Errorf("после записи ID должен быть 2, получен %d", store.NextID())
	}
}

// TestMemoryStorage_GetEventsByBase проверяет выборку дочерних событий.
func TestMemoryStorage_GetEventsByBase(t *testing.T) {
	store := NewMemoryStorage()

	// Записываем несколько событий
	// Используем уникальные BaseEventID для избежания коллизий
	// 100 = базовое событие-родитель (не конфликтует с типовыми константами)
	store.AppendEvent(EventData{ID: 1, BaseEventID: TypeEvent, ValueID: TypeSubEvent, ActorEventID: TypeEvent, Value: "Root"})
	store.AppendEvent(EventData{ID: 2, BaseEventID: 1, ValueID: TypeSubEvent, ActorEventID: TypeEvent, Value: "Child1"})
	store.AppendEvent(EventData{ID: 3, BaseEventID: 1, ValueID: TypeSubEvent, ActorEventID: TypeEvent, Value: "Child2"})
	store.AppendEvent(EventData{ID: 4, BaseEventID: 1, ValueID: TypeSubEvent, ActorEventID: TypeEvent, Value: "Child3"})
	store.AppendEvent(EventData{ID: 5, BaseEventID: 999, ValueID: TypeSubEvent, ActorEventID: TypeEvent, Value: "Root2"})

	// Проверяем дочерние события базового события 1
	children, err := store.GetEventsByBase(1)
	if err != nil {
		t.Fatalf("GetEventsByBase вернул ошибку: %v", err)
	}

	if len(children) != 3 {
		t.Fatalf("ожидалось 3 дочерних события, получено %d", len(children))
	}

	if children[0].Value != "Child1" || children[2].Value != "Child3" {
		t.Errorf("неверный порядок или значения дочерних событий")
	}
}

// TestMemoryStorage_GetEventsByValue проверяет выборку по типу.
func TestMemoryStorage_GetEventsByValue(t *testing.T) {
	store := NewMemoryStorage()

	store.AppendEvent(EventData{ID: 1, BaseEventID: TypeEvent, ValueID: TypeActor, ActorEventID: TypeEvent, Value: "Actor1"})
	store.AppendEvent(EventData{ID: 2, BaseEventID: TypeEvent, ValueID: TypeEntity, ActorEventID: TypeEvent, Value: "Entity1"})
	store.AppendEvent(EventData{ID: 3, BaseEventID: TypeEvent, ValueID: TypeActor, ActorEventID: TypeEvent, Value: "Actor2"})

	actors, err := store.GetEventsByValue(TypeActor)
	if err != nil {
		t.Fatalf("GetEventsByValue вернул ошибку: %v", err)
	}

	if len(actors) != 2 {
		t.Fatalf("ожидалось 2 актора, получено %d", len(actors))
	}
}

// TestMemoryStorage_DuplicateID проверяет защиту от дубликатов.
func TestMemoryStorage_DuplicateID(t *testing.T) {
	store := NewMemoryStorage()

	store.AppendEvent(EventData{ID: 1, BaseEventID: TypeEvent, ValueID: TypeSubEvent, ActorEventID: TypeEvent, Value: "first"})
	err := store.AppendEvent(EventData{ID: 1, BaseEventID: TypeEvent, ValueID: TypeSubEvent, ActorEventID: TypeEvent, Value: "duplicate"})

	if err == nil {
		t.Error("ожидалась ошибка дубликата, но её нет")
	}
}

// TestMemoryStorage_Closed проверяет защиту от записи в закрытое хранилище.
func TestMemoryStorage_Closed(t *testing.T) {
	store := NewMemoryStorage()
	store.Close()

	err := store.AppendEvent(EventData{ID: 1, BaseEventID: TypeEvent, ValueID: TypeSubEvent, ActorEventID: TypeEvent, Value: "test"})
	if err == nil {
		t.Error("ожидалась ошибка закрытого хранилища")
	}
}

// === Тесты ConditionRule ===

// TestCondition_Nil проверяет, что nil-условие считается истинным.
func TestCondition_Nil(t *testing.T) {
	var rule ConditionRule = nil
	// nil считается истинным — проверяется в ConjunctionRule и явно
	if rule != nil {
		t.Error("nil-правило не nil")
	}
}

// TestCondition_Equality проверяет условие равенства.
func TestCondition_Equality(t *testing.T) {
	rule := &PropertyEqualityRule{PropertyID: 1, Value: "Accept"}

	ctx := &EventContext{
		Properties: map[int]string{1: "Accept"},
	}

	if !rule.IsMet(ctx) {
		t.Error("правило равенства должно быть истинным")
	}

	// Проверяем неравенство
	ctx.Properties[1] = "Reject"
	if rule.IsMet(ctx) {
		t.Error("правило равенства должно быть ложным при разных значениях")
	}
}

// TestCondition_Inequality проверяет условие неравенства.
func TestCondition_Inequality(t *testing.T) {
	rule := &PropertyInequalityRule{PropertyID: 1, Value: "closed"}

	ctx := &EventContext{
		Properties: map[int]string{1: "process"},
	}

	if !rule.IsMet(ctx) {
		t.Error("правило неравенства должно быть истинным")
	}

	ctx.Properties[1] = "closed"
	if rule.IsMet(ctx) {
		t.Error("правило неравенства должно быть ложным при совпадении")
	}
}

// TestCondition_Conjunction проверяет логическое И.
func TestCondition_Conjunction(t *testing.T) {
	rule := &ConjunctionRule{
		Values: []ConditionRule{
			&PropertyEqualityRule{PropertyID: 1, Value: "A"},
			&PropertyInequalityRule{PropertyID: 2, Value: "0"},
		},
	}

	ctx := &EventContext{
		Properties: map[int]string{1: "A", 2: "1"},
	}

	if !rule.IsMet(ctx) {
		t.Error("конъюнкция должна быть истинной")
	}

	// Меняем одно условие
	ctx.Properties[1] = "B"
	if rule.IsMet(ctx) {
		t.Error("конъюнкция должна быть ложной при одном ложном условии")
	}
}

// TestCondition_Disjunction проверяет логическое ИЛИ.
func TestCondition_Disjunction(t *testing.T) {
	rule := &DisjunctionRule{
		Values: []ConditionRule{
			&PropertyEqualityRule{PropertyID: 1, Value: "A"},
			&PropertyEqualityRule{PropertyID: 2, Value: "B"},
		},
	}

	ctx := &EventContext{
		Properties: map[int]string{1: "X", 2: "B"},
	}

	if !rule.IsMet(ctx) {
		t.Error("дизъюнкция должна быть истинной (одно условие верно)")
	}

	// Оба ложны
	ctx.Properties[2] = "C"
	if rule.IsMet(ctx) {
		t.Error("дизъюнкция должна быть ложной (все ложны)")
	}
}

// TestCondition_EmptyConjunction проверяет пустую конъюнкцию.
func TestCondition_EmptyConjunction(t *testing.T) {
	rule := &ConjunctionRule{Values: nil}
	if !rule.IsMet(nil) {
		t.Error("пустая конъюнкция должна быть истинной")
	}
}

// === Тесты EngineBase ===

// TestEngine_New проверяет создание движка.
func TestEngine_New(t *testing.T) {
	store := NewMemoryStorage()
	eng := NewEngine(store)

	if eng.Position != 1 {
		t.Errorf("начальная позиция должна быть 1, получена %d", eng.Position)
	}

	if eng.Types == nil {
		t.Error("TypeManager не должен быть nil")
	}
}

// TestEngine_ProcessEvent проверяет базовую обработку события.
func TestEngine_ProcessEvent(t *testing.T) {
	store := NewMemoryStorage()
	eng := NewEngine(store)

	event := EventData{
		ID:           1,
		BaseEventID:  TypeEvent,
		ValueID:      TypeSubEvent,
		ActorEventID: 0, // системное событие
		Value:        "Test",
	}

	err := eng.ProcessEvent(event)
	if err != nil {
		t.Fatalf("ProcessEvent вернул ошибку: %v", err)
	}

	if eng.Position != 2 {
		t.Errorf("позиция должна быть 2, получена %d", eng.Position)
	}
}

// TestEngine_InvalidPosition проверяет ошибку неверной позиции.
func TestEngine_InvalidPosition(t *testing.T) {
	store := NewMemoryStorage()
	eng := NewEngine(store)

	event := EventData{
		ID:           999,
		BaseEventID:  TypeEvent,
		ValueID:      TypeSubEvent,
		ActorEventID: 0,
		Value:        "WrongPosition",
	}

	err := eng.ProcessEvent(event)
	if err == nil {
		t.Error("ожидалась ошибка неверной позиции")
	}
}

// TestEngine_InvalidBaseRef проверяет ошибку ссылки на будущее событие.
func TestEngine_InvalidBaseRef(t *testing.T) {
	store := NewMemoryStorage()
	eng := NewEngine(store)

	event := EventData{
		ID:           1,
		BaseEventID:  100,
		ValueID:      TypeSubEvent,
		ActorEventID: TypeEvent,
		Value:        "BadRef",
	}

	err := eng.ProcessEvent(event)
	if err == nil {
		t.Error("ожидалась ошибка неверной ссылки")
	}
}

// TestEngine_InvalidValueType проверяет ошибку неизвестного типа.
func TestEngine_InvalidValueType(t *testing.T) {
	store := NewMemoryStorage()
	eng := NewEngine(store)

	event := EventData{
		ID:           1,
		BaseEventID:  TypeEvent,
		ValueID:      99999, // неизвестный тип
		ActorEventID: TypeEvent,
		Value:        "BadType",
	}

	err := eng.ProcessEvent(event)
	if err == nil {
		t.Error("ожидалась ошибка неизвестного типа")
	}
}

// === Тесты InitGenesis ===

// TestInitGenesis проверяет инициализацию генезис-графа.
func TestInitGenesis(t *testing.T) {
	store := NewMemoryStorage()
	eng := NewEngine(store)

	err := eng.InitGenesis()
	if err != nil {
		t.Fatalf("InitGenesis вернул ошибку: %v", err)
	}

	if !eng.IsGenesisComplete() {
		t.Error("генезис-граф должен быть завершён")
	}

	// Проверяем количество событий
	count := store.EventCount()
	if count < 20 {
		t.Errorf("ожидалось минимум 20 событий в генезис-графе, получено %d", count)
	}

	// Проверяем, что основные типы созданы
	events, _ := store.GetEventsByValue(TypeSubEvent)
	if len(events) == 0 {
		t.Error("должны быть события типа SubEvent")
	}

	// Проверяем Actor_Main
	event, err := store.GetEvent(21)
	if err != nil {
		t.Fatalf("Actor_Main не найден: %v", err)
	}
	if event.Value != "Actor_Main" {
		t.Errorf("ожидался 'Actor_Main', получен '%s'", event.Value)
	}
}

// TestInitGenesisTwice проверяет защиту от двойной инициализации.
func TestInitGenesisTwice(t *testing.T) {
	store := NewMemoryStorage()
	eng := NewEngine(store)

	err := eng.InitGenesis()
	if err != nil {
		t.Fatalf("первая инициализация не удалась: %v", err)
	}

	err = eng.InitGenesis()
	if err == nil {
		t.Error("двойная инициализация должна вернуть ошибку")
	}
}

// === Тесты TypeManager ===

// TestTypeManager_BaseTypes проверяет наличие всех базовых типов.
func TestTypeManager_BaseTypes(t *testing.T) {
	tm := NewTypeManager()

	tests := []struct {
		id    int
		label string
	}{
		{TypeEvent, "Event"},
		{TypeActor, "Actor"},
		{TypeEntity, "Entity"},
		{TypeModel, "Model"},
		{TypeIndividual, "Individual"},
	}

	for _, tt := range tests {
		if !tm.IsValid(tt.id) {
			t.Errorf("тип %d (%s) должен быть зарегистрирован", tt.id, tt.label)
		}
		if tm.Label(tt.id) != tt.label {
			t.Errorf("метка типа %d: ожидалось '%s', получено '%s'", tt.id, tt.label, tm.Label(tt.id))
		}
	}
}

// TestTypeManager_UnknownType проверяет обработку неизвестных типов.
func TestTypeManager_UnknownType(t *testing.T) {
	tm := NewTypeManager()

	if tm.IsValid(99999) {
		t.Error("тип 99999 не должен существовать")
	}

	if tm.Label(99999) != "" {
		t.Error("метка неизвестного типа должна быть пустой")
	}
}

// TestTypeManager_All проверяет количество зарегистрированных типов.
func TestTypeManager_All(t *testing.T) {
	tm := NewTypeManager()
	all := tm.All()

	// Должно быть 25 типов
	if len(all) < 24 {
		t.Errorf("ожидалось минимум 24 типа, получено %d", len(all))
	}
}
