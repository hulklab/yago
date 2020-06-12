package treelib

import (
	"errors"
	"reflect"
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

func ptrVerify(p interface{}) error {
	value := reflect.ValueOf(p)
	if value.Kind() != reflect.Ptr {
		return errors.New("needs a pointer to a node")
	} else if value.Elem().Kind() == reflect.Ptr {
		return errors.New("a pointer to a pointer is not allowed")
	}
	return nil
}

func GenTree(root INode, list Nodes) error {
	if err := ptrVerify(root); err != nil {
		return err
	}

	var children Nodes
	for _, v := range list {
		if err := ptrVerify(v); err != nil {
			return err
		}
		if v.GetParentId() == root.GetId() {
			if err := GenTree(v, list); err != nil {
				return err
			}
			children = append(children, v)
		}
	}

	sort.Sort(children)

	root.SetChildren(children)

	return nil
}
