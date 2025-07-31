package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"

	"dijibill/database"
)

// Session represents a user session
type Session struct {
	ID        string    `json:"id"`
	UserID    int       `json:"user_id"`
	CompanyID int       `json:"company_id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

// SessionManager manages user sessions
type SessionManager struct {
	sessions map[string]*Session
	mutex    sync.RWMutex
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		sessions: make(map[string]*Session),
	}
	
	// Start cleanup goroutine
	go sm.cleanupExpiredSessions()
	
	return sm
}

// CreateSession creates a new session for a user
func (sm *SessionManager) CreateSession(user *database.User) (*Session, error) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Generate session ID
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("error generating session ID: %v", err)
	}

	// Create session
	session := &Session{
		ID:        sessionID,
		UserID:    user.ID,
		CompanyID: user.CompanyID,
		Username:  user.Username,
		Role:      user.Role,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hour session
	}

	sm.sessions[sessionID] = session
	return session, nil
}

// GetSession retrieves a session by ID
func (sm *SessionManager) GetSession(sessionID string) (*Session, bool) {
	sm.mutex.RLock()
	defer sm.mutex.RUnlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return nil, false
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		delete(sm.sessions, sessionID)
		return nil, false
	}

	return session, true
}

// UpdateSessionCompany updates the company for a session (for company switching)
func (sm *SessionManager) UpdateSessionCompany(sessionID string, companyID int) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found")
	}

	session.CompanyID = companyID
	return nil
}

// DeleteSession removes a session
func (sm *SessionManager) DeleteSession(sessionID string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	delete(sm.sessions, sessionID)
}

// ExtendSession extends the expiration time of a session
func (sm *SessionManager) ExtendSession(sessionID string) error {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	session, exists := sm.sessions[sessionID]
	if !exists {
		return fmt.Errorf("session not found")
	}

	session.ExpiresAt = time.Now().Add(24 * time.Hour)
	return nil
}

// cleanupExpiredSessions removes expired sessions periodically
func (sm *SessionManager) cleanupExpiredSessions() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sm.mutex.Lock()
			now := time.Now()
			for sessionID, session := range sm.sessions {
				if now.After(session.ExpiresAt) {
					delete(sm.sessions, sessionID)
				}
			}
			sm.mutex.Unlock()
		}
	}
}

// generateSessionID generates a random session ID
func generateSessionID() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// AuthContext represents the current authentication context
type AuthContext struct {
	SessionID string `json:"session_id"`
	UserID    int    `json:"user_id"`
	CompanyID int    `json:"company_id"`
	Username  string `json:"username"`
	Role      string `json:"role"`
}

// GetAuthContext returns the current authentication context
func (a *App) GetAuthContext() *AuthContext {
	if a.currentSession == nil {
		return nil
	}

	return &AuthContext{
		SessionID: a.currentSession.ID,
		UserID:    a.currentSession.UserID,
		CompanyID: a.currentSession.CompanyID,
		Username:  a.currentSession.Username,
		Role:      a.currentSession.Role,
	}
}