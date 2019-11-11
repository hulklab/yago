## 模型 Model

这一层比较简单，没有特别的要求，就是业务逻辑。在相应的函数内调用组件的来完成数据操作。只不过给http控制器调用的函数返回值的 error 尽量采用 yago.Err。

```go
package homemodel

import (
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/example/app/g"
	"github.com/hulklab/yago/libs/date"
	"github.com/hulklab/yago/coms/orm"

	"github.com/hulklab/yago/example/app/modules/home/homedao"
)

type HomeModel struct {
}

func NewHomeModel() *HomeModel {
	return &HomeModel{}
}

func (m *HomeModel) Add(name string, options map[string]interface{}) (int64, yago.Err) {

	// 判断 name 是否已存在
	exist := &homedao.HomeDao{Name: name}

	orm.Ins().Get(exist)

	if exist.Id != 0 {
		return 0, yago.NewErr("用户名 " + name + " 已存在")
	}

	// 添加用户
	user := &homedao.HomeDao{
		Name:  name,
		Ctime: date.Now(),
	}

	_, err := orm.Ins().Insert(user)

	return user.Id, yago.NewErr(err)

}
```