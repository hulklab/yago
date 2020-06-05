package elastic

import (
	"github.com/hulklab/yago"
	"github.com/hulklab/yago/coms/logger"
	"github.com/olivere/elastic"
	"log"
)

type ElasticV6 struct {
	*elastic.Client
}

func InsV6(id ...string) *ElasticV6 {

	var name string

	if len(id) == 0 {
		name = "elastic"
	} else if len(id) > 0 {
		name = id[0]
	}

	v := yago.Component.Ins(name, func() interface{} {
		config := yago.Config.GetStringMap(name)
		options := make([]elastic.ClientOptionFunc, 0)

		urls := yago.Config.GetStringSlice(name + ".urls")
		if len(urls) < 1 {
			log.Fatalln("Fatal error elastic urls not config")
		}

		options = append(options, elastic.SetURL(urls...))

		username := config["username"].(string)
		password := config["password"].(string)

		if username != "" && password != "" {
			options = append(options, elastic.SetBasicAuth(username, password))
		}

		if infolog_enable, ok := config["infolog_enable"]; ok && infolog_enable.(bool) {
			infolog := logger.Ins().Category("[ELASTIC_INFO " + name + "]")
			options = append(options, elastic.SetInfoLog(infolog))
		}

		if tracelog_enable, ok := config["tracelog_enable"]; ok && tracelog_enable.(bool) {
			tracelog := logger.Ins().Category("[ELASTIC_TRACE " + name + "]")
			options = append(options, elastic.SetTraceLog(tracelog))
		}

		errlog := logger.Ins().Category("[ELASTIC_ERROR " + name + "]")
		options = append(options, elastic.SetErrorLog(errlog))

		if sniff_enable, ok := config["sniff_enable"]; ok && !sniff_enable.(bool) {
			options = append(options, elastic.SetSniff(false))
		}

		client, err := elastic.NewClient(options...)
		if err != nil {
			//如果报错 no Elasticsearch node available  可能是用户名密码不正确，或者 sniff_enable没有置为false
			log.Fatalf("Fatal error url: %s", err)
		}
		return &ElasticV6{
			Client: client,
		}
	})

	return v.(*ElasticV6)
}
