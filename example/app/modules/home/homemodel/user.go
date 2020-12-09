package homemodel

import (
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/base/basemodel"
	"github.com/hulklab/yago/coms/orm"
	"github.com/hulklab/yago/example/app/g"
	"github.com/hulklab/yago/example/app/modules/home/homedao"
	"github.com/hulklab/yago/libs/date"
)

type UserModel struct {
	basemodel.BaseModel
}

func NewUserModel(opts ...basemodel.Option) *UserModel {
	model := &UserModel{}

	opts = append(
		opts,
		basemodel.WithDefaultPageSize(1, 20),
		basemodel.WithDefaultOrder(
			basemodel.Order{"created_at": -1},
		),
	)

	model.Init(opts...)

	return model
}

func (m *UserModel) Add(username, phone string, options map[string]interface{}) (int64, error) {
	// 判断 name 是否已存在
	exist := &homedao.UserDao{Username: username}

	_, _ = m.GetSession().Get(exist)

	if exist.Id != 0 {
		return 0, yago.NewErr("用户名 " + username + " 已存在")
	}

	// 添加用户
	user := &homedao.UserDao{
		Username:  username,
		Phone:     phone,
		CreatedAt: date.Now(),
	}

	_, err := m.GetSession().Insert(user)
	if err != nil {
		return 0, yago.WrapErr(yago.ErrSystem, err)
	}

	return user.Id, nil
}

func (m *UserModel) UpdateById(id int64, options map[string]interface{}) (*homedao.UserDao, error) {
	user := &homedao.UserDao{Id: id}

	exist, err := orm.Ins().Get(user)
	if err != nil {
		return nil, yago.WrapErr(yago.ErrSystem, err)
	}

	if !exist {
		return nil, yago.NewErr("用户不存在")
	}

	attrs := g.Hash{}

	// 更新姓名
	username, ok := options["username"]
	if ok {
		user.Username = username.(string)
		attrs["username"] = username
	}

	if len(attrs) > 0 {
		_, err := orm.Ins().Table(user.TableName()).ID(id).Update(attrs)
		if err != nil {
			return nil, yago.WrapErr(yago.ErrSystem, err)
		}
	}

	return user, nil
}

func (m *UserModel) DeleteById(id int64) (int64, error) {
	user := &homedao.UserDao{Id: id}

	deleted, err := orm.Ins().Delete(user)
	if err != nil {
		return 0, yago.WrapErr(yago.ErrSystem, err)
	}

	return deleted, nil
}

func (m *UserModel) GetDetail(id int64) (*homedao.UserDao, error) {
	user := &homedao.UserDao{Id: id}

	_, err := orm.Ins().Get(user)
	if err != nil {
		return nil, yago.WrapErr(yago.ErrSystem, err)
	}

	return user, nil
}

func (m *UserModel) GetList(q string, page, pageSize int) (int64, []*homedao.UserDao) {
	var users []*homedao.UserDao
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
