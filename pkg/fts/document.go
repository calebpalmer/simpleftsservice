package fts

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type DocumentJson struct {
	Id       string      `json:"id"`
	Document interface{} `json:"document"`
}

type Document struct {
	Id   string `json:"id,omitempty"`
	Path string `json:"path"`
}

// Json returns a json encoded Document
func (d *Document) Json() (interface{}, error) {
	if _, err := os.Stat(d.Path); os.IsNotExist(err) {
		return nil, err
	}

	jsonFile, err := os.Open(d.Path)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	bytes, _ := ioutil.ReadAll(jsonFile)

	var parsedJson interface{}
	err = json.Unmarshal(bytes, &parsedJson)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]interface{})
	ret["id"] = d.Id
	ret["contents"] = parsedJson
	return ret, nil
}

// Destroy deletes the persisted document
func (d *Document) destroy() error {
	if err := os.Remove(d.Path); err != nil {
		return err
	}
	return nil
}
