package handlers

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"

	"github.com/gorilla/mux"

	"github.com/calebpalmer/simpleftsservice/pkg/fts"
)

// IndexesHandler represents the handler for Indexes
type IndexesHandler struct {
	IndexManager *fts.IndexManager
}

// ServeHTTP is the handler for the indexes entities.
func (h *IndexesHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/indexes" {
		http.NotFound(w, req)
		return
	}

	switch req.Method {
	case http.MethodGet:
		h.getIndexes(w, req)
	case http.MethodPost:
		h.postIndexes(w, req)
	default:
		http.Error(w, "Method not allowed.", http.StatusMethodNotAllowed)
	}

}

// getIndexes returns the list of indexes.
func (h *IndexesHandler) getIndexes(w http.ResponseWriter, req *http.Request) {
	bytes, err := json.Marshal(h.IndexManager)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bytes))
}

// postIndexes is the handler for posting and index.
func (h *IndexesHandler) postIndexes(w http.ResponseWriter, req *http.Request) {
	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "Expected json body", http.StatusBadRequest)
		return
	}

	var newIndex fts.Index
	err := json.NewDecoder(req.Body).Decode(&newIndex)
	if err != nil {
		log.Println(err)
		msg := fmt.Sprintf("Error parsing json: %s", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// validate the index
	if err = newIndex.Validate(); err != nil {
		msg := fmt.Sprintf("Invalid index: %s", err)
		log.Println(msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// check to make sure index doesn't already exist
	for _, index := range h.IndexManager.Indexes {
		if index.Id == newIndex.Id {
			http.Error(w, "Index exists.", http.StatusConflict)
			return
		}
	}

	h.IndexManager.AddIndex(&newIndex)
	w.WriteHeader(http.StatusCreated)
}

// RegisterIndexesHandlers registeres the index handlers.
func RegisterIndexesHandlers(router *mux.Router, indexManager *fts.IndexManager) error {

	router.Handle("/indexes", &IndexesHandler{indexManager}).Methods("GET", "POST")
	return nil
}
