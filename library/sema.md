# goroutine 限流
goroutine 虽好，有时候使用过当也会造成 goroutine 泄露问题，semalib 的作用是用来控制并发 goroutine 的个数，
这在 for 循环里使用 goroutine 时特别管用。

使用方式如下：

```go
// 构造一个并发 3 个 goroutine 的 sema
sema := semalib.New(3)
for i := 0; i < 10; i++ {
    sema.Acquire() // 数量不足，阻塞等待
    go func() {
        defer sema.Release()
        time.Sleep(time.Second)
    }()
}
// 等待所有的 goroutine 任务执行结束
sema.Wait()
```
