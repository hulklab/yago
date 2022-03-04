package demomodel

import (
	"github.com/hulklab/yago/example/app/g/gmodel"

	"fmt"
	"github.com/hulklab/yago/example/app/g"
	"github.com/hulklab/yago/example/app/modules/demo/demodao"
)

type userModel struct {
	gmodel.BaseModel
}

func NewUserModel(opts ...gmodel.Option) *userModel {
	m := &userModel{}
	m.Init(opts...)
	return m
}

func (m *userModel) tableName() string {
	return new(demodao.UserDao).TableName()
}

func (m *userModel) InsertOne(dao *demodao.UserDao) (err error) {
	_, err = m.GetSession().InsertOne(dao)

	return
}

func (m *userModel) GetById(id int64) (b bool, dao *demodao.UserDao, err error) {
	dao = &demodao.UserDao{}

	b, err = m.GetSession().Where("id = ?", id).Get(dao)
	if err != nil || !b {
		return b, nil, err
	}

	return b, dao, err
}

func (m *userModel) MustGetById(id int64) (dao *demodao.UserDao, err error) {
	dao = &demodao.UserDao{}

	b, err := m.GetSession().Where("id = ?", id).Get(dao)
	if err != nil {
		return nil, err
	}

	if !b {
		return nil, fmt.Errorf("user %d is not exist", id)
	}

	return dao, err
}

func (m *userModel) GetListByIds(ids []int64) (list []*demodao.UserDao, err error) {
	list = make([]*demodao.UserDao, 0)

	err = m.GetSession().In("id", ids).Find(&list)

	return list, err
}

func (m *userModel) GetAll() (list []*demodao.UserDao, err error) {
	list = make([]*demodao.UserDao, 0)
	err = m.GetSession().Find(&list)

	return list, err
}

func (m *userModel) GetAllMap() (allMap map[int64]*demodao.UserDao, err error) {
	allMap = make(map[int64]*demodao.UserDao)
	err = m.GetSession().Find(&allMap)

	return allMap, err
}

func (m *userModel) DeleteById(id int64) (err error) {
	session := m.GetSession()
	session.Where("id=?", id)

	_, err = session.Delete(new(demodao.UserDao))
	return
}

func (m *userModel) DeleteByIds(ids []int64) (err error) {
	session := m.GetSession()
	session.In("id", ids)

	_, err = session.Delete(new(demodao.UserDao))
	return
}

func (m *userModel) UpdateAttrsById(id int64, attrs g.Hash) (err error) {
	if attrs == nil {
		return fmt.Errorf("attrs cannot be nil, id:%d", id)
	}

	_, err = m.GetSession().Table(m.tableName()).Where("id=?", id).Update(attrs)

	return
}
