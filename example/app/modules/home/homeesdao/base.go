package homeesdao

import "github.com/hulklab/yago/coms/elastic"

type BaseEsDao struct {
}

func (b *BaseEsDao) GetEs() *elastic.ElasticV6 {
	return elastic.InsV6()
}
