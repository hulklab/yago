# 组件
组件解决的是项目中依赖开源社区第三方包的问题，yago 组件的中心思想是在完整保留第三方包的所有功能基础上，也能扩展，并对访问者提供简便统一的操作接口。
应该说组件是 yago 最核心的思想之一。

### 组件原理
为了方便用户使用组件，我们定义了一个全局的变量：`yago.Component`

```
type components struct {
	m sync.Map
}

func (c *components) Ins(key string, f func() interface{}) interface{} {
	v, ok := c.m.Load(key)
	if !ok {
		val := f()
		v, _ = c.m.LoadOrStore(key, val)
	}
	return v

}

// 全局变量 Component 供外部调用
var Component = new(components)
```

可以看到, components 的本质是一个安全 Map，我们把所有的组件都放到 map 里面，每个组件都会有一个对应的 key，
这个 key，通常采用配置文件中的块名称，如下面日志组件 logger 的配置:
```
# logger 会做为组件的 Key，logger 也是日志组件的默认名称，修改 logger 名字的话，调用时需要指定修改后的名称
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
```

### 实现组件
每个组件需要实现 Ins() 方法，在自己的 Ins 方法中调用 yago.Component.Ins()，注册到全局组件中。
以 logger 组件为例：
```
// Logger 是对 logrus.Logger 的扩展
type Logger struct {
    *logrus.Logger
}

func Ins(id ...string) *Logger {
    var name string

    if len(id) == 0 {
        name = "logger"
    } else if len(id) > 0 {
        name = id[0]
    }

    // 调用 yago.Component.Ins() 注册到全局组件容器中
    v := yago.Component.Ins(name, func() interface{} {
        conf := yago.Config.GetStringMap(name)
        
        formatter := conf["formatter"].(string)
        filePath := conf["file_path"].(string)
        maxSize := int(conf["max_size"].(int64))
        maxBackups := int(conf["max_backups"].(int64))
        maxAge := int(conf["max_age"].(int64))
        level := logrus.Level(conf["level"].(int64))
        compress := conf["compress"].(bool)

        val := &Logger{logrus.New()}
        // 设置最低log level
        val.SetLevel(level)

        // 日志中显示记录的文件和函数名, 注意：textField 中需要避开 file 和 func 字段
        val.SetReportCaller(true)

        if formatter == "json" {
            val.Formatter = &logrus.JSONFormatter{CallerPrettyfier: CallerPretty}
        } else {
            val.Formatter = &logrus.TextFormatter{CallerPrettyfier: CallerPretty}
        }
        val.Out = &lumberjack.Logger{
            Filename:   filePath,
            MaxSize:    maxSize,
            MaxBackups: maxBackups,
            MaxAge:     maxAge,
            Compress:   compress,
        }

        return val
    })

    return v.(*Logger)
}
```


### 使用组件
组件只有第一次调用 Ins() 时才会初始化组件，所以可以理解为每个组件都是一个单例。
还是以 logger 组件为例：
```
logger.Ins().Info("this is a test msg")
```
logger.Ins() 返回的是 yago 定义的 Logger，同时它又可以完整的使用 logrus.Logger 的方法。从上面定义 logger 的 Ins() 方法可以看出，Ins() 是可以传参的，参数就是全局组件容器里面的 key，
这在项目需要多个同类型组件时是十分有用的，比如使用多个数据库连接，我们可以通过调用不同的 key 获取不同的数据库连接。

### 组件组成
yago 的组件分两部分，一部分十分常用的组件我们放在 `yago/coms` 下，
还有一部分我们认为并不是每个项目都需要，单独开源了一个项目，`https://github.com/hulklab/yago-coms`，里面的每个组件可以单独下载。

 