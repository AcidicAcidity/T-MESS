package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/tyler-smith/go-bip39"
)

type Identity struct {
	PrivateKey crypto.PrivKey
	PublicKey  crypto.PubKey
	PeerID     string
	Mnemonic   string
}

// GenerateNewIdentity создаёт новую идентичность
func GenerateNewIdentity() (*Identity, error) {
	// Генерируем энтропию для 12 слов
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return nil, err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}

	seed := bip39.NewSeed(mnemonic, "")

	// Создаём Ed25519 ключ из seed (детерминированно)
	privateKey, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return nil, err
	}

	// TODO: Использовать seed для детерминированной генерации
	// Для прототипа пока случайный ключ

	peerID, err := peer.IDFromPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	return &Identity{
		PrivateKey: privateKey,
		PublicKey:  privateKey.GetPublic(),
		PeerID:     peerID.String(),
		Mnemonic:   mnemonic,
	}, nil
}

// LoadOrCreateIdentity загружает существующий ключ или создаёт новый
func LoadOrCreateIdentity(path string) (*Identity, error) {
	data, err := os.ReadFile(path)
	if err == nil {
		// Загружаем существующий ключ
		block, _ := pem.Decode(data)
		if block == nil {
			return nil, fmt.Errorf("failed to decode PEM block")
		}

		privKey, err := crypto.UnmarshalPrivateKey(block.Bytes)
		if err != nil {
			return nil, err
		}

		peerID, err := peer.IDFromPrivateKey(privKey)
		if err != nil {
			return nil, err
		}

		return &Identity{
			PrivateKey: privKey,
			PublicKey:  privKey.GetPublic(),
			PeerID:     peerID.String(),
		}, nil
	}

	// Создаём новый
	identity, err := GenerateNewIdentity()
	if err != nil {
		return nil, err
	}

	// Сохраняем ключ
	keyBytes, err := identity.PrivateKey.Raw()
	if err != nil {
		return nil, err
	}

	pemBlock := &pem.Block{
		Type:  "ED25519 PRIVATE KEY",
		Bytes: keyBytes,
	}

	if err := os.WriteFile(path, pem.EncodeToMemory(pemBlock), 0600); err != nil {
		return nil, err
	}

	// Сохраняем мнемонику отдельно (опционально)
	mnemonicPath := path + ".mnemonic"
	if err := os.WriteFile(mnemonicPath, []byte(identity.Mnemonic), 0600); err != nil {
		return nil, err
	}

	return identity, nil
}

// Fingerprint возвращает короткий отпечаток ключа
func (i *Identity) Fingerprint() string {
	pubBytes, _ := i.PublicKey.Raw()
	return hex.EncodeToString(pubBytes[:8])
}
