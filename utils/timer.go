package utils

import (
	"fmt"
	"time"

	"github.com/x14n/evgateway/internal/gateway"
)

func StartSessionCleaner(gw *gateway.Gateway, ttl time.Duration) {
	go func() {
		ticker := time.NewTicker(time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			now := time.Now()
			sessions := gw.ListSessions()
			for _, s := range sessions {
				if now.Sub(s.Lastseen) > ttl {
					fmt.Printf("[session_cleaner] remove expired session: %s", s.ID)
					s.Close()
					gw.RemoveSession(s.ID)
				}
			}
		}
	}()
}
