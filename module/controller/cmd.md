# Cmd 控制器

home 模块提供了一个 cmd 的 hello world 程序，cmd 中我们使用了 [cobra](https://github.com/spf13/cobra) 命令行类库。所以 [cobra](https://github.com/spf13/cobra) 对象的使用可以直接参考官方文档，我们在这里主要帮你完成了路由注册和参数注册。

## 路由注册

```go
func init() {
	homeCmd := new(HomeCmd)
	// 注册路由
	yago.AddCmdRouter("demo", "Demo action", homeCmd.DemoAction, yago.CmdArg{
		Name: "arg", Shorthand: "a", Value: "value", Usage: "参数", Required: true,
	})
}
```

init 函数内完成了 Action 级别的路由注册。cmd 路由注册主要通过 yago.AddCmdRouter 这个函数来完成。参数含义如下表格

AddCmdRouter 参数说明

| 参数位置 | 参数类型 | 说明 |
| ------- | ------- | ------- |
| 1 | String | 命令行路径 |
| 2 | String | 命令简要说明 |
| 3 | Func | 命令行对应的 Action Func |
| >=4 | yago.AddCmdRouter | 命令行flag参数对象<br>可以配置参数的名称，缩写，默认值，参数含义，是否必须 |


## CmdAction

Cmd Action 接收两个参数一个 cmd *cobra.Command 用来获取 flag 参数，另外一个 args 用来获取非 flag 参数。

```go
func (c *HomeCmd) DemoAction(cmd *cobra.Command, args []string) {
	if arg, err := cmd.Flags().GetString("arg"); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("arg: " + arg)
	}
}

```