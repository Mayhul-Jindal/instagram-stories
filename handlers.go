package main

import (
	"encoding/json"
	"net/http"
	"time"
)

type HealthResponse struct {
	Status string `json:"status"`
}

func (a *api) health(w http.ResponseWriter, r *http.Request) (any, error) {
	return HealthResponse{Status: "ok"}, nil
}

func (a *api) signup(w http.ResponseWriter, r *http.Request) (any, error) {
	var req SignUpUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, ErrBadRequest
	}

	err = a.validator.Struct(req)
	if err != nil {
		return nil, err
	}

	resp, err := a.auth.SignUp(r.Context(), req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *api) login(w http.ResponseWriter, r *http.Request) (any, error) {
	var req LoginUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, ErrBadRequest
	}

	err = a.validator.Struct(req)
	if err != nil {
		return nil, err
	}

	resp, err := a.auth.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *api) getUsersProfiles(w http.ResponseWriter, r *http.Request) (any, error) {
	var req GetUsersProfilesRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, ErrBadRequest
	}

	err = a.validator.Struct(req)
	if err != nil {
		return nil, err
	}

	resp, err := a.bussiness.GetUsersProfiles(r.Context(), req.UserID, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *api) followUser(w http.ResponseWriter, r *http.Request) (any, error) {
	var req FollowUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, ErrBadRequest
	}

	err = a.validator.Struct(req)
	if err != nil {
		return nil, err
	}

	err = a.bussiness.FollowUser(r.Context(), req.UserID, req.Email)
	if err != nil {
		return nil, err
	}

	return map[string]string{"status": "ok"}, nil
}

func (a *api) getFollowing(w http.ResponseWriter, r *http.Request) (any, error) {
	var req GetUsersProfilesRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, ErrBadRequest
	}

	err = a.validator.Struct(req)
	if err != nil {
		return nil, err
	}

	resp, err := a.bussiness.GetFollowingProfiles(r.Context(), req.UserID, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *api) getFollowers(w http.ResponseWriter, r *http.Request) (any, error) {
	var req GetUsersProfilesRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, ErrBadRequest
	}

	err = a.validator.Struct(req)
	if err != nil {
		return nil, err
	}

	resp, err := a.bussiness.GetFollowersProfiles(r.Context(), req.UserID, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type CreateStoryRequest struct {
	UserID int64  `json:"user_id" validate:"required,number,min=1"`
	Data   string `json:"data" validate:"required,min=5"`
}

func (a *api) createStory(w http.ResponseWriter, r *http.Request) (any, error) {
	var req CreateStoryRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, ErrBadRequest
	}

	err = a.validator.Struct(req)
	if err != nil {
		return nil, err
	}

	err = a.bussiness.CreateStory(r.Context(), req.UserID, req.Data, time.Now())
	if err != nil {
		return nil, err
	}

	return map[string]string{"status": "ok"}, nil
}

type GetStoriesTimelineRequest struct {
	UserID int64 `json:"user_id" validate:"required,number,min=1"`
}

func (a *api) getStoriesTimeline(w http.ResponseWriter, r *http.Request) (any, error) {
	var req GetStoriesTimelineRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, ErrBadRequest
	}

	err = a.validator.Struct(req)
	if err != nil {
		return nil, err
	}

	resp, err := a.bussiness.GetStoriesTimeline(r.Context(), req.UserID)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type WatchStoryRequest struct {
	UserID  int64 `json:"user_id" validate:"required,number,min=1"`
	StoryID string `json:"story_id" validate:"required,number,min=1"`
}

func (a *api) watchStory(w http.ResponseWriter, r *http.Request) (any, error) {
	var req WatchStoryRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, ErrBadRequest
	}

	err = a.validator.Struct(req)
	if err != nil {
		return nil, err
	}

	resp, err := a.bussiness.WatchStoryById(r.Context(), req.UserID, req.StoryID)
	if err != nil {
		return nil, err
	}

	return resp, nil
}