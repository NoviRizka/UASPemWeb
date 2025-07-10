package models

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Username string
	Password string
}

var users = map[string]User{
	"admin": {Username: "admin", Password: "password123"},
}

type Session struct {
	Username string
	Expiry   time.Time
}

var sessions = make(map[string]Session)
var sessionsMutex sync.Mutex // Mutex untuk melindungi akses ke peta sesi

func IsAuthenticated(username, password string) bool {
	user, ok := users[username]
	return ok && user.Password == password
}

func CreateSession(username string) (string, error) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	sessionToken := uuid.New().String()
	expiresAt := time.Now().Add(1 * time.Hour)

	sessions[sessionToken] = Session{
		Username: username,
		Expiry:   expiresAt,
	}
	return sessionToken, nil
}

func GetSession(sessionToken string) (*Session, bool) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()

	session, ok := sessions[sessionToken]
	if !ok {
		return nil, false
	}
	if session.Expiry.Before(time.Now()) {
		delete(sessions, sessionToken)
		return nil, false
	}
	return &session, true
}

func DeleteSession(sessionToken string) {
	sessionsMutex.Lock()
	defer sessionsMutex.Unlock()
	delete(sessions, sessionToken)
}

func CleanUpExpiredSessions() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		sessionsMutex.Lock()
		for token, session := range sessions {
			if session.Expiry.Before(time.Now()) {
				delete(sessions, token)
				log.Printf("Sesi kadaluarsa dihapus: %s", token)
			}
		}
		sessionsMutex.Unlock()
	}
}
