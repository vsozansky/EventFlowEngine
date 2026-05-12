package engine

import (
	"encoding/xml"
	"fmt"
	"time"
)

// BackupData — корневой элемент XML-дампов графа.
// Соответствует формату Parallax/AuroraCore.
type BackupData struct {
	XMLName xml.Name      `xml:"BackupData"`
	XSI     string        `xml:"xmlns:xsi,attr"`
	XSD     string        `xml:"xmlns:xsd,attr"`
	Events  []BackupEvent `xml:"Events>BackupEvent"`
}

// BackupEvent — событие в формате XML-дампа.
type BackupEvent struct {
	ID           int             `xml:"ID"`
	BaseEventID  int             `xml:"BaseEventID"`
	ValueID      int             `xml:"ValueID"`
	Conditions   BackupCondition `xml:"Conditions"`
	ActorEventID int             `xml:"ActorEventID"`
	Value        string          `xml:"Value"`
	Date         string          `xml:"Date"`
}

// BackupCondition — условие в XML-формате.
type BackupCondition struct {
	XsiType string `xml:"xsi:type,attr"`
	EventID int    `xml:"EventID"`
}

// ImportBackupFromXML загружает события из XML-дампа в хранилище.
//
// Параметры:
//   - data: содержимое XML-файла (backup.xml)
//   - store: хранилище, в которое загружаются события
//
// Возвращает количество загруженных событий.
//
// Особенности:
//   - Сохраняет оригинальные ID событий (включая ID=0)
//   - Использует ImportEvents для прямой загрузки (без проверки ID)
//   - Используется для roundtrip-тестирования
func ImportBackupFromXML(data []byte, store *MemoryStorage) (int, error) {
	var backup BackupData
	backup.XSI = "http://www.w3.org/2001/XMLSchema-instance"
	backup.XSD = "http://www.w3.org/2001/XMLSchema"

	if err := xml.Unmarshal(data, &backup); err != nil {
		return 0, fmt.Errorf("ошибка парсинга XML: %w", err)
	}

	if len(backup.Events) == 0 {
		return 0, fmt.Errorf("XML не содержит событий")
	}

	// Преобразуем BackupEvent → EventData
	events := make([]EventData, 0, len(backup.Events))
	for _, be := range backup.Events {
		// Парсим дату
		var eventTime time.Time
		if be.Date != "" {
			var err error
			eventTime, err = time.Parse(time.RFC3339Nano, be.Date)
			if err != nil {
				eventTime, err = time.Parse("2006-01-02T15:04:05.000Z", be.Date)
				if err != nil {
					eventTime = time.Now().UTC()
				}
			}
		}

		event := EventData{
			ID:           be.ID,
			BaseEventID:  be.BaseEventID,
			ValueID:      be.ValueID,
			Conditions:   &EventConditionRule{EventID: be.Conditions.EventID},
			ActorEventID: be.ActorEventID,
			Value:        be.Value,
			Date:         eventTime,
		}
		events = append(events, event)
	}

	// Загружаем через ImportEvents (без проверки ID, с сохранением ID=0)
	if err := store.ImportEvents(events); err != nil {
		return 0, fmt.Errorf("ошибка загрузки событий: %w", err)
	}

	return len(events), nil
}

// ExportBackupToXML выгружает все события из хранилища в XML-формат.
//
// Параметры:
//   - store: хранилище с событиями
//   - events: список событий для выгрузки (если nil, выгружаются все события из хранилища)
//
// Возвращает XML-дамп в виде байтового среза.
func ExportBackupToXML(store *MemoryStorage, events []*EventData) ([]byte, error) {
	if events == nil {
		events = make([]*EventData, 0)
		for i := 0; ; i++ {
			e, err := store.GetEvent(i)
			if err != nil {
				break
			}
			events = append(events, e)
		}
	}

	if len(events) == 0 {
		return nil, fmt.Errorf("нет событий для выгрузки")
	}

	backup := BackupData{
		XSI: "http://www.w3.org/2001/XMLSchema-instance",
		XSD: "http://www.w3.org/2001/XMLSchema",
	}

	for _, e := range events {
		dateStr := e.Date.Format(time.RFC3339Nano)

		// Получаем EventID из условия
		condEventID := 0
		if ecr, ok := e.Conditions.(*EventConditionRule); ok {
			condEventID = ecr.EventID
		}

		be := BackupEvent{
			ID:           e.ID,
			BaseEventID:  e.BaseEventID,
			ValueID:      e.ValueID,
			Conditions: BackupCondition{
				XsiType: "EventConditionRule",
				EventID: condEventID,
			},
			ActorEventID: e.ActorEventID,
			Value:        e.Value,
			Date:         dateStr,
		}
		backup.Events = append(backup.Events, be)
	}

	data, err := xml.MarshalIndent(backup, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("ошибка сериализации XML: %w", err)
	}

	// Добавляем XML-заголовок
	result := append([]byte(xml.Header), data...)
	return result, nil
}
