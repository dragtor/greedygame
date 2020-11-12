package main

import (
	"encoding/json"
	"net/http"

	"github.com/dragtor/greedygame/pkg/inmemorytree"
	"github.com/gorilla/mux"
)

type App struct {
	InMemTree *inmemorytree.InMemTree
}

type InsertHTTPRequest struct {
	Dim []struct {
		Key string `json:"key"`
		Val string `json:"val"`
	} `json:"dim"`
	Metrics []struct {
		Key string `json:"key"`
		Val int    `json:"val"`
	} `json:"metrics"`
}

type QueryHTTPRequest struct {
	Dim []struct {
		Key string `json:"key"`
		Val string `json:"val"`
	} `json:"dim"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func convertHTTPRequestToInMemTreeRequest(req InsertHTTPRequest) inmemorytree.Record {
	var record inmemorytree.Record
	for _, r := range req.Dim {
		record.Dimensions = append(record.Dimensions, inmemorytree.VarKeyValue{r.Key, r.Val})
	}
	for _, r := range req.Metrics {
		record.Metrics = append(record.Metrics, inmemorytree.MatrixKeyValue{r.Key, r.Val})
	}
	return record
}

func convertHTTPQueryRequestToMemTreeQueryRequest(req QueryHTTPRequest) inmemorytree.Dimension {
	var dimension inmemorytree.Dimension
	for _, r := range req.Dim {
		dimension = append(dimension, inmemorytree.VarKeyValue{r.Key, r.Val})
	}
	return dimension
}

func (t *App) insertData(w http.ResponseWriter, r *http.Request) {
	var p InsertHTTPRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	rec := convertHTTPRequestToInMemTreeRequest(p)
	t.InMemTree.Insert(rec)
	respondWithJSON(w, http.StatusOK, map[string]interface{}{"status": "success"})
}
func (t *App) getData(w http.ResponseWriter, r *http.Request) {
	var p QueryHTTPRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	reqdim := convertHTTPQueryRequestToMemTreeQueryRequest(p)
	res, err := t.InMemTree.Query(reqdim)
	if err != nil {
		panic(err)
	}
	respondWithJSON(w, http.StatusOK, res)
}

func main() {
	r := mux.NewRouter()
	tree := inmemorytree.Tree()
	t := App{InMemTree: tree}
	r.HandleFunc("/v1/insert", t.insertData)
	r.HandleFunc("/v1/query", t.getData)
	http.ListenAndServe(":8091", r)
}
