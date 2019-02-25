package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/tommbee/go-article-ingest/model"
)

// MongoArticleRepository interfaces with a mongo db instance
type MongoArticleRepository struct {
	Server       string
	DatabaseName string
	AuthDatabase string
	DBSSL        string
	Collection   string
	Username     string
	Password     string
	ReplicaSet   string
	db           *mongo.Database
}

type ctx *context.Context

// type key string

// const (
// 	hostKey         = key("hostKey")
// 	usernameKey     = key("usernameKey")
// 	passwordKey     = key("passwordKey")
// 	databaseKey     = key("databaseKey")
// 	authDatabaseKey = key("authDatabaseKey")
// 	dBSSL           = key("dBSSL")
// )

// Connect to the db instance
// func (r *MongoArticleRepository) Connect() {
// 	log.Print("Connecting...")
// 	ctx := context.Background()
// 	ctx, cancel := context.WithCancel(ctx)
// 	defer cancel()
// 	ctx = context.WithValue(ctx, hostKey, r.Server)
// 	ctx = context.WithValue(ctx, usernameKey, r.Username)
// 	ctx = context.WithValue(ctx, passwordKey, r.Password)
// 	ctx = context.WithValue(ctx, databaseKey, r.DatabaseName)
// 	ctx = context.WithValue(ctx, authDatabaseKey, r.AuthDatabase)
// 	ctx = context.WithValue(ctx, dBSSL, r.DBSSL)
// 	db, err := configDB(ctx)
// 	if err != nil {
// 		log.Fatalf("todo: database configuration failed: %v", err)
// 	}
// 	r.db = db
// }

// func configDB(ctx context.Context) (*mongo.Database, error) {
// 	uri := fmt.Sprintf(`mongodb://%s:%s@%s/%s?authSource=%s&ssl=%s`,
// 		ctx.Value(usernameKey),
// 		ctx.Value(passwordKey),
// 		ctx.Value(hostKey),
// 		ctx.Value(databaseKey),
// 		ctx.Value(authDatabaseKey),
// 		ctx.Value(dBSSL),
// 	)
// 	log.Print(uri)
// 	client, err := mongo.NewClient(uri)
// 	if err != nil {
// 		return nil, fmt.Errorf("todo: couldn't connect to mongo: %v", err)
// 	}
// 	err = client.Connect(ctx)
// 	if err != nil {
// 		return nil, fmt.Errorf("todo: mongo client couldn't connect with background context: %v", err)
// 	}
// 	dbName := ctx.Value(databaseKey).(string)
// 	articleDB := client.Database(dbName)
// 	return articleDB, nil
// }

// Insert a record to the db
// func (r *MongoArticleRepository) Insert(a model.Article) (model.Article, error) {
// 	log.Printf("Attempting insert: %s", a.URL)
// 	r.Connect()
// 	collection := r.db.Collection(r.Collection)
// 	cur, err := collection.CountDocuments(context.Background(), bson.M{"url": a.URL})
// 	if err != nil {
// 		return a, err
// 	}
// 	if cur > 0 {
// 		return a, fmt.Errorf("resource %s already exists", a.URL)
// 	}
// 	a.CreatedAt = time.Now()

// 	res, err := collection.InsertOne(context.Background(), bson.M{
// 		"title":      a.Title,
// 		"url":        a.URL,
// 		"published":  a.Published,
// 		"created_at": a.CreatedAt,
// 	})

// 	if err != nil {
// 		return a, err
// 	}

// 	log.Printf("Inserted record: %s", res.InsertedID)
// 	return a, nil
// }

// // Connect to the db instance
// func (r *MongoArticleRepository) Connect() (*mongo.Client, error) {
// 	log.Print("Connecting...")

// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	uri := fmt.Sprintf(`mongodb://%s:%s@%s/%s?authSource=%s&ssl=%s`,
// 		r.Username,
// 		r.Password,
// 		r.Server,
// 		r.DatabaseName,
// 		r.AuthDatabase,
// 		r.DBSSL,
// 	)

// 	log.Print(uri)

// 	client, err := mongo.Connect(ctx, uri)

// 	return client, err
// }

// Insert a record to the db
func (r *MongoArticleRepository) Insert(a model.Article) (model.Article, error) {
	//client, err := r.Connect()
	log.Print("Connecting...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := fmt.Sprintf(`mongodb://%s:%s@%s/%s?authSource=%s&ssl=%s`,
		r.Username,
		r.Password,
		r.Server,
		r.DatabaseName,
		r.AuthDatabase,
		r.DBSSL,
	)

	log.Print(uri)

	client, err := mongo.Connect(ctx, uri)

	if err != nil {
		log.Fatalf("todo: database configuration failed: %v", err)
	}

	collection := client.Database(r.DatabaseName).Collection(r.Collection)
	matches, err := collection.CountDocuments(context.Background(), bson.M{"url": a.URL})
	if err != nil {
		return a, err
	}
	if matches > 0 {
		return a, fmt.Errorf("resource %s already exists", a.URL)
	}

	a.CreatedAt = time.Now()

	res, err := collection.InsertOne(context.Background(), bson.M{
		"title":      a.Title,
		"url":        a.URL,
		"published":  a.Published,
		"created_at": a.CreatedAt,
	})

	if err != nil {
		return a, err
	}

	err = client.Disconnect(ctx)

	if err != nil {
		return a, err
	}

	log.Printf("Inserted record: %s", res.InsertedID)
	return a, nil
}
