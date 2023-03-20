package main

import (
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
		//找到小于maxKey的key
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
