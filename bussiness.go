package main

import (
	"context"
	"errors"
	"time"

	"github.com/Mayhul-Jindal/instagram-stories/mongodb"
	"github.com/Mayhul-Jindal/instagram-stories/postgres"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bussiness interface {
	GetUsersProfiles(ctx context.Context, userID int64, limit int32, offset int32) ([]string, error)
	FollowUser(ctx context.Context, userID int64, email string) error
	GetFollowingProfiles(ctx context.Context, userID int64, limit int32, offset int32) ([]string, error)
	GetFollowersProfiles(ctx context.Context, userID int64, limit int32, offset int32) ([]string, error)

	CreateStory(ctx context.Context, userID int64, data string, date time.Time) error
	GetStoriesTimeline(ctx context.Context, userID int64) ([]mongodb.UserTimelineData, error)
	WatchStoryById(ctx context.Context, userID int64, storyID string) (mongodb.GetStoryByIdResponse, error)
}

type bussiness struct {
	postgres postgres.Querier
	mongo    mongodb.Querier
}

func NewBussiness(postgres postgres.Querier, mongo mongodb.Querier) Bussiness {
	return &bussiness{
		postgres: postgres,
		mongo:    mongo,
	}
}

// currently sends all them emails
func (b *bussiness) GetUsersProfiles(ctx context.Context, userID int64, limit int32, offset int32) ([]string, error) {
	resp, err := b.postgres.GetUsersEmails(ctx, postgres.GetUsersEmailsParams{
		ID:     userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (b *bussiness) FollowUser(ctx context.Context, userID int64, email string) error {
	followingUser, err := b.postgres.GetUserByEmail(ctx, email)
	if err != nil {
		return ErrUserDoesNotExist
	}

	if userID == followingUser.ID {
		return ErrCannotFollowSelf
	}

	// two cases
	// 1. following the same person twice
	// 2. following someone who doesn't exist
	err = b.postgres.FollowUser(ctx, postgres.FollowUserParams{
		FollowerID:  userID,
		FollowingID: followingUser.ID,
	})
	if err != nil {
		return pgErrorHandler(err, map[string]error{
			pgerrcode.ForeignKeyViolation: ErrUserDoesNotExist,
			pgerrcode.UniqueViolation:     ErrUserAlreayFollowed,
		})
	}

	return nil
}

func (b *bussiness) GetFollowingProfiles(ctx context.Context, userID int64, limit int32, offset int32) ([]string, error) {
	resp, err := b.postgres.GetFollowingEmails(ctx, postgres.GetFollowingEmailsParams{
		FollowerID: userID,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (b *bussiness) GetFollowersProfiles(ctx context.Context, userID int64, limit int32, offset int32) ([]string, error) {
	resp, err := b.postgres.GetFollowersEmails(ctx, postgres.GetFollowersEmailsParams{
		FollowingID: userID,
		Limit:       limit,
		Offset:      offset,
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}

/*
1. mongo mein ekh story update kardo ok
2. Get all your followers id
3. call update timeline
*/
func (b *bussiness) CreateStory(ctx context.Context, userID int64, data string, date time.Time) error {
	storyID, err := b.mongo.CreateStory(ctx, mongodb.Story{
		UserID:    userID,
		Data:      data,
		CreatedAt: date,
	})
	if err != nil {
		return err
	}

	followers, err := b.postgres.GetFollowersIDs(ctx, userID)
	if err != nil {
		return err
	}

	followers = append(followers, userID)

	err = b.mongo.UpdateTimelineOfFollowers(ctx, followers, mongodb.UserTimelineData{
		StoryID:   storyID,
		CreatedAt: date,
	})
	if err != nil {
		return err
	}

	return nil
}

// get timeline of stories
func (b *bussiness) GetStoriesTimeline(ctx context.Context, userID int64) ([]mongodb.UserTimelineData, error) {
	timeline, err := b.mongo.GetUserTimeline(ctx, userID)
	if err != nil {
		return nil, err
	}

	return timeline, nil
}

// watch the story one by one
func (b *bussiness) WatchStoryById(ctx context.Context, userID int64, storyID string) (mongodb.GetStoryByIdResponse, error) {
	storyid, err := primitive.ObjectIDFromHex(storyID)
	if err != nil {
		return mongodb.GetStoryByIdResponse{}, ErrInvalidStoryID
	}
	
	count, err := b.mongo.RemoveStoryFromUserTimeline(ctx, userID, storyid)
	if err != nil {
		return mongodb.GetStoryByIdResponse{}, err
	}

	story, err := b.mongo.GetStoryById(ctx, storyid)
	if err != nil {
		return mongodb.GetStoryByIdResponse{}, ErrStoryNotFound
	}

	if count == 0 {
		return mongodb.GetStoryByIdResponse{}, ErrStoryAlreadyWatched
	}


	return story, nil
}

// util function to handle postgres errors
func pgErrorHandler(err error, errMap map[string]error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		for k, v := range errMap {
			if pgErr.Code == k {
				return v
			}
		}
	}

	return err
}
