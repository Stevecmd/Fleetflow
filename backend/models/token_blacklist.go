package models

// TokenBlacklist represents a list of blacklisted tokens

import (
	"sync"
	"time"
)

type TokenBlacklist struct {
	tokens map[string]time.Time
	mu     sync.RWMutex
}

func NewTokenBlacklist() *TokenBlacklist {
	return &TokenBlacklist{
		tokens: make(map[string]time.Time),
	}
}

func (tb *TokenBlacklist) Add(token string) {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.tokens[token] = time.Now().Add(time.Hour * 24)
}

func (tb *TokenBlacklist) IsBlacklisted(token string) bool {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	expiry, exists := tb.tokens[token]
	return exists && time.Now().Before(expiry)
}

func (tb *TokenBlacklist) CleanupExpired() {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	now := time.Now()
	for token, expiry := range tb.tokens {
		if now.After(expiry) {
			delete(tb.tokens, token)
		}
	}
}
