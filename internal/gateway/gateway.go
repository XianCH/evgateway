package gateway

import "sync"

type Gateway struct {
	mu      sync.RWMutex
	session map[string]*Session
}

func NewGateway() *Gateway {
	return &Gateway{
		session: make(map[string]*Session),
	}
}

func (g *Gateway) AddSession(s *Session) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.session[s.ID] = s
}

func (g *Gateway) GetSession(id string) (*Session, bool) {
	g.mu.Lock()
	defer g.mu.Unlock()
	s, ok := g.session[id]
	return s, ok
}

func (g *Gateway) RemoveSession(id string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.session, id)
}

func (g *Gateway) ListSessions() []*Session {
	g.mu.RLock()
	defer g.mu.RUnlock()
	out := make([]*Session, 0, len(g.session))
	for _, s := range g.session {
		out = append(out, s)
	}
	return out
}
