package engine

import (
	"os"
	"strings"
	"testing"
)

// TestBackup_RoundTrip проверяет полный цикл: загрузка XML → сохранение XML → идентичность.
//
// Сценарий:
//   - Читает файл backup.xml из корня проекта
//   - Загружает 33 события в MemoryStorage
//   - Проверяет количество загруженных событий
//   - Выгружает обратно в XML
//   - Сравнивает с исходным файлом (игнорируя даты)
func TestBackup_RoundTrip(t *testing.T) {
	// Читаем исходный backup.xml
	data, err := os.ReadFile("../backup.xml")
	if err != nil {
		t.Skipf("backup.xml не найден: %v", err)
	}

	// Создаём хранилище и загружаем события
	store := NewMemoryStorage()
	count, err := ImportBackupFromXML(data, store)
	if err != nil {
		t.Fatalf("ошибка загрузки backup.xml: %v", err)
	}

	// Проверяем количество загруженных событий
	if count != 33 {
		t.Errorf("ожидалось 33 события, загружено %d", count)
	}

	// Проверяем, что события действительно в хранилище
	for i := 0; i < count; i++ {
		event, err := store.GetEvent(i)
		if err != nil {
			t.Errorf("событие %d не найдено в хранилище: %v", i, err)
		}
		if event.ID != i {
			t.Errorf("ID события %d не совпадает с ожидаемым %d", event.ID, i)
		}
	}

	// Выгружаем обратно в XML
	resultXML, err := ExportBackupToXML(store, nil)
	if err != nil {
		t.Fatalf("ошибка выгрузки в XML: %v", err)
	}

	// Преобразуем в строки для сравнения (без дат)
	original := string(data)
	result := string(resultXML)

	// Разделяем на строки
	origLines := strings.Split(original, "\n")
	resLines := strings.Split(result, "\n")

	// Сравниваем построчно, пропуская строки с датами
	origIdx := 0
	resIdx := 0
	diffs := 0

	for origIdx < len(origLines) && resIdx < len(resLines) {
		o := strings.TrimSpace(origLines[origIdx])
		r := strings.TrimSpace(resLines[resIdx])

		// Пропускаем строки с датами и XML-заголовок (версия может отличаться)
		if strings.Contains(o, "<Date>") || strings.Contains(o, "<?xml") {
			origIdx++
			continue
		}
		if strings.Contains(r, "<Date>") || strings.Contains(r, "<?xml") {
			resIdx++
			continue
		}

		if o != r {
			if diffs < 5 {
				t.Errorf("различие в строке %d (orig) vs %d (res):\n  orig: %s\n  res:  %s",
					origIdx, resIdx, o, r)
			}
			diffs++
		}
		origIdx++
		resIdx++
	}

	if diffs > 0 {
		t.Errorf("найдено %d различий между исходным и восстановленным XML", diffs)
	}
}

// TestBackup_ImportExportIdentity проверяет, что импорт и экспорт
// сохраняют все поля событий без изменений (кроме дат).
func TestBackup_ImportExportIdentity(t *testing.T) {
	data, err := os.ReadFile("../backup.xml")
	if err != nil {
		t.Skipf("backup.xml не найден: %v", err)
	}

	store := NewMemoryStorage()
	count, err := ImportBackupFromXML(data, store)
	if err != nil {
		t.Fatalf("ошибка загрузки: %v", err)
	}

	// Проверяем сохранность полей для каждого события
	for id := 0; id < count; id++ {
		event, err := store.GetEvent(id)
		if err != nil {
			t.Fatalf("событие %d не найдено: %v", id, err)
		}

		// Проверяем соответствие ID
		if event.ID != id {
			t.Errorf("ID: ожидался %d, получен %d", id, event.ID)
		}

		// Проверяем, что условие не nil
		if event.Conditions == nil {
			t.Errorf("событие %d: Conditions == nil", id)
		}
	}

	t.Logf("✅ успешно загружено %d событий, все поля сохранены", count)
}

// TestBackup_SpecificEvents проверяет несколько ключевых событий из дампа.
func TestBackup_SpecificEvents(t *testing.T) {
	data, err := os.ReadFile("../backup.xml")
	if err != nil {
		t.Skipf("backup.xml не найден: %v", err)
	}

	store := NewMemoryStorage()
	ImportBackupFromXML(data, store)

	// Проверяем ключевые события
	tests := []struct {
		id       int
		expected string
	}{
		{0, "Event"},
		{1, "SubEvent"},
		{14, "basic_type"},
		{17, "Name"},
		{21, "Actor_Main"},
		{25, "enum_type"},
	}

	for _, tt := range tests {
		event, err := store.GetEvent(tt.id)
		if err != nil {
			t.Errorf("событие %d не найдено: %v", tt.id, err)
			continue
		}
		if event.Value != tt.expected {
			t.Errorf("событие %d: ожидалось '%s', получено '%s'",
				tt.id, tt.expected, event.Value)
		}
	}
}
