# Kafka 组件
kafka 消息队列组件我们依赖的开源包是 `github.com/Shopify/sarama`。

按照组件的设计，我们定义了自己的 kafka 结构对其进行了组合，在保留其原生的功能之外，以便扩展。

```go
// yago/coms/kafka/kafka.go
type Kafka struct {
	connect string
	config  *cluster.Config
}
```

所以你可以查看 [sarama 官方文档](https://github.com/Shopify/sarama) 来获取所有支持的 api。

本文中仅介绍部分常用的 api 以及扩展的 api。

## 配置 kafka 组件
```toml
[kafka]
# 多个 broker 用逗号分隔
cluster = "127.0.0.1:9092"
topic = "demo"
```
我们在模版 app.toml 中默认配置开启了 kafka 组件，可根据实际情况进行调整。

## 使用 kafka 组件
* 生产消息
```go
kafka.Ins().Product("topic","value")
```

* 消费消息
```go
consumer, err := k.NewConsumer("topic","group")
if err == nil {
    // Consume 会一直阻塞等待消息，直到 err
    err := consumer.Consume(func(bytes []byte) {
        fmt.Println(string(bytes))
    })
    
    if err != nil {
        log.Println(err.Error())
    }
}
```
