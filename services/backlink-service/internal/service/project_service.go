package service

import (
	"context"

	"github.com/link-tracker/backlink-service/internal/model"
	"github.com/link-tracker/backlink-service/internal/repository"
)

type ProjectService struct {
	projectRepo *repository.ProjectRepository
}

func NewProjectService(projectRepo *repository.ProjectRepository) *ProjectService {
	return &ProjectService{projectRepo: projectRepo}
}

func (s *ProjectService) Create(ctx context.Context, userID int64, req *model.CreateProjectRequest) (*model.Project, error) {
	project := &model.Project{
		Name:          req.Name,
		UserID:        userID,
		GoogleSheetID: req.GoogleSheetID,
	}

	if err := s.projectRepo.Create(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) GetByID(ctx context.Context, userID, projectID int64) (*model.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if project.UserID != userID {
		return nil, ErrUnauthorized
	}

	return project, nil
}

func (s *ProjectService) List(ctx context.Context, userID int64) ([]*model.Project, error) {
	return s.projectRepo.GetByUserID(ctx, userID)
}

func (s *ProjectService) Update(ctx context.Context, userID, projectID int64, req *model.UpdateProjectRequest) (*model.Project, error) {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	if project.UserID != userID {
		return nil, ErrUnauthorized
	}

	if req.Name != nil {
		project.Name = *req.Name
	}
	if req.GoogleSheetID != nil {
		project.GoogleSheetID = req.GoogleSheetID
	}

	if err := s.projectRepo.Update(ctx, project); err != nil {
		return nil, err
	}

	return project, nil
}

func (s *ProjectService) Delete(ctx context.Context, userID, projectID int64) error {
	project, err := s.projectRepo.GetByID(ctx, projectID)
	if err != nil {
		return err
	}

	if project.UserID != userID {
		return ErrUnauthorized
	}

	return s.projectRepo.Delete(ctx, projectID)
}
