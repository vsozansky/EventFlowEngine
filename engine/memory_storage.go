package engine

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// MemoryStorage — реализация хранилища в оперативной памяти.
//
// Все события хранятся в map[int]*EventData. Потокобезопасна
// через sync.RWMutex. Подходит для тестирования, прототипирования
// и embedded-сценариев без персистентности.
//
// Не подходит для:
//   - Многогигабайтных графов
//   - Распределённых систем
//   - Сценариев, требующих долговременного хранения
type MemoryStorage struct {
	mu     sync.RWMutex
	events map[int]*EventData
	nextID int
	closed bool
}

// NewMemoryStorage создаёт новое in-memory хранилище.
// Счётчик ID начинается с 1.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		events: make(map[int]*EventData),
		nextID: 1,
	}
}

// Close закрывает хранилище. После закрытия все операции записи
// будут возвращать ErrStorageClosed.
func (s *MemoryStorage) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.closed = true
}

// NextID возвращает следующий свободный ID для события.
// ID монотонно возрастает, начиная с 1.
func (s *MemoryStorage) NextID() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.nextID
}

// AppendEvent добавляет новое событие в граф.
//
// Проверки:
//   - Хранилище не закрыто
//   - ID события совпадает с nextID (защита от дубликатов)
//   - Событие с таким ID ещё не существует
func (s *MemoryStorage) AppendEvent(event EventData) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return ErrStorageClosed
	}

	if event.ID != s.nextID {
		return fmt.Errorf("%w: ожидался ID %d, получен %d",
			ErrInvalidPosition, s.nextID, event.ID)
	}

	if _, exists := s.events[event.ID]; exists {
		return fmt.Errorf("%w: событие с ID %d уже существует",
			ErrDuplicateEvent, event.ID)
	}

	// Устанавливаем время, если не задано
	if event.Date.IsZero() {
		event.Date = time.Now().UTC()
	}

	s.events[event.ID] = &event
	s.nextID++
	return nil
}

// GetEvent возвращает событие по его ID.
// Возвращает ErrEventNotFound, если события нет.
func (s *MemoryStorage) GetEvent(id int) (*EventData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, ok := s.events[id]
	if !ok {
		return nil, fmt.Errorf("%w: ID %d", ErrEventNotFound, id)
	}
	return event, nil
}

// GetEventsByBase возвращает все события, зафиксированные на указанном
// базовом событии. Исключает само базовое событие.
// Возвращает пустой срез, если дочерних событий нет.
func (s *MemoryStorage) GetEventsByBase(baseEventID int) ([]*EventData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*EventData
	for _, e := range s.events {
		if e.BaseEventID == baseEventID && e.ID != baseEventID {
			result = append(result, e)
		}
	}

	// Сортируем по ID для детерминированного порядка
	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result, nil
}

// GetEventsByValue возвращает все события указанного семантического типа.
func (s *MemoryStorage) GetEventsByValue(valueID int) ([]*EventData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*EventData
	for _, e := range s.events {
		if e.ValueID == valueID {
			result = append(result, e)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result, nil
}

// GetEventsByActor возвращает все события, созданные указанным актором.
func (s *MemoryStorage) GetEventsByActor(actorEventID int) ([]*EventData, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []*EventData
	for _, e := range s.events {
		if e.ActorEventID == actorEventID {
			result = append(result, e)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})
	return result, nil
}

// ImportEvents загружает события напрямую, минуя проверку ID.
// Используется для восстановления графа из дампа, где ID событий
// уже предопределены. Не увеличивает nextID — хранилище переходит
// в режим «известных ID».
//
// Внимание: после ImportEvents nextID устанавливается на max(загруженные ID)+1.
func (s *MemoryStorage) ImportEvents(events []EventData) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		return ErrStorageClosed
	}

	maxID := 0
	for _, event := range events {
		if _, exists := s.events[event.ID]; exists {
			return fmt.Errorf("%w: событие с ID %d уже существует", ErrDuplicateEvent, event.ID)
		}
		if event.Date.IsZero() {
			event.Date = time.Now().UTC()
		}
		s.events[event.ID] = &event
		if event.ID > maxID {
			maxID = event.ID
		}
	}

	s.nextID = maxID + 1
	return nil
}

// EventCount возвращает общее количество событий в графе.
func (s *MemoryStorage) EventCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.events)
}

// GetEntities возвращает все сущности в графе.
// Заглушка: будет реализована в следующей фазе.
func (s *MemoryStorage) GetEntities() ([]*Entity, error) {
	return nil, nil
}

// GetEntity собирает сущность из событий графа.
// Заглушка: будет реализована в следующей фазе.
func (s *MemoryStorage) GetEntity(entityID int) (*Entity, error) {
	return nil, nil
}

// GetModel собирает модель из событий графа.
// Заглушка: будет реализована в следующей фазе.
func (s *MemoryStorage) GetModel(modelID int) (*Model, error) {
	return nil, nil
}

// GetIndividual собирает индивида из событий графа.
// Заглушка: будет реализована в следующей фазе.
func (s *MemoryStorage) GetIndividual(individualID int) (*Individual, error) {
	return nil, nil
}

// GetIndividualsByValue возвращает индивидов указанного типа.
// Заглушка: будет реализована в следующей фазе.
func (s *MemoryStorage) GetIndividualsByValue(valueID int) ([]*Individual, error) {
	return nil, nil
}

// GetPropertyProvider возвращает провайдер свойств.
// Заглушка: будет реализована в следующей фазе.
func (s *MemoryStorage) GetPropertyProvider(providerID int) (*PropertyProvider, error) {
	return nil, nil
}

// GetPropertyProviderAttribute возвращает привязанный атрибут.
// Заглушка: будет реализована в следующей фазе.
func (s *MemoryStorage) GetPropertyProviderAttribute(providerID, attrID int) (*AttachedAttribute, error) {
	return nil, nil
}
