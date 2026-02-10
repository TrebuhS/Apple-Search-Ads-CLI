package services

import (
	"fmt"
	"net/url"

	"github.com/trebuhs/asa-cli/internal/api"
	"github.com/trebuhs/asa-cli/internal/models"
)

type AppService struct {
	Client *api.Client
}

func NewAppService(client *api.Client) *AppService {
	return &AppService{Client: client}
}

func (s *AppService) Search(query string, limit, offset int, returnOwnedApps bool) ([]models.AppInfo, *models.PageDetail, error) {
	q := url.QueryEscape(query)
	path := fmt.Sprintf("/search/apps?query=%s&limit=%d&offset=%d&returnOwnedApps=%t", q, limit, offset, returnOwnedApps)
	var apps []models.AppInfo
	page, err := s.Client.Get(path, &apps)
	return apps, page, err
}

func (s *AppService) SearchGeo(query string, limit, offset int, entity string, countryCode string) ([]models.GeoEntity, *models.PageDetail, error) {
	q := url.QueryEscape(query)
	path := fmt.Sprintf("/search/geo?query=%s&limit=%d&offset=%d", q, limit, offset)
	if entity != "" {
		path += "&entity=" + url.QueryEscape(entity)
	}
	if countryCode != "" {
		path += "&countrycode=" + url.QueryEscape(countryCode)
	}
	var geos []models.GeoEntity
	page, err := s.Client.Get(path, &geos)
	return geos, page, err
}
