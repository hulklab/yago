package homecmd

import (
	"fmt"
	"time"

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
		Name: "age", Shorthand: "e", Value: 1, Usage: "年龄", Required: false,
	}, yago.CmdDurationArg{
		Name: "time", Shorthand: "t", Value: time.Second, Usage: "时间", Required: false,
	}, yago.CmdBoolArg{
		Name: "force", Shorthand: "f", Value: true, Usage: "强制", Required: false,
	})
}

func (c *HomeCmd) DemoAction(cmd *cobra.Command, args []string) {

	if arg, err := cmd.Flags().GetString("arg"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("arg: " + arg)
	}

	age, err := cmd.Flags().GetInt("age")
	fmt.Println(age, err)

	t, err := cmd.Flags().GetDuration("time")
	fmt.Println(t, err)

	f, err := cmd.Flags().GetBool("force")
	fmt.Println(f, err)
}
