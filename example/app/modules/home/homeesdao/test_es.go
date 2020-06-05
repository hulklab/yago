package homeesdao

type testEsDao struct {
	BaseEsDao
}

func NewDefaultEsDao() *testEsDao {
	return new(testEsDao)
}

func (d *testEsDao) IndexName() string {
	return "default_v1"
}

func (d *testEsDao) AliasName() string {
	return "default"
}

func (d *testEsDao) Mapping() string {
	return `
{
    "settings":{
        "number_of_shards":6,
        "number_of_replicas":1,
        "index":{
            "max_result_window":"100000000"
        },
        "analysis":{
            "normalizer":{
                "caseSensitive":{
                    "type":"custom",
                    "char_filter":[

                    ],
                    "filter":[
                        "lowercase",
                        "asciifolding"
                    ]
                }
            },
            "analyzer":{
                "comma":{
                    "type":"pattern",
                    "pattern":","
                }
            }
        }
    },
    "mappings":{
        "properties":{
            "id":{
                "type":"keyword"
            },
            "title":{
                "type":"text",
                "analyzer":"ik_max_word",
                "search_analyzer":"ik_smart",
                "fields":{
                    "suggest":{
                        "type":"completion",
                        "analyzer":"ik_max_word"
                    }
                }
            },
            "category":{
                "type":"keyword"
            },
            "content":{
                "type":"text",
                "analyzer":"ik_max_word",
                "search_analyzer":"ik_smart"
            },
            "author":{
                "type":"keyword"
            },
            "ctime":{
                "type":"date",
                "format":"yyyy-MM-dd HH:mm:ss"
            },
            "utime":{
                "type":"date",
                "format":"yyyy-MM-dd HH:mm:ss"
            },
            "attachment":{
                "type":"text",
                "analyzer":"ik_max_word",
                "search_analyzer":"ik_smart"
            },
            "origin_url":{
                "type":"keyword"
            },
			"keyword":{
				"type": "text",
        		"analyzer": "comma",
        		"search_analyzer": "comma"
			}
        }
    }
}`
}
