package homemiddleware

import (
	"fmt"

	"github.com/hulklab/yago"
)

func CheckUserName(c *yago.Ctx) {
	name := c.Param("name")
	if name == "devil" {
		c.SetError(yago.NewErr("path param name can not be devil"))
		c.Abort()
	}
}

func Compute(c *yago.Ctx) {

	// before request
	c.Set("number", 1)

	c.Next()

	// after request
	number := c.GetInt("number")

	c.SetData(fmt.Sprintf("the number is %d", number))
}
