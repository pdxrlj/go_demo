package main

type BPItem struct {
	Key int64
	Val any
}

type BPNode struct {
	MaxKey int64     //存储子树的最大关键数
	Nodes  []*BPNode //节点的子树
	Items  []BPItem  //叶子节点的记录数据
	Next   *BPNode
}

func NewLeafNode(width int) *BPNode {
	items := make([]BPItem, 0, width+1)
	return &BPNode{
		Nodes: nil,
		Items: items,
	}
}

func NewIndexNode(width int) *BPNode {
	nodes := make([]*BPNode, 0, width+1)
	return &BPNode{
		Nodes: nodes,
		Items: nil,
	}
}
