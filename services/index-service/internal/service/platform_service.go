package service

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/link-tracker/index-service/internal/model"
	"github.com/link-tracker/index-service/internal/repository"
)

var (
	ErrPlatformNotFound = errors.New("platform not found")
	ErrNotOwner         = errors.New("not platform owner")
	ErrBulkLimitExceeded = errors.New("bulk operation limit exceeded (max 100)")
)

type PlatformService struct {
	repo       *repository.PlatformRepository
	httpClient *http.Client
}

func NewPlatformService(repo *repository.PlatformRepository) *PlatformService {
	return &PlatformService{
		repo: repo,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *PlatformService) Create(ctx context.Context, userID int64, req *model.CreatePlatformRequest) (*model.Platform, error) {
	return s.repo.Create(ctx, userID, req)
}

func (s *PlatformService) GetByID(ctx context.Context, userID, platformID int64) (*model.Platform, error) {
	platform, err := s.repo.GetByID(ctx, platformID)
	if err != nil {
		return nil, err
	}
	if platform == nil {
		return nil, ErrPlatformNotFound
	}
	if platform.UserID != userID {
		return nil, ErrNotOwner
	}
	return platform, nil
}

func (s *PlatformService) List(ctx context.Context, userID int64, filters *model.PlatformFilters) ([]model.Platform, int64, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PerPage < 1 || filters.PerPage > 100 {
		filters.PerPage = 20
	}
	return s.repo.List(ctx, userID, filters)
}

func (s *PlatformService) Update(ctx context.Context, userID, platformID int64, req *model.UpdatePlatformRequest) (*model.Platform, error) {
	platform, err := s.repo.GetByID(ctx, platformID)
	if err != nil {
		return nil, err
	}
	if platform == nil {
		return nil, ErrPlatformNotFound
	}
	if platform.UserID != userID {
		return nil, ErrNotOwner
	}
	return s.repo.Update(ctx, platformID, req)
}

func (s *PlatformService) Delete(ctx context.Context, userID, platformID int64) error {
	platform, err := s.repo.GetByID(ctx, platformID)
	if err != nil {
		return err
	}
	if platform == nil {
		return ErrPlatformNotFound
	}
	if platform.UserID != userID {
		return ErrNotOwner
	}
	return s.repo.Delete(ctx, platformID)
}

func (s *PlatformService) BulkCreate(ctx context.Context, userID int64, req *model.BulkCreatePlatformsRequest) *model.BulkOperationResponse {
	response := &model.BulkOperationResponse{
		Errors:  make([]model.BulkError, 0),
		Created: make([]model.Platform, 0),
	}

	if len(req.Platforms) > 100 {
		response.Failed = len(req.Platforms)
		response.Errors = append(response.Errors, model.BulkError{
			Index:   0,
			Message: ErrBulkLimitExceeded.Error(),
		})
		return response
	}

	for i, createReq := range req.Platforms {
		platform, err := s.repo.Create(ctx, userID, &createReq)
		if err != nil {
			response.Failed++
			response.Errors = append(response.Errors, model.BulkError{
				Index:   i,
				URL:     createReq.URL,
				Message: err.Error(),
			})
		} else {
			response.Success++
			response.Created = append(response.Created, *platform)
		}
	}

	return response
}

func (s *PlatformService) CheckIndex(ctx context.Context, userID, platformID int64) (*model.IndexCheckResult, error) {
	platform, err := s.repo.GetByID(ctx, platformID)
	if err != nil {
		return nil, err
	}
	if platform == nil {
		return nil, ErrPlatformNotFound
	}
	if platform.UserID != userID {
		return nil, ErrNotOwner
	}

	result := &model.IndexCheckResult{
		PlatformID: platformID,
		URL:        platform.URL,
		CheckedAt:  time.Now(),
	}

	// MVP: Simple HTTP check
	req, err := http.NewRequestWithContext(ctx, "GET", platform.URL, nil)
	if err != nil {
		result.IndexStatus = model.IndexStatusError
		result.Error = err.Error()
		_ = s.repo.UpdateIndexStatus(ctx, platformID, model.IndexStatusError, false)
		return result, nil
	}

	req.Header.Set("User-Agent", "LinkTracker/1.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		result.IndexStatus = model.IndexStatusError
		result.Error = err.Error()
		_ = s.repo.UpdateIndexStatus(ctx, platformID, model.IndexStatusError, false)
		return result, nil
	}
	defer resp.Body.Close()

	result.HTTPStatus = resp.StatusCode

	// If status is 200, mark as potentially indexed
	if resp.StatusCode == http.StatusOK {
		result.IsIndexed = true
		result.IndexStatus = model.IndexStatusIndexed
		_ = s.repo.UpdateIndexStatus(ctx, platformID, model.IndexStatusIndexed, true)
	} else if resp.StatusCode >= 400 {
		result.IsIndexed = false
		result.IndexStatus = model.IndexStatusNotIndexed
		_ = s.repo.UpdateIndexStatus(ctx, platformID, model.IndexStatusNotIndexed, false)
	} else {
		result.IsIndexed = false
		result.IndexStatus = model.IndexStatusPending
		_ = s.repo.UpdateIndexStatus(ctx, platformID, model.IndexStatusPending, false)
	}

	return result, nil
}
