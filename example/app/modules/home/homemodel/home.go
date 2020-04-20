package homemodel

import (
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/coms/orm"
	"github.com/hulklab/yago/example/app/g"
	"github.com/hulklab/yago/libs/date"

	"github.com/hulklab/yago/example/app/modules/home/homedao"
)

type HomeModel struct {
}

func NewHomeModel() *HomeModel {
	return &HomeModel{}
}

func (m *HomeModel) Add(name string, options map[string]interface{}) (int64, error) {

	// 判断 name 是否已存在
	exist := &homedao.HomeDao{Name: name}

	_, _ = orm.Ins().Get(exist)

	if exist.Id != 0 {
		return 0, yago.NewErr("用户名 " + name + " 已存在")
	}

	// 添加用户
	user := &homedao.HomeDao{
		Name:  name,
		Ctime: date.Now(),
	}

	_, err := orm.Ins().Insert(user)

	return user.Id, yago.WrapErr(yago.ErrSystem, err)

}

func (m *HomeModel) UpdateById(id int64, options map[string]interface{}) (*homedao.HomeDao, error) {
	user := &homedao.HomeDao{Id: id}
	b, e := orm.Ins().Get(user)
	if e != nil {
		return nil, yago.WrapErr(yago.ErrSystem, e)
	}

	if !b {
		return nil, yago.NewErr("用户不存在")
	}

	attrs := g.Hash{}

	// 更新姓名
	name, ok := options["name"]
	if ok {
		user.Name = name.(string)
		attrs["name"] = name
	}

	if len(attrs) > 0 {
		_, err := orm.Ins().Table(user.TableName()).Where("id=?", id).Update(attrs)
		return user, yago.WrapErr(yago.ErrSystem, err)
	}

	return user, nil

}

func (m *HomeModel) DeleteById(id int64) (int64, error) {
	user := &homedao.HomeDao{Id: id}
	n, err := orm.Ins().Delete(user)
	if err != nil {
		return 0, yago.WrapErr(yago.ErrSystem, err)
	}
	return n, nil
}

func (m *HomeModel) GetDetail(id int64) *homedao.HomeDao {
	user := &homedao.HomeDao{Id: id}
	_, err := orm.Ins().Get(user)
	if err != nil {
		return nil
	}
	return user
}

func (m *HomeModel) GetList(q string, page, pageSize int) (int64, []*homedao.HomeDao) {

	var users []*homedao.HomeDao
	var total int64
	query := orm.Ins().NewSession()

	if q != "" {
		query.Where("name LIKE ?", "%"+q+"%")
	}

	query.Limit(pageSize, (page-1)*pageSize)
	query.OrderBy("id desc")

	total, _ = query.FindAndCount(&users)
	return total, users
}
