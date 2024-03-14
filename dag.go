package merkledag

import (
	"hash"
)

func Add(store KVStore, node Node, h hash.Hash) []byte {
	// 将分片写入到KVStore中
	writeToStore(store, node)

	// 计算Merkle Root
	root := calculateMerkleRoot(node, h)

	return root
}

func writeToStore(store KVStore, node Node) {
	switch n := node.(type) {
	case File:
		_ = store.Put([]byte("file"), n.Bytes()) // 假设这里使用 "file" 作为键
	case Dir:
		iter := n.It()
		for iter.Next() {
			childNode := iter.Node()
			writeToStore(store, childNode)
		}
	}
}

func calculateMerkleRoot(node Node, h hash.Hash) []byte {
	switch n := node.(type) {
	case File:
		h.Write(n.Bytes())
	case Dir:
		iter := n.It()
		for iter.Next() {
			childNode := iter.Node()
			childHash := calculateMerkleRoot(childNode, h)
			h.Write(childHash)
		}
	}
	return h.Sum(nil)
}
