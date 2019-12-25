# Locker 组件
Locker 组件实现了分布式锁，我们没有直接集成到 yago 里面，而是放在 `github.com/hulklab/yago-coms` 仓库里，
yago-coms 中的组件都是独立的 go.mod，你可以使用 `go get github.com/hulklab/yago-coms/locker` 单独加载此组件

locker 组件采专用 driver 的机制，目前支持的 driver 有 redis。

每个 driver 都需要实现下面 ILocker 接口的两个方法：

```go
// yago-coms/locker/locker.go
type ILocker struct {
	Lock(key string, timeout int64) error
	Unlock()
}
```


下面介绍具体的使用方式。

## 配置 locker 组件
### 使用 redis 做 locker 驱动的组件
```toml
[redis]
addr = "127.0.0.1:6379"
auth = "yourpass"
db = 0
max_idle = 5
idle_timeout = 30

[locker]
driver = "redis"
driver_instance_id = "redis"
```
locker 配置中只需要配置 driver 名称和 driver 实例 ID，
上例中 driver_instance_id 为 redis 表示使用配置文件中 section 为 redis 的组件实例作为 locker 的 driver 驱动对象

### 使用 etcd 做 locker 驱动的组件
```toml
[etcd]
endpoints = ["127.0.0.1:2379"]

[locker]
driver = "etcd"
driver_instance_id = "etcd"
```
上例中 driver_instance_id 为 etcd 表示使用配置文件中 section 为 etcd 的组件实例作为 locker 的 driver 驱动对象

## 使用 Locker 组件

```go
key := "locker_name"
// 锁的失效时间
timeout := 10

// 创建锁对象
lo := locker.New()

// Lock() will block until get the locker
// Get locker
err := lo.Lock(key, timeout)
if err != nil {
	// err handler
}
// Release locker
defer lo.Unlock()
// your code here

```
