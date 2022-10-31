package homeapi

import (
	"fmt"
	"log"

	"github.com/hulklab/yago"
	"github.com/hulklab/yago/base/basethird"
	"github.com/levigross/grequests"
)

type homeApi struct {
	basethird.HttpThird
}

// Usage: Ins().GetUserById()
func Ins() *homeApi {
	name := "home_api"
	v := yago.Component.Ins(name, func() interface{} {
		api := new(homeApi)

		// http 配置
		err := api.InitConfig(name)
		if err != nil {
			log.Fatal("init home api config error")
		}

		// 添加中间件
		api.AddInterceptor(func(method, uri string, ro *grequests.RequestOptions, call basethird.Caller) (response *basethird.Response, e error) {
			fmt.Println("before caller....", uri, method)
			// 注意：在 call 之前 return 的话，将不会真正执行接口调用，也没有日志

			resp, err := call(method, uri, ro)

			fmt.Println("after caller....", resp.StatusCode)

			return resp, err
		})

		return api
	})
	return v.(*homeApi)
}

func (a *homeApi) Hello(name string) {
	params := map[string]interface{}{}

	if name != "" {
		params["name"] = name
	}

	req, err := a.Get("/hello", params)

	fmt.Println("req:", req, "err:", err)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// } else {
	// 	s, _ := req.String()
	// 	fmt.Println(s)
	// }
}

func (a *homeApi) GetUserById(id int64) string {
	params := map[string]interface{}{
		"id": id,
	}

	req, err := a.Get("/home/user/detail", params)
	if err != nil {
		return err.Error()
	} else {
		s, _ := req.String()
		return s
	}
}

func (a *homeApi) UploadFile(filepath string) string {
	params := map[string]interface{}{
		"file": basethird.PostFile(filepath),
	}

	req, err := a.Post("/upload", params)
	if err != nil {
		return err.Error()
	} else {
		s, _ := req.String()
		return s
	}
}
