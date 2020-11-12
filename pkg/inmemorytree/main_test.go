package inmemorytree

import (
	"fmt"
	"testing"
)

type TestCase struct {
	RecordList             []Record
	ExpectedRootNodeMatrix int
}

var TestCasesInsert = []TestCase{
	TestCase{
		RecordList:             []Record{record1},
		ExpectedRootNodeMatrix: record1.Metrics[1].Value,
	},
	TestCase{
		RecordList:             []Record{record1, record2, record3},
		ExpectedRootNodeMatrix: record1.Metrics[1].Value + record2.Metrics[1].Value + record3.Metrics[1].Value,
	},
	TestCase{
		RecordList:             []Record{record2, record3},
		ExpectedRootNodeMatrix: record2.Metrics[1].Value + record3.Metrics[1].Value,
	},
}

var record1 = Record{
	Dimensions: []VarKeyValue{
		VarKeyValue{Key: "device", Value: "mobile"},
		VarKeyValue{Key: "country", Value: "IN"},
	},
	Metrics: []MatrixKeyValue{
		MatrixKeyValue{Key: "webreq", Value: 10},
		MatrixKeyValue{Key: "timespent", Value: 30},
	},
}

var record2 = Record{
	Dimensions: []VarKeyValue{
		VarKeyValue{Key: "device", Value: "tablet"},
		VarKeyValue{Key: "country", Value: "USA"},
	},
	Metrics: []MatrixKeyValue{
		MatrixKeyValue{Key: "webreq", Value: 400},
		MatrixKeyValue{Key: "timespent", Value: 200},
	},
}

var record3 = Record{
	Dimensions: []VarKeyValue{
		VarKeyValue{Key: "device", Value: "tablet"},
		VarKeyValue{Key: "country", Value: "IN"},
	},
	Metrics: []MatrixKeyValue{
		MatrixKeyValue{Key: "webreq", Value: 80},
		MatrixKeyValue{Key: "timespent", Value: 100},
	},
}

func TestInsert(t *testing.T) {
	for idx, tc := range TestCasesInsert {
		tree := Tree()
		for _, rc := range tc.RecordList {
			tree.Insert(rc)
		}
		if tree.RootNode.Metrics["timespent"] != tc.ExpectedRootNodeMatrix {
			t.Errorf(fmt.Sprintf("testcase_no : %d, Expected timespent : %d , Actual timespent : %d \n", idx+1, tc.ExpectedRootNodeMatrix, tree.RootNode.Metrics["timespent"]))
		}
	}

}

type TestCaseQry struct {
	InsertOperation []Record
	QueryDimension  Dimension
}

var TestCasesQuery = []TestCaseQry{
	TestCaseQry{
		InsertOperation: []Record{record1},
		QueryDimension:  Dimension{VarKeyValue{Key: "country", Value: "IN"}},
	},
}

func TestQuery(t *testing.T) {
	tc := Tree()
	r := Record{
		Dimensions: []VarKeyValue{
			VarKeyValue{Key: "device", Value: "mobile"},
			VarKeyValue{Key: "country", Value: "IN"},
		},
		Metrics: []MatrixKeyValue{
			MatrixKeyValue{Key: "webreq", Value: 0},
			MatrixKeyValue{Key: "timespent", Value: 0},
		},
	}
	r1 := Record{
		Dimensions: []VarKeyValue{
			VarKeyValue{Key: "device", Value: "web"},
			VarKeyValue{Key: "country", Value: "USA"},
		},
		Metrics: []MatrixKeyValue{
			MatrixKeyValue{Key: "webreq", Value: 100},
			MatrixKeyValue{Key: "timespent", Value: 100},
		},
	}
	r2 := Record{
		Dimensions: []VarKeyValue{
			VarKeyValue{Key: "device", Value: "tablet"},
			VarKeyValue{Key: "country", Value: "USA"},
		},
		Metrics: []MatrixKeyValue{
			MatrixKeyValue{Key: "webreq", Value: 100},
			MatrixKeyValue{Key: "timespent", Value: 100},
		},
	}
	r3 := Record{
		Dimensions: []VarKeyValue{
			VarKeyValue{Key: "device", Value: "tablet"},
			VarKeyValue{Key: "country", Value: "IN"},
		},
		Metrics: []MatrixKeyValue{
			MatrixKeyValue{Key: "webreq", Value: 0},
			MatrixKeyValue{Key: "timespent", Value: 0},
		},
	}
	tc.Insert(r)
	tc.Insert(r1)
	tc.Insert(r2)
	tc.Insert(r3)
	query1 := Dimension{VarKeyValue{Key: "country", Value: "USA"}}
	query2 := Dimension{VarKeyValue{Key: "country", Value: "USA"}}
	query3 := Dimension{VarKeyValue{Key: "country", Value: "IN"}}
	qry := []Dimension{query1, query2, query3}
	for _, q := range qry {
		rs, _ := tc.Query(q)
		if rs.Res[0].Value+rs.Res[1].Value != 0 && rs.Dimension[0].Value == "IN" {
			t.Errorf("Failed test case")
		}
	}
}
