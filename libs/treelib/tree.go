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
	SetChildren(n INodes)
}

type INodes []INode

func (nodes INodes) Len() int {
	return len(nodes)
}
func (nodes INodes) Swap(i, j int) {
	nodes[i], nodes[j] = nodes[j], nodes[i]
}
func (nodes INodes) Less(i, j int) bool {
	return nodes[i].GetSeq() < nodes[j].GetSeq()
}

func GenTree(root INode, list INodes) {
	var children INodes
	for _, v := range list {
		if v.GetParentId() == root.GetId() {
			GenTree(v, list)
			children = append(children, v)
		}
	}

	sort.Sort(children)

	root.SetChildren(children)
}
