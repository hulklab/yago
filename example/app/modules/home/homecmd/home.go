package homecmd

import (
	"fmt"
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/base/basecmd"
	"github.com/hulklab/yago/base/basethird"
	"github.com/hulklab/yago/example/app/g"
	"github.com/spf13/cobra"
)

type HomeCmd struct {
	basecmd.BaseCmd
}

func init() {
	homeCmd := new(HomeCmd)
	// 注册路由
	yago.AddCmdRouter("demo", "Demo action", homeCmd.DemoAction, yago.CmdArg{
		Name: "arg", Shorthand: "a", Value: "value", Usage: "参数", Required: false,
	})
}

func (c *HomeCmd) DemoAction(cmd *cobra.Command, args []string) {
	//client := &http.Client{
	//	Transport: &http.Transport{
	//		MaxIdleConnsPerHost: 100,
	//	},
	//	Timeout: time.Duration(5) * time.Second,
	//}
	//ro := &grequests.RequestOptions{
	//	HTTPClient: client,
	//}

	t := &basethird.HttpThird{
		Domain: "http://notexistdomain.org",
	}

	res, err := t.Get("/get", nil)
	fmt.Println(res, err)
	return

	for i := 1; i <= 200; i++ {
		resp, err := t.Get("/get", g.Hash{"status": i})
		//resp, err := grequests.Get("http://httpbin.org/get", nil)
		//grequests.DoRegularRequest()
		//resp, err := client.Get("http://httpbin.org/get")
		//resp, err := http.Get("http://httpbin.org/get")
		//resp, err := httplib.Get("http://httpbin.org/get").Response()
		fmt.Println(resp, err)
	}

	//if arg, err := cmd.Flags().GetString("arg"); err != nil {
	//	fmt.Println(err)
	//} else {
	//	fmt.Println("arg: " + arg)
	//}
}
