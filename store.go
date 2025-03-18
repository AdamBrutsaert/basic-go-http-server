package main

import "sync"

type Item struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type Store struct {
	mu     sync.RWMutex
	nextId int
	items  map[int]Item
}

func NewStore() *Store {
	return &Store{
		mu:     sync.RWMutex{},
		items:  map[int]Item{},
		nextId: 1,
	}
}

func (s *Store) GetItems() []Item {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := []Item{}
	for _, item := range s.items {
		items = append(items, item)
	}

	return items
}

func (s *Store) GetItem(id int) (Item, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[id]
	return item, ok
}

func (s *Store) AddItem(item Item) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[s.nextId] = item
	s.nextId++
}

func (s *Store) UpdateItem(id int, item Item) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return false
	}

	s.items[id] = item
	return true
}

func (s *Store) DeleteItem(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return false
	}

	delete(s.items, id)
	return true
}
