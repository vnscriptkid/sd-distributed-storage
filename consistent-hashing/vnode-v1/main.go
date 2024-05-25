package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sort"
)

// HashFunction to map values to uint32
func HashFunction(value string) uint32 {
	h := sha256.New()
	h.Write([]byte(value))
	hash := h.Sum(nil)
	return binary.BigEndian.Uint32(hash[:4])
}

// Ring represents a consistent hashing ring
type Ring struct {
	nodes   []uint32
	nodeMap map[uint32]string
	vnodes  int
}

// NewRing creates a new Ring
func NewRing(vnodes int) *Ring {
	return &Ring{
		nodeMap: make(map[uint32]string),
		vnodes:  vnodes,
	}
}

// AddNode adds a node to the ring
func (r *Ring) AddNode(node string) {
	for i := 0; i < r.vnodes; i++ {
		virtualNode := fmt.Sprintf("%s-%d", node, i)
		hash := HashFunction(virtualNode)
		r.nodes = append(r.nodes, hash)
		r.nodeMap[hash] = node
	}
	sort.Slice(r.nodes, func(i, j int) bool { return r.nodes[i] < r.nodes[j] })
}

// GetNode returns the closest node for a given key
func (r *Ring) GetNode(key string) string {
	hash := HashFunction(key)
	idx := sort.Search(len(r.nodes), func(i int) bool { return r.nodes[i] >= hash })
	if idx == len(r.nodes) {
		idx = 0
	}
	return r.nodeMap[r.nodes[idx]]
}

func main() {
	ring := NewRing(3)
	ring.AddNode("S1")
	ring.AddNode("S2")
	ring.AddNode("S3")

	keys := []string{"K1", "K2", "K3"}

	for _, key := range keys {
		fmt.Printf("Key %s is assigned to Node %s\n", key, ring.GetNode(key))
	}
}
