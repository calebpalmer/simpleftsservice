package cache

import (
	"fmt"

	"github.com/bradfitz/gomemcache/memcache"
)

type Cache struct {
	ConnString                string
	Client                    *memcache.Client
	ItemExpirationTimeSeconds int32 `default:60`
}

func NewCache(connString string) *Cache {
	client := memcache.New(connString)
	return &Cache{Client: client, ConnString: connString}
}

func makeKey(indexName string, searchValue string, extra string) string {
	return fmt.Sprintf("%s_%s_%s", indexName, searchValue, extra)
}

func (c *Cache) Get(indexName string, searchValue string, extra string) ([]byte, error) {
	key := makeKey(indexName, searchValue, extra)

	item, err := c.Client.Get(key)
	if err != nil {
		// this a normal event
		if err == memcache.ErrCacheMiss {
			return nil, nil
		}
		return []byte{}, err
	}

	return item.Value, nil
}

func (c *Cache) Add(indexName string, searchValue string, extra string, results []byte) error {
	key := makeKey(indexName, searchValue, extra)

	err := c.Client.Set(&memcache.Item{Key: key, Value: results, Expiration: c.ItemExpirationTimeSeconds})
	if err != nil {
		return err
	}

	return nil
}
