package infra

import (
	"context"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/surajjain36/assignment_service/misc"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

//Mongo Struct
type Mongo struct {
	db *mongo.Database
}

//NewMongo is a function which establishes connection to MongoDB.
func NewMongo(conf *misc.MongoConfig) (*Mongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", conf.Username, conf.Password, conf.Host, conf.Port, conf.DB)))

	if err != nil {
		log.Println("mongo.Connect() ERROR:", err)
		os.Exit(1)
	}

	err = client.Ping(ctx, readpref.Primary())
	log.Println(err)
	if err == nil {
		return &Mongo{
			db: client.Database(conf.DB),
		}, nil
	}

	return nil, err
}

//Insert is a generic function to do insert call on DB
//Author: Suraj
//date: 04/01/2020
func (mgo *Mongo) Insert(collection string, data interface{}) (interface{}, error) {
	var err error
	var res interface{}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if res, err = mgo.db.Collection(collection).InsertOne(ctx, data); err != nil {
		log.Println("Error while inserting record: ", err)
	}
	return res, err
}

//FindOne : Generic query to find one record in MongoDb
//Author: Suraj
//date: 04/01/2020
func (mgo *Mongo) FindOne(collection string, filters interface{}, result interface{}) error {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err = mgo.db.Collection(collection).FindOne(ctx, filters).Decode(result); err != nil {
		log.Println("Error while retrieving record: ", err)
	}
	return err
}

//Aggregate : Generic function to get aggregated data from DB
//Author: Suraj
//date: 04/01/2020
func (mgo *Mongo) Aggregate(collection string, pipeline interface{}, result interface{}) error {
	var err error
	var cursor *mongo.Cursor
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	opts := options.Aggregate()
	opts.SetAllowDiskUse(true)
	opts.SetBatchSize(5)
	if cursor, err = mgo.db.Collection(collection).Aggregate(ctx, pipeline, opts); err == nil && cursor != nil {
		defer cursor.Close(ctx)
		err = cursor.All(ctx, result)
		if err != nil {
			log.Println(err)
		}
	}
	return err
}
