package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/link-tracker/backlink-service/internal/model"
)

var (
	ErrProjectNotFound = errors.New("project not found")
)

type ProjectRepository struct {
	db *pgxpool.Pool
}

func NewProjectRepository(db *pgxpool.Pool) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, project *model.Project) error {
	query := `
		INSERT INTO projects (name, user_id, google_sheet_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(ctx, query,
		project.Name,
		project.UserID,
		project.GoogleSheetID,
	).Scan(&project.ID, &project.CreatedAt)

	return err
}

func (r *ProjectRepository) GetByID(ctx context.Context, id int64) (*model.Project, error) {
	query := `
		SELECT id, name, user_id, google_sheet_id, created_at
		FROM projects
		WHERE id = $1
	`

	project := &model.Project{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&project.ID,
		&project.Name,
		&project.UserID,
		&project.GoogleSheetID,
		&project.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	return project, nil
}

func (r *ProjectRepository) GetByUserID(ctx context.Context, userID int64) ([]*model.Project, error) {
	query := `
		SELECT id, name, user_id, google_sheet_id, created_at
		FROM projects
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*model.Project
	for rows.Next() {
		project := &model.Project{}
		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.UserID,
			&project.GoogleSheetID,
			&project.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	return projects, rows.Err()
}

func (r *ProjectRepository) Update(ctx context.Context, project *model.Project) error {
	query := `
		UPDATE projects
		SET name = $1, google_sheet_id = $2
		WHERE id = $3
	`

	result, err := r.db.Exec(ctx, query,
		project.Name,
		project.GoogleSheetID,
		project.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrProjectNotFound
	}

	return nil
}

func (r *ProjectRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM projects WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrProjectNotFound
	}

	return nil
}

func (r *ProjectRepository) IsOwner(ctx context.Context, projectID, userID int64) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1 AND user_id = $2)`
	var exists bool
	err := r.db.QueryRow(ctx, query, projectID, userID).Scan(&exists)
	return exists, err
}
