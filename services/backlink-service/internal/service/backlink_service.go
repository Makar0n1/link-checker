package service

import (
	"context"
	"errors"

	"github.com/link-tracker/backlink-service/internal/model"
	"github.com/link-tracker/backlink-service/internal/repository"
)

var (
	ErrUnauthorized    = errors.New("unauthorized access to resource")
	ErrValidation      = errors.New("validation error")
	ErrProjectRequired = errors.New("project_id is required")
)

type BacklinkService struct {
	backlinkRepo *repository.BacklinkRepository
	projectRepo  *repository.ProjectRepository
}

func NewBacklinkService(
	backlinkRepo *repository.BacklinkRepository,
	projectRepo *repository.ProjectRepository,
) *BacklinkService {
	return &BacklinkService{
		backlinkRepo: backlinkRepo,
		projectRepo:  projectRepo,
	}
}

func (s *BacklinkService) Create(ctx context.Context, userID int64, req *model.CreateBacklinkRequest) (*model.Backlink, error) {
	// Verify project ownership
	isOwner, err := s.projectRepo.IsOwner(ctx, req.ProjectID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, ErrUnauthorized
	}

	backlink := &model.Backlink{
		ProjectID:  req.ProjectID,
		SourceURL:  req.SourceURL,
		TargetURL:  req.TargetURL,
		AnchorText: req.AnchorText,
		Status:     model.LinkStatusPending,
		LinkType:   req.LinkType,
	}

	if backlink.LinkType == "" {
		backlink.LinkType = model.LinkTypeDoFollow
	}

	if err := s.backlinkRepo.Create(ctx, backlink); err != nil {
		return nil, err
	}

	return backlink, nil
}

func (s *BacklinkService) GetByID(ctx context.Context, userID, backlinkID int64) (*model.Backlink, error) {
	backlink, err := s.backlinkRepo.GetByID(ctx, backlinkID)
	if err != nil {
		return nil, err
	}

	// Verify project ownership
	isOwner, err := s.projectRepo.IsOwner(ctx, backlink.ProjectID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, ErrUnauthorized
	}

	return backlink, nil
}

func (s *BacklinkService) List(ctx context.Context, userID int64, filters *model.BacklinkFilters) ([]*model.Backlink, int64, error) {
	// If project_id is specified, verify ownership
	if filters.ProjectID != nil {
		isOwner, err := s.projectRepo.IsOwner(ctx, *filters.ProjectID, userID)
		if err != nil {
			return nil, 0, err
		}
		if !isOwner {
			return nil, 0, ErrUnauthorized
		}
	} else {
		// Get all user's projects and filter by them
		projects, err := s.projectRepo.GetByUserID(ctx, userID)
		if err != nil {
			return nil, 0, err
		}
		if len(projects) == 0 {
			return []*model.Backlink{}, 0, nil
		}
		// For simplicity, require project_id filter
		return nil, 0, ErrProjectRequired
	}

	// Set defaults
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PerPage < 1 || filters.PerPage > 100 {
		filters.PerPage = 20
	}

	return s.backlinkRepo.List(ctx, filters)
}

func (s *BacklinkService) Update(ctx context.Context, userID, backlinkID int64, req *model.UpdateBacklinkRequest) (*model.Backlink, error) {
	backlink, err := s.backlinkRepo.GetByID(ctx, backlinkID)
	if err != nil {
		return nil, err
	}

	// Verify project ownership
	isOwner, err := s.projectRepo.IsOwner(ctx, backlink.ProjectID, userID)
	if err != nil {
		return nil, err
	}
	if !isOwner {
		return nil, ErrUnauthorized
	}

	// Apply updates
	if req.SourceURL != nil {
		backlink.SourceURL = *req.SourceURL
	}
	if req.TargetURL != nil {
		backlink.TargetURL = *req.TargetURL
	}
	if req.AnchorText != nil {
		backlink.AnchorText = *req.AnchorText
	}
	if req.Status != nil {
		backlink.Status = *req.Status
	}
	if req.LinkType != nil {
		backlink.LinkType = *req.LinkType
	}

	if err := s.backlinkRepo.Update(ctx, backlink); err != nil {
		return nil, err
	}

	return backlink, nil
}

func (s *BacklinkService) Delete(ctx context.Context, userID, backlinkID int64) error {
	projectID, err := s.backlinkRepo.GetProjectID(ctx, backlinkID)
	if err != nil {
		return err
	}

	// Verify project ownership
	isOwner, err := s.projectRepo.IsOwner(ctx, projectID, userID)
	if err != nil {
		return err
	}
	if !isOwner {
		return ErrUnauthorized
	}

	return s.backlinkRepo.Delete(ctx, backlinkID)
}

func (s *BacklinkService) BulkCreate(ctx context.Context, userID int64, req *model.BulkCreateBacklinksRequest) (*model.BulkOperationResponse, error) {
	if len(req.Backlinks) == 0 {
		return &model.BulkOperationResponse{Success: 0, Failed: 0}, nil
	}

	// Verify ownership for all projects
	projectIDs := make(map[int64]bool)
	for _, b := range req.Backlinks {
		projectIDs[b.ProjectID] = true
	}

	for projectID := range projectIDs {
		isOwner, err := s.projectRepo.IsOwner(ctx, projectID, userID)
		if err != nil {
			return nil, err
		}
		if !isOwner {
			return nil, ErrUnauthorized
		}
	}

	// Create backlinks
	backlinks := make([]*model.Backlink, len(req.Backlinks))
	for i, b := range req.Backlinks {
		linkType := b.LinkType
		if linkType == "" {
			linkType = model.LinkTypeDoFollow
		}
		backlinks[i] = &model.Backlink{
			ProjectID:  b.ProjectID,
			SourceURL:  b.SourceURL,
			TargetURL:  b.TargetURL,
			AnchorText: b.AnchorText,
			Status:     model.LinkStatusPending,
			LinkType:   linkType,
		}
	}

	success, errs := s.backlinkRepo.BulkCreate(ctx, backlinks)

	var errStrings []string
	for _, e := range errs {
		errStrings = append(errStrings, e.Error())
	}

	return &model.BulkOperationResponse{
		Success: success,
		Failed:  len(errs),
		Errors:  errStrings,
	}, nil
}

func (s *BacklinkService) BulkDelete(ctx context.Context, userID int64, req *model.BulkDeleteBacklinksRequest) (*model.BulkOperationResponse, error) {
	if len(req.IDs) == 0 {
		return &model.BulkOperationResponse{Success: 0, Failed: 0}, nil
	}

	// Verify ownership for all backlinks
	for _, id := range req.IDs {
		projectID, err := s.backlinkRepo.GetProjectID(ctx, id)
		if err != nil {
			if errors.Is(err, repository.ErrBacklinkNotFound) {
				continue
			}
			return nil, err
		}

		isOwner, err := s.projectRepo.IsOwner(ctx, projectID, userID)
		if err != nil {
			return nil, err
		}
		if !isOwner {
			return nil, ErrUnauthorized
		}
	}

	deleted, err := s.backlinkRepo.BulkDelete(ctx, req.IDs)
	if err != nil {
		return nil, err
	}

	return &model.BulkOperationResponse{
		Success: deleted,
		Failed:  len(req.IDs) - deleted,
	}, nil
}
