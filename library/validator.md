# 验证器

验证器主要是为了在控制器层进行参数校验的，目前主要在 Http 控制器中进行了集成。我们实现了几个常用的验证器功能：
* [Required](#Required)：必传且不为空
* [String](#String)：字符串验证器
* [Number](#Number)：数字验证器
* [Float](#Float)：浮点数验证器
* [JSON](#JSON)：JSON 格式验证器
* [IP](#IP)：IP 格式验证器
* [Match](#Match)：正则验证器
* [Custom](#Custom)：自定义验证器

在 Http 控制器中，我们覆写 Rules 函数，在 Rules 函数中我们定义好验证器内容，然后返回一个验证器规则数组，yago 会在 Http 控制器的 BeforeAction 之后开始进行参数校验。

```go
func (h *HomeHttp) Rules() []validator.Rule {
	return []validator.Rule{
		{
			Params: []string{"name"},
			Method: validator.Required,
			On:     []string{"add"},
        },
    }
}
```

此外还有一个 `Label` 函数，用来对验证错误信息的参数名称进行替换和映射的。比如如果验证器错误信息为：name 参数不能为空，那么如果配置了 Label： `name=姓名`，则报错信息显示为 `姓名不能为空`。

```go
func (h *HomeHttp) Labels() validator.Label {
	return map[string]string{
		"id":       "ID",
		"name":     "姓名",
		"page":     "页码",
		"pagesize": "页内数量",
	}
}
```

## 规则 Rule 的结构

```go
type Rule struct {
	Params   []string // 参数名称，包括 URI 参数和表单参数
	Method   interface{} // 指定需要用到的验证器方法
	On       []string // 指定哪些 Action 会用到该条规则 Rule。选用当前 Action 在路由地址的末端名称。如果 on 为空则会对该控制器中的所有
	Min      float64 // 最小值验证
	Max      float64 // 最大值验证
    Pattern  string // 正则验证
	Message  string // 自定义验证错误信息，为空则展示默认错误信息
}
```

## 样例
### Required

验证该控制器中的 add，update Action 中的参数 name 和 score 是否必传且不为空

```go
validator.Rule{
    Params: []string{"name", "score"},
    Method: validator.Required,
    On:     []string{"add", "update"},
},

```

### String

验证该控制器中的 add Action 中的参数 name 是否为字符串，并且最短2个字符，最长10个字符。

```go
validator.Rule{
    Params: []string{"name"},
    Method: validator.String,
    On:     []string{"add"},
    Min:    2,
    Max:    10,
},

```

### Number

验证该控制器中的 add Action 中的参数 age 是否为数字，并且最大为200。

```go
validator.Rule{
    Params: []string{"age"},
    Method: validator.Int,
    On:     []string{"add"},
    Max:    200,
},

```

### Float

验证该控制器中的所有 Action 中的参数 score 是否为浮点数，并且最大为100。

```go
validator.Rule{
    Params: []string{"score"},
    Method: validator.Float,
    Max:    100,
},

```

### JSON

验证该控制器中的所有 Action 中的参数 extend 是否为json格式。

```go
validator.Rule{
    Params: []string{"extend"},
    Method: validator.JSON,
},

```

### IP

验证该控制器中的所有 Action 中的参数 ip 是否为ip格式。

```go
validator.Rule{
    Params: []string{"ip"},
    Method: validator.IP,
},

```

### Match

验证该控制器中的所有 Action 中的参数 email 格式是否正确。

```go
validator.Rule{
    Params: []string{"email"},
    Method: validator.Match,
    Pattern: `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`,
},

```

### Custom

自定义验证器支持我们自己定义验证函数

```go
func (h *HomeHttp) CheckNameExist(c *yago.Ctx, p string) (bool, error) {
	val, _ := c.Get(p)
	// check param p is exist
	var exists bool

	if val == "zhangsan" {
		exists = true
	}

	if exists {
		return false, fmt.Errorf("name %s is exists", val)
	}
	return true, nil

}
```

验证该控制器中的所有 Action 中的参数 name 是否为 zhangsan。

```go
validator.Rule{
    Params: []string{"name"},
    Method: h.CheckNameExist,
},

```