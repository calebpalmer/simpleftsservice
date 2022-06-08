package handlers

import (
	"github.com/calebpalmer/simpleftsservice/internal/cache"
	"github.com/calebpalmer/simpleftsservice/internal/config"
	"github.com/calebpalmer/simpleftsservice/pkg/fts"
	"github.com/gorilla/mux"
)

func RegisterHandlers(router *mux.Router, config config.Config) error {
	err := RegisterStatusHandlers(router)
	if err != nil {
		return err
	}

	var maybeCache *cache.Cache
	if config.CacheConfig != nil {
		maybeCache = cache.NewCache(config.CacheConfig.ConnString)
	}

	indexManager := fts.NewIndexManager("indexes.json", maybeCache)
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
