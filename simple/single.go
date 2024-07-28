package simple

import (
	"sync"
)

type Singleton[T any] struct {
	name           string
	provider       Provider[T]
	value          T
	state          State
	decDescriptors []DecDescriptor[T]
	getMutex       sync.Mutex
}

type State int
type Provider[T any] func() T
type Decorator[T any] func(Provider[T]) Provider[T]
type DecDescriptor[T any] struct {
	Decorator Decorator[T]
	Tag       string
}

const (
	StateEmpty State = iota
	StateValue
)

func NewSingleton[T any](name string, provider Provider[T]) Singleton[T] {
	return Singleton[T]{
		name:     name,
		provider: provider,
		state:    StateEmpty,
	}
}

func (s *Singleton[T]) GetState() State {
	return s.state
}

func (s *Singleton[T]) GetName() string {
	return s.name
}

func (s *Singleton[T]) doGet() T {
	provider := s.provider
	for _, desc := range s.decDescriptors {
		provider = desc.Decorator(s.provider)
	}
	return provider()
}

func (s *Singleton[T]) Get() T {
	if s.state == StateEmpty {
		s.getMutex.Lock()
		defer s.getMutex.Unlock()
		if s.state == StateEmpty {
			s.value = s.doGet()
			s.state = StateValue
		} else {
			return s.value
		}
	}
	return s.value
}

func (s *Singleton[T]) ResetToEmpty() {
	s.state = StateEmpty
}

func (s *Singleton[T]) AddDecorator(tag string, decorator Decorator[T]) {
	if s.state != StateEmpty {
		panic("single is not empty. but you can reset")
	}
	s.decDescriptors = append(s.decDescriptors, DecDescriptor[T]{
		Tag:       tag,
		Decorator: decorator,
	})
}

func (s *Singleton[T]) RemoveDecoratorsByTag(tag string) {
	if s.state != StateEmpty {
		panic("single is not empty. but you can reset")
	}
	decs := make([]DecDescriptor[T], 0)
	for _, d := range s.decDescriptors {
		if d.Tag != tag {
			decs = append(decs, d)
		}
	}
	s.decDescriptors = decs
}
