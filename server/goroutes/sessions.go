package goroutes

import (
	"github.com/xtaci/kcp-go"
	"sync"
)

// Controls the concurrent implementation of each session job

var sessionID = 0

// session struct

var Sessions = &sessions{
	active: map[int]*Session{}, // active session
	mutex:  &sync.RWMutex{},    // read and write lock
}

type Session struct {
	SessionID    int         `json:"session_id"`    // the session id
	Name         string      `json:"name"`          // the session name
	Description  string      `json:"description"`   // the session description
	PersistentID string      `json:"persistent_id"` // the session persistent id
	Debug        bool        `json:"debug"`         // determines whether the session is in use
	Implant      SessionType `json:"implant"`       // implant conn config
}

type SessionInfo struct {
	SessionID    int    `json:"session_id"`    // the session id
	Name         string `json:"name"`          // the session name
	Description  string `json:"description"`   // the session description
	PersistentID string `json:"persistent_id"` // the session persistent id
	Debug        bool   `json:"debug"`         // determines whether the session is in use
	ImplantType  string `json:"implant_type"`  // implant type
	ImplantUser  string `json:"implant_user"`  // implant user
}

type SessionType struct {
	ImplantType        string // implant data return type
	ImplantUser        string
	ImplantResponseURL string
	KcpConn            *kcp.UDPSession
}

// create a background job model

type sessions struct {
	active map[int]*Session // active session
	mutex  *sync.RWMutex    // read and write lock
}

func (s *sessions) All() []*Session {
	// return all active sessions
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var all []*Session
	//	loop through all sessions and add them to the all slice
	for _, session := range s.active {
		all = append(all, session)
	}
	return all
}

func (s *sessions) Add(session *Session) {
	// add a session
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.active[session.SessionID] = session
	// publish a start session event
	EventBroker.Publish(Event{
		Session:   session,
		EventType: "start-session",
	})
}

func (s *sessions) Remove(session *Session) {
	// remove a session
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.active, session.SessionID)
	// publish a stop session event
	EventBroker.Publish(Event{
		Session:   session,
		EventType: "stop-session",
	})
}

func (s *sessions) Get(SessionID int) *Session {
	// if SessionID < 0 not find to return nil
	if SessionID <= 0 {
		return nil
	}
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	// find session id
	return s.active[SessionID]
}

func NextSessionID() int {
	// return a new session id
	NewSessionID := sessionID + 1
	sessionID++
	return NewSessionID
}
