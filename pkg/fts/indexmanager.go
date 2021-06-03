package fts

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type IndexManager struct {
	mu      sync.Mutex        `json:"-"`
	Path    string            `json:"-"`
	Indexes map[string]*Index `json:"indexes"`
}

// NewIndexManager creates a new index manager object
func NewIndexManager(path string) *IndexManager {

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return &IndexManager{Path: path, Indexes: make(map[string]*Index)}
	}

	jsonFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var indexManager IndexManager
	json.Unmarshal(byteValue, &indexManager)

	indexManager.Path = path

	return &indexManager
}

// Save saves the indexes to persistent storage
func (indexManager *IndexManager) Save() error {
	bytes, err := json.Marshal(indexManager)
	if err != nil {
		log.Println(err)
		return err
	}

	jsonFile, err := os.Create(indexManager.Path)
	if err != nil {
		log.Println(err)
		return err
	}
	defer jsonFile.Close()

	jsonFile.Write(bytes)

	return nil
}

// AddIndex adds an index
func (indexManager *IndexManager) AddIndex(index *Index) error {
	indexManager.mu.Lock()
	defer indexManager.mu.Unlock()

	indexManager.Indexes[index.Id] = index
	if err := indexManager.Save(); err != nil {
		return err
	}

	// Create parent folder for the index
	if _, err := os.Stat("indexes/" + index.Id); os.IsNotExist(err) {
		os.MkdirAll("indexes/"+index.Id, os.ModePerm)
	}

	return nil
}

// GetIndex returns an index
func (indexManager *IndexManager) GetIndex(indexId string) (*Index, bool) {
	index, ok := indexManager.Indexes[indexId]
	return index, ok
}

// DeleteIndex deletes an index
func (indexManager *IndexManager) DeleteIndex(indexId string) error {
	indexManager.mu.Lock()
	defer indexManager.mu.Unlock()

	index, ok := indexManager.Indexes[indexId]
	if ok {
		index.Destroy()
		delete(indexManager.Indexes, indexId)
	}

	if err := indexManager.Save(); err != nil {
		return err
	}

	return nil
}
