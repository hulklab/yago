# 字符串函数

* md5 加密字符串
```go
// output: 29ddc288099264c17b07baf44d3f0adc
md5Value := str.Md5("kate")
```

* sha1 加密字符串
```go
// 986641af5eabfb409d9ba99ba07f746ce8fdc2af
sha1Value := str.Sha1("poly")
```

* 生成全局唯一 id
```go
// output: 36d9b8ef728ff3c09c15cfbbe0fdcfee
uniqId := str.UniqueId()
```

* 生成全局唯一短 id
```go
// output: 3LwHM98LvrT
shortId := str.UniqueIdShort()
```

* 下划线转驼峰
```go
// output: AaBbCc
camel := str.CamelString("aa_bb_cc")
```

* 首字母大写
```go
// output: Hello
uc := str.Ucfirst("hello")
```

* map[string]string -> string
```go
// output: k1=v1&k2=v2
m := map[string]string{"k1":"v1","k2":"v2"}
s := str.Kv2str(m,"=","&")
```

* string -> map[string]string
```go
// output: map[k1:v1 k2:v2]
s := "k1=v1&k2=v2"
str.Str2kv(s,"=","&")
```

* 按逗号，空格，换行，tab 等分隔字符串
```go
// output: [a,b,c,d]
s := "a, b,c\n, d\t"
ss := Split(s)
fmt.Println(ss)
```


