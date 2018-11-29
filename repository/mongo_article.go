package repository

import (
	"fmt"
	"log"
	"time"

	"github.com/tommbee/go-article-ingest/model"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoArticleRepository interfaces with a mongo db instance
type MongoArticleRepository struct {
	Server       string
	DatabaseName string
	Collection   string
	session      *mgo.Session
}

// Connect to the db instance
func (r *MongoArticleRepository) Connect() {
	session, err := mgo.Dial(r.Server)
	if err != nil {
		log.Fatal(err)
	}
	r.session = session
}

// Insert a record to the db
func (r *MongoArticleRepository) Insert(a model.Article) error {
	log.Printf("Attempting insert: %s", a.URL)
	r.Connect()

	defer r.session.Close()
	c := r.session.DB(r.DatabaseName).C(r.Collection)
	count, err := c.Find(bson.M{"url": a.URL}).Limit(1).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("resource %s already exists", a.URL)
	}
	a.CreatedAt = time.Now()

	return c.Insert(a)
}
