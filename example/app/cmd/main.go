package main

import (
	"github.com/hulklab/yago"

	_ "github.com/hulklab/yago/example/app/app/routes/cmdroute"
)

func main() {
	yago.NewCmd().RunCmd()
}
