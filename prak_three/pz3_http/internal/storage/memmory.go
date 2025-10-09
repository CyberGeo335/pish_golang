package storage

import (
	"errors"
	"fmt"
	"sync"
)

type Task struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

type MemoryStore struct {
	mu    sync.RWMutex
	auto  int64
	tasks map[int64]*Task
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks: make(map[int64]*Task),
	}
}

func (s *MemoryStore) Create(title string) *Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.auto++
	t := &Task{ID: s.auto, Title: title, Done: false}
	s.tasks[t.ID] = t
	return t
}

type TaskUpdatePayload struct {
	Done bool `json:"done"`
}

func (s *MemoryStore) Update(id int64, payload TaskUpdatePayload) *Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	t, ok := s.tasks[id]
	if !ok {
		return nil
	}
	fmt.Print(t.ID, t.Done)

	t.Done = payload.Done
	return t
}

func (s *MemoryStore) Delete(id int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.tasks, id)
}

func (s *MemoryStore) Get(id int64) (*Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.tasks[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return t, nil
}

func (s *MemoryStore) List() []*Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		out = append(out, t)
	}
	return out
}
