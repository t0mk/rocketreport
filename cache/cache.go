package cache

import (
	"sync"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

var Cache *ttlcache.Cache[string, interface{}]

var LoadingKeyedMutex = KeyedMutex{}

func init() {
	Cache = ttlcache.New(
		ttlcache.WithTTL[string, interface{}](time.Minute),
	)
	go Cache.Start()

}

type KeyedMutex struct {
    mutexes sync.Map // Zero value is empty and ready for use
}

func (m *KeyedMutex) Lock(key string) func() {
    value, _ := m.mutexes.LoadOrStore(key, &sync.Mutex{})
    mtx := value.(*sync.Mutex)
    mtx.Lock()

    return func() { mtx.Unlock() }
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

func Interface(key string, refresher func() (interface{}, error)) (interface{}, error) {
	unlock := LoadingKeyedMutex.Lock(key)
	defer unlock()

	item, found := Get(key)
	if found {
		return item, nil
	}
	i, err := refresher()
	if err != nil {
		return nil, err
	}
	Set(key, i)
	return i, nil
} 
	

func Float(key string, refresher func() (float64, error)) (float64, error) {
	unlock := LoadingKeyedMutex.Lock(key)
	defer unlock()

	item, found := Get(key)
	if found {
		return item.(float64), nil
	}
	f, err := refresher()
	if err != nil {
		return 0, err
	}
	Set(key, f)
	return f, nil
}

func FloatWrap(key string, refresher func() (float64, error)) func(...interface{}) (interface{}, error) {
	return func(...interface{}) (interface{}, error) {
		unlock := LoadingKeyedMutex.Lock(key)
		defer unlock()
		item, found := Get(key)
		if found {
			return item, nil
		}
		f, err := refresher()
		if err != nil {
			return nil, err
		}
		Set(key, f)
		return f, nil
	}
}

func TimeWrap(key string, refresher func() (time.Time, error)) func(...interface{}) (interface{}, error) {
	return func(...interface{}) (interface{}, error) {
		unlock := LoadingKeyedMutex.Lock(key)
		defer unlock()
		item, found := Get(key)
		if found {
			return item, nil
		}
		t, err := refresher()
		if err != nil {
			return nil, err
		}
		Set(key, t)
		return t, nil
	}
}

func Time(key string, refresher func() (time.Time, error)) (time.Time, error) {
	unlock := LoadingKeyedMutex.Lock(key)
	defer unlock()

	item, found := Get(key)
	if found {
		return item.(time.Time), nil
	}
	t, err := refresher()
	if err != nil {
		return time.Time{}, err
	}
	Set(key, t)
	return t, nil
}