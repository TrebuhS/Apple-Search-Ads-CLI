package models

// Keyword represents a targeting keyword.
type Keyword struct {
	ID               int64  `json:"id,omitempty"`
	CampaignID       int64  `json:"campaignId,omitempty"`
	AdGroupID        int64  `json:"adGroupId,omitempty"`
	Text             string `json:"text"`
	MatchType        string `json:"matchType"` // BROAD or EXACT
	Status           string `json:"status,omitempty"`
	BidAmount        *Money `json:"bidAmount,omitempty"`
	Deleted          bool   `json:"deleted,omitempty"`
	ModificationTime string `json:"modificationTime,omitempty"`
}

// NegativeKeyword represents a negative keyword (campaign or ad-group level).
type NegativeKeyword struct {
	ID               int64  `json:"id,omitempty"`
	CampaignID       int64  `json:"campaignId,omitempty"`
	AdGroupID        int64  `json:"adGroupId,omitempty"`
	Text             string `json:"text"`
	MatchType        string `json:"matchType"` // BROAD or EXACT
	Status           string `json:"status,omitempty"`
	Deleted          bool   `json:"deleted,omitempty"`
	ModificationTime string `json:"modificationTime,omitempty"`
}

// KeywordUpdate contains fields that can be updated on a keyword.
type KeywordUpdate struct {
	ID        int64  `json:"id"`
	Status    string `json:"status,omitempty"`
	BidAmount *Money `json:"bidAmount,omitempty"`
}
