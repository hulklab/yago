package homehttp

import (
	"fmt"
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/base/basehttp"
	"github.com/hulklab/yago/libs/validator"
	"net/http"

	"github.com/hulklab/yago/example/app/modules/home/homemodel"
)

type HomeHttp struct {
	basehttp.BaseHttp
}

func init() {
	homeHttp := new(HomeHttp)
	yago.AddHttpRouter("/home/hello", http.MethodGet, homeHttp.HelloAction, homeHttp)
	yago.AddHttpRouter("/home/add", http.MethodPost, homeHttp.AddAction, homeHttp)
	yago.AddHttpRouter("/home/delete", http.MethodPost, homeHttp.DeleteAction, homeHttp)
	yago.AddHttpRouter("/home/detail", http.MethodGet, homeHttp.DetailAction, homeHttp)
	yago.AddHttpRouter("/home/update", http.MethodPost, homeHttp.UpdateAction, homeHttp)
	yago.AddHttpRouter("/home/list", http.MethodPost, homeHttp.ListAction, homeHttp)
	yago.AddHttpRouter("/home/upload", http.MethodPost, homeHttp.UploadAction, homeHttp)
}

func (h *HomeHttp) Labels() validator.Label {
	return map[string]string{
		"id":       "ID",
		"name":     "姓名",
		"page":     "页码",
		"pagesize": "页内数量",
	}
}

func (h *HomeHttp) CheckNameExist(c *yago.Ctx, p string) (bool, error) {
	fmt.Println("here")
	val, _ := c.Get(p)
	// check param p is exist
	var exists bool

	if val == "zhangsan" {
		exists = true
	}

	if exists {
		return false, fmt.Errorf("name %s is exists", val)
	}
	return true, nil

}

func (h *HomeHttp) Rules() []validator.Rule {
	return []validator.Rule{
		{
			Params: []string{"name"},
			Method: validator.Required,
			On:     []string{"add"},
		},
		{
			Params: []string{"name"},
			Method: h.CheckNameExist,
			On:     []string{"add"},
		},
		//{
		//	Params: []string{"id"},
		//	Method: validator.Required,
		//	On:     []string{"delete", "detail", "update"},
		//},
		//{
		//	Params: []string{"page"},
		//	Method: validator.Int,
		//	Min:    1,
		//	On:     []string{"list"},
		//},
		//{
		//	Params: []string{"pagesize"},
		//	Method: validator.Int,
		//	Max:    100,
		//	On:     []string{"list"},
		//},
	}
}

func (h *HomeHttp) HelloAction(c *yago.Ctx) {
	name := c.RequestString("name")

	c.SetData("hello " + name)

	return
}

func (h *HomeHttp) AddAction(c *yago.Ctx) {
	name := c.RequestString("name")

	model := homemodel.NewHomeModel()
	id, err := model.Add(name, nil)
	if err.HasErr() {
		c.SetError(err)
		return
	}

	c.SetData(map[string]interface{}{"id": id})
	return
}

func (h *HomeHttp) DeleteAction(c *yago.Ctx) {

	uid, _ := c.RequestInt64("id")

	model := homemodel.NewHomeModel()

	n, err := model.DeleteById(uid)
	if err.HasErr() {
		c.SetError(err)
		return
	}

	c.SetData(n)
}

func (h *HomeHttp) DetailAction(c *yago.Ctx) {

	uid, _ := c.RequestInt64("id")

	model := homemodel.NewHomeModel()

	data := model.GetDetail(uid)

	c.SetData(data)
}

func (h *HomeHttp) UpdateAction(c *yago.Ctx) {
	uid, _ := c.RequestInt64("id")

	name := c.RequestString("name")

	model := homemodel.NewHomeModel()

	var options = make(map[string]interface{})

	if name != "" {
		options["name"] = name
	}

	user, err := model.UpdateById(uid, options)
	if err.HasErr() {
		c.SetError(err)
		return
	}
	c.SetData(user)
}

func (h *HomeHttp) ListAction(c *yago.Ctx) {

	q := c.RequestString("q")
	page, _ := c.RequestInt("page", 1)
	pageSize, _ := c.RequestInt("pagesize", 10)

	model := homemodel.NewHomeModel()
	total, users := model.GetList(q, page, pageSize)
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
