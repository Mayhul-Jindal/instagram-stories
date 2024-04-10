package main

import "errors"

var (
	ErrTokenExpired        = errors.New("token has expired")
	ErrInvalidKeySize      = errors.New("invalid key size for paseto token")
	ErrInvalidToken        = errors.New("validation failed for paseto token. invalid token")
	ErrUserAlreadyExists   = errors.New("user already exists")
	ErrUserAlreayFollowed  = errors.New("user is already followed")
	ErrNotAuthorized       = errors.New("not authorized")
	ErrCannotFollowSelf    = errors.New("cannot follow self")
	ErrUserDoesNotExist    = errors.New("user does not exist")
	ErrNoAuthHeader        = errors.New("no authorization header")
	ErrInvalidAuthHeader   = errors.New("invalid authorization header")
	ErrUnsupportedAuthType = errors.New("unsupported authorization type")
	ErrUserIdMissing       = errors.New("user id is missing")
	ErrBadRequest          = errors.New("bad request")
	ErrInvalidStoryID      = errors.New("invalid story id")
	ErrStoryNotFound       = errors.New("story not found")
	ErrStoryAlreadyWatched = errors.New("story already watched")
)
