# ORM 组件
ORM 组件我们依赖的开源包是 `github.com/go-xorm/xorm`。

按照组件的设计，我们定义了自己的 ORM 结构对其进行了组合，在保留其原生的功能之外，以便扩展。

```go
// yago/coms/orm/orm.go
type Orm struct {
	*xorm.Engine
}
```

所以你可以查看 [xorm 官方文档](http://gobook.io/read/gitea.com/xorm/manual-zh-CN/) 来获取所有支持的 api。

本文中仅介绍部分常用的 api 以及扩展的 api。


## 配置 ORM 组件
```toml
[db]
host = "127.0.0.1"
user = "user"
password = "password"
port = "3306"
database = "db"
prefix =""
timezone = "Asia/Shanghai"
charset = "utf8"
max_life_time = 8 # 连接最大生命周期 8s
max_idle_conn = 20 # 最多空闲连接数 20
max_open_conn = 500 # 最大打开连接数
show_log = true # 是否开启日志
```
我们在模版 app.toml 中默认配置开启了 ORM 组件，可根据实际情况进行调整。

## 使用 ORM 组件
### 定义表结构
```go
// 下文以 user 表为例
type UserDao struct {
	Id    int64  `json:"id" xorm:"autoincr"`
	Name  string `json:"name"`
	Ctime string `json:"ctime" xorm:"created"`
	Utime string `json:"utime" xorm:"updated"`
}

func (d *UserDao) TableName() string {
	return "user"
}
```

### 新增记录
* 新增一行

```go
user := &homedao.UserDao{
    Name:  "zhangsan",
}

_, err := orm.Ins().Insert(user)
```
> 调用 orm.Ins() 会返回 orm.Orm 实例，因其对 xorm.Engine 进行了扩展，我们可以使用 xorm.Engine 的所有 api。

* 新增多行

```go
users := []*homedao.UserDao{
    {Name:"lisi"},
    {Name:"wangwu"},
}
_, err := orm.Ins().Insert(users)
```

### 查询单行记录
```go
user := &homedao.UserDao{Id: id}
exists, err := orm.Ins().Get(user)
```

### 删除
* 删除一行

```go
user := &homedao.UserDao{Id: id}
n, err := orm.Ins().Delete(user)
```

* 删除多行

```go
n,err := orm.Ins().Where("id > ?", 100).Delete(new(homedao.UserDao))
```

### 修改
* 推荐用法

```go
// 推荐用法
attrs := g.Hash{}
attrs["name"] = "hello"

tableName := new(homedao.UserDao).TableName()

_, err := orm.Ins().Table(tableName).Where("id=?",id).Update(attrs)
```

* 简单用法

```go
user := &homedao.UserDao{Id: id}
user.Name = "kowloon"
n, err := orm.Ins().Update(user)
```
> 注意：简单用法里面，如果赋值为字段的零值时（整数的 0, 字符串的 空串）, orm 会忽略掉不更新

### 查询
* 查询多行多列

```go
users := make([]*homedao.UserDao, 0)
err := orm.Ins().Table(new(homedao.UserDao).TableName()).
		Select("*").
		Where("id > ?", 10).
		Find(&users)
```

* 查询多行单列

```go
userIds := make([]int64, 0)
err := orm.Ins().Table(new(homedao.UserDao).TableName()).
		Cols("id").
		Find(&userIds)

```

* 查询总数

```go
total, err := orm.Ins().Table(new(homedao.UserDao)).
            Where("id  > ?", 10).
            Count()
```

* 使用原生 sql

```go
sql := "select * from user"
results, err := orm.Ins().Query(sql)
```
当调用 Query 时，第一个返回值 results 为 []map[string][]byte 的形式。

### 分页 + Join + 排序查询
* 新增表结构

```go
// 此处新增一个表结构 user_detail, user_id 与 user 表关联
type DetailDao struct {
	Id    int64  `json:"id" xorm:"autoincr"`
	UserId  int64 `json:"user_id"`
	Phone string `json:"phone"`
}

func (d *DetailDao) TableName() string {
	return "user_detail"
}

// 再扩展一个新的结构体，来用存放连表查询的结构
type UserDetail struct {
    UserDao `xorm:"extends"`
    Phone string `json:"phone"`
}
```

* 查询

```go
userDetails := make([]*homedao.UserDetail, 0)

query := orm.Ins().Table(new(homedao.UserDao)).Alias("t").
		Select("t.*,d.phone").
		Join("INNER", []string{new(homedao.DetailDao).TableName(), "d"}, "d.user_id=t.id").
		Where("t.id > 10", 0)

query.Limit(pagesize,(page-1)*pagesize)
query.OrderBy("t.id desc")

// total 为总页数，userDetails 为结果集列表
total, err := query.FindAndCount(&userDetails)
```


### 事务
yago 扩展了 xorm 的事务功能，增加了一个 Transactional 方法
```go
err := orm.Ins().Transactional(func(session *xorm.Session) error {
    user := &homedao.UserDao{
        Name: "law",
    }
    if _, err := session.Insert(user);err != nil{
        // log
        return err
    }
    userDetail := &homedao.DetailDao{
        UserId: user.Id,
        Phone: "13800000000",
    }
    
    if _, err := session.Insert(userDetail);err != nil{
        // log
        return err
    }
)
``` 