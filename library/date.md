## 日期
go 的时间库已经很强大了，yago 基于 go 的时间库封装了类似 php 的 date 和 strtotime 函数

使用方式如下：

```go
// 打印当前时间的 2006-01-02 15:04:05 格式
dateTime := date.Date("Y-m-d H:i:s")
dateTime2 := date.Now()

// 打印当前时间的 2006-01-02 15:04:05.000 格式
dateTime := date.Date("Y-m-d H:i:J")

// 打印当前年
year := date.Date("Y")

// 打印当前月
month := date.Date("m")
monthShort := date.Date("n")

// 打印当前日
day := date.Date("d")
dayShort := date.Date("j")

// 打印当前小时
hour := date.Date("H")

// 打印当前分钟
minute := date.Date("i")

// 打印当前秒
second := date.Date("s")

// 打印当前微秒
millis := date.Date("Q")

// 打印当前纳秒
nana := date.Date("K")

// 打印某个时间戳的年月日
someDay := date.Date("Y-m-d H:i:s", 1573113481)

// 返回某个日期对应的时间戳
timestamp1 := date.Strtotime("2019-11-07 16:08:00","Y-m-d H:i:s")
timestamp2 := date.Strtotime("2019-11-07","Y-m-d")
```
