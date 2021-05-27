package fts

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
)

// Index struct
type Index struct {
	Id               string              `json:"id"`
	SearchProperties []string            `json:"searchProperties"`
	Documents        []Document          `json:"documents,omitempty"`
	InvertedIndex    map[string][]string `json:"-"`
	mu               sync.Mutex          `json:"-"`
}

// MakeIndex initializes and Index
func MakeIndex(name string, searchProperties []string) Index {
	return Index{Id: name, SearchProperties: searchProperties, Documents: make([]Document, 0, 10), InvertedIndex: make(map[string][]string)}
}

func (i *Index) Validate() error {
	if i.Id == "" {
		return errors.New("Index must have Id property.")
	}

	if len(i.SearchProperties) == 0 {
		return errors.New("Index must have searchProperties property.")
	}

	return nil
}

// AddDocument
func (i *Index) indexDocument(docId string, doc map[string]interface{}) error {
	if i.InvertedIndex == nil {
		i.InvertedIndex = make(map[string][]string)
	}
	for _, property := range i.SearchProperties {
		value, ok := doc[property]
		if !ok {
			return errors.New(fmt.Sprintf("Document %s does not have search property %s", docId, property))
		}

		stringValue, ok := value.(string)
		if !ok {
			return errors.New(fmt.Sprintf("Error generating index.  Document %s could not convert propery %s to string", docId, property))
		}

		for _, token := range getFilteredTokens(stringValue) {
			ids := i.InvertedIndex[token]
			found := false
			for _, id := range ids {
				if id == docId {
					found = true
				}
			}
			if !found {
				i.InvertedIndex[token] = append(ids, docId)
			}
		}
	}

	return nil
}

// AddDocument adds a document to the index
func (i *Index) AddDocument(id string, doc map[string]interface{}) (string, error) {
	i.mu.Lock()
	defer i.mu.Unlock()

	// create an id
	if id == "" {
		id = fmt.Sprintf("%s", uuid.New())
	}

	// save the contents to a file
	filePath := fmt.Sprintf("indexes/%s/%s.json", i.Id, id)

	bytes, err := json.Marshal(doc)
	if err != nil {
		panic(err)
	}

	jsonFile, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}

	jsonFile.Write(bytes)

	// add the document to the index
	d := Document{id, filePath}
	i.Documents = append(i.Documents, d)

	// index the document
	if err := i.indexDocument(id, doc); err != nil {
		return id, err
	}

	return id, nil
}

// Build builds the index
func (i *Index) Build() error {
	i.mu.Lock()
	defer i.mu.Unlock()

	for _, document := range i.Documents {
		if _, err := os.Stat(document.Path); os.IsNotExist(err) {
			log.Printf("Error during index generation.  File not found: %s", document.Path)
		}

		jsonFile, err := os.Open(document.Path)
		if err != nil {
			panic(err)
		}

		bytes, _ := ioutil.ReadAll(jsonFile)
		var jsonData interface{}
		json.Unmarshal(bytes, &jsonData)

		jsonMap := jsonData.(map[string]interface{})

		err = i.indexDocument(document.Id, jsonMap)
		if err != nil {
			log.Printf("Error indexing document %s, Error: %s", document.Id, err)
		}
	}

	return nil
}

// Destroy destroys the data assoicated with the index
func (i *Index) Destroy() {
	os.RemoveAll(fmt.Sprintf("indexes/%s", i.Id))
}

// GetDocument gets a document from the index.
func (i *Index) GetDocument(documentId string) (Document, bool) {
	for _, document := range i.Documents {
		if document.Id == documentId {
			return document, true
		}
	}
	return Document{}, false
}

// DeleteDocument deletes a document from the index
func (i *Index) DeleteDocument(documentId string) {
	// TODO this is probably not going to scale.
	i.mu.Lock()
	defer i.mu.Unlock()

	found := false
	foundIndex := 0
	for j, document := range i.Documents {
		if document.Id == documentId {
			found = true
			foundIndex = j
		}

		if found {
			break
		}
	}
	ret := make([]Document, 0)

	ret = append(ret, i.Documents[:foundIndex]...)
	ret = append(ret, i.Documents[foundIndex+1:]...)
	i.Documents = ret

	i.Build()
}

func (i *Index) SearchValue(value string) []string {
	set := make(map[string]struct{})

	for _, token := range getFilteredTokens(value) {
		if ids, ok := i.InvertedIndex[token]; ok {
			for _, id := range ids {
				set[id] = struct{}{}
			}
		}
	}

	r := []string{}
	for k := range set {
		r = append(r, k)
	}

	return r
}
