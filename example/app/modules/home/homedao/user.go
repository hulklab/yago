package homedao

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
