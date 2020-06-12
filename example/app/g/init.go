package g

import (
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/base/basemiddleware"
)

// load global lib, middleware or data
func init() {
	// load global middleware
	yago.HttpUse(basemiddleware.BizLog)
}
