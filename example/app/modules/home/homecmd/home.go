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
