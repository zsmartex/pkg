package vault

import (
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/vault/api"
)

type TransitService struct {
	vault_service          *VaultService
	vault_application_name string
}

func NewTransitService(vault_service *VaultService, vault_application_name string) *TransitService {
	return &TransitService{
		vault_service:          vault_service,
		vault_application_name: vault_application_name,
	}
}

func (s *TransitService) Encrypt(key, value string) (*api.Secret, error) {
	valueString := base64.StdEncoding.EncodeToString([]byte(value))

	return s.vault_service.Write(fmt.Sprintf("transit/encrypt/%s_%s", s.vault_application_name, key), map[string]interface{}{
		"plaintext": valueString,
	})
}

func (s *TransitService) Decrypt(key, ciphertext string) (*api.Secret, error) {
	return s.vault_service.Write(fmt.Sprintf("transit/encrypt/%s_%s", s.vault_application_name, key), map[string]interface{}{
		"ciphertext": ciphertext,
	})
}

func (s *TransitService) Delete(key string) (*api.Secret, error) {
	return s.vault_service.Delete(fmt.Sprintf("transit/encrypt/%s_%s", s.vault_application_name, key))
}
