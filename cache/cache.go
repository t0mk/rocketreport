package cache

import (
	"github.com/jellydator/ttlcache/v3"
)

var Cache = ttlcache.New[string, interface{}]()