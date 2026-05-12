package engine

import "fmt"

// EngineBase — ядро событийного движка.
//
// EngineBase управляет позицией в графе (Position), обрабатывает новые
// события через ProcessEvent и предоставляет высокоуровневые транзакции
// для работы с онтологией.
//
// Жизненный цикл:
//  1. NewEngine(storage) — создание
//  2. InitGenesis() — инициализация генезис-графа
//  3. ProcessEvent() / транзакции — работа
type EngineBase struct {
	// Position — текущая позиция в графе (следующий ID события).
	// Монотонно возрастает с каждым записанным событием.
	Position int

	// Storage — хранилище событийного графа.
	Storage Storage

	// Types — менеджер семантических типов.
	Types *TypeManager
}

// NewEngine создаёт новый экземпляр движка с указанным хранилищем.
//
// Параметры:
//   - storage — реализация Storage (например, MemoryStorage)
//
// Возвращает готовый к инициализации движок. Position начинается с 1.
func NewEngine(storage Storage) *EngineBase {
	return &EngineBase{
		Position: 1,
		Storage:  storage,
		Types:    NewTypeManager(),
	}
}

// ProcessEvent обрабатывает и записывает событие в граф.
//
// Процесс обработки:
//  1. Проверка ID: должен совпадать с Position (защита от дубликатов)
//  2. Проверка BaseEventID: должен быть меньше ID (ссылка в прошлое)
//  3. Проверка ValueID: должен быть известен TypeManager'у
//  4. Проверка ActorEventID: актор должен существовать в графе
//  5. Установка времени (UTC), если не задано
//  6. Запись в хранилище (AppendEvent)
//  7. Увеличение Position
func (e *EngineBase) ProcessEvent(event EventData) error {
	// Шаг 1: проверка ID
	if event.ID != e.Position {
		return fmt.Errorf("%w: ожидался ID %d, получен %d",
			ErrInvalidPosition, e.Position, event.ID)
	}

	// Шаг 2: проверка BaseEventID
	// Для корневых событий BaseEventID может равняться ID
	// (самоссылка на собственный тип). Ссылка в будущее (BaseEventID > ID)
	// всегда является ошибкой.
	if event.BaseEventID > event.ID {
		return fmt.Errorf("%w: BaseEventID %d > ID %d",
			ErrInvalidBaseRef, event.BaseEventID, event.ID)
	}

	// Шаг 3: проверка ValueID
	if !e.Types.IsValid(event.ValueID) {
		return fmt.Errorf("%w: ValueID %d не зарегистрирован",
			ErrInvalidValueType, event.ValueID)
	}

	// Шаг 4: проверка ActorEventID
	// ActorEventID = 0 означает системное событие (без конкретного актора).
	// Для не-системных событий актор должен существовать в графе.
	if event.ID > 1 && event.ActorEventID > 0 {
		// Проверяем, что событие актора уже существует
		_, err := e.Storage.GetEvent(event.ActorEventID)
		if err != nil {
			return fmt.Errorf("%w: актор с ID %d не найден: %w",
				ErrInvalidActor, event.ActorEventID, err)
		}
	}

	// Шаг 5-7: запись, увеличение счётчика и возврат
	err := e.Storage.AppendEvent(event)
	if err != nil {
		return fmt.Errorf("ошибка записи события: %w", err)
	}

	e.Position++
	return nil
}

// MustProcessEvent — упрощённый вариант ProcessEvent.
// Автоматически устанавливает ID = Position и вызывает ProcessEvent.
// Паникует при ошибке — используйте только в тестах и при инициализации.
func (e *EngineBase) MustProcessEvent(baseEventID, valueID int, conditions ConditionRule, actorEventID int, value string) int {
	pos := e.Position
	event := EventData{
		ID:           pos,
		BaseEventID:  baseEventID,
		ValueID:      valueID,
		Conditions:   conditions,
		ActorEventID: actorEventID,
		Value:        value,
	}
	if err := e.ProcessEvent(event); err != nil {
		panic(fmt.Sprintf("MustProcessEvent: %v", err))
	}
	return pos
}

// GetEvent возвращает событие из графа по его ID.
func (e *EngineBase) GetEvent(id int) (*EventData, error) {
	return e.Storage.GetEvent(id)
}

// GetEntity собирает сущность из событий графа.
func (e *EngineBase) GetEntity(entityID int) (*Entity, error) {
	return e.Storage.GetEntity(entityID)
}

// GetEntities возвращает все сущности в графе.
func (e *EngineBase) GetEntities() ([]*Entity, error) {
	return e.Storage.GetEntities()
}

// GetModel собирает модель из событий графа.
func (e *EngineBase) GetModel(modelID int) (*Model, error) {
	return e.Storage.GetModel(modelID)
}

// GetIndividual возвращает индивида из графа.
func (e *EngineBase) GetIndividual(individualID int) (*Individual, error) {
	return e.Storage.GetIndividual(individualID)
}
