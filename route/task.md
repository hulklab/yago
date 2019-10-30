## 6. 控制是否需要在此机器上开启 task 任务，有两种方式

* 修改配置文件中的 app.task_enable，默认为开启
* 修改环境变量 export {{配置文件中的app_name}}_APP_TASK_ENABLE=1, 1 表示开启，0 表示关闭，配置文件与环境变量同时存在时环境变量生效