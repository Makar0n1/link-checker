package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/link-tracker/backlink-service/internal/model"
)

var (
	ErrBacklinkNotFound = errors.New("backlink not found")
)

type BacklinkRepository struct {
	db *pgxpool.Pool
}

func NewBacklinkRepository(db *pgxpool.Pool) *BacklinkRepository {
	return &BacklinkRepository{db: db}
}

func (r *BacklinkRepository) Create(ctx context.Context, backlink *model.Backlink) error {
	query := `
		INSERT INTO backlinks (project_id, source_url, target_url, anchor_text, status, link_type)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query,
		backlink.ProjectID,
		backlink.SourceURL,
		backlink.TargetURL,
		backlink.AnchorText,
		backlink.Status,
		backlink.LinkType,
	).Scan(&backlink.ID, &backlink.CreatedAt)

	return err
}

func (r *BacklinkRepository) GetByID(ctx context.Context, id int64) (*model.Backlink, error) {
	query := `
		SELECT id, project_id, source_url, target_url, anchor_text, status, link_type, http_status, last_checked_at, created_at
		FROM backlinks
		WHERE id = $1
	`

	backlink := &model.Backlink{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&backlink.ID,
		&backlink.ProjectID,
		&backlink.SourceURL,
		&backlink.TargetURL,
		&backlink.AnchorText,
		&backlink.Status,
		&backlink.LinkType,
		&backlink.HTTPStatus,
		&backlink.LastCheckedAt,
		&backlink.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrBacklinkNotFound
		}
		return nil, err
	}

	return backlink, nil
}

func (r *BacklinkRepository) List(ctx context.Context, filters *model.BacklinkFilters) ([]*model.Backlink, int64, error) {
	var conditions []string
	var args []interface{}
	argNum := 1

	if filters.ProjectID != nil {
		conditions = append(conditions, fmt.Sprintf("project_id = $%d", argNum))
		args = append(args, *filters.ProjectID)
		argNum++
	}

	if filters.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argNum))
		args = append(args, *filters.Status)
		argNum++
	}

	if filters.LinkType != nil {
		conditions = append(conditions, fmt.Sprintf("link_type = $%d", argNum))
		args = append(args, *filters.LinkType)
		argNum++
	}

	if filters.SourceURL != nil {
		conditions = append(conditions, fmt.Sprintf("source_url ILIKE $%d", argNum))
		args = append(args, "%"+*filters.SourceURL+"%")
		argNum++
	}

	if filters.TargetURL != nil {
		conditions = append(conditions, fmt.Sprintf("target_url ILIKE $%d", argNum))
		args = append(args, "%"+*filters.TargetURL+"%")
		argNum++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM backlinks %s", whereClause)
	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (filters.Page - 1) * filters.PerPage
	query := fmt.Sprintf(`
		SELECT id, project_id, source_url, target_url, anchor_text, status, link_type, http_status, last_checked_at, created_at
		FROM backlinks
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argNum, argNum+1)

	args = append(args, filters.PerPage, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var backlinks []*model.Backlink
	for rows.Next() {
		backlink := &model.Backlink{}
		err := rows.Scan(
			&backlink.ID,
			&backlink.ProjectID,
			&backlink.SourceURL,
			&backlink.TargetURL,
			&backlink.AnchorText,
			&backlink.Status,
			&backlink.LinkType,
			&backlink.HTTPStatus,
			&backlink.LastCheckedAt,
			&backlink.CreatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		backlinks = append(backlinks, backlink)
	}

	return backlinks, total, rows.Err()
}

func (r *BacklinkRepository) Update(ctx context.Context, backlink *model.Backlink) error {
	query := `
		UPDATE backlinks
		SET source_url = $1, target_url = $2, anchor_text = $3, status = $4, link_type = $5
		WHERE id = $6
	`

	result, err := r.db.Exec(ctx, query,
		backlink.SourceURL,
		backlink.TargetURL,
		backlink.AnchorText,
		backlink.Status,
		backlink.LinkType,
		backlink.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrBacklinkNotFound
	}

	return nil
}

func (r *BacklinkRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM backlinks WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrBacklinkNotFound
	}

	return nil
}

func (r *BacklinkRepository) BulkCreate(ctx context.Context, backlinks []*model.Backlink) (int, []error) {
	var errs []error
	success := 0

	for _, backlink := range backlinks {
		if err := r.Create(ctx, backlink); err != nil {
			errs = append(errs, err)
		} else {
			success++
		}
	}

	return success, errs
}

func (r *BacklinkRepository) BulkDelete(ctx context.Context, ids []int64) (int, error) {
	if len(ids) == 0 {
		return 0, nil
	}

	query := `DELETE FROM backlinks WHERE id = ANY($1)`
	result, err := r.db.Exec(ctx, query, ids)
	if err != nil {
		return 0, err
	}

	return int(result.RowsAffected()), nil
}

func (r *BacklinkRepository) GetProjectID(ctx context.Context, backlinkID int64) (int64, error) {
	query := `SELECT project_id FROM backlinks WHERE id = $1`
	var projectID int64
	err := r.db.QueryRow(ctx, query, backlinkID).Scan(&projectID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrBacklinkNotFound
		}
		return 0, err
	}
	return projectID, nil
}
