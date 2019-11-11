# 控制器 Controller

 * [Http](controller/http.md)：Http 服务的控制器

 * [RPC](controller/rpc.md)：RPC 服务的控制器

 * [Cmd](controller/cmd.md)：命令行工具的控制器

 * [Task](controller/task.md)：定时任务和常驻进程的控制器

 控制器内的各个接口我们定义为 Action，并且对外暴露的各个接口的命名都是以 Action 作为后缀用来表明它是一个需要被外部调用的方法。这是一个不成文的规定，我们希望你也采用这种方式来使你的代码更清晰。