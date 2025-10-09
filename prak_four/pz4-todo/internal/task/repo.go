// internal/task/repo.go
package task

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("task not found")

// Repo — потокобезопасное файловое хранилище задач (JSON).
type Repo struct {
	mu       sync.RWMutex
	filePath string
}

func NewRepo(filePath string) *Repo {
	_ = os.MkdirAll(filepath.Dir(filePath), 0o755)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_ = os.WriteFile(filePath, []byte("{}\n"), 0o644)
	}
	return &Repo{filePath: filePath}
}

// --- низкоуровневые операции ---

func (r *Repo) readAll() (map[string]Task, error) {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}
	tasks := make(map[string]Task)
	if len(strings.TrimSpace(string(data))) == 0 {
		return tasks, nil
	}
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *Repo) writeAll(tasks map[string]Task) error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')

	dir := filepath.Dir(r.filePath)
	tmp, err := os.CreateTemp(dir, "tasks-*.tmp")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()

	if _, err := tmp.Write(data); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return err
	}
	if err := tmp.Sync(); err != nil {
		_ = tmp.Close()
		_ = os.Remove(tmpName)
		return err
	}
	if err := tmp.Close(); err != nil {
		_ = os.Remove(tmpName)
		return err
	}
	return os.Rename(tmpName, r.filePath)
}

// --- публичный API (совместим с handler.go) ---

// List: фильтр по подстроке title (case-insensitive), сортировка по UpdatedAt (убыв.), пагинация.
func (r *Repo) List(title string, page, limit int) ([]Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	all, err := r.readAll()
	if err != nil {
		return nil, err
	}

	title = strings.TrimSpace(title)
	titleLower := strings.ToLower(title)

	list := make([]Task, 0, len(all))
	for _, t := range all {
		if title == "" || strings.Contains(strings.ToLower(t.Title), titleLower) {
			list = append(list, t)
		}
	}

	sort.SliceStable(list, func(i, j int) bool {
		if list[i].UpdatedAt.Equal(list[j].UpdatedAt) {
			return list[i].CreatedAt.After(list[j].CreatedAt)
		}
		return list[i].UpdatedAt.After(list[j].UpdatedAt)
	})

	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	start := (page - 1) * limit
	if start >= len(list) {
		return []Task{}, nil
	}
	end := start + limit
	if end > len(list) {
		end = len(list)
	}
	return list[start:end], nil
}

func (r *Repo) Get(id string) (*Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	all, err := r.readAll()
	if err != nil {
		return nil, err
	}
	t, ok := all[id]
	if !ok {
		return nil, ErrNotFound
	}
	return &t, nil
}

func (r *Repo) Create(title string) (*Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	all, err := r.readAll()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	t := Task{
		ID:        uuid.NewString(),
		Title:     title,
		Done:      false,
		CreatedAt: now,
		UpdatedAt: now,
	}
	all[t.ID] = t

	if err := r.writeAll(all); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repo) Update(id, title string, done bool) (*Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	all, err := r.readAll()
	if err != nil {
		return nil, err
	}

	t, ok := all[id]
	if !ok {
		return nil, ErrNotFound
	}

	t.Title = title
	t.Done = done
	t.UpdatedAt = time.Now()
	all[id] = t

	if err := r.writeAll(all); err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *Repo) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	all, err := r.readAll()
	if err != nil {
		return err
	}
	if _, ok := all[id]; !ok {
		return ErrNotFound
	}
	delete(all, id)
	return r.writeAll(all)
}
