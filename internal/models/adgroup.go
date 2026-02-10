package models

// AdGroup represents an Apple Search Ads ad group.
type AdGroup struct {
	ID                    int64    `json:"id,omitempty"`
	CampaignID            int64    `json:"campaignId,omitempty"`
	OrgID                 int64    `json:"orgId,omitempty"`
	Name                  string   `json:"name"`
	Status                string   `json:"status,omitempty"`
	ServingStatus         string   `json:"servingStatus,omitempty"`
	ServingStateReasons   []string `json:"servingStateReasons,omitempty"`
	DisplayStatus         string   `json:"displayStatus,omitempty"`
	DefaultBidAmount      *Money   `json:"defaultBidAmount,omitempty"`
	CpaGoal               *Money   `json:"cpaGoal,omitempty"`
	AutomatedKeywordsOptIn bool   `json:"automatedKeywordsOptIn,omitempty"`
	StartTime             string   `json:"startTime,omitempty"`
	EndTime               string   `json:"endTime,omitempty"`
	ModificationTime      string   `json:"modificationTime,omitempty"`
	TargetingDimensions   *TargetingDimensions `json:"targetingDimensions,omitempty"`
	PaymentModel          string   `json:"paymentModel,omitempty"`
	PricingModel          string   `json:"pricingModel,omitempty"`
}

// TargetingDimensions for ad group targeting.
type TargetingDimensions struct {
	Age            *TargetingDimension `json:"age,omitempty"`
	Gender         *TargetingDimension `json:"gender,omitempty"`
	DeviceClass    *TargetingDimension `json:"deviceClass,omitempty"`
	Locality       *TargetingDimension `json:"locality,omitempty"`
	AdminArea      *TargetingDimension `json:"adminArea,omitempty"`
	Country        *TargetingDimension `json:"country,omitempty"`
	AppDownloaders *TargetingDimension `json:"appDownloaders,omitempty"`
	DayPart        *TargetingDimension `json:"daypart,omitempty"`
}

// TargetingDimension is a single targeting dimension.
type TargetingDimension struct {
	Included []interface{} `json:"included,omitempty"`
	Excluded []interface{} `json:"excluded,omitempty"`
}

// AdGroupUpdate contains fields that can be updated on an ad group.
type AdGroupUpdate struct {
	Name                   string `json:"name,omitempty"`
	Status                 string `json:"status,omitempty"`
	DefaultBidAmount       *Money `json:"defaultBidAmount,omitempty"`
	CpaGoal                *Money `json:"cpaGoal,omitempty"`
	AutomatedKeywordsOptIn *bool  `json:"automatedKeywordsOptIn,omitempty"`
	StartTime              string `json:"startTime,omitempty"`
	EndTime                string `json:"endTime,omitempty"`
}
