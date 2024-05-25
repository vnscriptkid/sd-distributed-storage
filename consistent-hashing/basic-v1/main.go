package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"sort"
)

// SimpleHashFunction maps a string to a uint32 value
func HashFunction(value string) uint32 {
	h := sha256.New()
	h.Write([]byte(value))
	// Sum returns the SHA-256 hash of the data (32 bytes - 256 bits)
	hash := h.Sum(nil)

	// Takes the first 4 bytes of the 32-byte hash (hash[:4])
	// Interprets them as a 32-bit unsigned integer in big-endian byte order using binary.BigEndian.Uint32.
	// Range: 0 to 999 (in decimal)
	return binary.BigEndian.Uint32(hash[:4]) % 1000
}

// Ring represents a consistent hashing ring
type Ring struct {
	nodes   []uint32
	nodeMap map[uint32]string
}

// NewRing creates a new Ring
func NewRing() *Ring {
	return &Ring{
		nodeMap: make(map[uint32]string),
	}
}

// AddNode adds a node to the ring
func (r *Ring) AddNode(node string) {
	hash := HashFunction(node)
	r.nodes = append(r.nodes, hash)
	r.nodeMap[hash] = node
	sort.Slice(r.nodes, func(i, j int) bool { return r.nodes[i] < r.nodes[j] })
	fmt.Printf("Node %s added with hash %d\n", node, hash)
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
	ring := NewRing()
	ring.AddNode("Node1")
	ring.AddNode("Node2")
	ring.AddNode("Node3")
	// Node Node1 added with hash 704
	// Node Node2 added with hash 566
	// Node Node3 added with hash 101

	keys := []string{"K1", "K2", "K3", "K4", "K5", "K6", "K7", "K8", "K9", "K10"}

	for _, key := range keys {
		fmt.Printf("Key %s is assigned to Node %s\n", key, ring.GetNode(key))
	}
}
