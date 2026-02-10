package services

import (
	"fmt"

	"github.com/trebuhs/asa-cli/internal/api"
	"github.com/trebuhs/asa-cli/internal/models"
)

type KeywordService struct {
	Client *api.Client
}

func NewKeywordService(client *api.Client) *KeywordService {
	return &KeywordService{Client: client}
}

// --- Targeting Keywords ---

func (s *KeywordService) List(campaignID, adGroupID int64, limit, offset int) ([]models.Keyword, *models.PageDetail, error) {
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/targetingkeywords?limit=%d&offset=%d", campaignID, adGroupID, limit, offset)
	var keywords []models.Keyword
	page, err := s.Client.Get(path, &keywords)
	return keywords, page, err
}

func (s *KeywordService) Get(campaignID, adGroupID, keywordID int64) (*models.Keyword, error) {
	var keyword models.Keyword
	_, err := s.Client.Get(fmt.Sprintf("/campaigns/%d/adgroups/%d/targetingkeywords/%d", campaignID, adGroupID, keywordID), &keyword)
	return &keyword, err
}

func (s *KeywordService) Find(campaignID, adGroupID int64, selector models.Selector) ([]models.Keyword, *models.PageDetail, error) {
	var keywords []models.Keyword
	page, err := s.Client.Post(fmt.Sprintf("/campaigns/%d/adgroups/%d/targetingkeywords/find", campaignID, adGroupID), &selector, &keywords)
	return keywords, page, err
}

func (s *KeywordService) FindAll(campaignID, adGroupID int64, selector models.Selector) ([]models.Keyword, error) {
	return api.PaginatedFetcher[models.Keyword](s.Client, fmt.Sprintf("/campaigns/%d/adgroups/%d/targetingkeywords/find", campaignID, adGroupID), selector)
}

func (s *KeywordService) Create(campaignID, adGroupID int64, keywords []models.Keyword) ([]models.Keyword, error) {
	var created []models.Keyword
	_, err := s.Client.Post(fmt.Sprintf("/campaigns/%d/adgroups/%d/targetingkeywords/bulk", campaignID, adGroupID), keywords, &created)
	return created, err
}

func (s *KeywordService) Update(campaignID, adGroupID int64, updates []models.KeywordUpdate) ([]models.Keyword, error) {
	var updated []models.Keyword
	_, err := s.Client.Put(fmt.Sprintf("/campaigns/%d/adgroups/%d/targetingkeywords/bulk", campaignID, adGroupID), updates, &updated)
	return updated, err
}

func (s *KeywordService) Delete(campaignID, adGroupID int64, keywordIDs []int64) error {
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/targetingkeywords/delete/bulk", campaignID, adGroupID)
	_, err := s.Client.Post(path, keywordIDs, nil)
	return err
}

// --- Campaign-level Negative Keywords ---

func (s *KeywordService) ListCampaignNegativeKeywords(campaignID int64, limit, offset int) ([]models.NegativeKeyword, *models.PageDetail, error) {
	path := fmt.Sprintf("/campaigns/%d/negativekeywords?limit=%d&offset=%d", campaignID, limit, offset)
	var keywords []models.NegativeKeyword
	page, err := s.Client.Get(path, &keywords)
	return keywords, page, err
}

func (s *KeywordService) GetCampaignNegativeKeyword(campaignID, keywordID int64) (*models.NegativeKeyword, error) {
	var keyword models.NegativeKeyword
	_, err := s.Client.Get(fmt.Sprintf("/campaigns/%d/negativekeywords/%d", campaignID, keywordID), &keyword)
	return &keyword, err
}

func (s *KeywordService) FindCampaignNegativeKeywords(campaignID int64, selector models.Selector) ([]models.NegativeKeyword, *models.PageDetail, error) {
	var keywords []models.NegativeKeyword
	page, err := s.Client.Post(fmt.Sprintf("/campaigns/%d/negativekeywords/find", campaignID), &selector, &keywords)
	return keywords, page, err
}

func (s *KeywordService) CreateCampaignNegativeKeywords(campaignID int64, keywords []models.NegativeKeyword) ([]models.NegativeKeyword, error) {
	var created []models.NegativeKeyword
	_, err := s.Client.Post(fmt.Sprintf("/campaigns/%d/negativekeywords/bulk", campaignID), keywords, &created)
	return created, err
}

func (s *KeywordService) DeleteCampaignNegativeKeywords(campaignID int64, keywordIDs []int64) error {
	path := fmt.Sprintf("/campaigns/%d/negativekeywords/delete/bulk", campaignID)
	_, err := s.Client.Post(path, keywordIDs, nil)
	return err
}

// --- Ad Group-level Negative Keywords ---

func (s *KeywordService) ListAdGroupNegativeKeywords(campaignID, adGroupID int64, limit, offset int) ([]models.NegativeKeyword, *models.PageDetail, error) {
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/negativekeywords?limit=%d&offset=%d", campaignID, adGroupID, limit, offset)
	var keywords []models.NegativeKeyword
	page, err := s.Client.Get(path, &keywords)
	return keywords, page, err
}

func (s *KeywordService) GetAdGroupNegativeKeyword(campaignID, adGroupID, keywordID int64) (*models.NegativeKeyword, error) {
	var keyword models.NegativeKeyword
	_, err := s.Client.Get(fmt.Sprintf("/campaigns/%d/adgroups/%d/negativekeywords/%d", campaignID, adGroupID, keywordID), &keyword)
	return &keyword, err
}

func (s *KeywordService) FindAdGroupNegativeKeywords(campaignID, adGroupID int64, selector models.Selector) ([]models.NegativeKeyword, *models.PageDetail, error) {
	var keywords []models.NegativeKeyword
	page, err := s.Client.Post(fmt.Sprintf("/campaigns/%d/adgroups/%d/negativekeywords/find", campaignID, adGroupID), &selector, &keywords)
	return keywords, page, err
}

func (s *KeywordService) CreateAdGroupNegativeKeywords(campaignID, adGroupID int64, keywords []models.NegativeKeyword) ([]models.NegativeKeyword, error) {
	var created []models.NegativeKeyword
	_, err := s.Client.Post(fmt.Sprintf("/campaigns/%d/adgroups/%d/negativekeywords/bulk", campaignID, adGroupID), keywords, &created)
	return created, err
}

func (s *KeywordService) DeleteAdGroupNegativeKeywords(campaignID, adGroupID int64, keywordIDs []int64) error {
	path := fmt.Sprintf("/campaigns/%d/adgroups/%d/negativekeywords/delete/bulk", campaignID, adGroupID)
	_, err := s.Client.Post(path, keywordIDs, nil)
	return err
}
