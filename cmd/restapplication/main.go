package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dragtor/greedygame/pkg/inmemorytree"
	"github.com/gorilla/mux"
)

type App struct {
	InMemTree *inmemorytree.InMemTree
	Carrier   chan inmemorytree.Record
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

type QueryResultHTTPResponse struct {
	Dim     []KVResponse  `json:"dim"`
	Metrics []MKVResponse `json:"metrics"`
}

type KVResponse struct {
	Key string `json:"key"`
	Val string `json:"val"`
}

type MKVResponse struct {
	Key string `json:"key"`
	Val int    `json:"val"`
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

func convertTreeQueryResultToHTTPResponse(qryRes inmemorytree.QueryResult) QueryResultHTTPResponse {
	var qryresult QueryResultHTTPResponse
	for _, qr := range qryRes.Dimension {
		qryresult.Dim = append(qryresult.Dim, KVResponse{Key: qr.Key, Val: qr.Value})
	}
	for _, qr := range qryRes.Res {
		qryresult.Metrics = append(qryresult.Metrics, MKVResponse{Key: qr.Key, Val: qr.Value})
	}
	return qryresult
}

func (t *App) insertData(w http.ResponseWriter, r *http.Request) {
	var p InsertHTTPRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		log.Printf("Error : Failed to decode request data ")
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	rec := convertHTTPRequestToInMemTreeRequest(p)
	log.Printf("Pushed data to channel")
	t.Carrier <- rec
	respondWithJSON(w, http.StatusOK, map[string]interface{}{"status": "success"})
}
func (t *App) getData(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request to fetch query")
	var p QueryHTTPRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		log.Printf("Error : Failed to decode request data ")
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()
	reqdim := convertHTTPQueryRequestToMemTreeQueryRequest(p)
	log.Printf("Retriving data from data store")

	res, err := t.InMemTree.Query(reqdim)
	if err == inmemorytree.ERROR_VALUE_NOT_PRESENT {
		respondWithJSON(w, http.StatusNotFound, map[string]string{"status": "Error: Value Not Present"})
		return
	}
	responseData := convertTreeQueryResultToHTTPResponse(*res)
	respondWithJSON(w, http.StatusOK, responseData)
}

func (t *App) asyncStoreInMemtree() {
	for {
		select {
		case rec := <-t.Carrier:
			t.InMemTree.Insert(rec)
		}
	}
}

func main() {
	log.Printf("Initializing server\n")
	r := mux.NewRouter()
	tree := inmemorytree.Tree()
	carrier := make(chan inmemorytree.Record)
	t := App{InMemTree: tree, Carrier: carrier}
	go t.asyncStoreInMemtree()
	r.HandleFunc("/v1/insert", t.insertData)
	r.HandleFunc("/v1/query", t.getData)
	log.Printf("Server listening on port :8080")
	http.ListenAndServe(":8080", r)
}
