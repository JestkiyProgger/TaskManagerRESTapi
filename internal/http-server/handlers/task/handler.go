package task

import (
	taskDomain "ProjectManagementAPI/internal/domain/task"
	resp "ProjectManagementAPI/internal/lib/api/response"
	"ProjectManagementAPI/internal/lib/logger/sl"
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, title, description, status string, assignees []uuid.UUID) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*taskDomain.Task, error)
}

type Handler struct {
	log     *slog.Logger
	service Service
}

func NewHandler(log *slog.Logger, service Service) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}

type CreateRequest struct {
	Title       string   `json:"title" validate:"required"`
	Description string   `json:"description" validate:"required"`
	Status      string   `json:"status" validate:"required"`
	Assignees   []string `json:"assignees" validate:"required,min=1,dive,uuid4"`
}

type CreateResponse struct {
	resp.Response
	ID string `json:"id"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "handlers/task.Create"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	var req CreateRequest

	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("decode error", sl.Err(err))
		render.JSON(w, r, resp.Error("invalid request"))
		return
	}

	if err := validator.New().Struct(req); err != nil {
		render.JSON(w, r, resp.Error("validation error"))
		return
	}

	var assigneeUUIDs []uuid.UUID
	for _, a := range req.Assignees {
		id, err := uuid.Parse(a)
		if err != nil {
			render.JSON(w, r, resp.Error("invalid assignee UUID: "+a))
			return
		}
		assigneeUUIDs = append(assigneeUUIDs, id)
	}

	id, err := h.service.Create(r.Context(), req.Title, req.Description, req.Status, assigneeUUIDs)

	if errors.Is(err, taskDomain.ErrNoAssignees) {
		render.JSON(w, r, resp.Error("task must have at least one assignee"))
		return
	}

	if errors.Is(err, taskDomain.ErrInvalidTitle) {
		render.JSON(w, r, resp.Error("invalid title"))
		return
	}

	if err != nil {
		log.Error("create failed", sl.Err(err))
		render.JSON(w, r, resp.Error("failed to create task"))
		return
	}

	render.JSON(w, r, CreateResponse{
		Response: resp.OK(),
		ID:       id.String(),
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "handlers/task.Delete"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		render.JSON(w, r, resp.Error("invalid id"))
		return
	}

	err = h.service.Delete(r.Context(), id)

	if errors.Is(err, taskDomain.ErrTaskNotFound) {
		render.JSON(w, r, resp.Error("task not found"))
		return
	}

	if err != nil {
		log.Error("delete failed", sl.Err(err))
		render.JSON(w, r, resp.Error("failed to delete"))
		return
	}

	render.JSON(w, r, resp.OK())
}

type GetByIDResponse struct {
	resp.Response
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Status      string   `json:"status"`
	Assignees   []string `json:"assignees"`
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	const op = "handlers/task.GetByID"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)

	idStr := chi.URLParam(r, "id")

	id, err := uuid.Parse(idStr)
	if err != nil {
		render.JSON(w, r, resp.Error("invalid id"))
		return
	}

	task, err := h.service.GetByID(r.Context(), id)

	if errors.Is(err, taskDomain.ErrTaskNotFound) {
		log.Info("task not found", slog.String("task_id", id.String()))
		render.JSON(w, r, resp.Error("task not found"))
		return
	}

	if err != nil {
		log.Error("get failed", sl.Err(err))
		render.JSON(w, r, resp.Error("failed to get"))
		return
	}

	assigneeIDs := make([]string, len(task.Assignees))
	for i, a := range task.Assignees {
		assigneeIDs[i] = a.String()
	}

	render.JSON(w, r, GetByIDResponse{
		Response:    resp.OK(),
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Assignees:   assigneeIDs,
	})
}
