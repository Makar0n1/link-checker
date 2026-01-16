package model

type CreatePlatformRequest struct {
	URL            string `json:"url"`
	PotentialScore int    `json:"potential_score,omitempty"`
	IsMustHave     bool   `json:"is_must_have,omitempty"`
	Notes          string `json:"notes,omitempty"`
}

type UpdatePlatformRequest struct {
	URL            *string `json:"url,omitempty"`
	PotentialScore *int    `json:"potential_score,omitempty"`
	IsMustHave     *bool   `json:"is_must_have,omitempty"`
	Notes          *string `json:"notes,omitempty"`
}

type BulkCreatePlatformsRequest struct {
	Platforms []CreatePlatformRequest `json:"platforms"`
}

type BulkOperationResponse struct {
	Success int           `json:"success"`
	Failed  int           `json:"failed"`
	Errors  []BulkError   `json:"errors"`
	Created []Platform    `json:"created,omitempty"`
}

type BulkError struct {
	Index   int    `json:"index"`
	URL     string `json:"url"`
	Message string `json:"message"`
}

type PlatformFilters struct {
	IndexStatus *IndexStatus `json:"index_status,omitempty"`
	IsIndexed   *bool        `json:"is_indexed,omitempty"`
	IsMustHave  *bool        `json:"is_must_have,omitempty"`
	Domain      string       `json:"domain,omitempty"`
	Page        int          `json:"page"`
	PerPage     int          `json:"per_page"`
}
