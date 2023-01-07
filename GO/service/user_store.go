package service

import (
	"fmt"
	"sync"
)

type UserStore interface {
	Save(user *User) error
	Find(userName string) (*User, error)
}

type InMemoryUserStore struct {
	mutex sync.RWMutex
	users map[string]*User
}

func NewInMemoryUserStore() *InMemoryUserStore {
	return &InMemoryUserStore{
		users: make(map[string]*User),
	}
}

func (store *InMemoryUserStore) Save(user *User) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.users[user.UserName] != nil {
		return fmt.Errorf("user %s already exists", user.UserName)
	}

	store.users[user.UserName] = user
	// log.Printf("user %s saved", user.UserName)
	return nil
}

func (store *InMemoryUserStore) Find(userName string) (*User, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	user := store.users[userName]
	if user == nil {
		return nil, fmt.Errorf("user %s not found", userName)
	}

	return user.Clone(), nil
}
