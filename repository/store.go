package repository

import (
	"context"
	"fmt"
)

var (
	ErrNotFound = fmt.Errorf("not found")
)

const (
	InMemory = "memory"
	PostGres = "postgres"
)

type Result struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

type Store interface {
	Insert(ctx context.Context, word string) error
	GetByPrefix(ctx context.Context, prefix string) (*Result, error)
	List(ctx context.Context) ([]Result, error)
}

func StoreFactory(kind string) (Store, error) {
	if kind == PostGres {
		return NewPGRepository("postgresql://postgres:secret@db:5432/postgres")
	}
	return NewMemRepository(), nil
}
