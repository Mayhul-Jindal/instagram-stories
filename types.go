package main

import "time"

type SignUpUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=7"`
}

type SignUpUserResponse struct {
	UserID    int64     `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=7"`
}

type LoginUserResponse struct {
	AccessToken          string             `json:"access_token"`
	AccessTokenExpiresAt time.Time          `json:"access_token_expires_at"`
	User                 SignUpUserResponse `json:"user"`
}

type GetUsersProfilesRequest struct {
	UserID int64 `json:"user_id" validate:"required,number,min=1"`
	Limit  int32 `json:"limit" validate:"required,number,min=1"`
	Offset int32 `json:"offset" validate:"number,min=0"`
}

type FollowUserRequest struct {
	UserID int64  `json:"user_id" validate:"required,number,min=1"`
	Email  string `json:"email" validate:"required,email"`
}

type CreateStoryRequest struct {
	UserID int64  `json:"user_id" validate:"required,number,min=1"`
	Data   string `json:"data" validate:"required,min=5"`
}

type GetStoriesTimelineRequest struct {
	UserID int64 `json:"user_id" validate:"required,number,min=1"`
}

type WatchStoryRequest struct {
	UserID  int64  `json:"user_id" validate:"required,number,min=1"`
	StoryID string `json:"story_id" validate:"required,min=1"`
}
