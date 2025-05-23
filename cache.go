package main

import (
	"net/http"
	"reflect"
	"time"

	"go.uber.org/zap"
)

type CacheElement struct {
	Response       http.Response
	RequestHeaders http.Header
	ExpireTime     time.Time
}

type Cache struct {
	store          map[string]CacheElement
	ExpireDuration time.Duration
}

func (c *Cache) Read(path string, headers http.Header) (*http.Response, bool) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()

	val, ok := c.store[path]
	if !ok {
		sugar.Infof("MISS - Path not found: %s", path)
		return nil, false
	} else if !reflect.DeepEqual(val.RequestHeaders, headers) {
		sugar.Info("MISS - Headers do not match")
		return nil, false
	} else if time.Now().After(val.ExpireTime) {
		sugar.Info("MISS - Cached response has expired")
		delete(c.store, path)
		return nil, false
	}
	sugar.Debug("HIT")
	return &val.Response, true
}

func (c *Cache) Write(path string, response *http.Response, requestHeaders http.Header) {
	c.store[path] = CacheElement{
		Response:       *response,
		RequestHeaders: requestHeaders,
		ExpireTime:     time.Now().Add(c.ExpireDuration),
	}
}

func NewCache(expireDuration time.Duration) Cache {
	return Cache{
		store:          make(map[string]CacheElement, 5),
		ExpireDuration: expireDuration,
	}
}
