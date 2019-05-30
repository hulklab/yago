package example // import "github.com/hulklab/yago/example/app"

import (
	"github.com/hulklab/yago"

	_ "github.com/hulklab/yago/example/app/app/routes/httproute"
	_ "github.com/hulklab/yago/example/app/app/routes/rpcroute"
	_ "github.com/hulklab/yago/example/app/app/routes/taskroute"
)

func main() {
	yago.NewApp().Run()
}
