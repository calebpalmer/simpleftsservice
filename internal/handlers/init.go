package handlers

import (
	"github.com/calebpalmer/simpleftsservice/pkg/fts"
	"github.com/gorilla/mux"
)

func RegisterHandlers(router *mux.Router) error {
	err := RegisterStatusHandlers(router)
	if err != nil {
		return err
	}

	indexManager := fts.NewIndexManager("indexes.json")
	for _, index := range indexManager.Indexes {
		index.Build()
	}

	err = RegisterIndexesHandlers(router, indexManager)
	if err != nil {
		return err
	}

	err = RegisterIndexHandlers(router, indexManager)
	if err != nil {
		return err
	}

	err = RegisterDocumentsHandlers(router, indexManager)
	if err != nil {
		return err
	}

	err = RegisterSearchHandlers(router, indexManager)
	if err != nil {
		return err
	}

	return nil
}
