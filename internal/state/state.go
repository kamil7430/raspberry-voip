package state

import "sync"

// State is a state struct providing concurrent-safe getters and setters (using mutexes).
type State struct {
	displayName      string
	displayNameMutex sync.Mutex
}

func (s *State) GetDisplayName() string {
	s.displayNameMutex.Lock()
	defer s.displayNameMutex.Unlock()
	return s.displayName
}

func (s *State) SetDisplayName(displayName string) {
	s.displayNameMutex.Lock()
	defer s.displayNameMutex.Unlock()
	s.displayName = displayName
}
