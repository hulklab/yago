package homemiddleware

import (
	"time"

	"github.com/hulklab/yago"
)

func CheckUserName(c *yago.Ctx) {
	name := c.Param("name")
	if name == "devil" {
		c.SetError(yago.NewErr("path param name can not be devil"))
		c.Abort()
	}
}

func ComputeConsume(c *yago.Ctx) {
	t := time.Now()

	// before request

	c.Next()

	// after request
	latency := time.Since(t)

	c.SetData("I'm awake and I slept for " + latency.String())
}
