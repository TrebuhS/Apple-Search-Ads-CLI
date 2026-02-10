package models

// AppInfo represents an app from the search API.
type AppInfo struct {
	AdamID       int64  `json:"adamId"`
	AppName      string `json:"appName"`
	DeveloperName string `json:"developerName"`
	CountryOrRegionCodes []string `json:"countryOrRegionCodes,omitempty"`
}

// GeoEntity represents a geographic location.
type GeoEntity struct {
	ID          string `json:"id"`
	Entity      string `json:"entity"`
	DisplayName string `json:"displayName"`
}
