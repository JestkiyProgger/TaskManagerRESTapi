package user

import (
	userDomain "ProjectManagementAPI/internal/domain/user"
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
	Create(ctx context.Context, email, name string) (uuid.UUID, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (*userDomain.User, error)
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
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required"`
}

type CreateResponse struct {
	resp.Response
	ID string `json:"id"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "handlers/user.Create"

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
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)
		render.JSON(w, r, resp.ValidationError(validateErr))
		return
	}

	id, err := h.service.Create(r.Context(), req.Email, req.Name)

	if errors.Is(err, userDomain.ErrEmailAlreadyExists) {
		log.Info("email already exists", slog.String("email", req.Email))
		render.JSON(w, r, resp.Error("email already exists"))
		return
	}

	if err != nil {
		log.Error("create failed", sl.Err(err))
		render.JSON(w, r, resp.Error("failed to create user"))
		return
	}

	render.JSON(w, r, CreateResponse{
		Response: resp.OK(),
		ID:       id.String(),
	})
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	const op = "handlers/user.Delete"

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

	if errors.Is(err, userDomain.ErrUserNotFound) {
		render.JSON(w, r, resp.Error("user not found"))
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
	Email string `json:"email"`
	Name  string `json:"name"`
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	const op = "handlers/user.GetByID"

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

	user, err := h.service.GetByID(r.Context(), id)

	if errors.Is(err, userDomain.ErrUserNotFound) {
		log.Info("user not found", slog.String("user_id", id.String()))
		render.JSON(w, r, resp.Error("user not found"))
		return
	}

	if err != nil {
		log.Error("get failed", sl.Err(err))
		render.JSON(w, r, resp.Error("failed to get"))
		return
	}

	render.JSON(w, r, GetByIDResponse{
		Response: resp.OK(),
		Email:    user.Email,
		Name:     user.Name,
	})
}
