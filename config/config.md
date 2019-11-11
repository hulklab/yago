# 配置

yago 模版生成的项目默认采用toml格式的配置文件，配置文件在程序启动后便开始解析，解析完成后存储在全局yago.Config中，yago.Config是 [viper](https://github.com/spf13/viper) 的扩展，原生采用 [viper](https://github.com/spf13/viper) 的方法来获取配置文件的值即可。

example:

```go
// 字符串
confString := yago.Config.GetString("conf_string")
// 整数
confInt := yago.Config.GetInt("conf_int")
// 浮点数
confFloat64 := yago.Config.GetFloat64("conf_float64")
// 持续时间
confDuration := yago.Config.GetDuration("conf_duration")
```

* [配置文件详细介绍](detail.md)
