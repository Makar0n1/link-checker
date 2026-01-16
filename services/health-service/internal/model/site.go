package model

import "time"

type MonitoredSite struct {
	ID              int64      `json:"id"`
	UserID          int64      `json:"user_id"`
	URL             string     `json:"url"`
	Domain          string     `json:"domain"`
	HTTPStatus      *int       `json:"http_status"`
	IsAlive         bool       `json:"is_alive"`
	ResponseTimeMs  *int       `json:"response_time_ms"`
	AllowsIndexing  *bool      `json:"allows_indexing"`
	RobotsTxtStatus string     `json:"robots_txt_status"`
	HasNoindex      bool       `json:"has_noindex"`
	PagesIndexed    int        `json:"pages_indexed"`
	LastCheckedAt   *time.Time `json:"last_checked_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

type SiteCheckHistory struct {
	ID             int64     `json:"id"`
	SiteID         int64     `json:"site_id"`
	HTTPStatus     *int      `json:"http_status"`
	IsAlive        bool      `json:"is_alive"`
	ResponseTimeMs *int      `json:"response_time_ms"`
	CheckedAt      time.Time `json:"checked_at"`
}

type SiteHealthCheck struct {
	SiteID          int64     `json:"site_id"`
	URL             string    `json:"url"`
	HTTPStatus      int       `json:"http_status"`
	IsAlive         bool      `json:"is_alive"`
	ResponseTimeMs  int       `json:"response_time_ms"`
	AllowsIndexing  bool      `json:"allows_indexing"`
	RobotsTxtStatus string    `json:"robots_txt_status"`
	HasNoindex      bool      `json:"has_noindex"`
	CheckedAt       time.Time `json:"checked_at"`
	Error           string    `json:"error,omitempty"`
}
