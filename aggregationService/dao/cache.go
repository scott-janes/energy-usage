package dao

import (
	"context"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/rs/zerolog/log"
	"github.com/scott-janes/energy-usage/shared/types"
)

type CacheValue struct {
	ChildProcesses map[string]string
}

func NewCache(ttl int) *ttlcache.Cache[string, CacheValue] {
	cache := ttlcache.New(
		ttlcache.WithTTL[string, CacheValue](time.Duration(ttl) * time.Minute),
	)

	cache.OnEviction(func(ctx context.Context, reason ttlcache.EvictionReason, item *ttlcache.Item[string, CacheValue]) {
		if reason == ttlcache.EvictionReasonExpired {
			log.Info().Msgf("Evicted %s because it expired", item.Key())
		}
	})

	return cache
}

func ShouldProcessRequest(ctx context.Context, cache *ttlcache.Cache[string, CacheValue], event *types.EnergyServiceEvent) bool {

	exists := cache.Has(event.ID)

	if !exists {
		cacheValue := CacheValue{ChildProcesses: make(map[string]string)}
		for _, process := range event.ChildProcesses {
			cacheValue.ChildProcesses[process] = "PENDING"
		}
		cache.Set(event.ID, cacheValue, ttlcache.DefaultTTL)
		log.Info().Msgf("Added event %s to cache for service %s", event.ID, event.Context.Service)
	}

	cacheValue := cache.Get(event.ID).Value()
	if event.Context.Status == "OK" {
		cacheValue.ChildProcesses[event.Context.Service] = "COMPLETED"
		cache.Set(event.ID, cacheValue, ttlcache.DefaultTTL)
		log.Info().Msgf("Updated event %s to cache for service %s", event.ID, event.Context.Service)
	}

	allCompleted := true
	for name, value := range cacheValue.ChildProcesses {
		log.Info().Msgf("Checking %s %s", name, value)
		if value != "COMPLETED" {
			allCompleted = false
		}
	}

	if allCompleted {
		log.Info().Msgf("All processes completed for event %s", event.ID)

		cache.Delete(event.ID)
		return true
	}
	return false
}
