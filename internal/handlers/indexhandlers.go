package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"log"

	"github.com/calebpalmer/simpleftsservice/pkg/fts"
	"github.com/gorilla/mux"
)

// IndexHandler represents the handler for an index.
type IndexHandler struct {
	IndexManager *fts.IndexManager
}

// ServeHTTP is the handler for the indexes entities.
func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		h.getIndexHandler(w, req)
		return
	case http.MethodDelete:
		h.deleteIndexHandler(w, req)
	// case http.MethodPost:
	//	h.postIndexHandler(w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

// getIndexHandler is the handler for getting an index
func (i *IndexHandler) getIndexHandler(w http.ResponseWriter, req *http.Request) {
	indexId := mux.Vars(req)["indexId"]
	index, ok := i.IndexManager.Indexes[indexId]
	if !ok {
		msg, err := json.Marshal(map[string]string{"error": "IndexNotFound"})
		if err != nil {
			log.Fatal(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, string(msg))
		return
	}

	jsonData, err := json.Marshal(index)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, string(jsonData))
}

// deleteIndexHandler is the handler for creating an index
func (i *IndexHandler) deleteIndexHandler(w http.ResponseWriter, req *http.Request) {
	indexId := mux.Vars(req)["indexId"]
	_, ok := i.IndexManager.Indexes[indexId]
	if !ok {
		msg, err := json.Marshal(map[string]string{"error": "IndexNotFound"})
		if err != nil {
			log.Fatal(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, string(msg))
		return
	}

	i.IndexManager.DeleteIndex(indexId)
	w.WriteHeader(http.StatusNoContent)
}

// RegisterIndexesHandlers registers the index handlers.
func RegisterIndexHandlers(router *mux.Router, indexManager *fts.IndexManager) error {

	router.Handle("/indexes/{indexId}", &IndexHandler{indexManager}).Methods("GET", "POST", "PUT", "DELETE")
	return nil
}
