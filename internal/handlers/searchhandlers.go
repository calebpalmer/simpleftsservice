package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/calebpalmer/simpleftsservice/pkg/fts"
	"github.com/gorilla/mux"
)

type SearchHandler struct {
	IndexManager *fts.IndexManager
}

// ServeHTTP is the handler for the indexes entities.
func (sh *SearchHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		sh.getSearchHandler(w, req)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (sh *SearchHandler) getSearchHandler(w http.ResponseWriter, req *http.Request) {
	indexId := mux.Vars(req)["indexId"]

	// get the index
	index, ok := sh.IndexManager.GetIndex(indexId)
	if !ok {
		msg, _ := json.Marshal(map[string]string{"error": "IndexNotFound"})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, string(msg))
		return
	}

	value := req.FormValue("value")
	ids := index.SearchValue(value)

	w.Header().Set("Content-Type", "application/json")

	if req.FormValue("documents") == "y" {
		// get documents instead of ids
		docs := []interface{}{}
		for _, id := range ids {
			doc, ok := index.GetDocument(id)
			if !ok {
				log.Printf("Document %s does not exist.", id)
				continue
			}
			docJson, err := doc.Json()
			if err != nil {
				writeInternalServerError(w, err)
			}
			docs = append(docs, docJson)
		}
		response, err := json.Marshal(map[string]interface{}{"results": docs})
		if err != nil {
			writeInternalServerError(w, err)
			return
		}
		fmt.Fprint(w, string(response))
	} else {
		response, _ := json.Marshal(map[string][]string{"results": ids})
		fmt.Fprint(w, string(response))
	}

}

// RegisterDocumentsesHandlers registers the index handlers.
func RegisterSearchHandlers(router *mux.Router, indexManager *fts.IndexManager) error {
	router.Handle("/indexes/{indexId}/search", &SearchHandler{indexManager}).
		Methods("GET")
	return nil
}
