package proxymanager

import (
	"math/rand"
	"time"

	"github.com/fsnotify/fsnotify"
)

// NextProxy will navigate the next proxy to use
func (p *ProxyManager) NextProxy() string {
	p.CurrentIndex++
	if p.CurrentIndex > len(p.Proxies)-1 {
		p.CurrentIndex = 0
	}

	proxy := p.Proxies[p.CurrentIndex]

	return proxy
}

// RandomProxy will choose a proxy randomly from the list
func (p *ProxyManager) RandomProxy() string {
	return p.Proxies[rand.Intn(len(p.Proxies))]
}

// Watch proxy file from events
func (p *ProxyManager) Watch() (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return watcher, err
	}

	if err := watcher.Add(p.filepath); err != nil {
		return watcher, err
	}

	return watcher, nil
}

// Reload proxy pool
func (p *ProxyManager) Reload() error {
	p.cleanupSessions(false)
	i := p.CurrentIndex

	p, err := New(p.filepath)
	if err != nil {
		return err
	}
	p.CurrentIndex = i

	return nil
}

func (p *ProxyManager) SessionProxy(sessionId string) string {
	p.cleanupSessions(true)

	p.RLock()
	session, isSessionExist := p.Sessions[sessionId]
	p.RUnlock()
	if isSessionExist {
		updatedSession := &Session{
			Proxy:     session.Proxy,
			Timestamp: time.Now(),
		}

		p.Lock()
		p.Sessions[sessionId] = updatedSession
		p.Unlock()

		return session.Proxy
	} else {
		proxy := p.NextProxy()
		p.Lock()

		newSession := &Session{
			Proxy:     proxy,
			Timestamp: time.Now(),
		}

		p.Sessions[sessionId] = newSession
		p.Unlock()
		return proxy
	}
}

func (p *ProxyManager) cleanupSessions(orphaned bool) {
	p.Lock()
	if orphaned {
		now := time.Now()
		for sessionId, session := range p.Sessions {
			diff := now.Sub(session.Timestamp)
			if diff.Minutes() > 10 {
				delete(p.Sessions, sessionId)
			}
		}
	} else {
		for sessionId := range p.Sessions {
			delete(p.Sessions, sessionId)
		}
	}
	p.Unlock()
}
