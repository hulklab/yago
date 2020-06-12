package treelib

import (
	"sort"
)

type INode interface {
	// node id
	GetId() int64
	// node pid
	GetParentId() int64
	// node seq
	GetSeq() int64
	// set node children
	SetChildren(n Nodes)
}

type Nodes []INode

func (nodes Nodes) Len() int {
	return len(nodes)
}
func (nodes Nodes) Swap(i, j int) {
	nodes[i], nodes[j] = nodes[j], nodes[i]
}
func (nodes Nodes) Less(i, j int) bool {
	return nodes[i].GetSeq() < nodes[j].GetSeq()
}

func GenTree(root INode, list Nodes) {
	var children Nodes
	for _, v := range list {
		if v.GetParentId() == root.GetId() {
			GenTree(v, list)
			children = append(children, v)
		}
	}

	sort.Sort(children)

	root.SetChildren(children)
}
