package authhandler

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"fmt"
	"gonotes/internal/api"
	"gonotes/internal/auth"
	"gonotes/internal/auth/dto"
	"gonotes/internal/auth/entity"
	"gonotes/internal/auth/storage"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/golang-jwt/jwt/v5"
)

type Handler struct {
	log     *slog.Logger
	storage storage.UserRepository
	secret  []byte
}

func NewHandler(log *slog.Logger, storage storage.UserRepository, secret string) *Handler {
	return &Handler{
		log:     log,
		storage: storage,
		secret:  []byte(secret),
	}
}
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	const op = "auth.handler.create"
	log := h.log.With(slog.String("op", op))

	var req dto.UserRequest
	if err := render.Bind(r, &req); err != nil {
		log.Error("invalid request body", slog.Any("err", err))
		render.Render(w, r, api.NewErrResponse(http.StatusBadRequest, err))
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", slog.Any("err", err))
		render.Render(w, r, api.NewErrResponse(http.StatusInternalServerError, err))
		return
	}

	id, err := h.storage.Create(&entity.User{Email: req.Email, Password: string(hashed)})
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			log.Error("failed to create user", slog.Any("err", err))
			render.Render(w, r, api.NewErrResponse(http.StatusConflict, err))
			return
		}
		log.Error("failed to create user", slog.Any("err", err))
		render.Render(w, r, api.NewErrResponse(http.StatusInternalServerError, err))
		return
	}

	render.Status(r, http.StatusCreated)
	render.Render(w, r, dto.NewUserResponse(&entity.User{Email: req.Email, ID: id}))

	log.Info("user created", slog.Int("user_id", id))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "auth.handler.login"
	log := h.log.With(slog.String("op", op))

	var req dto.UserRequest
	if err := render.Bind(r, &req); err != nil {
		log.Error("invalid request", slog.Any("err", err))
		render.Render(w, r, api.NewErrResponse(http.StatusBadRequest, err))
		return
	}

	user, err := h.storage.Get(req.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		log.Error("invalid credentials", slog.Any("err", err))
		render.Render(w, r, api.NewErrResponse(http.StatusUnauthorized, fmt.Errorf("invalid credentials")))
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenStr, err := token.SignedString(h.secret)
	if err != nil {
		log.Error("failed to sign token", slog.Any("err", err))
		render.Render(w, r, api.NewErrResponse(http.StatusInternalServerError, err))
		return
	}

	render.Render(w, r, dto.NewLoginResponse(user.Email, tokenStr))
	log.Info("user logged in", slog.Int("user_id", user.ID))
}
