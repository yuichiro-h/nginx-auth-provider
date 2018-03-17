package handler

import (
	"time"

	"github.com/wunderlist/ttlcache"
)

var stateCache = ttlcache.NewCache(3 * time.Minute)
