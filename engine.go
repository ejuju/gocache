package gdb

import "sync"

// Engine stores all the lists
type Engine struct {
	mu    sync.RWMutex
	lists map[string]*List
}

// NewEngine instanciates a new engine
func NewEngine() *Engine {
	return &Engine{
		lists: make(map[string]*List),
	}
}
