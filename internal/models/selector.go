package models

// Selector is the request body for Find endpoints.
type Selector struct {
	Conditions []Condition     `json:"conditions,omitempty"`
	Fields     []string        `json:"fields,omitempty"`
	OrderBy    []OrderByItem   `json:"orderBy,omitempty"`
	Pagination SelectorPagination `json:"pagination"`
}

// Condition represents a single filter condition.
type Condition struct {
	Field    string   `json:"field"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}

// OrderByItem represents a sorting criterion.
type OrderByItem struct {
	Field     string `json:"field"`
	SortOrder string `json:"sortOrder"` // ASCENDING or DESCENDING
}

// SelectorPagination controls result pagination.
type SelectorPagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

// NewSelector creates a Selector with default pagination.
func NewSelector(limit, offset int) Selector {
	if limit <= 0 {
		limit = 20
	}
	return Selector{
		Pagination: SelectorPagination{
			Offset: offset,
			Limit:  limit,
		},
	}
}

// ParseFilterOperator converts a shorthand operator to the API operator string.
func ParseFilterOperator(op string) string {
	switch op {
	case "=":
		return "EQUALS"
	case "~":
		return "CONTAINS"
	case "@":
		return "IN"
	case ">":
		return "GREATER_THAN"
	case "<":
		return "LESS_THAN"
	case ">=":
		return "GREATER_THAN_OR_EQUAL"
	case "<=":
		return "LESS_THAN_OR_EQUAL"
	case "!~":
		return "NOT_CONTAINS"
	default:
		return op // Allow passing raw operator names
	}
}

// ParseSortOrder converts sort direction shorthand.
func ParseSortOrder(dir string) string {
	switch dir {
	case "asc", "ASC", "ASCENDING":
		return "ASCENDING"
	case "desc", "DESC", "DESCENDING":
		return "DESCENDING"
	default:
		return "ASCENDING"
	}
}
