// +build go1.16

package demohttp

import (
	"net/http"

	"github.com/hulklab/yago"
	"github.com/hulklab/yago/example/app/g"
	"github.com/hulklab/yago/example/app/g/ghttp"
	"github.com/hulklab/yago/libs/date"
)

type IndexHttp struct {
	ghttp.BaseHttp
}

func init() {
	h := new(IndexHttp)

	ghttp.Root.Get("/", h.IndexAction)
	ghttp.Root.Get("/hello", h.HelloAction)
	ghttp.Root.Get("/hello2/:name", h.Hello2Action)
	ghttp.Root.Post("/upload", h.UploadAction)
	ghttp.Root.Get("/cookie", h.CookieAction)
	ghttp.Root.Get("/metadata", h.MetadataAction).WithMetadata(ghttp.HttpMetadata{
		Label: "自定义HTTP名称",
	})



	yago.SetHttpNoRouter(h.NoRouterAction)
}

func (h *IndexHttp) IndexAction(c *yago.Ctx) {
	c.HTML(http.StatusOK, "index.html", g.Hash{"date": date.Now()})
}

// curl -X GET 'http://127.0.0.1:8080/hello?username=zhangsan'
func (h *IndexHttp) HelloAction(c *yago.Ctx) {
	var p struct {
		Username string `json:"username" validate:"omitempty,max=20" form:"username" label:"姓名"`
	}

	err := c.ShouldBind(&p)
	if err != nil {
		c.SetError(err)
		return
	}

	data := "hello " + p.Username

	c.SetData(data)
}

// curl -X GET 'http://127.0.0.1:8080/hello2?name=zhangsan'
func (h *IndexHttp) Hello2Action(c *yago.Ctx) {
	name := c.Param("name")

	c.SetData("hello " + name)
}

func (h *IndexHttp) UploadAction(c *yago.Ctx) {

	file, _ := c.FormFile("file")

	// Upload the file to specific dst.
	if err := c.SaveUploadedFile(file, "/Users/xxx/Downloads/upload_test.png"); err != nil {
		c.SetError(err)
		return
	}

	c.SetData(file.Filename)
}

func (h *IndexHttp) CookieAction(c *yago.Ctx) {
	cookie, err := c.Cookie("user")

	if err != nil {
		c.SetError(err)
		return
	}

	c.SetData("hello " + cookie)
}

func (h *IndexHttp) NoRouterAction(c *yago.Ctx) {
	c.JSON(http.StatusNotFound, g.Hash{
		"error": "404, page not exists",
	})
}

func (h *IndexHttp) MetadataAction(c *yago.Ctx) {
	data := "get label from metadata:"

	for _, router := range yago.GetHttpRouters() {
		if router.Url() == "/metadata" {
			v, ok := router.Metadata.(ghttp.HttpMetadata)
			if ok {
				data = data + v.Label
			}
			break
		}
	}

	c.SetData(data)
}
