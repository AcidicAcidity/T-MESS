package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AcidicAcidity/t-mess/internal/crypto"
	"github.com/AcidicAcidity/t-mess/internal/p2p"
	"github.com/AcidicAcidity/t-mess/internal/storage"
	"github.com/AcidicAcidity/t-mess/internal/tui"
)

func main() {
	// Настройка логгирования в файл
	dataDir, err := storage.GetDataDir()
	if err != nil {
		log.Fatalf("Failed to get data dir: %v", err)
	}
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		log.Fatalf("Failed to create data dir: %v", err)
	}

	logFile, err := os.OpenFile(dataDir+"/t-mess.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// 1. Загрузка идентичности
	identityPath := storage.GetIdentityPath()
	identity, err := crypto.LoadOrCreateIdentity(identityPath)
	if err != nil {
		log.Fatalf("Failed to load/create identity: %v", err)
	}
	log.Printf("Node ID: %s", identity.PeerID)

	// 2. База данных
	dbPath := storage.GetDBPath()
	db, err := storage.NewDatabase(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()
	log.Printf("Database initialized at %s", dbPath)

	// 3. Запуск P2P узла
	ctx := context.Background()
	p2pNode, err := p2p.NewNode(ctx, identity.PrivateKey)
	if err != nil {
		log.Fatalf("Failed to start P2P node: %v", err)
	}
	defer p2pNode.Close()

	log.Printf("P2P node started: %s", p2pNode.ID())
	log.Printf("Listening on: %v", p2pNode.Addrs())

	// 4. Обработка сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down...")
		p2pNode.Close()
		os.Exit(0)
	}()

	// 5. Запуск TUI с передачей P2P узла
	app := tui.NewApp(identity, db, p2pNode)
	if err := app.Run(); err != nil {
		log.Fatalf("TUI error: %v", err)
	}
}
