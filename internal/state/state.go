package state

import (
	"context"
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"strconv"
	"sync"
)

const (
	minCodeValue = 100_000
	maxCodeValue = 999_999
)

// State is a state struct providing concurrent-safe getters and setters (using mutexes).
type State struct {
	displayName                 string
	displayNameMutex            sync.Mutex
	verificationCode            *string
	verificationCodeMutex       sync.Mutex
	connectionContext           context.Context
	connectionContextCancelFunc context.CancelFunc
	connectionContextMutex      sync.Mutex
}

func NewState() State {
	return State{
		displayName:                 "DefaultUser",
		verificationCode:            nil,
		connectionContext:           nil,
		connectionContextCancelFunc: nil,
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

// CheckAndConsumeVerificationCode checks whether the provided code is correct.
// If it is correct, consumes the verification code from the State struct and returns nil.
// Otherwise, returns an error and the code remains untouched.
func (s *State) CheckAndConsumeVerificationCode(codeToCheck string) error {
	s.verificationCodeMutex.Lock()
	defer s.verificationCodeMutex.Unlock()

	if s.verificationCode == nil {
		return errors.New("verification code is not set")
	}

	if *s.verificationCode != codeToCheck {
		return errors.New("verification code is wrong")
	}

	s.verificationCode = nil
	return nil
}

// CreateVerificationCode returns a verification code. If it was unused, the code is unchanged.
// Otherwise, creates a new (random) code.
func (s *State) CreateVerificationCode() string {
	s.verificationCodeMutex.Lock()
	defer s.verificationCodeMutex.Unlock()

	if s.verificationCode == nil {
		code, err := rand.Int(rand.Reader, big.NewInt(maxCodeValue-minCodeValue+1))
		if err != nil {
			log.Fatal(err)
		}
		codeString := strconv.Itoa(int(code.Int64() + minCodeValue))
		s.verificationCode = &codeString
	}

	return *s.verificationCode
}

// GetConnectionContext should not be used for checking whether the context is nil.
// Use TrySetConnectionContext instead.
func (s *State) GetConnectionContext() context.Context {
	s.connectionContextMutex.Lock()
	defer s.connectionContextMutex.Unlock()

	return s.connectionContext
}

// TrySetConnectionContext atomically checks whether the connection context is nil. If it is, sets the connection
// context to the provided argument. Returns the boolean value that indicates whether the connection context was set.
func (s *State) TrySetConnectionContext(ctx context.Context, cancelFun context.CancelFunc) bool {
	s.connectionContextMutex.Lock()
	defer s.connectionContextMutex.Unlock()

	if s.connectionContext != nil {
		return false
	}

	s.connectionContext = ctx
	s.connectionContextCancelFunc = cancelFun
	return true
}

func (s *State) TerminateConnection() {
	s.connectionContextMutex.Lock()
	defer s.connectionContextMutex.Unlock()

	s.connectionContextCancelFunc()
	s.connectionContext = nil
}
