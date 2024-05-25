package main

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math"
	"sort"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
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

// PlotRing visualizes the ring with nodes and keys
func PlotRing(ring *Ring, keys []string) {
	p := plot.New()
	p.Title.Text = "Consistent Hashing Ring"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	ringPoints := plotter.XYs{}
	ringCircle := plotter.XYs{}
	labels := plotter.XYs{}
	labelTexts := []string{}

	// Draw the circle
	for i := 0; i < 360; i++ {
		angle := float64(i) * math.Pi / 180.0
		ringCircle = append(ringCircle, plotter.XY{
			X: math.Cos(angle),
			Y: math.Sin(angle),
		})
	}

	// Plot nodes
	for _, node := range ring.nodes {
		angle := float64(node) / float64(math.MaxUint32) * 2.0 * math.Pi
		x := math.Cos(angle)
		y := math.Sin(angle)
		ringPoints = append(ringPoints, plotter.XY{X: x, Y: y})
		labels = append(labels, plotter.XY{X: x, Y: y})
		labelTexts = append(labelTexts, fmt.Sprintf("Node %s", ring.nodeMap[node]))
	}

	// Plot keys
	for _, key := range keys {
		hash := HashFunction(key)
		angle := float64(hash) / float64(math.MaxUint32) * 2.0 * math.Pi
		x := math.Cos(angle)
		y := math.Sin(angle)
		ringPoints = append(ringPoints, plotter.XY{X: x, Y: y})
		labels = append(labels, plotter.XY{X: x, Y: y})
		node := ring.GetNode(key)
		labelTexts = append(labelTexts, fmt.Sprintf("Key %s -> Node %s", key, node))
	}

	// Create scatter plots for the ring and points
	ringPlot, err := plotter.NewScatter(ringPoints)
	if err != nil {
		panic(err)
	}
	ringCirclePlot, err := plotter.NewScatter(ringCircle)
	if err != nil {
		panic(err)
	}

	// Create labels for nodes and keys
	nodeLabels, err := plotter.NewLabels(plotter.XYLabels{
		XYs:    labels,
		Labels: labelTexts,
	})
	if err != nil {
		panic(err)
	}

	// Add the scatter plots and labels to the plot
	p.Add(ringCirclePlot, ringPlot, nodeLabels)

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, 6*vg.Inch, "consistent_hashing_ring.png"); err != nil {
		panic(err)
	}
}

func main() {
	ring := NewRing()
	ring.AddNode("S1")
	ring.AddNode("S2")
	ring.AddNode("S3")

	keys := []string{"K1", "K2", "K3"}

	for _, key := range keys {
		fmt.Printf("Key %s is assigned to Node %s\n", key, ring.GetNode(key))
	}

	PlotRing(ring, keys)
}
