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
