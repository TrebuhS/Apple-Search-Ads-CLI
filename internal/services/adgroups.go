package services

import (
	"fmt"

	"github.com/trebuhs/asa-cli/internal/api"
	"github.com/trebuhs/asa-cli/internal/models"
)

type AdGroupService struct {
	Client *api.Client
}

func NewAdGroupService(client *api.Client) *AdGroupService {
	return &AdGroupService{Client: client}
}

func (s *AdGroupService) List(campaignID int64, limit, offset int) ([]models.AdGroup, *models.PageDetail, error) {
	path := fmt.Sprintf("/campaigns/%d/adgroups?limit=%d&offset=%d", campaignID, limit, offset)
	var adgroups []models.AdGroup
	page, err := s.Client.Get(path, &adgroups)
	return adgroups, page, err
}

func (s *AdGroupService) Get(campaignID, adGroupID int64) (*models.AdGroup, error) {
	var adgroup models.AdGroup
	_, err := s.Client.Get(fmt.Sprintf("/campaigns/%d/adgroups/%d", campaignID, adGroupID), &adgroup)
	return &adgroup, err
}

func (s *AdGroupService) Find(campaignID int64, selector models.Selector) ([]models.AdGroup, *models.PageDetail, error) {
	var adgroups []models.AdGroup
	page, err := s.Client.Post(fmt.Sprintf("/campaigns/%d/adgroups/find", campaignID), &selector, &adgroups)
	return adgroups, page, err
}

func (s *AdGroupService) FindAll(campaignID int64, selector models.Selector) ([]models.AdGroup, error) {
	return api.PaginatedFetcher[models.AdGroup](s.Client, fmt.Sprintf("/campaigns/%d/adgroups/find", campaignID), selector)
}

func (s *AdGroupService) Create(campaignID int64, adgroup *models.AdGroup) (*models.AdGroup, error) {
	var created models.AdGroup
	_, err := s.Client.Post(fmt.Sprintf("/campaigns/%d/adgroups", campaignID), adgroup, &created)
	return &created, err
}

func (s *AdGroupService) Update(campaignID, adGroupID int64, update *models.AdGroupUpdate) (*models.AdGroup, error) {
	var updated models.AdGroup
	_, err := s.Client.Put(fmt.Sprintf("/campaigns/%d/adgroups/%d", campaignID, adGroupID), update, &updated)
	return &updated, err
}

func (s *AdGroupService) Delete(campaignID, adGroupID int64) error {
	return s.Client.Delete(fmt.Sprintf("/campaigns/%d/adgroups/%d", campaignID, adGroupID))
}
