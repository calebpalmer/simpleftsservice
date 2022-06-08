package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/calebpalmer/simpleftsservice/pkg/fts"
	"github.com/gorilla/mux"
)

// IndexHandler represents the handler for an index.
type DocumentsHandler struct {
	IndexManager *fts.IndexManager
}

// ServeHTTP is the handler for the indexes entities.
func (d *DocumentsHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		d.getDocumentsHandler(w, req)
		return
	case http.MethodPost:
		d.postDocumentsHandler(w, req)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (d *DocumentsHandler) getDocumentsHandler(w http.ResponseWriter, req *http.Request) {
	indexId := mux.Vars(req)["indexId"]

	// get the index
	index, ok := d.IndexManager.GetIndex(indexId)
	if !ok {
		msg, err := json.Marshal(map[string]string{"error": "IndexNotFound"})
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, string(msg))
		return
	}

	documents := make([]interface{}, 0)
	for _, document := range index.Documents {
		docJson, err := document.Json()
		if err != nil {
			log.Println(err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		documents = append(documents, docJson)
	}

	results := map[string][]interface{}{"documents": documents}
	bytes, err := json.Marshal(results)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bytes))
}

func (d *DocumentsHandler) postDocumentsHandler(w http.ResponseWriter, req *http.Request) {
	indexId := mux.Vars(req)["indexId"]

	// get the index
	index, ok := d.IndexManager.GetIndex(indexId)
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

	// get the doc to add from the body
	var body map[string]interface{}
	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		log.Println(err)
		msg := fmt.Sprintf("Error parsing json: %s", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	// add the doc to the index
	var id string
	var doc interface{}

	id, ok = body["id"].(string)
	if !ok {
		id = ""
	}

	doc, ok = body["document"]
	if !ok {
		msg, _ := json.Marshal(map[string]string{"error": "\"document\" property is required."})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, string(msg))
		return
	}

	_, err = index.AddDocument(id, doc.(map[string]interface{}))

	if err != nil {
		msg, _ := json.Marshal(map[string]string{"error": "Internal Server Error"})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, string(msg))
		return
	}

	if err = d.IndexManager.Save(); err != nil {
		msg, _ := json.Marshal(map[string]string{"error": "Internal Server Error"})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, string(msg))
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// IndexHandler represents the handler for an index.
type DocumentHandler struct {
	IndexManager *fts.IndexManager
}

// ServeHTTP is the handler for the indexes entities.
func (d *DocumentHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		d.getDocumentHandler(w, req)
		return
	// case http.MethodPost:
	//	d.postDocumentsHandler(w, req)
	//	return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func (d *DocumentHandler) getDocumentHandler(w http.ResponseWriter, req *http.Request) {
	indexId := mux.Vars(req)["indexId"]

	// get the index
	index, ok := d.IndexManager.GetIndex(indexId)
	if !ok {
		msg, _ := json.Marshal(map[string]string{"error": "IndexNotFound"})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, string(msg))
		return
	}

	documentId := mux.Vars(req)["documentId"]
	document, ok := index.GetDocument(documentId)
	if !ok {
		msg, _ := json.Marshal(map[string]string{"error": "DocumentNotFound"})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, string(msg))
		return
	}

	docJson, err := document.Json()
	if err != nil {
		msg, _ := json.Marshal(map[string]string{"error": "InternalServerError"})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, string(msg))
		return
	}

	bytes, err := json.Marshal(docJson)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(bytes))
}

func (d *DocumentHandler) putDocumentHandler(w http.ResponseWriter, req *http.Request) {
	indexId := mux.Vars(req)["indexId"]

	// get the index
	index, ok := d.IndexManager.GetIndex(indexId)
	if !ok {
		msg, _ := json.Marshal(map[string]string{"error": "IndexNotFound"})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, string(msg))
		return
	}

	documentId := mux.Vars(req)["documentId"]
	document, ok := index.GetDocument(documentId)
	if !ok {
		msg, _ := json.Marshal(map[string]string{"error": "DocumentNotFound"})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, string(msg))
		return
	}

	// get the doc to add from the body
	var newDocument fts.DocumentJson
	err := json.NewDecoder(req.Body).Decode(&newDocument)
	if err != nil {
		log.Println(err)
		msg := fmt.Sprintf("Error parsing json: %s", err)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	if newDocument.Id != document.Id {
		http.Error(w, "Document id does not match.", http.StatusBadRequest)
		return
	}
}

// RegisterDocumentsesHandlers registers the index handlers.
func RegisterDocumentsHandlers(router *mux.Router, indexManager *fts.IndexManager) error {

	router.Handle("/indexes/{indexId}/documents", &DocumentsHandler{indexManager}).Methods("GET", "POST")
	router.Handle("/indexes/{indexId}/documents/{documentId}", &DocumentHandler{indexManager}).Methods("GET", "DELETE")
	return nil
}
