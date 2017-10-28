package counter

import (
	"sync"
	"testing"
)

func TestGen(t *testing.T) {

	if Incr() != 1 {
		t.Error("Incr() error")
	}

	Reset()

	if Incr() != 1 {
		t.Error("Incr() error")
	}
}

func TestCurrGen(t *testing.T) {

	Reset()
	wg := sync.WaitGroup{}
	num := 10000

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			Incr()
		}()
	}

	wg.Wait()

	n := Get()
	if n != uint64(num) {
		t.Errorf("Incr is not support concurrency: %d != %d", n, uint64(num))
	}

}
