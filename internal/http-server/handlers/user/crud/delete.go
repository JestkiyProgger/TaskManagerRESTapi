package crud

import (
	resp "ProjectManagementAPI/internal/lib/api/response"
	"ProjectManagementAPI/internal/lib/logger/sl"
	errors2 "ProjectManagementAPI/internal/user"
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type DeleteRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

type DeleteResponse struct {
	resp.Response
}

type DeleteService interface {
	Delete(ctx context.Context, id uuid.UUID) error
}

func NewDeleteHandler(log *slog.Logger, service DeleteService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers/user/create.NewDeleteHandler"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req DeleteRequest

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

		err = service.Delete(r.Context(), req.ID)
		if errors.Is(err, errors2.ErrUserNotFound) {
			log.Info("user not found", slog.String("id", req.ID.String()))

			render.JSON(w, r, resp.Error("user not found"))

			return
		}
		if err != nil {
			log.Info("failed to user not found", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to user not found"))

			return
		}

		log.Info("user delete", slog.String("id", req.ID.String()))

		render.JSON(w, r, DeleteResponse{
			Response: resp.OK(),
		})
	}
}
