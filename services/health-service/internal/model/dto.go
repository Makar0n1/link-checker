package model

type CreateSiteRequest struct {
	URL string `json:"url"`
}

type UpdateSiteRequest struct {
	URL          *string `json:"url,omitempty"`
	PagesIndexed *int    `json:"pages_indexed,omitempty"`
}

type SiteFilters struct {
	IsAlive    *bool  `json:"is_alive,omitempty"`
	HasNoindex *bool  `json:"has_noindex,omitempty"`
	Domain     string `json:"domain,omitempty"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
}

type HistoryFilters struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}
