package mongodb

import (
	"context"
	"errors"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Querier interface {
	CreateStory(ctx context.Context, story Story) (primitive.ObjectID, error)
	GetStoryById(ctx context.Context, storyID primitive.ObjectID) (GetStoryByIdResponse, error)
	CreateUserTimeline(ctx context.Context, userTimeline UserTimeline) error
	UpdateTimelineOfFollowers(ctx context.Context, userIDs []int64, userTimelineData UserTimelineData) error
	GetUserTimeline(ctx context.Context, userID int64) ([]UserTimelineData, error)
	RemoveStoryFromUserTimeline(ctx context.Context, userID int64, storyID primitive.ObjectID) (int64, error)
}

// Interface can be made if we want to have a common method for all the databases with same functionality as Mongo.
// Also inteface can be useful if want to add onion layers on top of this mongo layer (logging, metrics, ..etc)
// Currently no such requirement thus going with struct.
type mongodb struct {
	db *mongo.Database
}

func New(ctx context.Context, url string, dbName string) (Querier, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	return &mongodb{
		db: client.Database(dbName),
	}, nil
}

type Story struct {
	UserID    int64     `bson:"user_id,omitempty"`
	Data      string    `bson:"data,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty"`
}

const StoriesCollection = "StoriesCollection"

func (m *mongodb) CreateStory(ctx context.Context, story Story) (primitive.ObjectID, error) {
	resp, err := m.db.Collection(StoriesCollection).InsertOne(ctx, story)
	if err != nil {
		return primitive.NilObjectID, err
	}

	// The inserted ID is an interface{}, so we'll need to type assert it to
	// an ObjectID.
	oid, ok := resp.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, errors.New("cannot convert to ObjectID")
	}

	return oid, nil
}

type GetStoryByIdResponse struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    int64              `bson:"user_id"`
	Data      string             `bson:"data"`
	CreatedAt time.Time          `bson:"created_at"`
}

func (m *mongodb) GetStoryById(ctx context.Context, storyID primitive.ObjectID) (GetStoryByIdResponse, error) {
	resp := m.db.Collection(StoriesCollection).FindOne(ctx, primitive.D{{Key: "_id", Value: storyID}})

	// Decode the response into a Story struct
	var story GetStoryByIdResponse
	if err := resp.Decode(&story); err != nil {
		return GetStoryByIdResponse{}, err
	}

	return story, nil
}

const TimelineCollection = "TimelineCollection"

type UserTimelineData struct {
	UserID int64              `bson:"user_id,omitempty"`
	StoryID   primitive.ObjectID `bson:"story_id,omitempty"`
	CreatedAt time.Time          `bson:"created_at,omitempty"`
}

type UserTimeline struct {
	UserID   int64              `bson:"user_id,omitempty"`
	Timeline []UserTimelineData `bson:"timeline,omitempty"`
}

// create timeline document
func (m *mongodb) CreateUserTimeline(ctx context.Context, userTimeline UserTimeline) error {
	resp, err := m.db.Collection(TimelineCollection).InsertOne(ctx, userTimeline)
	if err != nil {
		return err
	}

	log.Println("Inserted timeline document with ID:", resp.InsertedID)
	return nil
}

func (m *mongodb) UpdateTimelineOfFollowers(ctx context.Context, userIDs []int64, userTimelineData UserTimelineData) error {
	resp, err := m.db.Collection(TimelineCollection).UpdateMany(ctx, primitive.D{{Key: "user_id", Value: primitive.D{{Key: "$in", Value: userIDs}}}}, primitive.D{{Key: "$push", Value: primitive.D{{Key: "timeline", Value: userTimelineData}}}})
	if err != nil {
		return err
	}

	log.Println("Updated timeline document with ID:", resp.UpsertedID)
	return nil
}

func (m *mongodb) GetUserTimeline(ctx context.Context, userID int64) ([]UserTimelineData, error) {
	resp := m.db.Collection(TimelineCollection).FindOne(ctx, primitive.D{{Key: "user_id", Value: userID}})

	var userTimeline UserTimeline
	if err := resp.Decode(&userTimeline); err != nil {
		return nil, err
	}

	return userTimeline.Timeline, nil
}

// It pulls story from the timeline array and they get data of story from GetStoryById
func (m *mongodb) RemoveStoryFromUserTimeline(ctx context.Context, userID int64, storyID primitive.ObjectID) (int64, error) {
	filter := bson.D{{Key: "user_id", Value: userID}} // Replace with your document's _id
	update := bson.D{{
		Key: "$pull",
		Value: bson.D{{
			Key: "timeline",
			Value: bson.D{{
				Key:   "story_id",
				Value: storyID,
			}},
		}},
	}}

	resp, err := m.db.Collection(TimelineCollection).UpdateOne(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	log.Printf("Removed %d documents\n", resp.ModifiedCount)
	log.Println("matched count", resp.MatchedCount)

	return resp.ModifiedCount, nil
}

// for debugging purposes

// func (m *Mongodb) GetStory(ctx context.Context)            {}
// func (m *Mongodb) UpdateStoryTimeline(ctx context.Context) {}
// func (m *Mongodb) GetStoryTimeline(ctx context.Context)    {}
