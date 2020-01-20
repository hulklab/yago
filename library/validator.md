# 验证器

验证器主要是为了在控制器层进行参数校验的，目前我们使用的是 gin 的验证器。gin 内部使用的包是开源的 `https://github.com/go-playground/validator`。



## 样例
比如有一个添加用户的接口，需要验证用户名必填，并且长度最大 20 个字符。
>此样例仅展示了 required 和 max 验证器，更多的验证器请参考 [validator 官方文档](https://godoc.org/github.com/go-playground/validator)


### 先定义一个结构体 p，用来接收请求的参数
```go
type par struct {
    Name string `json:"name" validate:"required,max=20" form:"name" label:"姓名"`
}

```

#### tag 说明
* validate 标签：yago 保留了原生的 validate 标签，替换掉了 gin 默认的 binding 标签
* json 标签：如果 Content-Type 是 application/json，gin 根据 json 标签的值来取值赋值
* form 标签：如果 Content-Type 不是 application/json，gin 根据 form 标签来取值赋值
* label 标签：是用来替换报错时的字段名信息的，上例中如果不指定默认是 Name

### 使用 ShouldBind 赋值并验证
通过 ctx 的 ShouldBind 方法，将请求参数赋值给变量 p，ShouldBind 方法内部会调用 validator 包做验证
调用 ctx 的 SetError 方法，yago 会自动处理验证的错误信息，并做相应的翻译。
```go
func (h *HomeHttp) AddAction(c *yago.Ctx) {
	p := par{}
	err := c.ShouldBind(&p)
    if err != nil {
        c.SetError(err)
        return
    }
	
	// your code here
	c.SetData(g.Hash{"name":p.Name})
	return
}

```

## 特别说明
### 默认值
目前未发现 validator 包设置默认值的方法，我们可以在初始化变量时给出，如下例：

```go
func (h *HomeHttp) ListAction(c *yago.Ctx) {
	type p struct {
		Q        string `json:"q" validate:"omitempty" form:"q"`
		Page     int    `json:"page" validate:"omitempty" form:"name" label:"当前页"`
		Pagesize int    `json:"pagesize" validate:"omitempty" form:"name" label:"页大小"`
	}
	
	// 设置默认值
	pi := &p{
		Page:     1,
		Pagesize: 10,
	}

	err := c.ShouldBind(&pi)
	if err != nil {
		c.SetError(err)
		return
	}
	
	// your code here
	c.SetData(g.Hash{})
	return
}
```

### 零值问题
go 的 Struct 字段是有零值的，比如上文中 Name 字段不传，它的值就是空串，我们不能根据他是否为空串
来判断用户是否有传这个字段，那么怎么处理咧，下文给出样例

```go

func (h *HomeHttp) AddAction(c *yago.Ctx) {
	// 采用字符串指针代替字符串类型
    type par struct {
        Name *string `json:"name" validate:"omitempty" form:"name" label:"姓名"`
    }
	p := par{}
	err := c.ShouldBind(&p)
    if err != nil {
        c.SetError(err)
        return
    }
	
	// 此处可以根据 par.Name 是否为 nil 来判断用户是否有传该值
	if p.Name == nil{
		// 说明用户没有传
		c.SetData("no name")
		return
	}
	
	c.SetData(*p.Name)
	return
}


```