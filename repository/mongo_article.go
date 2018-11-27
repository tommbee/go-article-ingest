package repository

import (
	"fmt"
	"log"

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

// Fetch all records from the article repository
func (r *MongoArticleRepository) Fetch(num int) ([]*model.Article, error) {
	r.Connect()
	defer r.session.Close()
	log.Print("Getting articles")
	var articles []*model.Article
	q := r.session.DB(r.DatabaseName).C(r.Collection).Find(bson.M{}).Limit(num)
	err := q.All(&articles)
	return articles, err
}

// GetByID an entity
func (r *MongoArticleRepository) GetByID(id int64) (*model.Article, error) {
	return nil, nil
}

// GetByTitle entity
func (r *MongoArticleRepository) GetByTitle(title string) (*model.Article, error) {
	return nil, nil
}

// GetByUrl entity
func (r *MongoArticleRepository) GetByUrl(URL string) (*model.Article, error) {
	return nil, nil
}

// Insert a record to the db
func (r *MongoArticleRepository) Insert(a model.Article) error {
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
	return c.Insert(a)
}
