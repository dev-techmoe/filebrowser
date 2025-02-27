package share

import (
	"time"

	"github.com/dev-techmoe/filebrowser/v2/errors"
)

// StorageBackend is the interface to implement for a share storage.
type StorageBackend interface {
	GetByHash(hash string) (*Link, error)
	GetPermanent(path string, id uint) (*Link, error)
	Gets(path string, id uint) ([]*Link, error)
	Save(s *Link) error
	Delete(hash string) error
}

// Storage is a storage.
type Storage struct {
	back StorageBackend
}

// NewStorage creates a share links storage from a backend.
func NewStorage(back StorageBackend) *Storage {
	return &Storage{back: back}
}

// GetByHash wraps a StorageBackend.GetByHash.
func (s *Storage) GetByHash(hash string) (*Link, error) {
	link, err := s.back.GetByHash(hash)
	if err != nil {
		return nil, err
	}

	if link.Expire != 0 && link.Expire <= time.Now().Unix() {
		s.Delete(link.Hash)
		return nil, errors.ErrNotExist
	}

	return link, nil
}

// GetPermanent wraps a StorageBackend.GetPermanent
func (s *Storage) GetPermanent(path string, id uint) (*Link, error) {
	return s.back.GetPermanent(path, id)
}

// Gets wraps a StorageBackend.Gets
func (s *Storage) Gets(path string, id uint) ([]*Link, error) {
	links, err := s.back.Gets(path, id)

	if err != nil {
		return nil, err
	}

	for i, link := range links {
		if link.Expire != 0 && link.Expire <= time.Now().Unix() {
			s.Delete(link.Hash)
			links = append(links[:i], links[i+1:]...)
		}
	}

	return links, nil
}

// Save wraps a StorageBackend.Save
func (s *Storage) Save(l *Link) error {
	return s.back.Save(l)
}

// Delete wraps a StorageBackend.Delete
func (s *Storage) Delete(hash string) error {
	return s.back.Delete(hash)
}
