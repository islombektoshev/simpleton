package simple

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type UserEntity struct {
	Id   int
	Name string
}

type DB interface {
	GetUserById(id int) (UserEntity, error)
}

type DBImpl struct {
}

func (DBImpl) GetUserById(id int) (UserEntity, error) {
	return UserEntity{1, "value from db"}, nil
}

var DbRepo = NewSingleton[DB]("db", func() DB {
	return DBImpl{}
})

type UserService interface {
	GetUser(id int) (UserEntity, error)
}

type UserServiceImpl struct {
}

func (UserServiceImpl) GetUserById(id int) (UserEntity, error) {
	return DbRepo.Get().GetUserById(1)
}

func TestName(t *testing.T) {
	user, err := DbRepo.Get().GetUserById(1)
	assert.Nil(t, err)
	assert.Equal(t, UserEntity{1, "value from db"}, user)
}

func TestSingleton_ResetToEmpty(t *testing.T) {

	providerCount := 0
	provider := func() int {
		providerCount++
		return 42
	}

	s := NewSingleton("MySingleton", provider)

	decoratingCount := 0
	// Add a decorator
	s.AddDecorator("log", func(next Provider[int]) Provider[int] {
		return func() int {
			decoratingCount++
			v := next()
			return v
		}
	})
	assert.Equal(t, 0, providerCount, "Expected not invoked yet")
	assert.Equal(t, 0, decoratingCount, "Expected not invoked yet")

	// Get the value and check it
	value := s.Get()
	assert.Equal(t, 42, value, "Expected value to be 42")
	assert.Equal(t, 1, providerCount, "Expected single invocation")
	assert.Equal(t, 1, decoratingCount, "Expected single invocation")

	value = s.Get()
	assert.Equal(t, 42, value, "Expected value to be 42")
	assert.Equal(t, 1, providerCount, "Expected single invocation")
	assert.Equal(t, 1, decoratingCount, "Expected single invocation")

	// Ensure the state is now StateValue
	assert.Equal(t, StateValue, s.GetState(), "Expected state to be StateValue")

	// Get the value again and check it
	value = s.Get()
	assert.Equal(t, 42, value, "Expected value to be 42")

	// Reset the singleton state
	s.ResetToEmpty()

	// Ensure the state is now StateEmpty
	assert.Equal(t, StateEmpty, s.GetState(), "Expected state to be StateEmpty")

	// Remove the decorator
	s.RemoveDecoratorsByTag("log")

	// Get the value again and check it
	value = s.Get()
	assert.Equal(t, 42, value, "Expected value to be 42 without decorator log")
	assert.Equal(t, 2, providerCount, "Expected twice invocation")
	assert.Equal(t, 1, decoratingCount, "Expected single invocation")
}
