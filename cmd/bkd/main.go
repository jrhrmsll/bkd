package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/oklog/run"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"bookmarks/cmd/bkd/http"
	"bookmarks/internal/repository/mongodb"
)

const (
	mongoDBURI     = "mongodb://localhost:27017"
	mongoDBTimeout = "10s"
)

func main() {
	ctx := context.Background()

	logger := log.New(os.Stdout, "INFO: ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	logger.Println("Starting Bookmark Service")

	app := run.Group{}
	app.Add(run.SignalHandler(ctx, os.Interrupt, os.Kill, syscall.SIGTERM))

	// MongoDB connection
	mgo, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoDBURI))
	if err != nil {
		logger.Fatal(err)
	}

	defer func() {
		err := mgo.Disconnect(context.Background())
		if err != nil {
			logger.Println(err)
		}
	}()

	if err := func() error { // Test MongoDB connection
		d, err := time.ParseDuration(mongoDBTimeout)
		if err != nil {
			return err
		}

		ctx, cancel := context.WithTimeout(ctx, d)
		defer cancel()

		err = mgo.Ping(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to connect to mongodb: %w", err)
		}

		return nil
	}(); err != nil {
		logger.Fatal(err)
	}

	// API
	bookmarksRepository := mongodb.NewBookmarksRepository(ctx, mgo, logger)
	bookmarksController := http.NewBookmarkController(ctx, bookmarksRepository, logger)

	bkd := http.NewAPI(ctx, logger, bookmarksController)
	app.Add(bkd.Execute, bkd.Interrupt)

	err = app.Run()
	if err != nil {
		logger.Println(err)
	}
}
