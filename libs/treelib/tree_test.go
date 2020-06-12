package treelib

import (
	"encoding/json"
	"fmt"
	"testing"
)

type ExampleNode struct {
	Id       int64   `json:"id"`
	Name     string  `json:"name"`
	ParentId int64   `json:"parent_id"`
	Children []INode `json:"children"`
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

func (n *ExampleNode) SetChildren(list Nodes) {
	n.Children = append(n.Children, list...)
}

func TestGenTree(t *testing.T) {
	var list Nodes
	list = Nodes{
		&ExampleNode{
			Id:       1,
			Name:     "一级节点",
			ParentId: 0,
		},
		&ExampleNode{
			Id:       2,
			Name:     "一级节点2",
			ParentId: 0,
		},
		&ExampleNode{
			Id:       3,
			Name:     "二级节点",
			ParentId: 1,
		},
		&ExampleNode{
			Id:       4,
			Name:     "二级节点2",
			ParentId: 2,
		},
		&ExampleNode{
			Id:       5,
			Name:     "三级节点",
			ParentId: 3,
		},
		&ExampleNode{
			Id:       6,
			Name:     "四级节点",
			ParentId: 5,
		},
		&ExampleNode{
			Id:       8,
			Name:     "三级节点a",
			ParentId: 4,
		},
		&ExampleNode{
			Id:       7,
			Name:     "三级节点b",
			ParentId: 4,
		},
	}

	root := &ExampleNode{
		Id:       0,
		Name:     "根节点",
		ParentId: -1,
	}

	err := GenTree(root, list)
	if err != nil {
		t.Error(err)
		return
	}

	jsonTree, _ := json.MarshalIndent(root, "", "  ")

	fmt.Println(string(jsonTree))
}
