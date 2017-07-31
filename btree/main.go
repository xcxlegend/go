package main

import (
	"fmt"
)

type TElemType int

/**
 * 二叉树节点
 */
type BiTNode struct {
	data   TElemType
	lchild *BiTNode
	rchild *BiTNode
}

/**
 * 前序遍历
 * @param {[type]} T *BiTNode [description]
 */
func PreOrderTraverse(T *BiTNode) {
	if T == nil {
		return
	}

	fmt.Println(T.data) // 前序先打印 然后遍历子节点 先左后右
	PreOrderTraverse(T.lchild)
	PreOrderTraverse(T.rchild)
}

func CreateBTree(datas []TElemType, T *BiTNode) {
	if len(datas) == 0 {
		return
	}
	var data = datas[0]
	if T == nil {
		T = new(BiTNode)
	}
	T.data = data
	CreateBTree(datas[1:len(datas)-1], T.lchild)
	if len(datas) > 2 {
		CreateBTree(datas[2:len(datas)-1], T.rchild)
	}
}

func main() {
	//Failed to continue: "Cannot find Delve debugger. Install from https://github.com/derekparker/delve & ensure it is in your "GOPATH/bin" or "PATH"."

	fmt.Println("s")
	// var Btree = new(BiTNode)
	// CreateBTree([]TElemType{1, 2, 3, 4, 5, 6, 7, 8}, Btree)
	// fmt.Println(Btree)
	// PreOrderTraverse(Btree)
}
