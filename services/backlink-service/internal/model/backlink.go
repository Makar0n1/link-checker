package model

import "time"

type LinkStatus string

const (
	LinkStatusPending  LinkStatus = "pending"
	LinkStatusActive   LinkStatus = "active"
	LinkStatusBroken   LinkStatus = "broken"
	LinkStatusRemoved  LinkStatus = "removed"
	LinkStatusNoFollow LinkStatus = "nofollow"
)

type LinkType string

const (
	LinkTypeDoFollow  LinkType = "dofollow"
	LinkTypeNoFollow  LinkType = "nofollow"
	LinkTypeSponsored LinkType = "sponsored"
	LinkTypeUGC       LinkType = "ugc"
)

type Project struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	UserID        int64     `json:"user_id"`
	GoogleSheetID *string   `json:"google_sheet_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

type Backlink struct {
	ID            int64      `json:"id"`
	ProjectID     int64      `json:"project_id"`
	SourceURL     string     `json:"source_url"`
	TargetURL     string     `json:"target_url"`
	AnchorText    string     `json:"anchor_text"`
	Status        LinkStatus `json:"status"`
	LinkType      LinkType   `json:"link_type"`
	HTTPStatus    *int       `json:"http_status,omitempty"`
	LastCheckedAt *time.Time `json:"last_checked_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}
