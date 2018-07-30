package mylog

import (
	"math/rand"
	"os"
	"testing"
	"time"
)

func BenchmarkLogger(b *testing.B) {
	nullf, err := os.OpenFile("/dev/null", os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		b.Fatalf("%v", err)
	}

	defer nullf.Close()
	log := New("debug", nullf, 0)
	entry := log.GetEntryWithFields(map[string]interface{}{
		"flag": "entry",
	})

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 设置并发数
	b.SetParallelism(100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			entry.Infof("benchmark test: %v", rand.Intn(b.N))
		}
	})
}
