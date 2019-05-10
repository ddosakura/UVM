package ast

import (
	"sync"
)

// Scope for Interpreter
type Scope interface {
	Reset()
	Put(k string, v Data)
	Get(k string) Data
}

// CommonScope for Interpreter
type CommonScope struct {
	scope map[string]Data
	mutex *sync.Mutex
}

// MakeScope for Common
func MakeScope() Scope {
	return &CommonScope{
		mutex: new(sync.Mutex),
	}
}

// Reset Scope
func (s *CommonScope) Reset() {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	s.reset()
}

func (s *CommonScope) reset() {
	s.scope = make(map[string]Data)
}

// Put Value
func (s *CommonScope) Put(k string, v Data) {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	if s.scope == nil {
		s.reset()
	}
	s.scope[k] = v
}

// Get Value
func (s *CommonScope) Get(k string) Data {
	defer s.mutex.Unlock()
	s.mutex.Lock()
	if s.scope == nil {
		s.reset()
	}
	return s.scope[k]
}
