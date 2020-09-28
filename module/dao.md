# 数据访问对象 Dao

Dao 用来存放数据库表的映射关系。除此之外还可以用来封装一些复杂的数据库操作方法。

mysql `table` 表的定义，关于 xorm 的使用请参考 [xorm 组件](../component/orm.md)

```go
type UserDao struct {
	Id         int64  `json:"user_id" xorm:"autoincr pk"`
	Username   string `json:"username"`
	UserType   int    `json:"user_type"`
	Password   string `json:"-"`
	Avatar     string `json:"avatar"`
	Phone      string `json:"phone"`
	PhoneState int    `json:"phone_state"`
	UserState  int    `json:"user_state"`
	CreatedAt  string `json:"created_at" xorm:"created"`
	UpdatedAt  string `json:"updated_at" xorm:"updated"`
}

func (d *UserDao) TableName() string {
	return "user"
}
```

