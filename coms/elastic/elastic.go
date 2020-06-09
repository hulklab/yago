package elastic

import (
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/coms/logger"
	"github.com/olivere/elastic/v7"
	"log"
)

type Elastic struct {
	*elastic.Client
}

func Ins(id ...string) *Elastic {

	var name string

	if len(id) == 0 {
		name = "elastic"
	} else if len(id) > 0 {
		name = id[0]
	}

	v := yago.Component.Ins(name, func() interface{} {
		config := yago.Config.GetStringMap(name)
		options := make([]elastic.ClientOptionFunc, 0)
		var username string
		var password string
		logLevel := int64(4)

		urls := yago.Config.GetStringSlice(name + ".urls")
		if len(urls) < 1 {
			log.Fatalln("Fatal error elastic urls not config")
		}

		options = append(options, elastic.SetURL(urls...))

		if config["username"] != nil {
			username = config["username"].(string)
		}

		if config["password"] != nil {
			password = config["password"].(string)
		}

		if username != "" && password != "" {
			options = append(options, elastic.SetBasicAuth(username, password))
		}

		if config["level"] != nil {
			logLevel = config["level"].(int64)
		}

		if logLevel >= 6 {
			tracelog := logger.Ins().Category("[ELASTIC_TRACE " + name + "]")
			options = append(options, elastic.SetTraceLog(tracelog))
		}

		if logLevel >= 4 {
			infolog := logger.Ins().Category("[ELASTIC_INFO " + name + "]")
			options = append(options, elastic.SetInfoLog(infolog))
		}

		if logLevel >= 2 {
			errlog := logger.Ins().Category("[ELASTIC_ERROR " + name + "]")
			options = append(options, elastic.SetErrorLog(errlog))
		}

		if sniffEnable, ok := config["sniff_enable"]; ok && !sniffEnable.(bool) {
			options = append(options, elastic.SetSniff(false))
		}

		client, err := elastic.NewClient(options...)
		if err != nil {
			// 如果报错 no Elasticsearch node available 可能是用户名密码不正确，或者 sniff_enable 没有置为 false。
			log.Fatalf("Fatal error url: %s", err)
		}
		return &Elastic{
			Client: client,
		}
	})

	return v.(*Elastic)
}
