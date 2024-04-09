package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Interface can be made if we want to have a common method for all the databases with same functionality as Mongo.
// Also inteface can be useful if want to add onion layers on top of this mongo layer (logging, metrics, ..etc)

type Mongodb struct {
	client   *mongo.Client
	database *mongo.Database
}

func New(ctx context.Context, url string, dbName string) (*Mongodb, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	return &Mongodb{
		client:   client,
		database: client.Database(dbName),
	}, nil
}

type Story struct {
	UserID int `bson:"user_id"`
	Data   string
}

func (m *Mongodb) InsertStory(ctx context.Context, story Story) error {
	_, err := m.database.Collection("allStories").InsertOne(ctx, story)
	if err != nil {
		return err
	}

	return nil
}

// for debugging purposes

func (m *Mongodb) GetStory(ctx context.Context)            {}
func (m *Mongodb) UpdateStoryTimeline(ctx context.Context) {}
func (m *Mongodb) GetStoryTimeline(ctx context.Context)    {}
