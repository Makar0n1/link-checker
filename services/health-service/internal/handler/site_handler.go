package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/link-tracker/health-service/internal/model"
	"github.com/link-tracker/health-service/internal/service"
	"github.com/link-tracker/shared/pkg/models"
	"github.com/link-tracker/shared/pkg/response"
)

type SiteHandler struct {
	service *service.SiteService
}

func NewSiteHandler(service *service.SiteService) *SiteHandler {
	return &SiteHandler{service: service}
}

func (h *SiteHandler) List(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*models.Claims)

	filters := &model.SiteFilters{
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
	if isAlive := r.URL.Query().Get("is_alive"); isAlive != "" {
		b := isAlive == "true"
		filters.IsAlive = &b
	}
	if hasNoindex := r.URL.Query().Get("has_noindex"); hasNoindex != "" {
		b := hasNoindex == "true"
		filters.HasNoindex = &b
	}
	if domain := r.URL.Query().Get("domain"); domain != "" {
		filters.Domain = domain
	}

	sites, total, err := h.service.List(r.Context(), claims.UserID, filters)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Paginated(w, sites, filters.Page, filters.PerPage, total)
}

func (h *SiteHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*models.Claims)

	var req model.CreateSiteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.URL == "" {
		response.Error(w, http.StatusBadRequest, "url is required")
		return
	}

	site, err := h.service.Create(r.Context(), claims.UserID, &req)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.Created(w, site)
}

func (h *SiteHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*models.Claims)

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid site id")
		return
	}

	site, err := h.service.GetByID(r.Context(), claims.UserID, id)
	if err != nil {
		switch err {
		case service.ErrSiteNotFound:
			response.Error(w, http.StatusNotFound, err.Error())
		case service.ErrNotOwner:
			response.Error(w, http.StatusForbidden, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.JSON(w, http.StatusOK, site)
}

func (h *SiteHandler) Update(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*models.Claims)

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid site id")
		return
	}

	var req model.UpdateSiteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	site, err := h.service.Update(r.Context(), claims.UserID, id, &req)
	if err != nil {
		switch err {
		case service.ErrSiteNotFound:
			response.Error(w, http.StatusNotFound, err.Error())
		case service.ErrNotOwner:
			response.Error(w, http.StatusForbidden, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.JSON(w, http.StatusOK, site)
}

func (h *SiteHandler) Delete(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*models.Claims)

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid site id")
		return
	}

	err = h.service.Delete(r.Context(), claims.UserID, id)
	if err != nil {
		switch err {
		case service.ErrSiteNotFound:
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

func (h *SiteHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*models.Claims)

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid site id")
		return
	}

	result, err := h.service.CheckHealth(r.Context(), claims.UserID, id)
	if err != nil {
		switch err {
		case service.ErrSiteNotFound:
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

func (h *SiteHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*models.Claims)

	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		response.Error(w, http.StatusBadRequest, "invalid site id")
		return
	}

	filters := &model.HistoryFilters{
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

	history, total, err := h.service.GetHistory(r.Context(), claims.UserID, id, filters)
	if err != nil {
		switch err {
		case service.ErrSiteNotFound:
			response.Error(w, http.StatusNotFound, err.Error())
		case service.ErrNotOwner:
			response.Error(w, http.StatusForbidden, err.Error())
		default:
			response.Error(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.Paginated(w, history, filters.Page, filters.PerPage, total)
}
