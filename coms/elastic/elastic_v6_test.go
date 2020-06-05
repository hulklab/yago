package elastic

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
	"testing"
)

//参考地址  https://olivere.github.io/elastic/
func TestElastic(t *testing.T) {

	indices, err := InsV6().CatIndices().Do(context.Background())

	if err != nil {
		fmt.Printf("cat index err :%s", err)
		return
	}

	for _, index := range indices {
		fmt.Printf("index name : %s,  docs total %d, size %s \n", index.Index, index.DocsCount, index.StoreSize)
	}

	// Search with a term query
	term_query := elastic.NewTermQuery("trace_id", "a522189c642e2f0c213e2cf55a0a7139")

	ret, err := InsV6().Search().
		Index("test"). // search in index "test"
		Query(term_query). // specify the query
		From(0).Size(10). // take documents 0-9
		Do(context.Background()) // execute
	if err != nil {
		fmt.Printf("term query err :%s", err)
		return
	}

	fmt.Printf("term query took %d milliseconds\n", ret.TookInMillis)

	fmt.Printf("found total of %d \n", ret.TotalHits())

	if ret.TotalHits() > 0 {
		for _, hit := range ret.Hits.Hits {
			s, err := hit.Source.MarshalJSON()
			fmt.Println(string(s), err)
		}
	}

}
