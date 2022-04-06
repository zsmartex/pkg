package services

import (
	"fmt"
)

type TOTPService struct {
	vault_service          *VaultService
	application_name       string
	vault_application_name string
}

func NewTOTPService(vault_service *VaultService, application_name string, vault_application_name string) *TOTPService {
	return &TOTPService{
		vault_service:          vault_service,
		application_name:       application_name,
		vault_application_name: vault_application_name,
	}
}

func (s *TOTPService) Create(uid, email string) (map[string]interface{}, error) {
	if result, err := s.vault_service.Write(s.totp_key(uid), map[string]interface{}{
		"generate":     true,
		"issuer":       s.application_name,
		"account_name": email,
		"qr_size":      100,
	}); err != nil {
		return nil, err
	} else {
		return result.Data, nil
	}
}

func (s *TOTPService) Validate(uid, code string) bool {
	secret, err := s.vault_service.Write(s.totp_code_key(uid), map[string]interface{}{
		"code": code,
	})
	if err != nil {
		return false
	}

	return secret.Data["valid"].(bool)
}

func (s *TOTPService) Delete(uid string) {
	s.vault_service.Delete(s.totp_key(uid))
}

func (s *TOTPService) Exist(uid string) bool {
	secret, err := s.vault_service.Read(s.totp_key(uid))
	if err != nil {
		return false
	}

	return secret != nil
}

func (s *TOTPService) totp_key(uid string) string {
	return fmt.Sprintf("totp/keys/%s_%s", s.vault_application_name, uid)
}

func (s *TOTPService) totp_code_key(uid string) string {
	return fmt.Sprintf("totp/codes/%s_%s", s.vault_application_name, uid)
}
