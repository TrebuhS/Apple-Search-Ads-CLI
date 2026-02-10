package services

import (
	"encoding/json"
	"fmt"

	"github.com/trebuhs/asa-cli/internal/api"
	"github.com/trebuhs/asa-cli/internal/models"
)

type ReportingService struct {
	Client *api.Client
}

func NewReportingService(client *api.Client) *ReportingService {
	return &ReportingService{Client: client}
}

func (s *ReportingService) GetCampaignReport(req *models.ReportRequest) (*models.ReportingDataResponse, error) {
	return s.getReport("/reports/campaigns", req)
}

func (s *ReportingService) GetAdGroupReport(campaignID int64, req *models.ReportRequest) (*models.ReportingDataResponse, error) {
	return s.getReport(fmt.Sprintf("/reports/campaigns/%d/adgroups", campaignID), req)
}

func (s *ReportingService) GetKeywordReport(campaignID int64, req *models.ReportRequest) (*models.ReportingDataResponse, error) {
	return s.getReport(fmt.Sprintf("/reports/campaigns/%d/keywords", campaignID), req)
}

func (s *ReportingService) GetSearchTermReport(campaignID int64, req *models.ReportRequest) (*models.ReportingDataResponse, error) {
	return s.getReport(fmt.Sprintf("/reports/campaigns/%d/searchterms", campaignID), req)
}

func (s *ReportingService) getReport(path string, req *models.ReportRequest) (*models.ReportingDataResponse, error) {
	var raw json.RawMessage
	_, err := s.Client.Post(path, req, &raw)
	if err != nil {
		return nil, err
	}

	var resp models.ReportResponse
	if err := json.Unmarshal(raw, &resp); err != nil {
		// Try direct unmarshal
		var direct models.ReportingDataResponse
		if err2 := json.Unmarshal(raw, &direct); err2 != nil {
			return nil, fmt.Errorf("parsing report response: %w", err)
		}
		return &direct, nil
	}

	return &resp.ReportingDataResponse, nil
}
