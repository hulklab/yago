### Cmd控制器

home模块提供了一个cmd的hello world程序，cmd中我们用了[cobra](https://github.com/spf13/cobra)一个命令行类库。所以[cobra](https://github.com/spf13/cobra)对象的使用可以直接参考官方文档，我们在这里主要帮你完成了路由注册和参数注册。

homecmd/home.go
```go
package homecmd

import (
	"fmt"
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/base/basecmd"
	"github.com/spf13/cobra"
)

type HomeCmd struct {
	basecmd.BaseCmd
}

func init() {
	homeCmd := new(HomeCmd)
	// 注册路由
	yago.AddCmdRouter("demo", "Demo action", homeCmd.DemoAction, yago.CmdArg{
		Name: "arg", Shorthand: "a", Value: "value", Usage: "参数", Required: true,
	})
}

func (c *HomeCmd) DemoAction(cmd *cobra.Command, args []string) {
	if arg, err := cmd.Flags().GetString("arg"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("arg: " + arg)
	}
}

```

init函数内完成了Action级别的路由注册。cmd路由注册主要通过yago.AddCmdRouter这个函数来完成。参数含义如下表格

#### AddCmdRouter参数说明

| 参数位置 | 参数类型 | 说明 |
| ------- | ------- | ------- |
| 1 | String | 命令行路径 |
| 2 | String | 命令简要说明 |
| 3 | Func | 命令行对应的Action Func |
| >=4 | yago.AddCmdRouter | 命令行flag参数对象<br>可以配置参数的名称，缩写，默认值，参数含义，是否必须 |

Cmd Action接收两个参数一个cmd *cobra.Command用来获取flag参数，另外一个args用来获取非flag参数。