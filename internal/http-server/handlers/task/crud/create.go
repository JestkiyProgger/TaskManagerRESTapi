package crud

import (
	resp "ProjectManagementAPI/internal/lib/api/response"
	"ProjectManagementAPI/internal/lib/logger/sl"
	errors2 "ProjectManagementAPI/internal/task"
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CreateRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	Status      string `json:"status" validate:"required"`
}

type CreateResponse struct {
	resp.Response
}

type CreateService interface {
	Create(ctx context.Context, title, description, status string) (uuid.UUID, error)
}

func CreateNew(log *slog.Logger, service CreateService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers/task/create.CreateNew"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req CreateRequest

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			errors.As(err, &validateErr)

			log.Error("failed to validate request", sl.Err(err))

			render.JSON(w, r, resp.ValidationError(validateErr))

			return
		}

		id, err := service.Create(r.Context(), req.Title, req.Description, req.Status)
		if errors.Is(err, errors2.ErrInvalidTitle) {
			log.Info("invalid title", slog.String("title", req.Title))

			render.JSON(w, r, resp.Error("invalid title"))

			return
		}
		if err != nil {
			log.Info("failed to create task", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to create task"))

			return
		}

		log.Info("task created", slog.String("id", id.String()))

		render.JSON(w, r, CreateResponse{
			Response: resp.OK(),
		})
	}
}
