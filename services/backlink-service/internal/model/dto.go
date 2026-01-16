package model

// Request DTOs

type CreateBacklinkRequest struct {
	ProjectID  int64      `json:"project_id"`
	SourceURL  string     `json:"source_url"`
	TargetURL  string     `json:"target_url"`
	AnchorText string     `json:"anchor_text"`
	LinkType   LinkType   `json:"link_type"`
}

type UpdateBacklinkRequest struct {
	SourceURL  *string     `json:"source_url,omitempty"`
	TargetURL  *string     `json:"target_url,omitempty"`
	AnchorText *string     `json:"anchor_text,omitempty"`
	Status     *LinkStatus `json:"status,omitempty"`
	LinkType   *LinkType   `json:"link_type,omitempty"`
}

type BulkCreateBacklinksRequest struct {
	Backlinks []CreateBacklinkRequest `json:"backlinks"`
}

type BulkDeleteBacklinksRequest struct {
	IDs []int64 `json:"ids"`
}

type ImportFromSheetsRequest struct {
	ProjectID     int64  `json:"project_id"`
	GoogleSheetID string `json:"google_sheet_id"`
	SheetName     string `json:"sheet_name"`
}

type CreateProjectRequest struct {
	Name          string  `json:"name"`
	GoogleSheetID *string `json:"google_sheet_id,omitempty"`
}

type UpdateProjectRequest struct {
	Name          *string `json:"name,omitempty"`
	GoogleSheetID *string `json:"google_sheet_id,omitempty"`
}

// Query parameters

type BacklinkFilters struct {
	ProjectID  *int64      `json:"project_id,omitempty"`
	Status     *LinkStatus `json:"status,omitempty"`
	LinkType   *LinkType   `json:"link_type,omitempty"`
	SourceURL  *string     `json:"source_url,omitempty"`
	TargetURL  *string     `json:"target_url,omitempty"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
}

// Response DTOs

type BacklinkResponse struct {
	ID            int64      `json:"id"`
	ProjectID     int64      `json:"project_id"`
	SourceURL     string     `json:"source_url"`
	TargetURL     string     `json:"target_url"`
	AnchorText    string     `json:"anchor_text"`
	Status        LinkStatus `json:"status"`
	LinkType      LinkType   `json:"link_type"`
	HTTPStatus    *int       `json:"http_status,omitempty"`
	LastCheckedAt *string    `json:"last_checked_at,omitempty"`
	CreatedAt     string     `json:"created_at"`
}

type ProjectResponse struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	UserID        int64   `json:"user_id"`
	GoogleSheetID *string `json:"google_sheet_id,omitempty"`
	CreatedAt     string  `json:"created_at"`
}

type BulkOperationResponse struct {
	Success int   `json:"success"`
	Failed  int   `json:"failed"`
	Errors  []string `json:"errors,omitempty"`
}

func BacklinkToResponse(b *Backlink) BacklinkResponse {
	resp := BacklinkResponse{
		ID:         b.ID,
		ProjectID:  b.ProjectID,
		SourceURL:  b.SourceURL,
		TargetURL:  b.TargetURL,
		AnchorText: b.AnchorText,
		Status:     b.Status,
		LinkType:   b.LinkType,
		HTTPStatus: b.HTTPStatus,
		CreatedAt:  b.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
	if b.LastCheckedAt != nil {
		formatted := b.LastCheckedAt.Format("2006-01-02T15:04:05Z")
		resp.LastCheckedAt = &formatted
	}
	return resp
}

func ProjectToResponse(p *Project) ProjectResponse {
	return ProjectResponse{
		ID:            p.ID,
		Name:          p.Name,
		UserID:        p.UserID,
		GoogleSheetID: p.GoogleSheetID,
		CreatedAt:     p.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
