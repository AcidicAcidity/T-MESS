package crypto

import (
	"encoding/hex"
	"encoding/pem"
	"os"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/tyler-smith/go-bip39"
	"golang.org/x/crypto/ed25519"
)

type Identity struct {
	PrivateKey crypto.PrivKey
	PublicKey  crypto.PubKey
	PeerID     string
	Mnemonic   string
}

// GenerateNewIdentity создаёт новую идентичность из мнемонической фразы
func GenerateNewIdentity() (*Identity, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return nil, err
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return nil, err
	}

	seed := bip39.NewSeed(mnemonic, "")
	privateKey := ed25519.NewKeyFromSeed(seed[:32])

	privKey, err := crypto.UnmarshalEd25519PrivateKey(privateKey)
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
		Mnemonic:   mnemonic,
	}, nil
}

// LoadOrCreateIdentity загружает существующий ключ или создаёт новый
func LoadOrCreateIdentity(path string) (*Identity, error) {
	data, err := os.ReadFile(path)
	if err == nil {
		// Пытаемся загрузить существующий ключ
		block, _ := pem.Decode(data)
		if block == nil {
			// Если не удалось распарсить PEM, удаляем повреждённый файл и создаём новый
			os.Remove(path)
			return createAndSaveIdentity(path)
		}

		privKey, err := crypto.UnmarshalPrivateKey(block.Bytes)
		if err != nil {
			// Если не удалось распарсить ключ, удаляем повреждённый файл и создаём новый
			os.Remove(path)
			return createAndSaveIdentity(path)
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

	// Файла нет — создаём новый
	return createAndSaveIdentity(path)
}

// createAndSaveIdentity создаёт новую идентичность и сохраняет
func createAndSaveIdentity(path string) (*Identity, error) {
	identity, err := GenerateNewIdentity()
	if err != nil {
		return nil, err
	}

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
