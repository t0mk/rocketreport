package cache

import (
	"time"

	"github.com/jellydator/ttlcache/v3"
)

var Cache *ttlcache.Cache[string, interface{}]

func init() {
	Cache = ttlcache.New(
		ttlcache.WithTTL[string, interface{}](time.Minute),
	)
	go Cache.Start()
}

func Get(key string) (interface{}, bool) {
	item := Cache.Get(key)
	if (item != nil) && (!item.IsExpired()) {
		return item.Value(), true
	}
	return nil, false
}

func Set(key string, value interface{}) {
	Cache.Set(key, value, ttlcache.DefaultTTL)
}

func FloatWrap(key string, refresher func() (float64, error)) func(...interface{}) (interface{}, error) {
	return func(...interface{}) (interface{}, error) {
		item := Cache.Get(key)
		if (item != nil) && (!item.IsExpired()) {
			return item.Value().(float64), nil
		}
		f, err := refresher()
		if err != nil {
			return nil, err
		}
		Cache.Set(key, f, ttlcache.DefaultTTL)
		return f, nil
	}
}

func TimeWrap(key string, refresher func() (time.Time, error)) func(...interface{}) (interface{}, error) {
	return func(...interface{}) (interface{}, error) {
		item := Cache.Get(key)
		if (item != nil) && (!item.IsExpired()) {
			return item.Value().(time.Time), nil
		}
		t, err := refresher()
		if err != nil {
			return nil, err
		}
		Cache.Set(key, t, ttlcache.DefaultTTL)
		return t, nil
	}
}
