package main

import (
	"fmt"
	"log"
	"os"

	"github.com/AcidicAcidity/t-mess/internal/crypto"
	"github.com/AcidicAcidity/t-mess/internal/storage"
	"github.com/AcidicAcidity/t-mess/internal/tui"
)

func main() {
	// 1. Получаем путь к данным
	dataDir, err := storage.GetDataDir()
	if err != nil {
		log.Fatalf("Failed to get data dir: %v", err)
	}

	// 2. Создаём папку, если нет
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		log.Fatalf("Failed to create data dir: %v", err)
	}

	// 3. Загружаем или создаём идентичность
	identityPath := storage.GetIdentityPath()
	identity, err := crypto.LoadOrCreateIdentity(identityPath)
	if err != nil {
		log.Fatalf("Failed to load/create identity: %v", err)
	}
	fmt.Printf("Node ID: %s\n", identity.PeerID)

	// 4. Инициализируем базу данных
	dbPath := storage.GetDBPath()
	db, err := storage.NewDatabase(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// 5. Запускаем TUI
	app := tui.NewApp(identity, db)
	if err := app.Run(); err != nil {
		log.Fatalf("TUI error: %v", err)
	}
}
