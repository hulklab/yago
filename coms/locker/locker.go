package locker

import (
	"log"

	"github.com/hulklab/yago"
	"github.com/hulklab/yago/coms/locker/lock"
	_ "github.com/hulklab/yago/coms/locker/redis"
)

func New(id ...string) lock.ILocker {
	var name string

	if len(id) == 0 {
		name = "locker"
	} else if len(id) > 0 {
		name = id[0]
	}

	driver := yago.Config.GetString(name + ".driver")
	driverInsId := yago.Config.GetString(name + ".driver_instance_id")

	if len(driverInsId) == 0 {
		log.Fatalln("driver_id is required in locker config")
	}

	newFunc, b := lock.LoadLocker(driver)
	if !b {
		log.Fatalf("unsupport driver %s, or driver is not register yet", driver)
	}

	return newFunc(name)
}
