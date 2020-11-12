package inmemorytree

import (
	"fmt"
	"testing"
)

func TestInsert(t *testing.T) {
	tc := Tree()
	r := Record{
		Dimensions: []VarKeyValue{
			VarKeyValue{Key: "device", Value: "mobile"},
			VarKeyValue{Key: "country", Value: "IN"},
		},
		Metrics: []MatrixKeyValue{
			MatrixKeyValue{Key: "webreq", Value: 70},
			MatrixKeyValue{Key: "timespent", Value: 30},
		},
	}
	r1 := Record{
		Dimensions: []VarKeyValue{
			VarKeyValue{Key: "device", Value: "web"},
			VarKeyValue{Key: "country", Value: "USA"},
		},
		Metrics: []MatrixKeyValue{
			MatrixKeyValue{Key: "webreq", Value: 900},
			MatrixKeyValue{Key: "timespent", Value: 30},
		},
	}
	r2 := Record{
		Dimensions: []VarKeyValue{
			VarKeyValue{Key: "device", Value: "tablet"},
			VarKeyValue{Key: "country", Value: "USA"},
		},
		Metrics: []MatrixKeyValue{
			MatrixKeyValue{Key: "webreq", Value: 70000},
			MatrixKeyValue{Key: "timespent", Value: 3000},
		},
	}
	r3 := Record{
		Dimensions: []VarKeyValue{
			VarKeyValue{Key: "device", Value: "tablet"},
			VarKeyValue{Key: "country", Value: "IN"},
		},
		Metrics: []MatrixKeyValue{
			MatrixKeyValue{Key: "webreq", Value: 52349},
			MatrixKeyValue{Key: "timespent", Value: 299},
		},
	}
	tc.Insert(r)
	tc.Insert(r1)
	tc.Insert(r2)
	tc.Insert(r3)
}

func TestQuery(t *testing.T) {
	tc := Tree()
	r := Record{
		Dimensions: []VarKeyValue{
			VarKeyValue{Key: "device", Value: "mobile"},
			VarKeyValue{Key: "country", Value: "IN"},
		},
		Metrics: []MatrixKeyValue{
			MatrixKeyValue{Key: "webreq", Value: 70},
			MatrixKeyValue{Key: "timespent", Value: 30},
		},
	}
	r1 := Record{
		Dimensions: []VarKeyValue{
			VarKeyValue{Key: "device", Value: "web"},
			VarKeyValue{Key: "country", Value: "USA"},
		},
		Metrics: []MatrixKeyValue{
			MatrixKeyValue{Key: "webreq", Value: 900},
			MatrixKeyValue{Key: "timespent", Value: 30},
		},
	}
	r2 := Record{
		Dimensions: []VarKeyValue{
			VarKeyValue{Key: "device", Value: "tablet"},
			VarKeyValue{Key: "country", Value: "USA"},
		},
		Metrics: []MatrixKeyValue{
			MatrixKeyValue{Key: "webreq", Value: 70000},
			MatrixKeyValue{Key: "timespent", Value: 3000},
		},
	}
	r3 := Record{
		Dimensions: []VarKeyValue{
			VarKeyValue{Key: "device", Value: "tablet"},
			VarKeyValue{Key: "country", Value: "IN"},
		},
		Metrics: []MatrixKeyValue{
			MatrixKeyValue{Key: "webreq", Value: 52349},
			MatrixKeyValue{Key: "timespent", Value: 299},
		},
	}
	tc.Insert(r)
	tc.Insert(r1)
	tc.Insert(r2)
	tc.Insert(r3)
	query := Dimension{VarKeyValue{Key: "country", Value: "USA"}} //, VarKeyValue{Key: "device", Value: "web"}}
	rs, _ := tc.Query(query)
	fmt.Println(rs)
}
