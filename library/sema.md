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

除了上面的问题外，当多个goroutine并发地执行时，一旦某个goroutine报错，我们有的时候需要的，不是静静地等待所有goroutine执行完成，而是关闭并发的goroutine并及时返回错误。
使用方式如下：

```go
// 构造一个并发 4 个 goroutine 的 sema
sema := New(4)
// 添加 3个不会触发err的goroutine
for i := 0; i < 3; i++ {
    sema.Add(func() error {
        time.Sleep(time.Second * 2)
        return nil
    })
}

// 添加一个可以触发errr的 goroutine
// 触发了err后，会关闭并发中goroutine并及时返回错误
sema.Add(func() error {
    return errors.New("occur error")
})

// 等待所有的 goroutine 任务执行结束
err := sema.Wait()
if err != nil {
    fmt.Println(err) //打印错误
}
```