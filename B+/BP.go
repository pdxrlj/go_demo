package main

import (
	"sort"
	"sync"
)

type BPTree struct {
	mux   sync.RWMutex
	root  *BPNode
	width int //B+树的阶层
	halfw int //ceil(M/2)
}

func NewBPTree(width int) *BPTree {
	return &BPTree{
		mux:   sync.RWMutex{},
		root:  nil,
		width: width,
		halfw: (width + 1) / 2,
	}
}

// Get 查找
func (b *BPTree) Get(key int64) any {
	b.mux.RLock()
	defer b.mux.RUnlock()
	node := b.root
	for i := 0; i < len(node.Nodes); i++ {
		// 找到小于maxKey的key
		// 找到当前节点的第一个大于key的节点，然后继续往下找，并且重置i=0
		if node.Nodes[i].MaxKey >= key {
			node = node.Nodes[i]
			i = 0
		}
	}

	//没有达到叶子节点
	if len(node.Nodes) > 0 {
		return nil
	}
	for i := 0; i < len(node.Items); i++ {
		if node.Items[i].Key == key {
			return node.Items[i].Val
		}
	}
	return nil
}

func (b *BPTree) Add(key int64, val any) {
	b.mux.Lock()
	defer b.mux.Unlock()
	if b.root == nil {
		b.root = NewLeafNode(b.width)
	}
	node := b.root

	// 找到叶子节点
	for i := 0; i < len(node.Nodes); i++ {
		if node.Nodes[i].MaxKey >= key {
			node = node.Nodes[i]
			i = 0
		}
	}

	// 没有达到叶子节点
	if len(node.Nodes) > 0 {
		return
	}

	// 更新
	for i := 0; i < len(node.Items); i++ {
		if node.Items[i].Key == key {
			node.Items[i].Val = val
			return
		}
	}

	//插入
	node.Items = append(node.Items, BPItem{Key: key, Val: val})
	//排序
	sort.Slice(node.Items, func(i, j int) bool {
		return node.Items[i].Key < node.Items[j].Key
	})
	//判断是否需要分裂
	if len(node.Items) > b.width {
		b.split(node)
	}
}

func (b *BPTree) split(node *BPNode) {
	// 分裂
	var left = NewLeafNode(b.width)
	var right = NewLeafNode(b.width)
	var mid = b.halfw
	left.Items = append(left.Items, node.Items[:mid]...)
	right.Items = append(right.Items, node.Items[mid:]...)

	// 更新父节点
	if node.Parent == nil {
		node.Parent = NewIndexNode(b.width)
		b.root = node.Parent
	}
	node.Parent.Nodes = append(node.Parent.Nodes, left)
	node.Parent.Nodes = append(node.Parent.Nodes, right)
	sort.Slice(node.Parent.Nodes, func(i, j int) bool {
		return node.Parent.Nodes[i].MaxKey < node.Parent.Nodes[j].MaxKey
	})
	// 更新maxKey
	for i := 0; i < len(node.Parent.Nodes); i++ {
		node.Parent.Nodes[i].MaxKey = node.Parent.Nodes[i].Items[len(node.Parent.Nodes[i].Items)-1].Key
	}

	// 判断是否需要继续分裂
	if len(node.Parent.Nodes) > b.width {
		b.split(node.Parent)
	}
}
