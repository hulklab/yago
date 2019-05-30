package homemodel

import (
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/libs/date"
	"github.com/hulklab/yago/libs/orm"

	"github.com/hulklab/yago/example/app/app/modules/home/homedao"
)

type HomeModel struct {
}

func NewHomeModel() *HomeModel {
	return &HomeModel{}
}

func (m *HomeModel) Add(name string, options map[string]interface{}) (int64, yago.Err) {

	// 判断 name 是否已存在
	exist := &homedao.HomeDao{Name: name}

	orm.Ins().Get(exist)

	if exist.Id != 0 {
		return 0, yago.NewErr("用户名 " + name + " 已存在")
	}

	// 添加用户
	user := &homedao.HomeDao{
		Name:  name,
		Ctime: date.Now(),
	}

	_, err := orm.Ins().Insert(user)

	return user.Id, yago.NewErr(err)

}

func (m *HomeModel) UpdateById(id int64, options map[string]interface{}) (*homedao.HomeDao, yago.Err) {
	user := &homedao.HomeDao{Id: id}
	b, e := orm.Ins().Get(user)
	if e != nil {
		return nil, yago.NewErr(e)
	}

	if !b {
		return nil, yago.ErrNotFound
	}

	attrs := make(map[string]interface{})

	// 更新姓名
	name, ok := options["name"]
	if ok {
		user.Name = name.(string)
		attrs["name"] = name
	}

	if len(attrs) > 0 {
		_, err := orm.Ins().Table(user.TableName()).Where("id=?", id).Update(attrs)
		return user, yago.NewErr(err)
	}

	return user, yago.OK

}

func (m *HomeModel) DeleteById(id int64) (int64, yago.Err) {
	user := &homedao.HomeDao{Id: id}
	n, err := orm.Ins().Delete(user)
	return n, yago.NewErr(err)
}

func (m *HomeModel) GetDetail(id int64) *homedao.HomeDao {
	user := &homedao.HomeDao{Id: id}
	orm.Ins().Get(user)

	return user
}

func (m *HomeModel) GetList(q string, page, pagesize int) (int64, []*homedao.HomeDao) {

	var users []*homedao.HomeDao
	var total int64
	query := orm.Ins().NewSession()

	if q != "" {
		query.Where("name LIKE ?", "%"+q+"%")
	}

	query.Limit(pagesize, (page-1)*pagesize)
	query.OrderBy("id desc")

	total, _ = query.FindAndCount(&users)
	return total, users
}
