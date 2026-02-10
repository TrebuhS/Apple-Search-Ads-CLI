package services

import (
	"github.com/trebuhs/asa-cli/internal/api"
	"github.com/trebuhs/asa-cli/internal/models"
)

type ACLService struct {
	Client *api.Client
}

func NewACLService(client *api.Client) *ACLService {
	return &ACLService{Client: client}
}

func (s *ACLService) GetACLs() ([]models.UserACL, error) {
	var acls []models.UserACL
	_, err := s.Client.Get("/acls", &acls)
	if err != nil {
		return nil, err
	}
	return acls, nil
}
