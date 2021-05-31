// +build go1.16

package homehttp

import (
	"net/http"

	"github.com/hulklab/yago"
	"github.com/hulklab/yago/base/basehttp"
	"github.com/hulklab/yago/example/app/g"
	"github.com/hulklab/yago/libs/date"
)

type IndexHttp struct {
	basehttp.BaseHttp
}

func init() {
	h := new(IndexHttp)

	group := yago.NewHttpGroupRouter("/")
	group.Get("/", h.IndexAction)
}

func (h IndexHttp) IndexAction(c *yago.Ctx) {
	c.HTML(http.StatusOK, "index.html", g.Hash{"date": date.Now()})
	return
}
