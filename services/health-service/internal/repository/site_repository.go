package repository

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/link-tracker/health-service/internal/model"
)

type SiteRepository struct {
	db *pgxpool.Pool
}

func NewSiteRepository(db *pgxpool.Pool) *SiteRepository {
	return &SiteRepository{db: db}
}

func (r *SiteRepository) Create(ctx context.Context, userID int64, req *model.CreateSiteRequest) (*model.MonitoredSite, error) {
	domain := extractDomain(req.URL)

	var site model.MonitoredSite
	err := r.db.QueryRow(ctx, `
		INSERT INTO monitored_sites (user_id, url, domain)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, url, domain, http_status, is_alive, response_time_ms,
		          allows_indexing, robots_txt_status, has_noindex, pages_indexed, last_checked_at, created_at
	`, userID, req.URL, domain).Scan(
		&site.ID, &site.UserID, &site.URL, &site.Domain, &site.HTTPStatus, &site.IsAlive,
		&site.ResponseTimeMs, &site.AllowsIndexing, &site.RobotsTxtStatus, &site.HasNoindex,
		&site.PagesIndexed, &site.LastCheckedAt, &site.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &site, nil
}

func (r *SiteRepository) GetByID(ctx context.Context, id int64) (*model.MonitoredSite, error) {
	var site model.MonitoredSite
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, url, domain, http_status, is_alive, response_time_ms,
		       allows_indexing, robots_txt_status, has_noindex, pages_indexed, last_checked_at, created_at
		FROM monitored_sites WHERE id = $1
	`, id).Scan(
		&site.ID, &site.UserID, &site.URL, &site.Domain, &site.HTTPStatus, &site.IsAlive,
		&site.ResponseTimeMs, &site.AllowsIndexing, &site.RobotsTxtStatus, &site.HasNoindex,
		&site.PagesIndexed, &site.LastCheckedAt, &site.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &site, nil
}

func (r *SiteRepository) List(ctx context.Context, userID int64, filters *model.SiteFilters) ([]model.MonitoredSite, int64, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
	args = append(args, userID)
	argIndex++

	if filters.IsAlive != nil {
		conditions = append(conditions, fmt.Sprintf("is_alive = $%d", argIndex))
		args = append(args, *filters.IsAlive)
		argIndex++
	}

	if filters.HasNoindex != nil {
		conditions = append(conditions, fmt.Sprintf("has_noindex = $%d", argIndex))
		args = append(args, *filters.HasNoindex)
		argIndex++
	}

	if filters.Domain != "" {
		conditions = append(conditions, fmt.Sprintf("domain ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Domain+"%")
		argIndex++
	}

	whereClause := strings.Join(conditions, " AND ")

	// Count total
	var total int64
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM monitored_sites WHERE %s", whereClause)
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (filters.Page - 1) * filters.PerPage
	args = append(args, filters.PerPage, offset)

	query := fmt.Sprintf(`
		SELECT id, user_id, url, domain, http_status, is_alive, response_time_ms,
		       allows_indexing, robots_txt_status, has_noindex, pages_indexed, last_checked_at, created_at
		FROM monitored_sites
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var sites []model.MonitoredSite
	for rows.Next() {
		var s model.MonitoredSite
		err := rows.Scan(
			&s.ID, &s.UserID, &s.URL, &s.Domain, &s.HTTPStatus, &s.IsAlive,
			&s.ResponseTimeMs, &s.AllowsIndexing, &s.RobotsTxtStatus, &s.HasNoindex,
			&s.PagesIndexed, &s.LastCheckedAt, &s.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		sites = append(sites, s)
	}

	return sites, total, nil
}

func (r *SiteRepository) Update(ctx context.Context, id int64, req *model.UpdateSiteRequest) (*model.MonitoredSite, error) {
	var setClauses []string
	var args []interface{}
	argIndex := 1

	if req.URL != nil {
		setClauses = append(setClauses, fmt.Sprintf("url = $%d", argIndex))
		args = append(args, *req.URL)
		argIndex++
		setClauses = append(setClauses, fmt.Sprintf("domain = $%d", argIndex))
		args = append(args, extractDomain(*req.URL))
		argIndex++
	}

	if req.PagesIndexed != nil {
		setClauses = append(setClauses, fmt.Sprintf("pages_indexed = $%d", argIndex))
		args = append(args, *req.PagesIndexed)
		argIndex++
	}

	if len(setClauses) == 0 {
		return r.GetByID(ctx, id)
	}

	args = append(args, id)
	query := fmt.Sprintf(`
		UPDATE monitored_sites SET %s
		WHERE id = $%d
		RETURNING id, user_id, url, domain, http_status, is_alive, response_time_ms,
		          allows_indexing, robots_txt_status, has_noindex, pages_indexed, last_checked_at, created_at
	`, strings.Join(setClauses, ", "), argIndex)

	var site model.MonitoredSite
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&site.ID, &site.UserID, &site.URL, &site.Domain, &site.HTTPStatus, &site.IsAlive,
		&site.ResponseTimeMs, &site.AllowsIndexing, &site.RobotsTxtStatus, &site.HasNoindex,
		&site.PagesIndexed, &site.LastCheckedAt, &site.CreatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &site, nil
}

func (r *SiteRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, "DELETE FROM monitored_sites WHERE id = $1", id)
	return err
}

func (r *SiteRepository) UpdateHealthCheck(ctx context.Context, id int64, check *model.SiteHealthCheck) error {
	_, err := r.db.Exec(ctx, `
		UPDATE monitored_sites
		SET http_status = $1, is_alive = $2, response_time_ms = $3, allows_indexing = $4,
		    robots_txt_status = $5, has_noindex = $6, last_checked_at = NOW()
		WHERE id = $7
	`, check.HTTPStatus, check.IsAlive, check.ResponseTimeMs, check.AllowsIndexing,
		check.RobotsTxtStatus, check.HasNoindex, id)
	return err
}

func (r *SiteRepository) AddCheckHistory(ctx context.Context, siteID int64, check *model.SiteHealthCheck) error {
	_, err := r.db.Exec(ctx, `
		INSERT INTO site_check_history (site_id, http_status, is_alive, response_time_ms)
		VALUES ($1, $2, $3, $4)
	`, siteID, check.HTTPStatus, check.IsAlive, check.ResponseTimeMs)
	return err
}

func (r *SiteRepository) GetHistory(ctx context.Context, siteID int64, filters *model.HistoryFilters) ([]model.SiteCheckHistory, int64, error) {
	// Count total
	var total int64
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM site_check_history WHERE site_id = $1", siteID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (filters.Page - 1) * filters.PerPage
	rows, err := r.db.Query(ctx, `
		SELECT id, site_id, http_status, is_alive, response_time_ms, checked_at
		FROM site_check_history
		WHERE site_id = $1
		ORDER BY checked_at DESC
		LIMIT $2 OFFSET $3
	`, siteID, filters.PerPage, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var history []model.SiteCheckHistory
	for rows.Next() {
		var h model.SiteCheckHistory
		err := rows.Scan(&h.ID, &h.SiteID, &h.HTTPStatus, &h.IsAlive, &h.ResponseTimeMs, &h.CheckedAt)
		if err != nil {
			return nil, 0, err
		}
		history = append(history, h)
	}

	return history, total, nil
}

func extractDomain(rawURL string) string {
	if !strings.HasPrefix(rawURL, "http://") && !strings.HasPrefix(rawURL, "https://") {
		rawURL = "https://" + rawURL
	}
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return parsed.Host
}
