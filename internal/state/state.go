package state

import (
	"errors"
	"sync"
)

// State is a state struct providing concurrent-safe getters and setters (using mutexes).
type State struct {
	displayName      string
	displayNameMutex sync.Mutex
}

func NewState() State {
	return State{
		displayName: "DefaultUser",
	}
}

func (s *State) GetDisplayName() string {
	s.displayNameMutex.Lock()
	defer s.displayNameMutex.Unlock()

	return s.displayName
}

func (s *State) SetDisplayName(displayName string) error {
	if len(displayName) <= 0 {
		return errors.New("nazwa wyświetlana nie może być pusta")
	}
	if len(displayName) > 16 {
		return errors.New("nazwa wyświetlana musi mieć maksymalnie 16 znaków")
	}

	s.displayNameMutex.Lock()
	defer s.displayNameMutex.Unlock()

	s.displayName = displayName
	return nil
}
