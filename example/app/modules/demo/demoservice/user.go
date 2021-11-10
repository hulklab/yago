package demoservice

import (
	"github.com/hulklab/yago/example/app/g/gservice"
	"github.com/hulklab/yago/example/app/libs/trace"

	"fmt"

	"github.com/hulklab/yago/example/app/g/gmodel"
	"github.com/hulklab/yago/example/app/modules/demo/demodao"
	"github.com/hulklab/yago/example/app/modules/demo/demodto"
	"github.com/hulklab/yago/example/app/modules/demo/demomodel"
)

type userService struct {
	gservice.BaseService
}

func NewUserService(ctx *trace.Context) *userService {
	s := &userService{}
	s.Init(ctx)
	return s
}

func (s *userService) GetList(req *demodto.UserListReq) (resp *demodto.UserListResp, err error) {
	resp = &demodto.UserListResp{
		List: make([]*demodto.User, 0),
	}

	session := demomodel.NewUserModel(gmodel.WithCtx(s.Ctx)).GetSession()

	list := make([]*demodao.UserDao, 0)

	if len(req.Q) > 0 {
		// session.Where("? like ?","%"+req.Q+"%")
	}

	session.OrderBy("id desc")

	if req.PageNum > 0 && req.PageSize > 0 {
		session.Limit(req.PageSize, req.PageSize*(req.PageNum-1))
	}

	total, err := session.FindAndCount(&list)
	if err != nil {
		return
	}

	resp.Total = total

	for _, dao := range list {
		info := &demodto.User{
			UpdatedAt: dao.UpdatedAt,
			Username:  dao.Username,
			Avatar:    dao.Avatar,
			CreatedAt: dao.CreatedAt,
			Id:        dao.Id,
			Name:      dao.Name,
			Phone:     dao.Phone,
			Status:    dao.Status,
		}

		resp.List = append(resp.List, info)
	}
	return
}

func (s *userService) GetDetail(req *demodto.UserDetailReq) (resp *demodto.User, err error) {

	dao, err := demomodel.NewUserModel(gmodel.WithCtx(s.Ctx)).MustGetById(req.Id)
	if err != nil {
		return nil, err
	}

	info := &demodto.User{
		Status:    dao.Status,
		UpdatedAt: dao.UpdatedAt,
		Username:  dao.Username,
		Avatar:    dao.Avatar,
		CreatedAt: dao.CreatedAt,
		Id:        dao.Id,
		Name:      dao.Name,
		Phone:     dao.Phone,
	}

	return info, nil
}

func (s *userService) AddUser(req *demodto.UserAddReq) (resp *demodto.UserAddResp, err error) {

	dao := &demodao.UserDao{
		Phone:    req.Phone,
		Username: req.Username,
		Avatar:   req.Avatar,
		Name:     req.Name,
	}
	err = demomodel.NewUserModel(gmodel.WithCtx(s.Ctx)).InsertOne(dao)

	if err != nil {
		return nil, fmt.Errorf("insert record err:%w", err)
	}

	resp = &demodto.UserAddResp{
		Id: dao.Id,
	}

	return
}

func (s *userService) UpdateById(req *demodto.UserUpdateReq) (err error) {
	_, err = demomodel.NewUserModel(gmodel.WithCtx(s.Ctx)).MustGetById(req.Id)
	if err != nil {
		return err
	}

	dao := &demodao.UserDao{
		Phone:    req.Phone,
		Status:   req.Status,
		Username: req.Username,
		Avatar:   req.Avatar,
		Name:     req.Name,
	}

	err = demomodel.NewUserModel(gmodel.WithCtx(s.Ctx)).UpdateById(req.Id, dao)

	return
}

func (s *userService) DeleteById(req *demodto.UserDeleteReq) (err error) {

	_, err = demomodel.NewUserModel(gmodel.WithCtx(s.Ctx)).MustGetById(req.Id)
	if err != nil {
		return err
	}

	err = demomodel.NewUserModel(gmodel.WithCtx(s.Ctx)).DeleteById(req.Id)

	return
}
