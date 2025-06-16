package service

import (
	"github.com/iamtvk/jsontransformer/internal/models"
	"github.com/patrickmn/go-cache"
	"sync"
	"time"
)

// TODO: implement invalidation of compiled cache
type CacheLayer struct {
	scriptCache   *cache.Cache
	compiledCache *sync.Map
}

func NewCacheLayer() *CacheLayer {
	return &CacheLayer{
		scriptCache:   cache.New(15*time.Minute, 30*time.Minute),
		compiledCache: &sync.Map{},
	}
}

func (c *CacheLayer) GetScript(identifier string) (models.TransformationScript, bool) {
	if script, found := c.scriptCache.Get(identifier); found {
		return script.(models.TransformationScript), found
	}
	return models.TransformationScript{}, false
}
func (c *CacheLayer) SetScript(identifier string, script models.TransformationScript) error {
	return c.scriptCache.Add(identifier, script, cache.DefaultExpiration)
}

func (c *CacheLayer) SetCompiledExpression(scriptHash string, expr any) {
	c.compiledCache.Store(scriptHash, expr)
}

func (c *CacheLayer) GetCompiledExpression(scriptHash string) (any, bool) {
	return c.compiledCache.Load(scriptHash)
}

// TODO: implement invalidation of compiled cache
func (c *CacheLayer) InvalidateCompiledCache() {}
