package model

import "time"

type IndexStatus string

const (
	IndexStatusPending    IndexStatus = "pending"
	IndexStatusIndexed    IndexStatus = "indexed"
	IndexStatusNotIndexed IndexStatus = "not_indexed"
	IndexStatusError      IndexStatus = "error"
)

type Platform struct {
	ID             int64       `json:"id"`
	UserID         int64       `json:"user_id"`
	URL            string      `json:"url"`
	Domain         string      `json:"domain"`
	IndexStatus    IndexStatus `json:"index_status"`
	IsIndexed      bool        `json:"is_indexed"`
	FirstIndexedAt *time.Time  `json:"first_indexed_at"`
	LastCheckedAt  *time.Time  `json:"last_checked_at"`
	CheckCount     int         `json:"check_count"`
	PotentialScore int         `json:"potential_score"`
	IsMustHave     bool        `json:"is_must_have"`
	Notes          string      `json:"notes"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

type IndexCheckResult struct {
	PlatformID  int64       `json:"platform_id"`
	URL         string      `json:"url"`
	HTTPStatus  int         `json:"http_status"`
	IsIndexed   bool        `json:"is_indexed"`
	IndexStatus IndexStatus `json:"index_status"`
	CheckedAt   time.Time   `json:"checked_at"`
	Error       string      `json:"error,omitempty"`
}
