package main // import "github.com/hulklab/yago/example"

import (
	"github.com/hulklab/yago"

	_ "github.com/hulklab/yago/example/app/route"
)

func main() {
	yago.NewCmd().RunCmd()
}
