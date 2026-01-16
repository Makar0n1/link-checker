package repository

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/link-tracker/index-service/internal/model"
)

type PlatformRepository struct {
	db *pgxpool.Pool
}

func NewPlatformRepository(db *pgxpool.Pool) *PlatformRepository {
	return &PlatformRepository{db: db}
}

func (r *PlatformRepository) Create(ctx context.Context, userID int64, req *model.CreatePlatformRequest) (*model.Platform, error) {
	domain := extractDomain(req.URL)

	var platform model.Platform
	err := r.db.QueryRow(ctx, `
		INSERT INTO platforms (user_id, url, domain, potential_score, is_must_have, notes)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, url, domain, index_status, is_indexed, first_indexed_at,
		          last_checked_at, check_count, potential_score, is_must_have, notes, created_at, updated_at
	`, userID, req.URL, domain, req.PotentialScore, req.IsMustHave, req.Notes).Scan(
		&platform.ID, &platform.UserID, &platform.URL, &platform.Domain, &platform.IndexStatus,
		&platform.IsIndexed, &platform.FirstIndexedAt, &platform.LastCheckedAt, &platform.CheckCount,
		&platform.PotentialScore, &platform.IsMustHave, &platform.Notes, &platform.CreatedAt, &platform.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &platform, nil
}

func (r *PlatformRepository) GetByID(ctx context.Context, id int64) (*model.Platform, error) {
	var platform model.Platform
	err := r.db.QueryRow(ctx, `
		SELECT id, user_id, url, domain, index_status, is_indexed, first_indexed_at,
		       last_checked_at, check_count, potential_score, is_must_have, notes, created_at, updated_at
		FROM platforms WHERE id = $1
	`, id).Scan(
		&platform.ID, &platform.UserID, &platform.URL, &platform.Domain, &platform.IndexStatus,
		&platform.IsIndexed, &platform.FirstIndexedAt, &platform.LastCheckedAt, &platform.CheckCount,
		&platform.PotentialScore, &platform.IsMustHave, &platform.Notes, &platform.CreatedAt, &platform.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &platform, nil
}

func (r *PlatformRepository) List(ctx context.Context, userID int64, filters *model.PlatformFilters) ([]model.Platform, int64, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	conditions = append(conditions, fmt.Sprintf("user_id = $%d", argIndex))
	args = append(args, userID)
	argIndex++

	if filters.IndexStatus != nil {
		conditions = append(conditions, fmt.Sprintf("index_status = $%d", argIndex))
		args = append(args, *filters.IndexStatus)
		argIndex++
	}

	if filters.IsIndexed != nil {
		conditions = append(conditions, fmt.Sprintf("is_indexed = $%d", argIndex))
		args = append(args, *filters.IsIndexed)
		argIndex++
	}

	if filters.IsMustHave != nil {
		conditions = append(conditions, fmt.Sprintf("is_must_have = $%d", argIndex))
		args = append(args, *filters.IsMustHave)
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
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM platforms WHERE %s", whereClause)
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (filters.Page - 1) * filters.PerPage
	args = append(args, filters.PerPage, offset)

	query := fmt.Sprintf(`
		SELECT id, user_id, url, domain, index_status, is_indexed, first_indexed_at,
		       last_checked_at, check_count, potential_score, is_must_have, notes, created_at, updated_at
		FROM platforms
		WHERE %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var platforms []model.Platform
	for rows.Next() {
		var p model.Platform
		err := rows.Scan(
			&p.ID, &p.UserID, &p.URL, &p.Domain, &p.IndexStatus,
			&p.IsIndexed, &p.FirstIndexedAt, &p.LastCheckedAt, &p.CheckCount,
			&p.PotentialScore, &p.IsMustHave, &p.Notes, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		platforms = append(platforms, p)
	}

	return platforms, total, nil
}

func (r *PlatformRepository) Update(ctx context.Context, id int64, req *model.UpdatePlatformRequest) (*model.Platform, error) {
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

	if req.PotentialScore != nil {
		setClauses = append(setClauses, fmt.Sprintf("potential_score = $%d", argIndex))
		args = append(args, *req.PotentialScore)
		argIndex++
	}

	if req.IsMustHave != nil {
		setClauses = append(setClauses, fmt.Sprintf("is_must_have = $%d", argIndex))
		args = append(args, *req.IsMustHave)
		argIndex++
	}

	if req.Notes != nil {
		setClauses = append(setClauses, fmt.Sprintf("notes = $%d", argIndex))
		args = append(args, *req.Notes)
		argIndex++
	}

	if len(setClauses) == 0 {
		return r.GetByID(ctx, id)
	}

	args = append(args, id)
	query := fmt.Sprintf(`
		UPDATE platforms SET %s
		WHERE id = $%d
		RETURNING id, user_id, url, domain, index_status, is_indexed, first_indexed_at,
		          last_checked_at, check_count, potential_score, is_must_have, notes, created_at, updated_at
	`, strings.Join(setClauses, ", "), argIndex)

	var platform model.Platform
	err := r.db.QueryRow(ctx, query, args...).Scan(
		&platform.ID, &platform.UserID, &platform.URL, &platform.Domain, &platform.IndexStatus,
		&platform.IsIndexed, &platform.FirstIndexedAt, &platform.LastCheckedAt, &platform.CheckCount,
		&platform.PotentialScore, &platform.IsMustHave, &platform.Notes, &platform.CreatedAt, &platform.UpdatedAt,
	)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &platform, nil
}

func (r *PlatformRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.Exec(ctx, "DELETE FROM platforms WHERE id = $1", id)
	return err
}

func (r *PlatformRepository) UpdateIndexStatus(ctx context.Context, id int64, status model.IndexStatus, isIndexed bool) error {
	query := `
		UPDATE platforms
		SET index_status = $1, is_indexed = $2, last_checked_at = NOW(), check_count = check_count + 1
		WHERE id = $3
	`
	if isIndexed {
		query = `
			UPDATE platforms
			SET index_status = $1, is_indexed = $2, last_checked_at = NOW(), check_count = check_count + 1,
			    first_indexed_at = COALESCE(first_indexed_at, NOW())
			WHERE id = $3
		`
	}
	_, err := r.db.Exec(ctx, query, status, isIndexed, id)
	return err
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
