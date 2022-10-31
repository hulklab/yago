package demohttp

import (
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/example/app/g/ghttp"

	"github.com/hulklab/cast"
	"github.com/hulklab/yago/example/app/g"
	"github.com/hulklab/yago/example/app/modules/demo/demodto"
	"github.com/hulklab/yago/example/app/modules/demo/demoservice"
)

type UserHttp struct {
	ghttp.BaseHttp
}

func init() {
	h := new(UserHttp)

	ghttp.Root.Post("/demo/user/add", h.AddAction)
	ghttp.Root.Post("/demo/user/delete", h.DeleteAction)
	ghttp.Root.Post("/demo/user/update", h.UpdateAction)
	ghttp.Root.Post("/demo/user/list", h.ListAction)
	ghttp.Root.Get("/demo/user/detail", h.DetailAction)

	// routing groups are recommended
	memberGroup := yago.NewHttpGroupRouter("/demo/user/member", ghttp.CheckUserName)
	{
		memberGroup.Post("/:name", h.UserSetAction)
		memberGroup.Get("/:name", h.UserGetAction)
		memberGroup.Put("/:name", h.UserUpdateAction)
		memberGroup.Delete("/:name", h.UserDeleteAction)

		consumeSubGroup := memberGroup.Group("/plus")
		consumeSubGroup.Patch("/number/:number", ghttp.Compute, h.PlusAction)
	}

}

// curl 'http://127.0.0.1:8080/user/list' -H "Content-type:application/json" -XPOST -d '{"page_size":10,"page_num":1}'
func (h *UserHttp) ListAction(c *yago.Ctx) {
	req := &demodto.UserListReq{}

	ctx := h.GetTraceCtx(c)

	if err := c.ShouldBind(req); err != nil {
		c.SetError(err)
		return
	}

	data, err := demoservice.NewUserService(ctx).GetList(req)
	c.SetDataOrErr(data, err)
}

// curl 'http://127.0.0.1:8080/demo/user/add' -H "Content-type:application/x-www-form-urlencoded" -XPOST -d "username=lisi&phone=13090001112"
func (h *UserHttp) AddAction(c *yago.Ctx) {
	req := &demodto.UserAddReq{}

	ctx := h.GetTraceCtx(c)

	if err := c.ShouldBind(req); err != nil {
		c.SetError(err)
		return
	}

	data, err := demoservice.NewUserService(ctx).AddUser(req)
	c.SetDataOrErr(data, err)
}

// curl 'http://127.0.0.1:8080/demo/user/update' -H "Content-type:application/json" -XPOST -d '{"id":2,"username":"zhangsan"}'
func (h *UserHttp) UpdateAction(c *yago.Ctx) {
	req := &demodto.UserUpdateReq{}

	ctx := h.GetTraceCtx(c)

	if err := c.ShouldBind(req); err != nil {
		c.SetError(err)
		return
	}

	err := demoservice.NewUserService(ctx).UpdateById(req)
	c.SetDataOrErr(g.Hash{}, err)
}

// curl 'http://127.0.0.1:8080/demo/user/delete' -H "Content-type:application/json" -XPOST -d '{"id":1}'
func (h *UserHttp) DeleteAction(c *yago.Ctx) {
	req := &demodto.UserDeleteReq{}

	ctx := h.GetTraceCtx(c)

	if err := c.ShouldBind(req); err != nil {
		c.SetError(err)
		return
	}

	err := demoservice.NewUserService(ctx).DeleteById(req)
	c.SetDataOrErr(g.Hash{}, err)
}

// curl 'http://127.0.0.1:8080/demo/user/detail?id=2' -H "Content-type:application/json" -XGET
func (h *UserHttp) DetailAction(c *yago.Ctx) {
	req := &demodto.UserDetailReq{}

	ctx := h.GetTraceCtx(c)

	if err := c.ShouldBind(req); err != nil {
		c.SetError(err)
		return
	}

	data, err := demoservice.NewUserService(ctx).GetDetail(req)
	c.SetDataOrErr(data, err)
}

func (h *UserHttp) UserSetAction(c *yago.Ctx) {
	name := c.Param("name")

	c.SetData("set " + name)
}

func (h *UserHttp) UserGetAction(c *yago.Ctx) {
	name := c.Param("name")

	c.SetData("get " + name)
}

func (h *UserHttp) UserUpdateAction(c *yago.Ctx) {
	name := c.Param("name")

	c.SetData("update " + name)
}

func (h *UserHttp) UserDeleteAction(c *yago.Ctx) {
	name := c.Param("name")

	c.SetData("delete " + name)
}

func (h *UserHttp) PlusAction(c *yago.Ctx) {
	plusNumber := c.Param("number")
	number := c.GetInt("number")
	number = number + cast.ToInt(plusNumber)
	c.Set("number", number)
}
