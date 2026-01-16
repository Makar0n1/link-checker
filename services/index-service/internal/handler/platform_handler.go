package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/link-tracker/index-service/internal/model"
	"github.com/link-tracker/index-service/internal/service"
	"github.com/link-tracker/shared/pkg/models"
	"github.com/link-tracker/shared/pkg/response"
)

type PlatformHandler struct {
	service *service.PlatformService
}

func NewPlatformHandler(service *service.PlatformService) *PlatformHandler {
	return &PlatformHandler{service: service}
}

func (h *PlatformHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*models.Claims)

	filters := &model.PlatformFilters{
		Page:    1,
		PerPage: 20,
	}

	if page := r.URL.Query().Get("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil {
			filters.Page = p
		}
	}
	if perPage := r.URL.Query().Get("per_page"); perPage != "" {
		if pp, err := strconv.Atoi(perPage); err == nil {
			filters.PerPage = pp
		}
	}
	if indexStatus := r.URL.Query().Get("index_status"); indexStatus != "" {
		status := model.IndexStatus(indexStatus)
		filters.IndexStatus = &status
	}
	if isIndexed := r.URL.Query().Get("is_indexed"); isIndexed != "" {
		b := isIndexed == "true"
		filters.IsIndexed = &b
	}
	if isMustHave := r.URL.Query().Get("is_must_have"); isMustHave != "" {
		b := isMustHave == "true"
		filters.IsMustHave = &b
	}
	if domain := r.URL.Query().Get("domain"); domain != "" {
		filters.Domain = domain
	}

	platforms, total, err := h.service.List(r.Context(), claims.UserID, filters)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Paginated(w, platforms, filters.Page, filters.PerPage, total)
}

func (h *PlatformHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*models.Claims)

	var req model.CreatePlatformRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.URL == "" {
		response.Error(w, http.StatusBadRequest, "url is required")
		return
	}

	platform, err := h.service.Create(r.Context(), claims.UserID, &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(w, platform)
}

func (h *PlatformHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*models.Claims)

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid platform id")
		return
	}

	platform, err := h.service.GetByID(r.Context(), claims.UserID, id)
	if err != nil {
		switch err {
		case service.ErrPlatformNotFound:
			response.Error(w, http.StatusNotFound, err.Error())
		case service.ErrNotOwner:
			response.Error(w, http.StatusForbidden, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.JSON(w, http.StatusOK, platform)
}

func (h *PlatformHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*models.Claims)

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid platform id")
		return
	}

	var req model.UpdatePlatformRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	platform, err := h.service.Update(r.Context(), claims.UserID, id, &req)
	if err != nil {
		switch err {
		case service.ErrPlatformNotFound:
			response.Error(w, http.StatusNotFound, err.Error())
		case service.ErrNotOwner:
			response.Error(w, http.StatusForbidden, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.JSON(w, http.StatusOK, platform)
}

func (h *PlatformHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*models.Claims)

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid platform id")
		return
	}

	err = h.service.Delete(r.Context(), claims.UserID, id)
	if err != nil {
		switch err {
		case service.ErrPlatformNotFound:
			response.Error(w, http.StatusNotFound, err.Error())
		case service.ErrNotOwner:
			response.Error(w, http.StatusForbidden, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.NoContent(w)
}

func (h *PlatformHandler) BulkCreate(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*models.Claims)

	var req model.BulkCreatePlatformsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.Platforms) == 0 {
		response.Error(w, http.StatusBadRequest, "platforms array is required")
		return
	}

	result := h.service.BulkCreate(r.Context(), claims.UserID, &req)
	response.JSON(w, http.StatusOK, result)
}

func (h *PlatformHandler) CheckIndex(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*models.Claims)

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid platform id")
		return
	}

	result, err := h.service.CheckIndex(r.Context(), claims.UserID, id)
	if err != nil {
		switch err {
		case service.ErrPlatformNotFound:
			response.Error(w, http.StatusNotFound, err.Error())
		case service.ErrNotOwner:
			response.Error(w, http.StatusForbidden, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.JSON(w, http.StatusOK, result)
}
