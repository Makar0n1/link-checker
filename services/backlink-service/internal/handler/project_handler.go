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

type ProjectHandler struct {
	projectService *service.ProjectService
}

func NewProjectHandler(projectService *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

// List handles GET /api/v1/projects
func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return
	}

	projects, err := h.projectService.List(r.Context(), userID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to list projects", "INTERNAL_ERROR")
		return
	}

	data := make([]model.ProjectResponse, len(projects))
	for i, p := range projects {
		data[i] = model.ProjectToResponse(p)
	}

	response.JSON(w, http.StatusOK, data)
}

// Create handles POST /api/v1/projects
func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return
	}

	var req model.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", "INVALID_REQUEST")
		return
	}

	if req.Name == "" {
		response.Error(w, http.StatusBadRequest, "name is required", "VALIDATION_ERROR")
		return
	}

	project, err := h.projectService.Create(r.Context(), userID, &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "failed to create project", "INTERNAL_ERROR")
		return
	}

	response.Created(w, model.ProjectToResponse(project))
}

// Get handles GET /api/v1/projects/:id
func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid project id", "INVALID_ID")
		return
	}

	project, err := h.projectService.GetByID(r.Context(), userID, id)
	if err != nil {
		if errors.Is(err, repository.ErrProjectNotFound) {
			response.Error(w, http.StatusNotFound, "project not found", "NOT_FOUND")
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			response.Error(w, http.StatusForbidden, "access denied", "FORBIDDEN")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to get project", "INTERNAL_ERROR")
		return
	}

	response.JSON(w, http.StatusOK, model.ProjectToResponse(project))
}

// Update handles PUT /api/v1/projects/:id
func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid project id", "INVALID_ID")
		return
	}

	var req model.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body", "INVALID_REQUEST")
		return
	}

	project, err := h.projectService.Update(r.Context(), userID, id, &req)
	if err != nil {
		if errors.Is(err, repository.ErrProjectNotFound) {
			response.Error(w, http.StatusNotFound, "project not found", "NOT_FOUND")
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			response.Error(w, http.StatusForbidden, "access denied", "FORBIDDEN")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to update project", "INTERNAL_ERROR")
		return
	}

	response.JSON(w, http.StatusOK, model.ProjectToResponse(project))
}

// Delete handles DELETE /api/v1/projects/:id
func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Error(w, http.StatusUnauthorized, "unauthorized", "UNAUTHORIZED")
		return
	}

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid project id", "INVALID_ID")
		return
	}

	err = h.projectService.Delete(r.Context(), userID, id)
	if err != nil {
		if errors.Is(err, repository.ErrProjectNotFound) {
			response.Error(w, http.StatusNotFound, "project not found", "NOT_FOUND")
			return
		}
		if errors.Is(err, service.ErrUnauthorized) {
			response.Error(w, http.StatusForbidden, "access denied", "FORBIDDEN")
			return
		}
		response.Error(w, http.StatusInternalServerError, "failed to delete project", "INTERNAL_ERROR")
		return
	}

	response.NoContent(w)
}
