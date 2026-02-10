package models

// UserACL represents an Access Control List entry.
type UserACL struct {
	OrgName    string   `json:"orgName"`
	OrgID      int64    `json:"orgId"`
	Currency   string   `json:"currency"`
	RoleNames  []string `json:"roleNames"`
	ParentOrgID *int64  `json:"parentOrgId,omitempty"`
}
