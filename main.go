package main

import (
	"context"
	"log"
	"time"

	"github.com/Mayhul-Jindal/instagram-stories/mongodb"
	"github.com/Mayhul-Jindal/instagram-stories/postgres"
	"github.com/go-playground/validator"
)

const (
	symmetricKey  = "12345678901234567890123456789012"
	postgresURL   = "postgresql://postgres:postgres@localhost:5432/instagram-stories?sslmode=disable"
	tokenDuration = 1 * time.Hour
	listenAddr    = ":3000"
	mondodbURL = "mongodb://mongo:mongo@localhost:27017"
)

func main() {
	tokenSvc, err := NewPasetoTokenSvc(symmetricKey)
	if err != nil {
		log.Fatalf("cannot create token service: %v", err)
	}

	postgres, conn, err := postgres.NewPostresDB(context.TODO(), postgresURL)
	if err != nil {
		log.Fatalf("cannot connect to postgres: %v", err)
	}
	defer conn.Close(context.TODO())

	mongo, err := mongodb.New(context.TODO(), mondodbURL, "instagram-stories")
	if err != nil {
		log.Fatalf("cannot connect to mongo: %v", err)
	}

	authSvc := NewAuthSvc(postgres, tokenSvc, tokenDuration, mongo)

	validator := validator.New()

	businessSvc := NewBussiness(postgres, mongo)

	api := NewAPI(listenAddr, tokenSvc, authSvc, validator, businessSvc)
	log.Println("starting server at", listenAddr)
	log.Fatalln(api.Run())
}

// mongo, err := mongodb.New(context.TODO(), "mongodb://mongo:mongo@localhost:27017", "instagram-stories")
// if err != nil {
// 	log.Fatalf("cannot connect to mongo: %v", err)
// }

// // story creating
// storyID, err := mongo.CreateStory(context.TODO(), mongodb.Story{
// 	UserID:    1,
// 	Data:      "hello world",
// 	CreatedAt: time.Now(),
// })
// if err != nil {
// 	log.Fatalf("cannot create story: %v", err)
// }
// log.Println("story created", storyID.Hex())

// // get story details by story id
// story, err := mongo.GetStoryById(context.TODO(), storyID)
// if err != nil {
// 	log.Fatalf("cannot get story: %v", err)
// }
// log.Println("story retrieved", story)

// // creating user timeline
// err = mongo.CreateUserTimeline(context.TODO(), mongodb.UserTimeline{
// 	UserID: 100,
// 	Timeline: []mongodb.UserTimelineData{},
// })
// if err != nil {
// 	log.Fatalf("cannot create user timeline: %v", err)
// }

// err = mongo.CreateUserTimeline(context.TODO(), mongodb.UserTimeline{
// 	UserID: 200,
// 	Timeline: []mongodb.UserTimelineData{},
// })
// if err != nil {
// 	log.Fatalf("cannot create user timeline: %v", err)
// }

// err = mongo.UpdateTimelineOfFollowers(context.TODO(), []int64{100,200}, mongodb.UserTimelineData{
// 	StoryID: storyID,
// 	CreatedAt: time.Now(),
// })
// if err != nil {
// 	log.Fatalf("cannot update timeline of followers: %v", err)
// }

// timeline, err := mongo.GetUserTimeline(context.TODO(), 100)
// if err != nil {
// 	log.Fatalf("cannot get user timeline: %v", err)
// }

// log.Println("user timeline", timeline)

// story, err = mongo.GetStoryFromUserTimeline(context.TODO(), 100, timeline[0].StoryID)
// if err != nil {
// 	log.Fatalf("cannot get story from user timeline: %v", err)
// }
// log.Println("story from user timeline", story)

// story, err = mongo.GetStoryFromUserTimeline(context.TODO(), 100, timeline[0].StoryID)
// if err != nil {
// 	log.Fatalf("cannot get story from user timeline: %v", err)
// }
// log.Println("story from user timeline", story)
