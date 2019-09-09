package mgo

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
)

// go test -v ./coms/mgo -run TestCollection_InsertOne -args "-c=${PWD}/example/conf/app.toml"
func TestCollection_InsertOne(t *testing.T) {
	mgoClient := Ins()
	insertResult, err := mgoClient.C("test").InsertOne(bson.M{"name": "tom"})
	if err != nil {
		t.Errorf("InsertOne error: %s", err)
		t.FailNow()
	}
	t.Log("InsertOne ok", insertResult)
}

// go test -v ./coms/mgo -run TestCollection_InsertMany -args "-c=${PWD}/example/conf/app.toml"
func TestCollection_InsertMany(t *testing.T) {
	mgoClient := Ins()
	insertResult, err := mgoClient.C("test").InsertMany(bson.A{
		bson.M{"name": "henry"},
		bson.M{"name": "lily"},
		bson.M{"name": "peter"},
		bson.M{"name": "sheldon"},
		bson.M{"name": "john"},
		bson.M{"name": "stark"},
	})
	if err != nil {
		t.Errorf("InsertMany error: %s", err)
		t.FailNow()
	}
	t.Log("InsertMany ok", insertResult)
}

// go test -v ./coms/mgo -run ^TestCollection_FindOne$ -args "-c=${PWD}/example/conf/app.toml"
func TestCollection_FindOne(t *testing.T) {
	mgoClient := Ins()
	findResult := mgoClient.C("test").FindOne(bson.M{"name": "henry"})
	if findResult.Err() != nil {
		t.Errorf("FindOne error: %s", findResult.Err())
		t.FailNow()
	}
	result := bson.D{}
	if err := findResult.Decode(&result); err != nil {
		t.Errorf("FindOne error: %s", err)
		t.FailNow()
	}
	t.Log("FindOne ok", result)
}

// go test -v ./coms/mgo -run ^TestCollection_Find$ -args "-c=${PWD}/example/conf/app.toml"
func TestCollection_Find(t *testing.T) {
	mgoClient := Ins()
	cursor, err := mgoClient.C("test").Find(bson.M{})
	if err != nil {
		t.Errorf("Find error: %s", err)
		t.FailNow()
	}
	defer cursor.Close()

	resultAll := bson.A{}
	if err := cursor.All(&resultAll); err != nil {
		t.Errorf("Find all error: %s", err)
		t.FailNow()
	}
	t.Log("Find all ok", resultAll)

	cursor, err = mgoClient.C("test").Find(bson.M{})
	if err != nil {
		t.Errorf("Find error: %s", err)
		t.FailNow()
	}
	defer cursor.Close()
	for {
		if !cursor.Next() {
			break
		}
		result := bson.D{}
		if err := cursor.Decode(&result); err != nil {
			t.Errorf("cursor decode error: %s", err)
			t.FailNow()
		}
		t.Log("Find ok", result)
	}
	t.Log("Find Cursor ok")
}

// go test -v ./coms/mgo -run TestCollection_Distinct -args "-c=${PWD}/example/conf/app.toml"
func TestCollection_Distinct(t *testing.T) {
	mgoClient := Ins()
	result, err := mgoClient.C("test").Distinct("name", bson.M{})
	if err != nil {
		t.Errorf("Distinct error: %s", err)
		t.FailNow()
	}
	t.Log("Distinct ok", result)
}

// go test -v ./coms/mgo -run TestCollection_Aggregate -args "-c=${PWD}/example/conf/app.toml"
func TestCollection_Aggregate(t *testing.T) {
	mgoClient := Ins()
	cursor, err := mgoClient.C("test").Aggregate(bson.A{bson.D{bson.E{Key: "$skip", Value: 1}}})
	if err != nil {
		t.Errorf("Aggregate error: %s", err)
		t.FailNow()
	}

	resultAll := bson.A{}
	if err := cursor.All(&resultAll); err != nil {
		t.Errorf("Aggregate error: %s", err)
		t.FailNow()
	}
	t.Log("Aggregate ok ", resultAll)
}

// go test -v ./coms/mgo -run TestCollection_BulkWrite -args "-c=${PWD}/example/conf/app.toml"
func TestCollection_BulkWrite(t *testing.T) {
	mgoClient := Ins()

	result, err := mgoClient.C("test").BulkWrite([]mongo.WriteModel{
		&mongo.InsertOneModel{Document: bson.M{"name": "x"}},
		&mongo.DeleteOneModel{Filter: bson.M{"name": "x"}},
	})
	if err != nil {
		t.Errorf("BulkWrite error: %s", err)
		t.FailNow()
	}
	t.Log("BulkWrite ok", bson.M{
		"inserted": result.InsertedCount,
		"deleted":  result.DeletedCount,
		"matched":  result.MatchedCount,
		"upserted": result.UpsertedCount,
		"modified": result.ModifiedCount,
	})
}

// go test -v ./coms/mgo -run TestCollection_CountDocuments -args "-c=${PWD}/example/conf/app.toml"
func TestCollection_CountDocuments(t *testing.T) {
	mgoClient := Ins()

	result, err := mgoClient.C("test").CountDocuments(bson.M{"name": "tom"})
	if err != nil {
		t.Errorf("CountDocuments error: %s", err)
		t.FailNow()
	}
	t.Log("CountDocuments ok, count: ", result)
}

// go test -v ./coms/mgo -run TestCollection_DeleteOne -args "-c=${PWD}/example/conf/app.toml"
func TestCollection_DeleteOne(t *testing.T) {
	mgoClient := Ins()

	result, err := mgoClient.C("test").DeleteOne(bson.M{"name": "tom"})
	if err != nil {
		t.Errorf("DeleteOne error: %s", err)
		t.FailNow()
	}
	t.Log("DeleteOne ok, delete count: ", result.DeletedCount)
}

// go test -v ./coms/mgo -run TestCollection_DeleteMany -args "-c=${PWD}/example/conf/app.toml"
func TestCollection_DeleteMany(t *testing.T) {
	mgoClient := Ins()

	result, err := mgoClient.C("test").DeleteMany(bson.M{"name": "henry"})
	if err != nil {
		t.Errorf("DeleteMany error: %s", err)
		t.FailNow()
	}
	t.Log("DeleteMany ok, delete count: ", result.DeletedCount)
}

// go test -v ./coms/mgo -run TestCollection_EstimatedDocumentCount -args "-c=${PWD}/example/conf/app.toml"
func TestCollection_EstimatedDocumentCount(t *testing.T) {
	mgoClient := Ins()

	result, err := mgoClient.C("test").EstimatedDocumentCount()
	if err != nil {
		t.Errorf("EstimatedDocumentCount error: %s", err)
		t.FailNow()
	}
	t.Log("EstimatedDocumentCount ok: ", result)
}

// go test -v ./coms/mgo -run TestCollection_FindOneAndUpdate -args "-c=${PWD}/example/conf/app.toml"
func TestCollection_FindOneAndUpdate(t *testing.T) {
	mgoClient := Ins()

	findResult := mgoClient.C("test").FindOneAndUpdate(bson.M{"name": "lily"}, bson.M{"$set": bson.M{"age": 17}})
	if findResult.Err() != nil {
		t.Errorf("FindOneAndUpdate error: %s", findResult.Err())
		t.FailNow()
	}
	result := bson.D{}
	if err := findResult.Decode(&result); err != nil {
		t.Errorf("FindOneAndUpdate decode error: %s", err)
		return
	}
	t.Log("FindOneAndUpdate ok: ", result)
}

// go test -v ./coms/mgo -run TestCollection_FindOneAndReplace -args "-c=${PWD}/example/conf/app.toml"
func TestCollection_FindOneAndReplace(t *testing.T) {
	mgoClient := Ins()

	findResult := mgoClient.C("test").FindOneAndReplace(bson.M{"name": "lily"}, bson.M{"name": "lily", "age": 18})
	if findResult.Err() != nil {
		t.Errorf("FindOneAndReplace error: %s", findResult.Err())
		t.FailNow()
	}
	result := bson.D{}
	if err := findResult.Decode(&result); err != nil {
		t.Errorf("FindOneAndReplace decode error: %s", err)
		t.FailNow()
	}
	t.Log("FindOneAndReplace ok: ", result)
}

// go test -v ./coms/mgo -run TestCollection_FindOneAndDelete -args "-c=${PWD}/example/conf/app.toml"
func TestCollection_FindOneAndDelete(t *testing.T) {
	mgoClient := Ins()

	findResult := mgoClient.C("test").FindOneAndDelete(bson.M{"name": "lily"})
	if findResult.Err() != nil {
		t.Errorf("FindOneAndDelete error: %s", findResult.Err())
		t.FailNow()
	}
	result := bson.D{}
	if err := findResult.Decode(&result); err != nil {
		t.Errorf("FindOneAndDelete decode error: %s", err)
		t.FailNow()
	}
	t.Log("FindOneAndDelete ok: ", result)
}

// go test -v ./coms/mgo -run TestCollection_ReplaceOne -args "-c=${PWD}/example/conf/app.toml"
func TestCollection_ReplaceOne(t *testing.T) {
	mgoClient := Ins()

	result, err := mgoClient.C("test").ReplaceOne(bson.M{"name": "sheldon"}, bson.M{"name": "lily", "age": 18})
	if err != nil {
		t.Errorf("ReplaceOne error: %s", err)
		t.FailNow()
	}
	t.Log("ReplaceOne ok: ", result)
}

// go test -v ./coms/mgo -run TestCollection_UpdateOne -args "-c=${PWD}/example/conf/app.toml"
func TestCollection_UpdateOne(t *testing.T) {
	mgoClient := Ins()

	result, err := mgoClient.C("test").UpdateOne(bson.M{"name": "sheldon"}, bson.M{"$set": bson.M{"age": 18}})
	if err != nil {
		t.Errorf("UpdateOne error: %s", err)
		t.FailNow()
	}
	t.Log("UpdateOne ok: ", result)
}

// go test -v ./coms/mgo -run TestCollection_UpdateMany -args "-c=${PWD}/example/conf/app.toml"
func TestCollection_UpdateMany(t *testing.T) {
	mgoClient := Ins()

	result, err := mgoClient.C("test").UpdateMany(bson.M{"name": bson.M{"$ne": ""}}, bson.M{"$set": bson.M{"age": 18}})
	if err != nil {
		t.Errorf("UpdateMany error: %s", err)
		t.FailNow()
	}
	t.Log("UpdateMany ok: ", result)
}

// go test -v ./coms/mgo -run TestCollection_Upsert -args "-c=${PWD}/example/conf/app.toml"
func TestCollection_Upsert(t *testing.T) {
	mgoClient := Ins()

	result, err := mgoClient.C("test").Upsert(bson.M{"name": "test"}, bson.M{"name": "test", "age": 18})
	if err != nil {
		t.Errorf("Upsert error: %s", err)
		t.FailNow()
	}
	t.Log("Upsert ok: ", result)
}
