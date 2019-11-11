# 模块 

Yago 在最上层采用了 Module（模块）的分层设计，用来应对复杂业务的场景上的逻辑功能的划分。你可以理解为，Module 是一个 yago 程序核心的业务逻辑的实现。所以我们就从一个 module 开始完整地介绍一个 yago 程序的编写。默认用 yago 生成的样板代码自带了一个 home 模块，里面集成了 yago 的一些编写规范样例，比如 controller，model，dao，third-api 的实现，是我们写代码的很好的参考。当然你也可以通过下面这条命令来创建一个你自己的模块。

```bash
$ yago new -m mymodule
2019/11/03 17:15:41 create module mymodule
```

## 模块内的 MVC

Yago 在 Module（模块）内采用了分层结构设计，并且提供了四种类型的 Controller（控制器），分别是 Http，GRPC，Cmd，Task。 Model 层用来实现主要的业务逻辑方法和数据访问。Dao 层用来实现数据库表的结构定义与映射。

下面我们分别介绍

* [控制器](controller.md)
* [模型](model.md)
* [数据访问对象](dao.md)