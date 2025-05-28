package tui

import (
	"fmt"
	"sync"
)

type Leaf[T any] struct {
	mutex *sync.Mutex
	items []T
}

type StringLeaf Leaf[string]

// Functions for generic Leaf

// Creates and return a new empty leaf
func NewLeaf[T any]() *Leaf[T] {
	return &Leaf[T]{
		mutex: &sync.Mutex{},
		items: []T{},
	}
}

// Appends the given element to the leaf
func (leaf *Leaf[T]) Append(item T) {
	leaf.mutex.Lock()
	defer leaf.mutex.Unlock()

	leaf.items = append(leaf.items, item)
}

// Removes all remaining items for leaf
func (leaf *Leaf[T]) Clear() {
	leaf.mutex.Lock()
	defer leaf.mutex.Unlock()

	leaf.items = []T{}
}

// Return a slice of all items in leaf
func (leaf *Leaf[T]) GetAll() []T {
	leaf.mutex.Lock()
	defer leaf.mutex.Unlock()

	return leaf.items
}

// Returns the oldest element and true, then removes it from the leaf.
// Return an empty var of type T and false if leaf has no elements
func (leaf *Leaf[T]) Harvest() (T, bool) {
	leaf.mutex.Lock()
	defer leaf.mutex.Unlock()

	if len(leaf.items) == 0 {
		var zeroValue T
		return zeroValue, false
	}

	harvest := leaf.items[0]
	leaf.items = leaf.items[1:]

	return harvest, true
}

// Returns all items from leaf and clears it.
func (leaf *Leaf[T]) HarvestAll() []T {
	leaf.mutex.Lock()
	defer leaf.mutex.Unlock()

	harvest := leaf.items
	leaf.items = []T{}

	return harvest
}

//
// Functions for StringLeaf
//

// Creates and return a new empty string leaf
func NewStringLeaf() *StringLeaf {
	return &StringLeaf{
		mutex: &sync.Mutex{},
		items: []string{},
	}
}

// Implements fmt.Println()
func (leaf *StringLeaf) Println(a ...any) {
	leaf.mutex.Lock()
	defer leaf.mutex.Unlock()

	leaf.items = append(leaf.items, fmt.Sprint(a...))
}

// Implements fmt.Printf() with a new line after
func (leaf *StringLeaf) Printlnf(format string, a ...any) {
	leaf.mutex.Lock()
	defer leaf.mutex.Unlock()

	leaf.items = append(leaf.items, fmt.Sprintf(format, a...))
}

// Removes all remaining items for leaf
func (leaf *StringLeaf) Clear() {
	leaf.mutex.Lock()
	defer leaf.mutex.Unlock()

	leaf.items = []string{}
}

// Return a slice of all items in leaf
func (leaf *StringLeaf) GetAll() []string {
	leaf.mutex.Lock()
	defer leaf.mutex.Unlock()

	return leaf.items
}

// Returns the oldest element and true, then removes it from the leaf.
// Return an empty string and false if leaf has no elements
func (leaf *StringLeaf) Harvest() (string, bool) {
	leaf.mutex.Lock()
	defer leaf.mutex.Unlock()

	if len(leaf.items) == 0 {
		var zeroValue string
		return zeroValue, false
	}

	harvest := leaf.items[0]
	leaf.items = leaf.items[1:]

	return harvest, true
}

// Returns all items from leaf and clears it.
func (leaf *StringLeaf) HarvestAll() []string {
	leaf.mutex.Lock()
	defer leaf.mutex.Unlock()

	harvest := leaf.items
	leaf.items = []string{}

	return harvest
}
