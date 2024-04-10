package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Mayhul-Jindal/instagram-stories/mongodb"
	"github.com/Mayhul-Jindal/instagram-stories/postgres"
	"golang.org/x/crypto/bcrypt"
)

// mein yaha peh hee jwt aur login signup banara hun maa chuday
type Auther interface {
	SignUp(ctx context.Context, email string, password string) (SignUpUserResponse, error)
	Login(ctx context.Context, email string, password string) (LoginUserResponse, error)
}

type auth struct {
	db       postgres.Querier
	token    Token
	tokenExp time.Duration
	mongo mongodb.Querier
}

func NewAuthSvc(db postgres.Querier, token Token, tokenExp time.Duration, mongo mongodb.Querier) Auther {
	return &auth{
		db:       db,
		token:    token,
		tokenExp: tokenExp,
		mongo: mongo,
	}
}


func (a *auth) SignUp(ctx context.Context, email string, password string) (SignUpUserResponse, error) {
	hashedPassword, err := hashPassword(password)
	if err != nil {
		return SignUpUserResponse{}, err
	}

	res, err := a.db.CreateUser(ctx, postgres.CreateUserParams{
		Email:          email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		return SignUpUserResponse{}, ErrUserAlreadyExists
	}

	// creating user timeline
	err = a.mongo.CreateUserTimeline(ctx, mongodb.UserTimeline{
		UserID: res.ID,
		Timeline: []mongodb.UserTimelineData{},
	})
	if err != nil {
		return SignUpUserResponse{}, err
	}

	return SignUpUserResponse{
		UserID:    res.ID,
		CreatedAt: res.CreatedAt,
	}, nil
}

func (a *auth) Login(ctx context.Context, email string, password string) (LoginUserResponse, error) {
	user, err := a.db.GetUserByEmail(ctx, email)
	if err != nil {
		return LoginUserResponse{}, err
	}

	err = checkPassword(password, user.HashedPassword)
	if err != nil {
		return LoginUserResponse{}, ErrNotAuthorized
	}

	accessToken, accessPayload, err := a.token.Create(
		user.ID,
		a.tokenExp,
	)
	if err != nil {
		return LoginUserResponse{}, err
	}

	return LoginUserResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: accessPayload.ExpiredAt,
		User:                 SignUpUserResponse{UserID: user.ID, CreatedAt: user.CreatedAt},
	}, nil
}

// hashPassword returns the bcrypt hash of the password
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// checkPassword checks if the provided password is correct or not
func checkPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
