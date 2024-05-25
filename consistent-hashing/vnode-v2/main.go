package main

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"sort"
)

// MD5HashFunction maps a string to a uint32 value using the MD5 hash and limits the range to 1000
func MD5HashFunction(value string) uint32 {
	h := md5.New()
	h.Write([]byte(value))
	hash := h.Sum(nil)
	return binary.BigEndian.Uint32(hash[:4]) % 1000
}

// VirtualNode represents a virtual node in the consistent hashing ring
type VirtualNode struct {
	hash uint32
	node string
}

// Ring represents a consistent hashing ring with virtual nodes
type Ring struct {
	virtualNodes []VirtualNode
	nodeMap      map[string]int
}

// NewRing creates a new Ring
func NewRing() *Ring {
	return &Ring{
		nodeMap: make(map[string]int),
	}
}

// AddNode adds a physical node to the ring with a specified number of virtual nodes
func (r *Ring) AddNode(node string, virtualNodeCount int) {
	for i := 0; i < virtualNodeCount; i++ {
		virtualNodeID := fmt.Sprintf("%s#%d", node, i)
		hash := MD5HashFunction(virtualNodeID)
		r.virtualNodes = append(r.virtualNodes, VirtualNode{hash: hash, node: node})
	}
	sort.Slice(r.virtualNodes, func(i, j int) bool { return r.virtualNodes[i].hash < r.virtualNodes[j].hash })
}

// GetNode returns the closest node for a given key
func (r *Ring) GetNode(key string) string {
	hash := MD5HashFunction(key)
	idx := sort.Search(len(r.virtualNodes), func(i int) bool { return r.virtualNodes[i].hash >= hash })
	if idx == len(r.virtualNodes) {
		idx = 0
	}
	return r.virtualNodes[idx].node
}

// PrintNodeRanges prints the ranges each node covers and the total number of values each node covers
func (r *Ring) PrintNodeRanges() {
	for i, vNode := range r.virtualNodes {
		start := vNode.hash
		end := uint32(999)
		if i+1 < len(r.virtualNodes) {
			end = r.virtualNodes[i+1].hash - 1
		}
		r.nodeMap[vNode.node] += int(end - start + 1)
		fmt.Printf("Node %s covers range %d - %d\n", vNode.node, start, end)
	}
	for node, count := range r.nodeMap {
		fmt.Printf("Node %s covers a total of %d values\n", node, count)
	}
}

func main() {
	ring := NewRing()
	ring.AddNode("S1", 20)
	ring.AddNode("S2", 5)
	ring.AddNode("S3", 20)

	ring.PrintNodeRanges()

	keys := []string{"K1", "K2", "K3", "K4", "K5", "K6", "K7", "K8", "K9", "K10"}

	for _, key := range keys {
		fmt.Printf("Key %s is assigned to Node %s\n", key, ring.GetNode(key))
	}
}
