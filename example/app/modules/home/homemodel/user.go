package homemodel

import (
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/base/basemodel"
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

func (m *UserModel) Add(username, phone string) (int64, error) {
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
