package cache

import (
	"os"
	"testing"
	"time"

	"beeme/util/mylog"
)

var cnt = 0
var clog = mylog.New("INFO", os.Stderr, mylog.Lfile)

// interval 0, no update
func TestCachePermanent(t *testing.T) {
	cnt = 0 // reset cnt

	c := NewCache("incr", incr, 0, clog)
	for i := 0; i < 3; i++ {
		v, _ := c.Load().(int)
		if v != 1 {
			t.Error("TestCachedPermanent fail")
		}
	}

}

// update each interval time
func TestCachePeriod(t *testing.T) {
	cnt = 0 // reset cnt

	interval := time.Second
	c := NewCache("incr", incr, interval, clog)

	time.Sleep(interval / 2) // sleep half interval to make sure update before each Load
	for i := 1; i < 5; i++ {
		v, _ := c.Load().(int)
		if v != i {
			t.Error("TestCachedPeriod fail")
		}
		time.Sleep(interval)
	}
}

func TestNewCache(t *testing.T) {
	name := "aa"
	c := NewCache(name, nil, 0, clog)
	if c != nil {
		t.Error("update is nil, NewCache should return nil")
	}

	update := func() interface{} {
		return nil
	}
	c = NewCache(name, update, 0, clog)
	if c != nil {
		t.Error("update return nil, NewCache should return nil")
	}
}

func TestGetName(t *testing.T) {
	name := "aa"
	c := NewCache(name, incr, 0, clog)
	if c.GetName() != name {
		t.Error("GetName fail")
	}
}

func TestGetInterval(t *testing.T) {
	name := "aa"
	interval := time.Minute
	c := NewCache(name, incr, interval, clog)
	if c.GetInterval() != interval {
		t.Errorf("GetInterval fail")
	}

}

func incr() interface{} {
	cnt++
	return cnt
}
