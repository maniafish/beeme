// Package cache provide cache that is updated periodically
package cache

import (
	"sync/atomic"
	"time"

	"beeme/util/mylog"
)

// Cache with update interval, run update method each interval time
type Cache struct {
	atomic.Value                    // 存储cache值
	name         string             // cache名称
	interval     time.Duration      // 更新时间间隔
	update       func() interface{} // 更新方法
	log          mylog.MyLogger
}

// NewCache return an cache instance
// update function is executed each interval time and value return
// will be stored in cache, if interval given is 0, cache will not be updated
func NewCache(name string, update func() interface{}, interval time.Duration, log mylog.MyLogger) *Cache {
	c := &Cache{
		name:     name,
		interval: interval,
		update:   update,
		log:      log,
	}

	if c.update == nil {
		log.Infof("update func is nil")
		return nil
	}

	// first update fail, return nil
	v := c.update()
	if v == nil {
		log.Errorf("%s update fail: nil return", name)
		return nil
	}
	c.Store(v)

	if interval == 0 {
		log.Infof("interval 0, no update")
		return c
	}

	go func() {
		for {
			time.Sleep(c.interval)
			v = c.update()
			if v == nil {
				log.Errorf("%s update fail: nil return", name)
				continue
			}
			c.Store(v)
		}
	}()
	return c
}

// GetName return cache name
func (c *Cache) GetName() string {
	return c.name
}

// GetInterval return cache interval
func (c *Cache) GetInterval() time.Duration {
	return c.interval
}
