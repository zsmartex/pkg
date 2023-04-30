package vault

import (
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/vault/api"
)

type TransitService struct {
	vaultService         *VaultService
	vaultApplicationName string
}

func NewTransitService(vault_service *VaultService, vault_application_name string) *TransitService {
	return &TransitService{
		vaultService:         vault_service,
		vaultApplicationName: vault_application_name,
	}
}

func (s *TransitService) Encrypt(key, value string) (*api.Secret, error) {
	valueString := base64.StdEncoding.EncodeToString([]byte(value))

	return s.vaultService.Write(fmt.Sprintf("transit/encrypt/%s_%s", s.vaultApplicationName, key), map[string]interface{}{
		"plaintext": valueString,
	})
}

func (s *TransitService) Decrypt(key, ciphertext string) (*api.Secret, error) {
	return s.vaultService.Write(fmt.Sprintf("transit/decrypt/%s_%s", s.vaultApplicationName, key), map[string]interface{}{
		"ciphertext": ciphertext,
	})
}

func (s *TransitService) Delete(key string) (*api.Secret, error) {
	return s.vaultService.Delete(fmt.Sprintf("transit/encrypt/%s_%s", s.vaultApplicationName, key))
}
