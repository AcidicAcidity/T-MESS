package main

import (
	"context"
	"fmt"
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
	// 1. Путь к данным
	dataDir, err := storage.GetDataDir()
	if err != nil {
		log.Fatalf("Failed to get data dir: %v", err)
	}
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		log.Fatalf("Failed to create data dir: %v", err)
	}

	// 2. Загрузка идентичности
	identityPath := storage.GetIdentityPath()
	identity, err := crypto.LoadOrCreateIdentity(identityPath)
	if err != nil {
		log.Fatalf("Failed to load/create identity: %v", err)
	}
	fmt.Printf("Node ID: %s\n", identity.PeerID)

	// 3. База данных
	dbPath := storage.GetDBPath()
	db, err := storage.NewDatabase(dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// 4. Запуск P2P узла
	ctx := context.Background()
	p2pNode, err := p2p.NewNode(ctx, identity.PrivateKey)
	if err != nil {
		log.Fatalf("Failed to start P2P node: %v", err)
	}
	defer p2pNode.Close()

	fmt.Printf("P2P node started: %s\n", p2pNode.ID())
	fmt.Printf("Listening on: %v\n", p2pNode.Addrs())

	// 5. Обработка сигналов завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nShutting down...")
		p2pNode.Close()
		os.Exit(0)
	}()

	// 6. Запуск TUI (пока без интеграции с P2P)
	app := tui.NewApp(identity, db)
	if err := app.Run(); err != nil {
		log.Fatalf("TUI error: %v", err)
	}
}
