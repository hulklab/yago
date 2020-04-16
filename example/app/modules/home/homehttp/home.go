package homehttp

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
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
	yago.AddHttpRouter("/home/hello", http.MethodGet, homeHttp.HelloAction, homeHttp)
	yago.AddHttpRouter("/home/add", http.MethodPost, homeHttp.AddAction, homeHttp)
	yago.AddHttpRouter("/home/delete", http.MethodPost, homeHttp.DeleteAction, homeHttp)
	yago.AddHttpRouter("/home/detail", http.MethodGet, homeHttp.DetailAction, homeHttp)
	yago.AddHttpRouter("/home/update", http.MethodPost, homeHttp.UpdateAction, homeHttp)
	yago.AddHttpRouter("/home/list", http.MethodPost, homeHttp.ListAction, homeHttp)
	yago.AddHttpRouter("/home/upload", http.MethodPost, homeHttp.UploadAction, homeHttp)
	yago.AddHttpRouter("/home/metadata", http.MethodGet, homeHttp.MetadataAction, homeHttp, HttpMetadata{
		Label: "自定义HTTP名称",
	})
	yago.SetHttpNoRouter(homeHttp.NoRouterAction)

	// 注册自定义验证器和翻译
	// 例：添加一个手机号验证器
	e := binding.Validator.Engine()
	if v, ok := e.(*validator.Validate); ok {
		// 添加手机号验证器
		_ = v.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
			regular := "^((13[0-9])|(14[5,7])|(15[0-3,5-9])|(17[0,3,5-8])|(18[0-9])|166|198|199|(147))\\d{8}$"

			reg := regexp.MustCompile(regular)
			return reg.MatchString(fl.Field().String())
		})

		// 添加手机号验证器的翻译
		_ = v.RegisterTranslation("phone", yago.GetTranslator(), func(ut ut.Translator) error {
			var e error
			switch ut.Locale() {
			case "zh":
				e = ut.Add("phone", "{0} 必须是一个有效的手机号", false)
			case "en":
				e = ut.Add("phone", "{0} must be a valid phone number", false)
			default:
				e = ut.Add("phone", "{0} 必须是一个有效的手机号", false)
			}
			return e
		}, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("phone", fe.Field())
			return t
		})
	}
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

// curl 'http://127.0.0.1:8080/home/add' -H "Content-type:application/x-www-form-urlencoded" -XPOST -d name=lisi&phone=13090001112
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
	if e.HasErr() {
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
	if e.HasErr() {
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
	if e.HasErr() {
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

func (h *HomeHttp) MetadataAction(c *yago.Ctx) {
	data := "get label from metadata:"

	for _, router := range yago.HttpRouterMap {
		if router.Url == "/home/metadata" {
			v, ok := router.Metadata.(HttpMetadata)
			if ok {
				data = data + v.Label
			}
			break
		}
	}

	c.SetData(data)
}
