package storage

import (
	"github.com/SudoQ/crisp/item"
	"errors"
)

type Store struct {
	items []*item.Item
	current uint
	size uint
}

func New() *Store {
	var default_size uint = 5
	return &Store{
		items: make([]*item.Item, 5),
		current: (default_size-1),
		size: default_size,
	}
}

func (this *Store) next() uint {
	return (this.current + 1)%this.size
}

func (this *Store) Add(newItem *item.Item) {
	this.items[this.next()] = newItem
	this.current = (this.current + 1) % this.size
}

func (this *Store) Get() (*item.Item, error) {
	if len(this.items) == 0 {
		return nil, errors.New("Empty storage")
	}
	return this.items[this.current], nil
}
