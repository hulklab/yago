package homedao

type HomeDao struct {
	Id        int64  `json:"id" xorm:"autoincr pk"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at" xorm:"created"`
	UpdatedAt string `json:"updated_at" xorm:"updated"`
}

func (d *HomeDao) TableName() string {
	return "home"
}
