package demodao

type UserDao struct {
	Avatar    string `json:"avatar"`                    //头像地址
	CreatedAt string `json:"created_at" xorm:"created"` //创建时间
	Id        int64  `json:"id" xorm:"autoincr pk"`     //
	Name      string `json:"name"`                      //名称
	Phone     string `json:"phone"`                     //手机
	Status    int    `json:"status"`                    //1可用
	UpdatedAt string `json:"updated_at" xorm:"updated"` //更新时间
	Username  string `json:"username"`                  //用户名
}

func (t *UserDao) TableName() string {
	return "user"
}
