# 日志组件
日志组件我们依赖开源的包 `github.com/sirupsen/logrus`，按照组件的设计，我们对其进行了组合，在保留其原生的功能之外，以便扩展
```
// yago/coms/logger/logger.go
type Logger struct {
	*logrus.Logger
}
```

## 配置日志组件
```
[logger]
# json | text, default text
formatter = "json"
# 日志最低等级 Panic = 0, Fatal = 1, Error = 2, Warn = 3, Info = 4, Debug = 5, Trace = 6
level = 5
# 文件路径
file_path = "./logs/app.log"
# 最大保留的备份数
max_backups = 20
# 日志最大保留天数
max_age = 30
# 文件最大大小(mb)
max_size = 500
# 是否开启压缩
compress = true
# 是否开启终端输出
# stdout_enable = true
```
我们在模版 app.toml 中默认配置开启了日志组件，可根据实际情况进行调整。

## 使用
```
logger.Ins().Info("this is a info level msg")
```
在任意地方执行 `logger.Ins()`，会实例化 logger 组件到 yago 的全局的 Components 安全 map 中，