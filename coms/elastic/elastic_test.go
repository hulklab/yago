package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic/v7"
	"testing"
)

// 参考地址 https://olivere.github.io/elastic/

// Tweet is a structure used for serializing/deserializing data in Elasticsearch.
type Tweet struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

// 获取所有的 index
// go test -v ./coms/elastic -run TestCatIndex -args "-c=${PWD}/example/conf/app.toml"
func TestCatIndex(t *testing.T) {
	indices, err := Ins().CatIndices().Do(context.Background())

	if err != nil {
		fmt.Printf("cat index err :%s", err)
		return
	}

	for _, index := range indices {
		fmt.Printf("index name : %s,  docs total %d, size %s \n", index.Index, index.DocsCount, index.StoreSize)
	}
}

// 创建 index
// go test -v ./coms/elastic -run TestCreateIndex -args "-c=${PWD}/example/conf/app.toml"
func TestCreateIndex(t *testing.T) {
	exist, err := Ins().IndexExists("es_test").Do(context.Background())
	if err != nil {
		panic(err)
	}

	if exist {
		fmt.Println("索引已存在")
		return
	}

	mapping := `
{
    "settings":{
        "number_of_shards":1,
        "number_of_replicas":0
    },
    "mappings":{
        "properties":{
            "user":{
                "type":"keyword"
            },
            "message":{
                "type":"keyword"
            }
        }
    }
}`

	ret, err := Ins().CreateIndex("es_test").BodyString(mapping).Do(context.Background())

	if err != nil {
		panic(err)
	}

	fmt.Println(ret.Acknowledged)
}

// 删除 index
// go test -v ./coms/elastic -run TestDeleteIndex -args "-c=${PWD}/example/conf/app.toml"
func TestDeleteIndex(t *testing.T) {
	exist, err := Ins().IndexExists("es_test").Do(context.Background())
	if err != nil {
		panic(err)
	}

	if !exist {
		fmt.Println("索引不存在")
		return
	}

	ret, err := Ins().DeleteIndex("es_test").Do(context.Background())

	if err != nil {
		panic(err)
	}

	fmt.Println(ret.Acknowledged)
}

// 插入数据
// go test -v ./coms/elastic -run TestInsert -args "-c=${PWD}/example/conf/app.toml"
func TestInsert(t *testing.T) {
	msg := Tweet{User: "bob", Message: "hello"}
	// 不存在就创建，存在就更新
	ret, err := Ins().Index().Index("es_test").Id("1").BodyJson(msg).Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(ret.Id)
	fmt.Println(ret.Result)
}

// 查询数据
// go test -v ./coms/elastic -run TestTermQuery -args "-c=${PWD}/example/conf/app.toml"
func TestTermQuery(t *testing.T) {
	// 利用elastic/v7 的包构建 term 查询
	termQuery := elastic.NewTermQuery("user", "bob")

	ret, err := Ins().Search().
		Index("es_test"). // use index
		Query(termQuery). // termQuery
		From(0).Size(10). // pagesize
		Do(context.Background()) // execute
	if err != nil {
		fmt.Printf("term query err :%s", err)
		return
	}

	// 耗时
	fmt.Printf("term query took %d milliseconds\n", ret.TookInMillis)

	// 总数
	fmt.Printf("found total of %d \n", ret.TotalHits())

	// list
	if ret.TotalHits() > 0 {
		for _, hit := range ret.Hits.Hits {
			s, err := hit.Source.MarshalJSON()
			fmt.Println(hit.Type, hit.Id, string(s), err)

			var d Tweet
			err = json.Unmarshal(hit.Source, &d)

			fmt.Printf("err %+v, %+v\n", err, d)
		}
	}
}

// 根据 id 更新
// go test -v ./coms/elastic -run TestUpdateById -args "-c=${PWD}/example/conf/app.toml"
func TestUpdateById(t *testing.T) {
	ret, err := Ins().Update().Index("es_test").Id("1").Doc(map[string]interface{}{"message": "haha"}).Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(ret.Id)
	fmt.Println(ret.Type)
	fmt.Println(ret.Result)
}

// 根据 id 删除数据
// go test -v ./coms/elastic -run TestDeleteById -args "-c=${PWD}/example/conf/app.toml"
func TestDeleteById(t *testing.T) {
	ret, err := Ins().Delete().Index("es_test").Id("1").Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(ret.Result)
}

// 根据查询语句删除数据
// go test -v ./coms/elastic -run TestDeleteByQuery -args "-c=${PWD}/example/conf/app.toml"
func TestDeleteByQuery(t *testing.T) {
	termQuery := elastic.NewTermQuery("user", "bob")
	ret, err := Ins().DeleteByQuery().Index("es_test").Query(termQuery).Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(ret.Total)
}
