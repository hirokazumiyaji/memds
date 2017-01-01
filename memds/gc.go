package memds

import (
	"context"
	"sync"
	"time"
)

func gc(ctx context.Context, wg *sync.WaitGroup, c *Config) {
	defer wg.Done()

	i, bl := 0, buckets.Len()
	duration := time.Duration(c.GCCycle) * time.Second
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(duration):
			bucket := buckets[i]
			for k, v := range bucket.value {
				value := Value{}
				err := value.Decode(v)
				if err != nil {
					Del(k)
					continue
				}
				if value.IsExpire() {
					Del(k)
				}
			}
			i = (i + 1) % bl
		}
	}
}
