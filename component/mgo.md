# Mongo 组件
Mongo 组件我们依赖的开源包是 `go.mongodb.org/mongo-driver/mongo`。

按照组件的设计，我们定义了自己的 Mgo 结构对其进行了组合，在保留其原生的功能之外，以便扩展。

```go
// yago/coms/mgo/mgo.go
type Mgo struct {
	*mongo.Database
}

type Collection struct {
	c *mongo.Collection
}

type Cursor struct {
	c *mongo.Cursor
}
```

所以你可以查看 [mongo 官方文档](https://godoc.org/go.mongodb.org/mongo-driver/mongo) 来获取所有支持的 api。

本文中仅介绍部分常用的 api 以及扩展的 api。

## 配置 Mongo 组件
```toml
[mongodb]
# https://docs.mongodb.com/manual/reference/connection-string/
mongodb_uri = "mongodb://user:password@127.0.0.1:27017/?connectTimeoutMS=5000&socketTimeoutMS=5000&maxPoolSize=100"
database = "test"
```
我们在模版 app.toml 中默认配置开启了 mgo 组件，可根据实际情况进行调整。

## 使用 Mongo 组件
* 添加一行记录

```go
import "go.mongodb.org/mongo-driver/bson"

insertResult, err := mgo.Ins().C("test").InsertOne(bson.M{"name": "tom"})
```

* 添加多行记录

```go
insertResult, err := mgo.Ins().C("test").InsertMany(bson.A{
    bson.M{"name": "henry"},
    bson.M{"name": "lily"},
    bson.M{"name": "sheldon"},
})

```

* 查询一行记录

```go
findResult := mgo.Ins().C("test").FindOne(bson.M{"name": "henry"})

result := bson.M{}
err := findResult.Decode(&result)

```

* 查询多行记录

```go
cursor, err := mgo.Ins().C("test").Find(bson.M{})
defer cursor.Close()
resultAll := bson.A{}
if err := cursor.All(&resultAll); err != nil {

}
```






