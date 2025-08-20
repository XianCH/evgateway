package gateway

import (
	"net"
	"sync"
	"time"
)

type Session struct {
	ID         string
	Addr       string
	Lastseen   time.Time
	Conn       net.Conn
	ConnClosed bool
	mu         sync.Mutex
}

func (s *Session) UpdateLastSeen() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Lastseen = time.Now()
}

func (s *Session) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ConnClosed {
		return nil
	}
	if err := s.Conn.Close(); err != nil {
		return err
	}
	s.ConnClosed = true
	return nil
}
