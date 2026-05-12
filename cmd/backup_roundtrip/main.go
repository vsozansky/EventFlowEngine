// Программа для тестирования загрузки/выгрузки backup.xml
// Запуск: go run cmd/backup_roundtrip/main.go
package main

import (
	"fmt"
	"os"

	"github.com/eventflow/engine/engine"
)

func main() {
	// Читаем исходный дамп
	data, err := os.ReadFile("backup.xml")
	if err != nil {
		panic(fmt.Errorf("ошибка чтения backup.xml: %w", err))
	}

	// Загружаем в хранилище
	store := engine.NewMemoryStorage()
	count, err := engine.ImportBackupFromXML(data, store)
	if err != nil {
		panic(fmt.Errorf("ошибка загрузки: %w", err))
	}
	fmt.Printf("Загружено событий: %d\n", count)

	// Выгружаем обратно
	result, err := engine.ExportBackupToXML(store, nil)
	if err != nil {
		panic(fmt.Errorf("ошибка выгрузки: %w", err))
	}

	// Сохраняем результат
	if err := os.WriteFile("backup-result.xml", result, 0644); err != nil {
		panic(fmt.Errorf("ошибка сохранения: %w", err))
	}
	fmt.Printf("Сохранено backup-result.xml (%d байт)\n", len(result))
	fmt.Println("✅ Roundtrip завершён успешно!")
}
