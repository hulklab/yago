package treelib

// ConvertToINodeArray 其他的结构体想要生成树，直接实现这个接口
type INode interface {
	// GetId获取id
	GetId() int64
	// GetParentId 获取父id
	GetParentId() int64
	// IsRoot 判断当前节点是否是顶层根节点
	IsRoot() bool
	// SetChildren 设置子节点
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
	return nodes[i].GetId() < nodes[j].GetId()
}
func GenerateTree(nodes []INode) (trees INodes) {
	trees = INodes{}
	// 定义顶层根和子节点
	var roots, childs []INode
	for _, v := range nodes {
		if v.IsRoot() {
			// 判断顶层根节点
			roots = append(roots, v)
		}
		childs = append(childs, v)
	}

	for _, v := range roots {
		// 递归
		RecursiveTree(v, childs)
		trees = append(trees, v)
	}
	return
}
func RecursiveTree(root INode, nodes INodes) {
	var children INodes
	for _, v := range nodes {
		if v.IsRoot() {
			// 如果当前节点是顶层根节点就跳过
			continue
		}
		if root.GetId() == v.GetParentId() {
			RecursiveTree(v, nodes)
			children = append(children, v)
		}
	}
	root.SetChildren(children)
}
