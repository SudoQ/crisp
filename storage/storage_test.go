package storage

import (
	"fmt"
	"github.com/SudoQ/satchel/item"
	"testing"
	"time"
)

func TestAdd(t *testing.T) {
	store := New()
	newItem := item.New(time.Now(), []byte("payload"))
	store.Add(newItem)
	if store.items[0] != newItem {
		t.Fail()
	}
}

func TestCurrent(t *testing.T) {
	store := New()
	if store.current != (store.size - 1) {
		t.Fail()
	}
}

func TestGet(t *testing.T) {
	store := New()
	newItem := item.New(time.Now(), []byte("payload"))
	store.Add(newItem)
	currentItem, err := store.Get()
	if err != nil {
		t.Fail()
	}
	if currentItem != newItem {
		t.Fail()
	}
}

func TestAddFull(t *testing.T) {
	store := New()
	times := 10
	for i := 0; i < times; i++ {
		newItem := item.New(time.Now(), []byte(fmt.Sprintf("%d", i)))
		store.Add(newItem)
		sameItem, err := store.Get()
		if err != nil {
			t.Fail()
		}
		if sameItem != newItem {
			t.Fail()
		}
	}
}
