# Elastic 组件

Elastic 组件我们依赖的开源包是 `github.com/olivere/elastic/v7`。

按照组件的设计，我们定义了自己的 Elastic 结构对其进行了组合，在保留其原生的功能之外，以便扩展。

```go
// yago/coms/elastic/elastic.go
type Elastic struct {
	*elastic.Client
}
```

所以你可以查看 [elastic 官方文档](https://github.com/olivere/elastic) 来获取所有支持的 api。

本文中仅介绍部分常用的 api 以及扩展的 api。

## 配置 Elastic 组件
```toml
[elastic]
urls = ["http://127.0.0.1:9200/", "http://127.0.0.1:9300"]
# username = ""
# password = ""
# sniff_enable = false
# 日志最低等级 Error = 2, Info = 4, Trace = 6
level = 6
```

## 使用 Elastic 组件

* 创建索引


```go
    exist, err := es.Ins().IndexExists("es_test").Do(context.Background())
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

	ret, err := es.Ins().CreateIndex("es_test").BodyString(mapping).Do(context.Background())

	if err != nil {
		panic(err)
	}

	fmt.Println(ret.Acknowledged)
```

* 删除索引

```go
	exist, err := es.Ins().IndexExists("es_test").Do(context.Background())
	if err != nil {
		panic(err)
	}

	if !exist {
		fmt.Println("索引不存在")
		return
	}

	ret, err := es.Ins().DeleteIndex("es_test").Do(context.Background())

	if err != nil {
		panic(err)
	}

	fmt.Println(ret.Acknowledged)
```

* 插入数据

```go
	msg := Tweet{User: "bob", Message: "hello"}
	// 不存在就创建，存在就更新
	ret, err := es.Ins().Index().Index("es_test").Id("1").BodyJson(msg).Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(ret.Id)
	fmt.Println(ret.Result)
```

* 查询数据

```go
	// 利用elastic/v7 的包构建 term 查询
	termQuery := elastic.NewTermQuery("user", "bob")

	ret, err := es.Ins().Search().
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
```

* 根据 id 更新

```go
	ret, err := es.Ins().Update().Index("es_test").Id("1").Doc(map[string]interface{}{"message": "haha"}).Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(ret.Id)
	fmt.Println(ret.Type)
	fmt.Println(ret.Result)
```

* 根据 id 删除数据

```go
	ret, err := es.Ins().Delete().Index("es_test").Id("1").Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(ret.Result)
```

* 根据查询语句删除数据

```go
	termQuery := elastic.NewTermQuery("user", "bob")
	ret, err := es.Ins().DeleteByQuery().Index("es_test").Query(termQuery).Do(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(ret.Total)
```

