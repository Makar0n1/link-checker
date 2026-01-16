package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/link-tracker/backlink-service/internal/model"
	"github.com/link-tracker/backlink-service/internal/repository"
	"github.com/link-tracker/backlink-service/internal/service"
	"github.com/link-tracker/shared/pkg/middleware"
	"github.com/link-tracker/shared/pkg/response"
)

type BacklinkHandler struct {
	backlinkService *service.BacklinkService
}

func NewBacklinkHandler(backlinkService *service.BacklinkService) *BacklinkHandler {
	return &BacklinkHandler{backlinkService: backlinkService}
}

// List handles GET /api/v1/backlinks
func (h *BacklinkHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return
	}

	filters := &model.BacklinkFilters{
		Page:    1,
		PerPage: 20,
	}

	// Parse query params
	if v := r.URL.Query().Get("project_id"); v != "" {
		if id, err := strconv.ParseInt(v, 10, 64); err == nil {
			filters.ProjectID = &id
		}
	}
	if v := r.URL.Query().Get("status"); v != "" {
		status := model.LinkStatus(v)
		filters.Status = &status
	}
	if v := r.URL.Query().Get("link_type"); v != "" {
		linkType := model.LinkType(v)
		filters.LinkType = &linkType
	}
	if v := r.URL.Query().Get("source_url"); v != "" {
		filters.SourceURL = &v
	}
	if v := r.URL.Query().Get("target_url"); v != "" {
		filters.TargetURL = &v
	}
	if v := r.URL.Query().Get("page"); v != "" {
		if page, err := strconv.Atoi(v); err == nil && page > 0 {
			filters.Page = page
		}
	}
	if v := r.URL.Query().Get("per_page"); v != "" {
		if perPage, err := strconv.Atoi(v); err == nil && perPage > 0 {
			filters.PerPage = perPage
		}
	}

	backlinks, total, err := h.backlinkService.List(r.Context(), userID, filters)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			response.Error(w, http.StatusForbidden, "access denied", "FORBIDDEN")
			return
		}
		if errors.Is(err, service.ErrProjectRequired) {
			response.Error(w, http.StatusBadRequest, "project_id query parameter is required", "VALIDATION_ERROR")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to list backlinks", "INTERNAL_ERROR")
		return
	}

	data := make([]model.BacklinkResponse, len(backlinks))
	for i, b := range backlinks {
		data[i] = model.BacklinkToResponse(b)
	}

	response.Paginated(w, data, filters.Page, filters.PerPage, total)
}

// Create handles POST /api/v1/backlinks
func (h *BacklinkHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return
	}

	var req model.CreateBacklinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", "INVALID_REQUEST")
		return
	}

	if err := validateCreateBacklinkRequest(&req); err != nil {
		response.Error(w, http.StatusBadRequest, err.Error(), "VALIDATION_ERROR")
		return
	}

	backlink, err := h.backlinkService.Create(r.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			response.Error(w, http.StatusForbidden, "access denied to project", "FORBIDDEN")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to create backlink", "INTERNAL_ERROR")
		return
	}

	response.Created(w, model.BacklinkToResponse(backlink))
}

// Get handles GET /api/v1/backlinks/:id
func (h *BacklinkHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid backlink id", "INVALID_ID")
		return
	}

	backlink, err := h.backlinkService.GetByID(r.Context(), userID, id)
	if err != nil {
		if errors.Is(err, repository.ErrBacklinkNotFound) {
			response.Error(w, http.StatusNotFound, "backlink not found", "NOT_FOUND")
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			response.Error(w, http.StatusForbidden, "access denied", "FORBIDDEN")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to get backlink", "INTERNAL_ERROR")
		return
	}

	response.JSON(w, http.StatusOK, model.BacklinkToResponse(backlink))
}

// Update handles PUT /api/v1/backlinks/:id
func (h *BacklinkHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid backlink id", "INVALID_ID")
		return
	}

	var req model.UpdateBacklinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", "INVALID_REQUEST")
		return
	}

	backlink, err := h.backlinkService.Update(r.Context(), userID, id, &req)
	if err != nil {
		if errors.Is(err, repository.ErrBacklinkNotFound) {
			response.Error(w, http.StatusNotFound, "backlink not found", "NOT_FOUND")
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			response.Error(w, http.StatusForbidden, "access denied", "FORBIDDEN")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to update backlink", "INTERNAL_ERROR")
		return
	}

	response.JSON(w, http.StatusOK, model.BacklinkToResponse(backlink))
}

// Delete handles DELETE /api/v1/backlinks/:id
func (h *BacklinkHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid backlink id", "INVALID_ID")
		return
	}

	err = h.backlinkService.Delete(r.Context(), userID, id)
	if err != nil {
		if errors.Is(err, repository.ErrBacklinkNotFound) {
			response.Error(w, http.StatusNotFound, "backlink not found", "NOT_FOUND")
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			response.Error(w, http.StatusForbidden, "access denied", "FORBIDDEN")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to delete backlink", "INTERNAL_ERROR")
		return
	}

	response.NoContent(w)
}

// BulkCreate handles POST /api/v1/backlinks/bulk
func (h *BacklinkHandler) BulkCreate(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return
	}

	var req model.BulkCreateBacklinksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", "INVALID_REQUEST")
		return
	}

	if len(req.Backlinks) == 0 {
		response.Error(w, http.StatusBadRequest, "backlinks array is required", "VALIDATION_ERROR")
		return
	}

	if len(req.Backlinks) > 100 {
		response.Error(w, http.StatusBadRequest, "maximum 100 backlinks per request", "VALIDATION_ERROR")
		return
	}

	result, err := h.backlinkService.BulkCreate(r.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			response.Error(w, http.StatusForbidden, "access denied to one or more projects", "FORBIDDEN")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to create backlinks", "INTERNAL_ERROR")
		return
	}

	response.JSON(w, http.StatusCreated, result)
}

// BulkDelete handles DELETE /api/v1/backlinks/bulk
func (h *BacklinkHandler) BulkDelete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return
	}

	var req model.BulkDeleteBacklinksRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", "INVALID_REQUEST")
		return
	}

	if len(req.IDs) == 0 {
		response.Error(w, http.StatusBadRequest, "ids array is required", "VALIDATION_ERROR")
		return
	}

	if len(req.IDs) > 100 {
		response.Error(w, http.StatusBadRequest, "maximum 100 ids per request", "VALIDATION_ERROR")
		return
	}

	result, err := h.backlinkService.BulkDelete(r.Context(), userID, &req)
	if err != nil {
		if errors.Is(err, service.ErrUnauthorized) {
			response.Error(w, http.StatusForbidden, "access denied to one or more backlinks", "FORBIDDEN")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to delete backlinks", "INTERNAL_ERROR")
		return
	}

	response.JSON(w, http.StatusOK, result)
}

// Import handles POST /api/v1/backlinks/import
func (h *BacklinkHandler) Import(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return
	}

	var req model.ImportFromSheetsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", "INVALID_REQUEST")
		return
	}

	// TODO: Implement Google Sheets import
	// This would require Google Sheets API integration
	_ = userID
	response.Error(w, http.StatusNotImplemented, "Google Sheets import not yet implemented", "NOT_IMPLEMENTED")
}

func validateCreateBacklinkRequest(req *model.CreateBacklinkRequest) error {
	if req.ProjectID == 0 {
		return errors.New("project_id is required")
	}
	if req.SourceURL == "" {
		return errors.New("source_url is required")
	}
	if req.TargetURL == "" {
		return errors.New("target_url is required")
	}
	return nil
}
