## 数据映射 Dao

dao 用来存放数据库表的映射关系。除此之外还可以用来封装一些复杂的数据库操作方法。

xorm mysql `table` 表的定义

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

