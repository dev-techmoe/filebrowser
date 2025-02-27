package users

import (
	"sync"
	"time"

	"github.com/dev-techmoe/filebrowser/v2/errors"
)

// StorageBackend is the interface to implement for a users storage.
type StorageBackend interface {
	GetBy(interface{}) (*User, error)
	Gets() ([]*User, error)
	Save(u *User) error
	Update(u *User, fields ...string) error
	DeleteByID(uint) error
	DeleteByUsername(string) error
}

// Storage is a users storage.
type Storage struct {
	back    StorageBackend
	updated map[uint]int64
	mux     sync.RWMutex
}

// NewStorage creates a users storage from a backend.
func NewStorage(back StorageBackend) *Storage {
	return &Storage{
		back:    back,
		updated: map[uint]int64{},
	}
}

// Get allows you to get a user by its name or username. The provided
// id must be a string for username lookup or a uint for id lookup. If id
// is neither, a ErrInvalidDataType will be returned.
func (s *Storage) Get(baseScope string, id interface{}) (user *User, err error) {
	user, err = s.back.GetBy(id)
	if err != nil {
		return
	}
	user.Clean(baseScope)
	return
}

// Gets gets a list of all users.
func (s *Storage) Gets(baseScope string) ([]*User, error) {
	users, err := s.back.Gets()
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		user.Clean(baseScope)
	}

	return users, err
}

// Update updates a user in the database.
func (s *Storage) Update(user *User, fields ...string) error {
	err := user.Clean("", fields...)
	if err != nil {
		return err
	}

	err = s.back.Update(user, fields...)
	if err != nil {
		return err
	}

	s.mux.Lock()
	s.updated[user.ID] = time.Now().Unix()
	s.mux.Unlock()
	return nil
}

// Save saves the user in a storage.
func (s *Storage) Save(user *User) error {
	if err := user.Clean(""); err != nil {
		return err
	}

	return s.back.Save(user)
}

// Delete allows you to delete a user by its name or username. The provided
// id must be a string for username lookup or a uint for id lookup. If id
// is neither, a ErrInvalidDataType will be returned.
func (s *Storage) Delete(id interface{}) (err error) {
	switch id := id.(type) {
	case string:
		err = s.back.DeleteByUsername(id)
	case uint:
		err = s.back.DeleteByID(id)
	default:
		err = errors.ErrInvalidDataType
	}

	return
}

// LastUpdate gets the timestamp for the last update of an user.
func (s *Storage) LastUpdate(id uint) int64 {
	s.mux.RLock()
	defer s.mux.RUnlock()
	if val, ok := s.updated[id]; ok {
		return val
	}
	return 0
}
