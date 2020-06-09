package homehttp

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/hulklab/yago/example/app/modules/home/homehttp/homemiddleware"

	"github.com/hulklab/yago"
	"github.com/hulklab/yago/base/basehttp"
	"github.com/hulklab/yago/example/app/g"
	"github.com/hulklab/yago/example/app/modules/home/homemodel"
)

type HomeHttp struct {
	basehttp.BaseHttp
}

type HttpMetadata struct {
	Label string `json:"label"`
}

func init() {
	homeHttp := new(HomeHttp)

	// simple route, not recommend
	yago.AddHttpRouter("/home/hello", http.MethodGet, homeHttp.HelloAction)
	yago.AddHttpRouter("/home/add", http.MethodPost, homeHttp.AddAction)
	yago.AddHttpRouter("/home/delete", http.MethodPost, homeHttp.DeleteAction)
	yago.AddHttpRouter("/home/detail", http.MethodGet, homeHttp.DetailAction)
	yago.AddHttpRouter("/home/update", http.MethodPost, homeHttp.UpdateAction)
	yago.AddHttpRouter("/home/list", http.MethodPost, homeHttp.ListAction)
	yago.AddHttpRouter("/home/upload", http.MethodPost, homeHttp.UploadAction)
	yago.AddHttpRouter("/home/hello/:name", http.MethodGet, homeHttp.Hello2Action)
	yago.AddHttpRouter("/home/cookie", http.MethodGet, homeHttp.CookieAction)
	yago.AddHttpRouter("/home/metadata", http.MethodGet, homeHttp.MetadataAction, HttpMetadata{
		Label: "自定义HTTP名称",
	})

	// routing groups are recommended
	userGroup := yago.NewHttpGroupRouter("/home/user")
	userGroup.Use(homemiddleware.CheckUserName)
	{
		userGroup.Post("/:name", homeHttp.UserSetAction)
		userGroup.Get("/:name", homeHttp.UserGetAction)
		userGroup.Put("/:name", homeHttp.UserUpdateAction)
		userGroup.Delete("/:name", homeHttp.UserDeleteAction)

		consumeSubGroup := userGroup.Group("/consume")
		consumeSubGroup.Use(homemiddleware.ComputeConsume)
		consumeSubGroup.Patch("/sleep/:name", homeHttp.ConsumeSleepAction)
	}

	yago.SetHttpNoRouter(homeHttp.NoRouterAction)
}

func (h *HomeHttp) NoRouterAction(c *yago.Ctx) {
	c.JSON(http.StatusNotFound, g.Hash{
		"error": "404, page not exists",
	})
}

// curl -X GET 'http://127.0.0.1:8080/home/hello?name=zhangsan'
func (h *HomeHttp) HelloAction(c *yago.Ctx) {
	var p struct {
		Name string `json:"name" validate:"omitempty,max=20" form:"name" label:"姓名"`
	}

	err := c.ShouldBind(&p)
	if err != nil {
		c.SetError(err)
		return
	}

	data := "hello " + p.Name

	c.SetData(data)
}

// curl 'http://127.0.0.1:8080/home/add' -H "Content-type:application/x-www-form-urlencoded" -XPOST -d "name=lisi&phone=13090001112"
func (h *HomeHttp) AddAction(c *yago.Ctx) {
	var p struct {
		Name  string `json:"name" validate:"required,max=20" form:"name" label:"姓名"`
		Phone string `json:"phone" validate:"required,phone" form:"phone" label:"手机号"`
	}

	err := c.ShouldBind(&p)
	if err != nil {
		c.SetError(err)
		return
	}

	model := homemodel.NewHomeModel()
	id, e := model.Add(p.Name, nil)
	if e != nil {
		c.SetError(e)
		return
	}

	c.SetData(map[string]interface{}{"id": id})
}

var p struct {
	Id int64 `json:"id" validate:"required" form:"id" label:"Id"`
}

// curl 'http://127.0.0.1:8080/home/delete' -H "Content-type:application/json" -XPOST -d '{"id":1}'
func (h *HomeHttp) DeleteAction(c *yago.Ctx) {

	err := c.ShouldBind(&p)
	if err != nil {
		c.SetError(err)
		return
	}

	model := homemodel.NewHomeModel()

	n, e := model.DeleteById(p.Id)
	if e != nil {
		c.SetError(err)
		return
	}

	c.SetData(n)
}

// curl 'http://127.0.0.1:8080/home/detail?id=2' -H "Content-type:application/json" -XGET
func (h *HomeHttp) DetailAction(c *yago.Ctx) {
	err := c.ShouldBind(&p)
	if err != nil {
		c.SetError(err)
		return
	}

	model := homemodel.NewHomeModel()

	data := model.GetDetail(p.Id)

	c.SetData(data)
}

// curl 'http://127.0.0.1:8080/home/update' -H "Content-type:application/json" -XPOST -d '{"id":2,"name":"zhangsan"}'
func (h *HomeHttp) UpdateAction(c *yago.Ctx) {
	var p struct {
		Id   int64  `json:"id" validate:"required" form:"id" label:"Id"`
		Name string `json:"name" validate:"required" form:"name" label:"姓名"`
	}

	err := c.ShouldBind(&p)
	if err != nil {
		c.SetError(err)
		return
	}

	model := homemodel.NewHomeModel()

	var options = make(map[string]interface{})

	if p.Name != "" {
		options["name"] = p.Name
	}

	user, e := model.UpdateById(p.Id, options)
	if e != nil {
		c.SetError(err)
		return
	}
	c.SetData(user)
}

// curl 'http://127.0.0.1:8080/home/list' -H "Content-type:application/json" -XPOST -d '{"pagesize":1}'
func (h *HomeHttp) ListAction(c *yago.Ctx) {
	type p struct {
		Q        string `json:"q" validate:"omitempty" form:"q"`
		Page     int    `json:"page" validate:"omitempty" form:"name" label:"当前页"`
		Pagesize int    `json:"pagesize" validate:"omitempty" form:"pagesize" label:"页大小"`
	}

	pi := &p{
		Page:     1,
		Pagesize: 10,
	}

	err := c.ShouldBind(&pi)
	if err != nil {
		c.SetError(err)
		return
	}

	model := homemodel.NewHomeModel()
	total, users := model.GetList(pi.Q, pi.Page, pi.Pagesize)
	c.SetData(map[string]interface{}{
		"total": total,
		"list":  users,
	})
}

func (h *HomeHttp) UploadAction(c *yago.Ctx) {

	file, _ := c.FormFile("file")

	// Upload the file to specific dst.
	if err := c.SaveUploadedFile(file, "/Users/xxx/Downloads/upload_test.png"); err != nil {
		c.SetError(yago.NewErr(err.Error()))
		return
	}

	c.SetData(file.Filename)
}

func (h *HomeHttp) Hello2Action(c *yago.Ctx) {
	name := c.Param("name")

	c.SetData("hello " + name)
}

func (h *HomeHttp) UserSetAction(c *yago.Ctx) {
	name := c.Param("name")

	c.SetData("set " + name)
}

func (h *HomeHttp) UserGetAction(c *yago.Ctx) {
	name := c.Param("name")

	c.SetData("get " + name)
}

func (h *HomeHttp) UserUpdateAction(c *yago.Ctx) {
	name := c.Param("name")

	c.SetData("update " + name)
}

func (h *HomeHttp) UserDeleteAction(c *yago.Ctx) {
	name := c.Param("name")

	c.SetData("delete " + name)
}

func (h *HomeHttp) CookieAction(c *yago.Ctx) {
	cookie, err := c.Cookie("user")

	if err != nil {
		c.SetError(err)
		return
	}

	c.SetData("hello " + cookie)
}

func (h *HomeHttp) ConsumeSleepAction(c *yago.Ctx) {
	c.SetData("I'm sleeping zzz.....")
	time.Sleep(time.Second * time.Duration(rand.Intn(5)))
}

func (h *HomeHttp) MetadataAction(c *yago.Ctx) {
	data := "get label from metadata:"

	for _, router := range yago.GetHttpRouters() {
		if router.Url() == "/home/metadata" {
			v, ok := router.Metadata.([]interface{})
			if ok {
				data = data + v[0].(HttpMetadata).Label
			}
			break
		}
	}

	c.SetData(data)
}
