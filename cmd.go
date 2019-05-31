package yago

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type Cmd struct {
	*cobra.Command
}

func NewCmd() *Cmd {
	cmd := &Cmd{&cobra.Command{
		Use:   Config.GetString("app.app_name") + "-cmd",
		Short: "cmd 命令执行工具",
	}}
	return cmd
}

func (c *Cmd) LoadCmdRouter() {
	if len(CmdRouterMap) == 0 {
		return
	}

	baseCmdMap := make(map[string]*cobra.Command)

	for use, router := range CmdRouterMap {
		useSlice := strings.Split(use, "/")
		length := len(useSlice)

		if length == 0 {
			log.Fatalf("add cmd router failed: %s", use)
		}

		var baseCmdSlice []string
		var cmd string
		if length == 1 {
			baseCmdSlice = useSlice
			cmd = ""
		}
		if length > 2 {
			baseCmdSlice = useSlice[:(length - 1)]
			cmd = useSlice[length-1]
		}

		baseCmdStr := strings.Join(baseCmdSlice, "/")

		var baseCmd *cobra.Command

		if _, ok := baseCmdMap[baseCmdStr]; !ok {
			baseCmd = &cobra.Command{
				Use: baseCmdStr,
			}
			baseCmdMap[baseCmdStr] = baseCmd
			c.AddCommand(baseCmd)
		} else {
			baseCmd = baseCmdMap[baseCmdStr]
		}

		var rootCmd *cobra.Command

		if cmd == "" {
			baseCmd.Short = router.Short
			baseCmd.Run = router.Action
			rootCmd = baseCmd
		} else {
			subCmd := &cobra.Command{
				Use:   cmd,
				Short: router.Short,
				Run:   router.Action,
			}
			baseCmd.AddCommand(subCmd)
			rootCmd = subCmd
		}
		if len(router.Args) > 0 {
			for _, arg := range router.Args {
				rootCmd.Flags().StringP(arg.Name, arg.Shorthand, arg.Value, arg.Usage)
				if arg.Required {
					if err := rootCmd.MarkFlagRequired(arg.Name); err != nil {
						log.Printf("cmd %s mark flag failed: %s", use, err.Error())
					}
				}
			}
		}
	}
}

func (c *Cmd) RunCmd() {
	c.LoadCmdRouter()

	if err := c.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
