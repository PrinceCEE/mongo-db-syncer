package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func main() {
	var mongoURI, sourcedb, destdb string

	flag.StringVar(&mongoURI, "mongo_uri", "", "The mongodb connection uri, without the database name")
	flag.StringVar(&sourcedb, "source_db", "", "The source database")
	flag.StringVar(&destdb, "dest_db", "", "The destination database")
	flag.Parse()

	if mongoURI == "" {
		fmt.Println("mongo_uri is required")
		os.Exit(1)
	}

	if sourcedb == "" {
		fmt.Println("source db is required")
		os.Exit(1)
	}
	if destdb == "" {
		fmt.Println("dest db is required")
		os.Exit(1)
	}

	if err := dbSyncer(mongoURI, sourcedb, destdb); err != nil {
		fmt.Println(fmt.Errorf("error syncing from source to destination: %w", err))
		os.Exit(1)
	}

	fmt.Println("Source DB has been synced to the destination DB successfully.")
}

func dbSyncer(uri, source, dest string) error {
	ctx := context.Background()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return fmt.Errorf("failed to ping: %w", err)
	}

	sourceDB := client.Database(source)
	collectionNames, err := sourceDB.ListCollectionNames(ctx, bson.D{})
	if err != nil {
		return err
	}

	worker := func(ctx context.Context, name string, errChan chan<- error, wg *sync.WaitGroup) {
		defer wg.Done()

		collection := sourceDB.Collection(name)
		_, err := collection.Aggregate(ctx, mongo.Pipeline{
			{{
				Key: "$out",
				Value: bson.D{
					{
						Key:   "db",
						Value: dest,
					},
					{
						Key:   "coll",
						Value: name,
					},
				},
			}},
		})

		if err != nil {
			errChan <- err
			return
		}
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 1)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, name := range collectionNames {
		fmt.Printf("Syncing collection: %s\n", name)
		wg.Add(1)
		go worker(ctx, name, errChan, &wg)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case err := <-errChan:
		cancel()
		return err
	case <-done:
		return nil
	}
}
