package engine

import "errors"

// Sentinel-ошибки (предопределённые константы ошибок).
// Возвращаются из функций движка при нарушении контрактов.
var (
	// ErrEventNotFound — событие с указанным ID не найдено в графе.
	ErrEventNotFound = errors.New("событие не найдено")

	// ErrInvalidPosition — ID события не совпадает с текущей позицией движка.
	ErrInvalidPosition = errors.New("неверная позиция события")

	// ErrInvalidBaseRef — BaseEventID ссылается на несуществующее событие.
	ErrInvalidBaseRef = errors.New("неверная ссылка на базовое событие")

	// ErrInvalidValueType — тип события (ValueID) не известен TypeManager'у.
	ErrInvalidValueType = errors.New("неизвестный тип события")

	// ErrInvalidActor — актор с указанным ID не найден или не валиден.
	ErrInvalidActor = errors.New("неверный актор")

	// ErrDuplicateEvent — событие с таким ID уже существует.
	ErrDuplicateEvent = errors.New("дубликат события")

	// ErrStorageClosed — хранилище закрыто для операций.
	ErrStorageClosed = errors.New("хранилище закрыто")

	// ErrImmutable — попытка изменить существующее событие (запрещено).
	ErrImmutable = errors.New("события неизменяемы")
)

// Storage — интерфейс хранилища событийного графа.
//
// Хранилище отвечает за персистентность событий и сборку агрегатов
// (сущностей, моделей, индивидов) из событий по запросу.
//
// Все реализации должны соблюдать принципы:
//   - Append-only: никаких Update/Delete
//   - Immutable: события не изменяются после записи
//   - Event Sourcing: полная воспроизводимость состояния
type Storage interface {
	// === Базовые операции с событиями ===

	// AppendEvent добавляет новое событие в граф.
	// Событие должно иметь ID = NextID(). После успешной записи
	// внутренний счётчик увеличивается.
	AppendEvent(event EventData) error

	// NextID возвращает следующий свободный ID для события.
	// ID — монотонно возрастающее целое число.
	NextID() int

	// GetEvent возвращает событие по его ID.
	// Возвращает ErrEventNotFound, если события нет.
	GetEvent(id int) (*EventData, error)

	// === Выборки событий ===

	// GetEventsByBase возвращает все события, зафиксированные на указанном
	// базовом событии (по BaseEventID). Исключает само базовое событие.
	// Порядок — по возрастанию ID.
	GetEventsByBase(baseEventID int) ([]*EventData, error)

	// GetEventsByValue возвращает все события указанного семантического типа
	// (по ValueID). Порядок — по возрастанию ID.
	GetEventsByValue(valueID int) ([]*EventData, error)

	// GetEventsByActor возвращает все события, созданные указанным актором.
	GetEventsByActor(actorEventID int) ([]*EventData, error)

	// === Агрегаты (собираются из событий) ===

	// GetEntity собирает сущность из событий графа по её ID.
	// Сущность — это событие с ValueID = TypeEntity.
	GetEntity(entityID int) (*Entity, error)

	// GetEntities возвращает все сущности в графе.
	GetEntities() ([]*Entity, error)

	// GetModel собирает модель из событий графа по её ID.
	// Модель — это событие с ValueID = TypeModel.
	GetModel(modelID int) (*Model, error)

	// GetIndividual собирает индивида из событий графа по его ID.
	// Индивид — это событие с ValueID = TypeIndividual.
	GetIndividual(individualID int) (*Individual, error)

	// GetIndividuals возвращает всех индивидов указанного типа.
	GetIndividualsByValue(valueID int) ([]*Individual, error)

	// === Провайдеры свойств ===

	// GetPropertyProvider возвращает провайдер свойств для указанного
	// провайдера (модели или контейнера). Провайдер содержит атрибуты,
	// отношения и события, привязанные к нему с ограничениями.
	GetPropertyProvider(providerID int) (*PropertyProvider, error)

	// GetPropertyProviderAttribute возвращает привязанный атрибут
	// из указанного провайдера по ID атрибута.
	GetPropertyProviderAttribute(providerID, attrID int) (*AttachedAttribute, error)
}
