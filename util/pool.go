package util

import (
	"errors"
	"fmt"
)

type Pool struct {
	Depth   []bool
	MaxSize int
}

func (p *Pool) IsFull() bool {
	if len(p.Depth) >= p.MaxSize {
		return true
	}
	return false
}

// Remove an entry from the end of the depth slice
func (p *Pool) Drain() {
	i := len(p.Depth) - 1
	if i < 0 {
		return
	}

	p.Depth = p.Depth[:i+copy(p.Depth[i:], p.Depth[i+1:])]
}

func (p *Pool) Fill() error {
	if p.IsFull() {
		return errors.New(
			fmt.Sprintf("Pool is full. Maximum depth: ", len(p.Depth)))
	}

	p.Depth = append(p.Depth, true)

	return nil
}

func NewPool(max_size int) *Pool {
	p := new(Pool)
	p.MaxSize = max_size

	return p
}
