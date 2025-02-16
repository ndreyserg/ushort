package storage

import (
	"context"
	"testing"
)

func BenchmarkStorage(b *testing.B) {

	urls := make([]string, 1000)
	for i := 0; i < 1000; i++ {
		urls[i] = getUniqKey()
	}

	b.Run("set mem storage", func(b *testing.B) {
		st := NewMemoryStorage()

		for i := 0; i < b.N; i++ {
			for i := 1; i < len(urls); i++ {
				st.Set(context.Background(), urls[i], "tt")
			}
		}

	})

	b.Run("get mm storage", func(b *testing.B) {
		b.StopTimer()
		st := NewMemoryStorage()
		short := make([]string, len(urls))
		for i := 0; i < len(urls); i++ {
			sh, _ := st.Set(context.Background(), urls[i], "tt")
			short[i] = sh
		}
		b.StartTimer()
		for i := 0; i < b.N; i++ {
			for _, sh := range short {
				st.Get(context.Background(), sh)
			}
		}
	})
}
