package homecmd

import (
	"fmt"
	"time"

	"github.com/hulklab/yago/example/app/third/homeapi"

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
	yago.AddCmdRouter("demo", "Demo action", homeCmd.DemoAction, yago.CmdStringArg{
		Name: "arg", Shorthand: "a", Value: "value", Usage: "参数", Required: true,
	}, yago.CmdIntArg{
		Name: "int", Shorthand: "i", Value: 1, Usage: "整数", Required: false,
	}, yago.CmdInt64Arg{
		Name: "int64", Shorthand: "x", Value: 64, Usage: "整数64位", Required: false,
	}, yago.CmdDurationArg{
		Name: "duration", Shorthand: "d", Value: time.Second, Usage: "时间", Required: false,
	}, yago.CmdBoolArg{
		Name: "bool", Shorthand: "b", Value: true, Usage: "布尔", Required: false,
	}, yago.CmdFloat64Arg{
		Name: "float", Shorthand: "f", Value: 0.1, Usage: "浮点", Required: false,
	}, yago.CmdStringSliceArg{
		Name: "string_slice", Shorthand: "y", Value: nil, Usage: "字串数组", Required: false,
	}, yago.CmdIntSliceArg{
		Name: "int_slice", Shorthand: "z", Value: nil, Usage: "整型数组", Required: false,
	})

	yago.AddCmdRouter("test", "test action", homeCmd.TestAction)
}

func (c *HomeCmd) TestAction(cmd *cobra.Command, args []string) {
	homeapi.New().RpcHello()
	time.Sleep(1 * time.Second)
	homeapi.New().RpcHelloStream()
}

// ./example demo -a=hello -i 11 -x 66 -d 10ms -b=false -y=java,php -z=1,2,3 -f 3.14
func (c *HomeCmd) DemoAction(cmd *cobra.Command, args []string) {

	if arg, err := cmd.Flags().GetString("arg"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("arg: " + arg)
	}

	i, err := cmd.Flags().GetInt("int")
	fmt.Println(i, err)

	x, err := cmd.Flags().GetInt64("int64")
	fmt.Println(x, err)

	d, err := cmd.Flags().GetDuration("duration")
	fmt.Println(d, err)

	b, err := cmd.Flags().GetBool("bool")
	fmt.Println(b, err)

	f, err := cmd.Flags().GetFloat64("float")
	fmt.Println(f, err)

	y, err := cmd.Flags().GetStringSlice("string_slice")
	fmt.Println(y, err)

	z, err := cmd.Flags().GetIntSlice("int_slice")
	fmt.Println(z, err)

}
