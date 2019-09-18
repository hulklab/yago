package mgo

import (
	"context"
	"github.com/hulklab/yago"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type Mgo struct {
	*mongo.Database
}

type Collection struct {
	c *mongo.Collection
}

type Cursor struct {
	c *mongo.Cursor
}

// 返回 mongodb 的一个数据库连接
func Ins(id ...string) *Mgo {

	var name string

	if len(id) == 0 {
		name = "mongodb"
	} else if len(id) > 0 {
		name = id[0]
	}

	v := yago.Component.Ins(name, func() interface{} {
		m := new(Mgo)
		conf := yago.Config.GetStringMap(name)
		uri := conf["mongodb_uri"].(string)

		client, err := mongo.NewClient(options.Client().ApplyURI(uri))
		if err != nil {
			log.Fatalf("Fatal error mongo: %s", err)
		}

		if err := client.Connect(defCtx()); err != nil {
			log.Fatalf("Fatal error mongo: %s", err)
		}

		if database, ok := conf["database"].(string); !ok {
			log.Fatal("Fatal error mongo: database is required")
		} else {
			m.Database = client.Database(database)
		}
		return m
	})

	return v.(*Mgo)
}

func defCtx() context.Context {
	return context.Background()
}

func (m *Mgo) DB(name string) *Mgo {
	return &Mgo{m.Client().Database(name)}
}

func (m *Mgo) C(name string) *Collection {
	return &Collection{m.Collection(name)}
}

// collection
func (c *Collection) InsertOne(document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return c.c.InsertOne(defCtx(), document, opts...)
}

func (c *Collection) InsertMany(document []interface{}, opts ...*options.InsertManyOptions) (*mongo.InsertManyResult, error) {
	return c.c.InsertMany(defCtx(), document, opts...)
}

func (c *Collection) FindOne(filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return c.c.FindOne(defCtx(), filter, opts...)
}

func (c *Collection) Find(filter interface{}, opts ...*options.FindOptions) (*Cursor, error) {
	cursor, err := c.c.Find(defCtx(), filter, opts...)
	return &Cursor{cursor}, err
}

func (c *Collection) Distinct(fieldName string, filter interface{}, opts ...*options.DistinctOptions) ([]interface{}, error) {
	return c.c.Distinct(defCtx(), fieldName, filter, opts...)
}

func (c *Collection) Aggregate(pipeline interface{}, opts ...*options.AggregateOptions) (*Cursor, error) {
	cursor, err := c.c.Aggregate(defCtx(), pipeline, opts...)
	return &Cursor{cursor}, err
}

func (c *Collection) BulkWrite(models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	return c.c.BulkWrite(defCtx(), models, opts...)
}

func (c *Collection) CountDocuments(filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return c.c.CountDocuments(defCtx(), filter, opts...)
}

func (c *Collection) DeleteOne(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return c.c.DeleteOne(defCtx(), filter, opts...)
}

func (c *Collection) DeleteMany(filter interface{}, opts ...*options.DeleteOptions) (*mongo.DeleteResult, error) {
	return c.c.DeleteMany(defCtx(), filter, opts...)
}

func (c *Collection) Drop() error {
	return c.c.Drop(defCtx())
}

func (c *Collection) EstimatedDocumentCount(opts ...*options.EstimatedDocumentCountOptions) (int64, error) {
	return c.c.EstimatedDocumentCount(defCtx(), opts...)
}

func (c *Collection) Clone(opts ...*options.CollectionOptions) (*Collection, error) {
	if collection, err := c.c.Clone(opts...); err != nil {
		return nil, err
	} else {
		return &Collection{collection}, nil
	}
}

func (c *Collection) FindOneAndDelete(filter interface{}, opts ...*options.FindOneAndDeleteOptions) *mongo.SingleResult {
	return c.c.FindOneAndDelete(defCtx(), filter, opts...)
}

func (c *Collection) FindOneAndReplace(filter, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) *mongo.SingleResult {
	return c.c.FindOneAndReplace(defCtx(), filter, replacement, opts...)
}

func (c *Collection) FindOneAndUpdate(filter, update interface{}, opts ...*options.FindOneAndUpdateOptions) *mongo.SingleResult {
	return c.c.FindOneAndUpdate(defCtx(), filter, update, opts...)
}

func (c *Collection) ReplaceOne(filter, replacement interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	return c.c.ReplaceOne(defCtx(), filter, replacement, opts...)
}

func (c *Collection) UpdateOne(filter, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.c.UpdateOne(defCtx(), filter, update, opts...)
}

func (c *Collection) UpdateMany(filter, update interface{}, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return c.c.UpdateMany(defCtx(), filter, update, opts...)
}

func (c *Collection) Upsert(filter, update interface{}, opts ...*options.ReplaceOptions) (*mongo.UpdateResult, error) {
	opt := options.ReplaceOptions{}
	opt.SetUpsert(true)

	if len(opts) > 0 {
		opts = append(opts, &opt)
	} else {
		opts = []*options.ReplaceOptions{&opt}
	}

	return c.c.ReplaceOne(defCtx(), filter, update, opts...)
}

func (c *Collection) Watch(pipeline interface{}, opts ...*options.ChangeStreamOptions) (*mongo.ChangeStream, error) {
	return c.c.Watch(defCtx(), pipeline, opts...)
}

func (c *Collection) Orig() *mongo.Collection {
	return c.c
}

// cursor
func (c *Cursor) Next() bool {
	return c.c.Next(defCtx())
}

func (c *Cursor) Close() error {
	return c.c.Close(defCtx())
}

func (c *Cursor) All(results interface{}) error {
	return c.c.All(defCtx(), results)
}

func (c *Cursor) Decode(val interface{}) error {
	return c.c.Decode(val)
}

func (c *Cursor) Orig() *mongo.Cursor {
	return c.c
}
