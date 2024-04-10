package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator"
)

type api struct {
	listenAddr string
	token      Token
	auth       Auther
	validator  *validator.Validate
	bussiness  Bussiness
}

func NewAPI(listenAddr string, token Token, auth Auther, validator *validator.Validate, bussiness Bussiness) *api {
	return &api{
		listenAddr: listenAddr,
		token:      token,
		auth:       auth,
		validator:  validator,
		bussiness:  bussiness,
	}
}

func (a *api) Run() error {
	r := chi.NewRouter()

	r.Get("/health", apiFnHandler(a.health))

	r.Post("/signup", apiFnHandler(a.signup))
	r.Post("/login", apiFnHandler(a.login))

	r.Get("/profiles", apiFnHandler(a.verifyTokenMiddleware(a.getUsersProfiles)))
	r.Post("/follow", apiFnHandler(a.verifyTokenMiddleware(a.followUser)))
	r.Get("/following", apiFnHandler(a.verifyTokenMiddleware(a.getFollowing)))
	r.Get("/followers", apiFnHandler(a.verifyTokenMiddleware(a.getFollowers)))

	r.Post("/story/create", apiFnHandler(a.verifyTokenMiddleware(a.createStory)))
	r.Get("/story/get-timeline", apiFnHandler(a.verifyTokenMiddleware(a.getStoriesTimeline)))
	r.Get("/story/watch", apiFnHandler(a.verifyTokenMiddleware(a.watchStory)))

	return http.ListenAndServe(a.listenAddr, r)
}
