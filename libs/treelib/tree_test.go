package treelib

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/tidwall/pretty"
)

type ExampleNode struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	ParentId int64  `json:"parent_id"`
	Children INodes `json:"children"` //子树
}

func (n ExampleNode) GetId() int64 {
	return n.Id
}

func (n ExampleNode) GetParentId() int64 {
	return n.ParentId
}

func (n ExampleNode) GetSeq() int64 {
	return n.Id
}

func (n ExampleNode) IsRoot() bool {
	// 这里通过ParentId等于0 或者 ParentId等于自身Id表示顶层根节点
	return n.ParentId == 0 || n.ParentId == n.Id
}

func (n *ExampleNode) SetChildren(nodes INodes) {
	n.Children = append(n.Children, nodes...)
}

// ConvertToINodeArray 将当前数组转换成父类 INode 接口 数组
func ConvertToINodeArray(ns []*ExampleNode) (nodes []INode) {
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
	resp := GenerateTree(ConvertToINodeArray(list))
	bytes, _ := json.MarshalIndent(resp, "", "\t")
	fmt.Println(string(pretty.Color(pretty.PrettyOptions(bytes, pretty.DefaultOptions), nil)))
}
