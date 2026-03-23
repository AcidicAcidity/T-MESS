package storage

import (
	"os"
	"path/filepath"
)

// GetDataDir возвращает путь к папке с данными приложения
func GetDataDir() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "t-mess"), nil
}

// GetDBPath возвращает путь к файлу базы данных
func GetDBPath() string {
	dataDir, _ := GetDataDir()
	return filepath.Join(dataDir, "t-mess.db")
}

// GetIdentityPath возвращает путь к файлу с приватным ключом
func GetIdentityPath() string {
	dataDir, _ := GetDataDir()
	return filepath.Join(dataDir, "identity.key")
}

// GetConfigPath возвращает путь к файлу конфигурации
func GetConfigPath() string {
	dataDir, _ := GetDataDir()
	return filepath.Join(dataDir, "config.json")
}
