package service

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/link-tracker/health-service/internal/model"
	"github.com/link-tracker/health-service/internal/repository"
)

var (
	ErrSiteNotFound = errors.New("site not found")
	ErrNotOwner     = errors.New("not site owner")
)

type SiteService struct {
	repo       *repository.SiteRepository
	httpClient *http.Client
}

func NewSiteService(repo *repository.SiteRepository) *SiteService {
	return &SiteService{
		repo: repo,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (s *SiteService) Create(ctx context.Context, userID int64, req *model.CreateSiteRequest) (*model.MonitoredSite, error) {
	return s.repo.Create(ctx, userID, req)
}

func (s *SiteService) GetByID(ctx context.Context, userID, siteID int64) (*model.MonitoredSite, error) {
	site, err := s.repo.GetByID(ctx, siteID)
	if err != nil {
		return nil, err
	}
	if site == nil {
		return nil, ErrSiteNotFound
	}
	if site.UserID != userID {
		return nil, ErrNotOwner
	}
	return site, nil
}

func (s *SiteService) List(ctx context.Context, userID int64, filters *model.SiteFilters) ([]model.MonitoredSite, int64, error) {
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PerPage < 1 || filters.PerPage > 100 {
		filters.PerPage = 20
	}
	return s.repo.List(ctx, userID, filters)
}

func (s *SiteService) Update(ctx context.Context, userID, siteID int64, req *model.UpdateSiteRequest) (*model.MonitoredSite, error) {
	site, err := s.repo.GetByID(ctx, siteID)
	if err != nil {
		return nil, err
	}
	if site == nil {
		return nil, ErrSiteNotFound
	}
	if site.UserID != userID {
		return nil, ErrNotOwner
	}
	return s.repo.Update(ctx, siteID, req)
}

func (s *SiteService) Delete(ctx context.Context, userID, siteID int64) error {
	site, err := s.repo.GetByID(ctx, siteID)
	if err != nil {
		return err
	}
	if site == nil {
		return ErrSiteNotFound
	}
	if site.UserID != userID {
		return ErrNotOwner
	}
	return s.repo.Delete(ctx, siteID)
}

func (s *SiteService) CheckHealth(ctx context.Context, userID, siteID int64) (*model.SiteHealthCheck, error) {
	site, err := s.repo.GetByID(ctx, siteID)
	if err != nil {
		return nil, err
	}
	if site == nil {
		return nil, ErrSiteNotFound
	}
	if site.UserID != userID {
		return nil, ErrNotOwner
	}

	result := &model.SiteHealthCheck{
		SiteID:    siteID,
		URL:       site.URL,
		CheckedAt: time.Now(),
	}

	// 1. HTTP GET request to main URL
	startTime := time.Now()
	req, err := http.NewRequestWithContext(ctx, "GET", site.URL, nil)
	if err != nil {
		result.Error = err.Error()
		_ = s.repo.UpdateHealthCheck(ctx, siteID, result)
		_ = s.repo.AddCheckHistory(ctx, siteID, result)
		return result, nil
	}
	req.Header.Set("User-Agent", "LinkTracker/1.0")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		result.Error = err.Error()
		_ = s.repo.UpdateHealthCheck(ctx, siteID, result)
		_ = s.repo.AddCheckHistory(ctx, siteID, result)
		return result, nil
	}
	defer resp.Body.Close()

	result.ResponseTimeMs = int(time.Since(startTime).Milliseconds())
	result.HTTPStatus = resp.StatusCode
	result.IsAlive = resp.StatusCode >= 200 && resp.StatusCode < 400

	// 2. Read HTML to check for noindex
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024*1024)) // Max 1MB
	bodyStr := strings.ToLower(string(body))
	result.HasNoindex = strings.Contains(bodyStr, `name="robots"`) && strings.Contains(bodyStr, "noindex")

	// 3. Check robots.txt
	robotsURL := s.buildRobotsURL(site.URL)
	robotsReq, _ := http.NewRequestWithContext(ctx, "GET", robotsURL, nil)
	if robotsReq != nil {
		robotsReq.Header.Set("User-Agent", "LinkTracker/1.0")
		robotsResp, err := s.httpClient.Do(robotsReq)
		if err != nil {
			result.RobotsTxtStatus = "error"
		} else {
			defer robotsResp.Body.Close()
			if robotsResp.StatusCode == http.StatusOK {
				robotsBody, _ := io.ReadAll(io.LimitReader(robotsResp.Body, 64*1024))
				robotsContent := strings.ToLower(string(robotsBody))

				if strings.Contains(robotsContent, "disallow: /") {
					result.RobotsTxtStatus = "disallow"
					result.AllowsIndexing = false
				} else {
					result.RobotsTxtStatus = "allow"
					result.AllowsIndexing = true
				}
			} else if robotsResp.StatusCode == http.StatusNotFound {
				result.RobotsTxtStatus = "not_found"
				result.AllowsIndexing = true
			} else {
				result.RobotsTxtStatus = "error"
			}
		}
	}

	// Save results
	_ = s.repo.UpdateHealthCheck(ctx, siteID, result)
	_ = s.repo.AddCheckHistory(ctx, siteID, result)

	return result, nil
}

func (s *SiteService) GetHistory(ctx context.Context, userID, siteID int64, filters *model.HistoryFilters) ([]model.SiteCheckHistory, int64, error) {
	site, err := s.repo.GetByID(ctx, siteID)
	if err != nil {
		return nil, 0, err
	}
	if site == nil {
		return nil, 0, ErrSiteNotFound
	}
	if site.UserID != userID {
		return nil, 0, ErrNotOwner
	}

	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PerPage < 1 || filters.PerPage > 100 {
		filters.PerPage = 20
	}

	return s.repo.GetHistory(ctx, siteID, filters)
}

func (s *SiteService) buildRobotsURL(siteURL string) string {
	if !strings.HasPrefix(siteURL, "http://") && !strings.HasPrefix(siteURL, "https://") {
		siteURL = "https://" + siteURL
	}

	// Extract base URL
	idx := strings.Index(siteURL[8:], "/")
	if idx == -1 {
		return siteURL + "/robots.txt"
	}
	return siteURL[:8+idx] + "/robots.txt"
}
