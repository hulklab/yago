package treelib

import (
	"encoding/json"
	"fmt"
	"github.com/tidwall/pretty"
	"testing"
)

type ExampleNode struct {
	Id       int64   `json:"id"`
	Name     string  `json:"name"`
	ParentId int64   `json:"parent_id"`
}

func (n ExampleNode) GetId() int64 {
	return n.Id
}

func (n ExampleNode) GetParentId() int64 {
	return n.ParentId
}

func (n ExampleNode) GetTitle() string {
	return n.Name
}

func (n ExampleNode) GetData() interface{} {
	return n
}

func (n ExampleNode) IsRoot() bool {
	// 这里通过ParentId等于0 或者 ParentId等于自身Id表示顶层根节点
	return n.ParentId == 0 || n.ParentId == n.Id
}

type ExampleNodes []*ExampleNode

// ConvertToINodeArray 将当前数组转换成父类 INode 接口 数组
func (ns ExampleNodes) ConvertToINodeArray() (nodes []INode) {
	for _, v := range ns {
		nodes = append(nodes, v)
	}
	return
}

func TestGenTree(t *testing.T) {
	list := []*ExampleNode{
		{
			Id:       1,
			Name:     "一级节点",
			ParentId: 0,
		},
		{
			Id:       2,
			Name:     "一级节点2",
			ParentId: 0,
		},
		{
			Id:       3,
			Name:     "二级节点",
			ParentId: 1,
		},
		{
			Id:       4,
			Name:     "二级节点2",
			ParentId: 2,
		},
		{
			Id:       5,
			Name:     "三级节点",
			ParentId: 3,
		},
		{
			Id:       6,
			Name:     "四级节点",
			ParentId: 5,
		},
		{
			Id:       8,
			Name:     "三级节点a",
			ParentId: 4,
		},
		{
			Id:       7,
			Name:     "三级节点b",
			ParentId: 4,
		},
	}

	// 生成完全树
	resp := GenerateTree(ExampleNodes.ConvertToINodeArray(list), nil)
	bytes, _ := json.MarshalIndent(resp, "", "\t")
	fmt.Println(string(pretty.Color(pretty.PrettyOptions(bytes, pretty.DefaultOptions), nil)))

	// 模拟从数据库中查询出 '三级节点'
	threeNode := []*ExampleNode{list[4]}
	// 查询 '三级节点' 的所有父节点
	respNodes := FindRelationNode(ExampleNodes.ConvertToINodeArray(threeNode), ExampleNodes.ConvertToINodeArray(list))
	resp = GenerateTree(respNodes, nil)
	bytes, _ = json.Marshal(resp)
	fmt.Println(string(pretty.Color(pretty.PrettyOptions(bytes, pretty.DefaultOptions), nil)))
}
