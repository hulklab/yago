# yago-locker

## config
```toml
[locker]
driver = "redis"
driver_instance_id = "redis"

[redis]
addr = "127.0.0.1:6379"
auth = "yourpass"
db = 0
max_idle = 5
idle_timeout = 30
```
>注: driver_instance_id 表示配置文件中的 redis 组件名
