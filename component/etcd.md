# Etcd 组件
etcd 分布式键值数据库组件我们依赖的开源包是 `go.etcd.io/etcd/clientv3`。

按照组件的设计，我们定义了自己的 etcd 结构对其进行了组合，在保留其原生的功能之外，以便扩展。

```go
// yago/coms/etcd/etcd.go
type Etcd struct {
	*clientv3.Client
}
```

所以你可以查看 [etcd-client 官方文档](https://godoc.org/go.etcd.io/etcd/clientv3) 来获取所有支持的 api。

本文中仅介绍部分常用的 api。

## 配置 etcd 组件

```toml
[etcd]
endpoints    = ["127.0.0.1:2379"]
dial_timeout = 2
# username = ""
# password = ""
# max_call_recv_msgsize_byte = 0  # default 2 * 1024 * 1024
# max_call_send_msgsize_byte = 0  # default math.MaxInt32
# cert_file = "conf/ssl/etcd_cert.pem"
# cert_key_file = "conf/ssl/etcd_cert_key.pem"
# cert_ca_file = "conf/ssl/etcd_ca.pem"
```

我们在模版 app.toml 中默认配置开启了 etcd 组件，可根据实际情况进行调整。

## 使用 etcd 组件
* 设置值 put

```go
_, err := etcd.Ins().Put(context.Background(),"key","value")
```

* 获取值 get

```go
_, err := etcd.Ins().Get(context.Background(),"key")

if err != nil {
    // err handler
}

if len(res.Kvs) == 0{
	return
}

for _, item := range res.Kvs {
    fmt.Println(string(item.Key), "=>", string(item.Value))
}

```

* 根据前缀获取值 get with prefix by desc order

```go
prefix := "k"
res, err := etcd.Ins().Get(context.Background(),prefix,clientv3.WithPrefix(),clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))

if err != nil {
    // err handler
}

if len(res.Kvs) == 0{
	return
}

for _, item := range res.Kvs {
    fmt.Println(string(item.Key), "=>", string(item.Value))
}
```

* 删除值 delete

```go
_, err := etcd.Ins().Delete(context.Background(),"key")
```

* 监听值的变化 watch

```go
watchChan := etcd.Ins().Watch(context.Background(), "", clientv3.WithPrefix()) // watch all key

go func() {
    for {
        msg := <-watchChan
        for _, event := range msg.Events {
            if event.Type == clientv3.EventTypePut {
                fmt.Println("watch:", string(event.Kv.Key), "=>put=> ", string(event.Kv.Value))
            } else if event.Type == clientv3.EventTypeDelete {
                fmt.Println("watch:", string(event.Kv.Key), "=>delete=> ", string(event.Kv.Value))
            }

        }
    }
}()
```
