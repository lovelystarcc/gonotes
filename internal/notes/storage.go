package notes

import (
	"fmt"
	"sync"
)

type Storage struct {
	mu     sync.RWMutex
	notes  map[int]Note
	nextID int
}

func NewStorage() *Storage {
	return &Storage{
		notes:  make(map[int]Note),
		nextID: 1,
	}
}

func (s *Storage) Create(note Note) Note {
	s.mu.Lock()
	defer s.mu.Unlock()

	note.ID = s.nextID
	s.nextID++
	s.notes[note.ID] = note
	return note
}

func (s *Storage) Get(id int) (*Note, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	n, ok := s.notes[id]
	if !ok {
		return nil, fmt.Errorf("note with id %d not found", id)
	}
	return &n, nil
}

func (s *Storage) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.notes[id]
	if !ok {
		return fmt.Errorf("note with id %d not found", id)
	}

	delete(s.notes, id)
	return nil
}

func (s *Storage) GetAll() []Note {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list := make([]Note, 0, len(s.notes))
	for _, note := range s.notes {
		list = append(list, note)
	}

	return list
}
