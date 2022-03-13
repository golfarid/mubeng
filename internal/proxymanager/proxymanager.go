package proxymanager

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"

	"ktbs.dev/mubeng/pkg/mubeng"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Session struct {
	Proxy     string
	Timestamp time.Time
}

// ProxyManager defines the proxy list and current proxy position
type ProxyManager struct {
	sync.RWMutex
	Proxies      []string
	CurrentIndex int
	Sessions     map[string]*Session
}

// New initialize ProxyManager
func New(filename string) (*ProxyManager, error) {
	keys := make(map[string]bool)

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	manager := &ProxyManager{}
	manager.CurrentIndex = -1
	manager.Sessions = make(map[string]*Session)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxy := scanner.Text()
		if _, value := keys[proxy]; !value {
			if _, err = mubeng.Transport(proxy); err == nil {
				keys[proxy] = true
				manager.Proxies = append(manager.Proxies, proxy)
			}
		}
	}

	if len(manager.Proxies) < 1 {
		return manager, fmt.Errorf("open %s: has no valid proxy URLs", filename)
	}

	return manager, scanner.Err()
}

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

func (p *ProxyManager) SessionProxy(sessionId string) string {
	p.cleanupOrphanedSessions()

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

func (p *ProxyManager) cleanupOrphanedSessions() {
	now := time.Now()
	p.Lock()
	for sessionId, session := range p.Sessions {
		diff := now.Sub(session.Timestamp)
		if diff.Minutes() > 1 {
			delete(p.Sessions, sessionId)
		}
	}
	p.Unlock()
}
