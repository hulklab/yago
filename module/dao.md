# 数据访问对象 Dao

Dao 用来存放数据库表的映射关系。除此之外还可以用来封装一些复杂的数据库操作方法。

mysql `table` 表的定义，关于 xorm 的使用请参考 [xorm 组件](../component/orm.md)

```go
package homedao

type HomeDao struct {
	Id    int64  `json:"id" xorm:"autoincr"`
	Name  string `json:"name"`
	Ctime string `json:"ctime" xorm:"created"`
	Utime string `json:"utime" xorm:"updated"`
}

func (d *HomeDao) TableName() string {
	return "table"
}

```

