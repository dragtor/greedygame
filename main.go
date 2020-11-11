package main

import (
	"encoding/json"
	"errors"
	"sort"
	"sync"
)

type InMemTree struct {
	RootNode *ANode
	mu       sync.Mutex
}

func (t *InMemTree) String() string {
	d, _ := json.MarshalIndent(t, " ", " ")
	return string(d)
}

func Tree() *InMemTree {
	return &InMemTree{
		RootNode: ANodeFactory("root"),
	}
}

const (
	DIM_NONE = iota
	DIM_COUNTRY
	DIM_DEVICE
)

func rankOfDimension(dim string) int {
	switch dim {
	case "country":
		return DIM_COUNTRY
	case "device":
		return DIM_DEVICE
	}
	return 0
}

type ANode struct {
	DimensionType string
	Metrics       map[string]int
	ChildNodeMap  map[string]*ANode
}

func ANodeFactory(dim string) *ANode {
	return &ANode{
		DimensionType: dim,
		Metrics:       make(map[string]int),
		ChildNodeMap:  make(map[string]*ANode),
	}
}

func (n *ANode) Dimension() string {
	return n.DimensionType
}

func (n *ANode) childNodes() []*ANode {
	return nil
}

var (
	ERROR_VALUE_NOT_PRESENT = errors.New("ERROR_VALUE_NOT_PRESENT")
)

func (n *ANode) childNodeOfDimension(dim string) (*ANode, error) {
	if val, present := n.ChildNodeMap[dim]; present {
		return val, nil
	}
	return nil, ERROR_VALUE_NOT_PRESENT
}

type Record struct {
	Dimensions Dimension
	Metrics    Metrics
}

type Dimension []VarKeyValue
type Metrics []MatrixKeyValue

func (d Dimension) Len() int {
	return len(d)
}

func (d Dimension) Less(a, b int) bool {
	return rankOfDimension(d[a].Key) < rankOfDimension(d[b].Key)
}

func (d Dimension) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

type MatrixKeyValue struct {
	Key   string
	Value int
}

type VarKeyValue struct {
	Key   string
	Value string
}

// Insert method is responsible to insert new entry in in-memory tree
// insert operation should take care of mutation in tree
func (t *InMemTree) Insert(r Record) error {
	return t.insert(r)
}

func (n *ANode) incrementMetrics(params Metrics) {
	for _, m := range params {
		n.Metrics[m.Key] += m.Value
	}
}

func (t *InMemTree) insert(r Record) error {
	sort.Sort(r.Dimensions)
	currNode := t.RootNode
	currNode.incrementMetrics(r.Metrics)
	for _, d := range r.Dimensions {
		child, err := currNode.childNodeOfDimension(d.Value)
		if err == ERROR_VALUE_NOT_PRESENT {
			node := ANodeFactory(d.Key)
			currNode.ChildNodeMap[d.Value] = node
			node.incrementMetrics(r.Metrics)
			currNode = node
			continue
		}
		child.incrementMetrics(r.Metrics)
		currNode = child
	}
	return nil
}

type QueryResult struct {
	Dimension Dimension
	Res       Metrics
}

func (t *QueryResult) String() string {
	d, _ := json.MarshalIndent(t, " ", " ")
	return string(d)
}

func metrixstructformat(m map[string]int) Metrics {
	var mtc Metrics
	for k, v := range m {
		m1 := MatrixKeyValue{
			Key:   k,
			Value: v,
		}
		mtc = append(mtc, m1)
	}
	return mtc
}

func queryResultGen(dim Dimension, mtc Metrics) *QueryResult {
	return &QueryResult{
		Dimension: dim,
		Res:       mtc,
	}
}

// Query method is responsible to fetch metrix for given dimension
// Assumption : we will be quering for only set of Dimension
func (t *InMemTree) Query(queryDimSet Dimension) (*QueryResult, error) {
	return t.query(queryDimSet)
}

func (t *InMemTree) query(queryDimSet Dimension) (*QueryResult, error) {
	sort.Sort(queryDimSet)
	currNode := t.RootNode
	for _, dim := range queryDimSet {
		child, err := currNode.childNodeOfDimension(dim.Value)
		if err == ERROR_VALUE_NOT_PRESENT {
			return nil, ERROR_VALUE_NOT_PRESENT
		}
		currNode = child
	}
	return queryResultGen(queryDimSet, metrixstructformat(currNode.Metrics)), nil
}
