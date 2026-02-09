package crud

import (
	resp "ProjectManagementAPI/internal/lib/api/response"
	"ProjectManagementAPI/internal/lib/logger/sl"
	errors2 "ProjectManagementAPI/internal/user"
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type DeleteRequest struct {
	ID uuid.UUID `json:"id"`
}

type DeleteResponse struct {
	resp.Response
}

type DeleteService interface {
	Delete(ctx context.Context, id uuid.UUID) error
}

func NewDeleteHandler(log *slog.Logger, service DeleteService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers/user/delete.NewDeleteHandler"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		idStr := chi.URLParam(r, "id")

		id, err := uuid.Parse(idStr)
		if err != nil {
			log.Error("invalid UUID format",
				slog.String("id", idStr),
				sl.Err(err))

			render.JSON(w, r, resp.Error("invalid user ID format"))
			return
		}

		log.Info("deleting user", slog.String("id", id.String()))

		err = service.Delete(r.Context(), id)
		if errors.Is(err, errors2.ErrUserNotFound) {
			log.Info("user not found", slog.String("id", id.String()))

			render.JSON(w, r, resp.Error("user not found"))

			return
		}
		if err != nil {
			log.Info("failed delete user", sl.Err(err))

			render.JSON(w, r, resp.Error("failed delete user"))

			return
		}

		log.Info("user delete", slog.String("id", id.String()))

		render.JSON(w, r, DeleteResponse{
			Response: resp.OK(),
		})
	}
}
