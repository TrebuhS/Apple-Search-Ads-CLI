package models

// ReportRequest is the request body for reporting endpoints.
type ReportRequest struct {
	StartTime        string   `json:"startTime"`
	EndTime          string   `json:"endTime"`
	Granularity      string   `json:"granularity,omitempty"` // HOURLY, DAILY, WEEKLY, MONTHLY
	GroupBy          []string `json:"groupBy,omitempty"`     // countryOrRegion, deviceClass, ageRange, gender, adminArea, locality
	Selector         *Selector `json:"selector,omitempty"`
	ReturnGrandTotals bool    `json:"returnGrandTotals,omitempty"`
	ReturnRecordsWithNoMetrics bool `json:"returnRecordsWithNoMetrics,omitempty"`
	ReturnRowTotals  bool    `json:"returnRowTotals,omitempty"`
	TimeZone         string  `json:"timeZone,omitempty"`
}

// ReportResponse wraps reporting response data.
type ReportResponse struct {
	ReportingDataResponse ReportingDataResponse `json:"reportingDataResponse"`
}

// ReportingDataResponse contains the actual report rows.
type ReportingDataResponse struct {
	Row        []ReportRow    `json:"row"`
	GrandTotals *ReportRow   `json:"grandTotals,omitempty"`
}

// ReportRow represents a single row in a report.
type ReportRow struct {
	Other    bool                   `json:"other,omitempty"`
	Total    *SpendRow              `json:"total,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Granularity []GranularityRow   `json:"granularity,omitempty"`
	Insights *InsightData           `json:"insights,omitempty"`
}

// SpendRow contains the metrics for a report row.
type SpendRow struct {
	Impressions    int64   `json:"impressions"`
	Taps           int64   `json:"taps"`
	Installs       int64   `json:"installs"`
	NewDownloads   int64   `json:"newDownloads"`
	Redownloads    int64   `json:"redownloads"`
	LatOnInstalls  int64   `json:"latOnInstalls"`
	LatOffInstalls int64   `json:"latOffInstalls"`
	TTR            float64 `json:"ttr"`
	AvgCPA         Money   `json:"avgCPA"`
	AvgCPT         Money   `json:"avgCPT"`
	LocalSpend     Money   `json:"localSpend"`
	ConversionRate float64 `json:"conversionRate"`
}

// GranularityRow is a time-bucketed metrics row.
type GranularityRow struct {
	Date    string    `json:"date"`
	Metrics *SpendRow `json:"metrics,omitempty"`
}

// InsightData contains keyword-level insights.
type InsightData struct {
	BidRecommendation *BidRecommendation `json:"bidRecommendation,omitempty"`
}

// BidRecommendation for keyword bid suggestions.
type BidRecommendation struct {
	SuggestedBidAmount *Money `json:"suggestedBidAmount,omitempty"`
}

// SearchTermReportRow is a row in the search terms report.
type SearchTermReportRow struct {
	Other    bool                   `json:"other,omitempty"`
	Total    *SpendRow              `json:"total,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}
