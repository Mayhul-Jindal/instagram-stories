package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
)

type apiFn func(w http.ResponseWriter, r *http.Request) (any, error)

type ApiError struct {
	Error string `json:"error"`
}

func writeJSON(ctx context.Context, w http.ResponseWriter, status int, resp any) {
	if status == http.StatusInternalServerError {
		slog.Error(
			"internal server error",
			"method", ctx.Value(method),
			"route", ctx.Value(route),
			"status", status,
			"msg", resp,
		)
	} else {
		slog.Info("request processed",
			"method", ctx.Value(method),
			"route", ctx.Value(route),
			"status", status,
			"msg", resp,
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

type Route string
const route Route = "route"

type Method string
const method Method = "method"

func apiFnHandler(apiFn apiFn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(context.WithValue(r.Context(), route, r.URL.Path))
		r = r.WithContext(context.WithValue(r.Context(), method, r.Method))

		resp, err := apiFn(w, r)
		if err == nil {
			writeJSON(r.Context(), w, http.StatusOK, resp)
			return
		}

		switch err {
		case ErrNotAuthorized, ErrNoAuthHeader, ErrInvalidAuthHeader, ErrUnsupportedAuthType, ErrTokenExpired, ErrInvalidKeySize, ErrInvalidToken:
			writeJSON(r.Context(), w, http.StatusUnauthorized, ApiError{Error: err.Error()})

		// todo better status code ?
		case ErrUserAlreadyExists, ErrUserAlreayFollowed, ErrCannotFollowSelf, ErrUserDoesNotExist:
			writeJSON(r.Context(), w, http.StatusBadRequest, ApiError{Error: err.Error()})

		case ErrUserIdMissing, ErrBadRequest:
			writeJSON(r.Context(), w, http.StatusBadRequest, ApiError{Error: err.Error()})

		default:
			var validationErr validator.ValidationErrors
			if errors.As(err, &validationErr) {
				writeJSON(r.Context(), w, http.StatusBadRequest, ApiError{Error: validationErr.Error()})
				return
			}
			writeJSON(r.Context(), w, http.StatusInternalServerError, ApiError{Error: err.Error()})

		}
	}
}

type Request struct {
	UserID int64 `json:"user_id"`
}

func (a *api) verifyTokenMiddleware(next apiFn) apiFn {
	return func(w http.ResponseWriter, r *http.Request) (any, error) {
		authorizationHeader := r.Header.Get("authorization")

		if len(authorizationHeader) == 0 {
			return nil, ErrNoAuthHeader
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			return nil, ErrInvalidAuthHeader
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != "bearer" {
			return nil, ErrUnsupportedAuthType
		}

		accessToken := fields[1]
		payload, err := a.token.Verify(accessToken)
		if err != nil {
			return nil, err
		}

		var req Request
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, &req)
		if err != nil {
			return nil, ErrUserIdMissing
		}
		if payload.UserID != req.UserID {
			return nil, ErrNotAuthorized
		}

		r.Body = io.NopCloser(bytes.NewReader(body))
		return next(w, r)
	}
}
